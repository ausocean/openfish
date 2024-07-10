/*
AUTHORS
  Alan Noble <alan@ausocean.org>
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023, The OpenFish Contributors.

  Redistribution and use in source and binary forms, with or without
  modification, are permitted provided that the following conditions are met:

  1. Redistributions of source code must retain the above copyright notice, this
     list of conditions and the following disclaimer.

  2. Redistributions in binary form must reproduce the above copyright notice,
     this list of conditions and the following disclaimer in the documentation
     and/or other materials provided with the distribution.

  3. Neither the name of The Australian Ocean Lab Ltd. ("AusOcean")
     nor the names of its contributors may be used to endorse or promote
     products derived from this software without specific prior written permission.

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
  SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
  CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
  OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
  OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package datastore offers common datastore API with multiple store implementations:
//
//   - CloudStore is a Google datastore implementation.
//   - FileStore is a simple file-based store implementation.
package datastore

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// CloudStore limits.
const (
	MaxKeys = 500     // Maximum number of datastore keys to retrieve in one go.
	MaxBlob = 1000000 // 1MB
)

// FileStore constants.
const (
	EpochStart  = 1483228800    // Start of the AusOcean epoch, namely 2017-01-01 00:00:00+Z.
	EpochEnd    = math.MaxInt64 // End of the epoch.
	SubTimeBits = 3             // Sub-second time bits.
)

var (
	ErrDecoding        = errors.New("decoding error")
	ErrUnimplemented   = errors.New("unimplemented feature")
	ErrInvalidField    = errors.New("invalid field")
	ErrInvalidFilter   = errors.New("invalid filter")
	ErrInvalidOperator = errors.New("invalid operator")
	ErrInvalidValue    = errors.New("invalid value")
	ErrOperatorMissing = errors.New("operator missing")
	ErrNoSuchEntity    = errors.New("no such entity")
	ErrEntityExists    = errors.New("entity exists")
	ErrWrongType       = errors.New("wrong type")
	ErrInvalidType     = errors.New("invalid type")
)

// We reuse some Google Datastore types.
type (
	Key      = datastore.Key      // Datastore key.
	Property = datastore.Property // Datastore property.
)

// Entity defines the common interface for our datastore entities.
type Entity interface {
	Encode() []byte                  // Encode an entity into bytes.
	Decode([]byte) error             // Decode bytes into an entity.
	Copy(dst Entity) (Entity, error) // Copy an entity to dst, or return a copy of the entity when dst is nil.
	GetCache() Cache                 // Returns a cache, or nil for no caching.
}

// newEntity maps entity type names to their respective constructor function.
// It is populated by RegisterEntity.
var newEntity = map[string]func() Entity{}

// Store defines the datastore interface. It is a blend and subset of
// Google Cloud datastore functions and datastore.Client methods.
//
// See also https://godoc.org/cloud.google.com/go.
type Store interface {
	IDKey(kind string, id int64) *Key                                        // Returns an ID key.
	NameKey(kind, name string) *Key                                          // Returns a name key.
	IncompleteKey(kind string) *Key                                          // Returns an incomplete key.
	NewQuery(kind string, keysOnly bool, keyParts ...string) Query           // Returns a new query.
	Get(ctx context.Context, key *Key, dst Entity) error                     // Gets a single entity by its key.
	GetAll(ctx context.Context, q Query, dst interface{}) ([]*Key, error)    // Runs a query and returns all matching entities.
	Create(ctx context.Context, key *Key, src Entity) error                  // Creates a single entity by its key.
	Put(ctx context.Context, key *Key, src Entity) (*Key, error)             // Put or creates a single entity by its key.
	Update(ctx context.Context, key *Key, fn func(Entity), dst Entity) error // Atomically updates a single entity by its key.
	Delete(ctx context.Context, key *Key) error                              // Deletes a single entity by its key.
	DeleteMulti(ctx context.Context, keys []*Key) error                      // Deletes multiple entities by their keys.
}

// Query defines the query interface, which is a subset of Google
// Cloud's datastore.Query. Note that unlike datastore.Query methods,
// these methods modify the Query value without returning it.
// A nil value denotes a wild card (match anything).
//
// See also Google Cloud datastore.Query.Filter and datastore.Query.Order.
type Query interface {
	Filter(filterStr string, value interface{}) error                       // Filters a query (deprecated).
	FilterField(fieldName string, operator string, value interface{}) error // Filters a query.
	Order(fieldName string)                                                 // Orders a query.
	Limit(limit int)                                                        // Limits the number of results returned.
	Offset(offset int)                                                      // How many keys to skip before returning results.
}

// RegisterEntity registers a new kind of entity and its constructor.
func RegisterEntity(kind string, construct func() Entity) {
	newEntity[kind] = construct
}

// NewEntity instantiates a new entity of the given kind, else returns an error.
func NewEntity(kind string) (Entity, error) {
	construct, ok := newEntity[kind]
	if !ok {
		return nil, ErrInvalidType
	}
	return construct(), nil
}

// NewStore returns a new Store. If kind is "cloud" a CloudStore is
// returned. If kind is "file" a FileStore is returned. The id is the
// Project ID of the requested datastore service. The url parameter is
// a URL to credentials for CloudStore, or a directory path for
// FileStore.
func NewStore(ctx context.Context, kind, id, url string) (s Store, err error) {
	switch kind {
	case "cloud":
		return newCloudStore(ctx, id, url)
	case "file":
		return newFileStore(ctx, id, url)
	default:
		return nil, errors.New("unexpected kind: " + kind)
	}
}

// CloudStore implements Store for the Google Cloud Datastore.
type CloudStore struct {
	client *datastore.Client
}

// newCloudStore returns a new CloudStore, using the given URL to
// retrieve credentials and authenticate. To obtain credentials from a
// Google storage bucket, URL takes the form gs://bucket_name/creds. A
// URL without a scheme is interpreted as a file. If the environment
// variable <ID>_CREDENTIALS is defined it overrides the supplied URL.
func newCloudStore(ctx context.Context, id, url string) (*CloudStore, error) {
	s := new(CloudStore)

	ev := strings.ToUpper(id) + "_CREDENTIALS"
	if os.Getenv(ev) != "" {
		url = os.Getenv(ev)
	}

	var err error
	if url == "" {
		// Attempt authentication using the default credentials.
		s.client, err = datastore.NewClient(ctx, id)
		if err != nil {
			log.Printf("datastore.NewCient failed: %v ", err)
			return nil, err
		}
		return s, nil
	}

	var creds []byte
	if strings.HasPrefix(url, "gs://") {
		// Obtain credentials from a Google Storage bucket.
		url = url[5:]
		sep := strings.IndexByte(url, '/')
		if sep == -1 {
			log.Printf("invalid gs bucket URL: %s", url)
			return nil, errors.New("invalid gs bucket URL")
		}
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Printf("storage.NewCient failed: %v ", err)
			return nil, err
		}
		bkt := client.Bucket(url[:sep])
		obj := bkt.Object(url[sep+1:])
		r, err := obj.NewReader(ctx)
		if err != nil {
			log.Printf("NewReader failed for gs bucket %s: %v", url, err)
			return nil, err
		}
		defer r.Close()
		creds, err = ioutil.ReadAll(r)
		if err != nil {
			log.Printf("cannot read gs bucket %s: %v ", url, err)
			return nil, err
		}

	} else {
		// Interpret url as a file name.
		creds, err = ioutil.ReadFile(url)
		if err != nil {
			log.Printf("cannot read file %s: %v", url, err)
			return nil, err
		}
	}

	s.client, err = datastore.NewClient(ctx, id, option.WithCredentialsJSON(creds))
	return s, err
}

// IDKey returns an ID key given a kind and an int64 ID.
func (s *CloudStore) IDKey(kind string, id int64) *Key {
	return datastore.IDKey(kind, id, nil)
}

// NameKey returns an name key given a kind and a (string) name.
func (s *CloudStore) NameKey(kind, name string) *Key {
	return datastore.NameKey(kind, name, nil)
}

// IncompleteKey returns an incomplete key given a kind.
func (s *CloudStore) IncompleteKey(kind string) *Key {
	return datastore.IncompleteKey(kind, nil)
}

// NewQuery returns a new CloudQuery and is a wrapper for
// datastore.NewQuery. If keysOnly is true the query is set to keys
// only, but keyParts are ignored.
func (s *CloudStore) NewQuery(kind string, keysOnly bool, keyParts ...string) Query {
	q := new(CloudQuery)
	q.query = datastore.NewQuery(kind)
	if keysOnly {
		q.query = q.query.KeysOnly()
	}
	return q
}

func (s *CloudStore) Get(ctx context.Context, key *Key, dst Entity) error {
	if cache := dst.GetCache(); cache != nil {
		err := cache.Get(key, dst)
		if err == nil {
			return nil
		}
	}
	err := s.client.Get(ctx, key, dst)
	if err == datastore.ErrNoSuchEntity {
		return ErrNoSuchEntity
	}
	return err
}

func (s *CloudStore) GetAll(ctx context.Context, query Query, dst interface{}) ([]*Key, error) {
	q, ok := query.(*CloudQuery)
	if !ok {
		return nil, errors.New("expected *CloudQuery type")
	}
	return s.client.GetAll(ctx, q.query, dst)
}

func (s *CloudStore) Create(ctx context.Context, key *Key, src Entity) error {
	_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err := tx.Get(key, src)
		if err == nil {
			return ErrEntityExists
		}
		if err != datastore.ErrNoSuchEntity {
			return err
		}
		_, err = tx.Put(key, src)
		return err
	})
	return err
}

func (s *CloudStore) Put(ctx context.Context, key *Key, src Entity) (*Key, error) {
	key, err := s.client.Put(ctx, key, src)
	if err != nil {
		return key, err
	}
	if cache := src.GetCache(); cache != nil {
		cache.Set(key, src)
	}
	return key, err
}

func (s *CloudStore) Update(ctx context.Context, key *Key, fn func(Entity), dst Entity) error {
	_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err := tx.Get(key, dst)
		if err != nil {
			return err
		}
		fn(dst)
		_, err = tx.Put(key, dst)
		return err
	})
	return err
}

func (s *CloudStore) DeleteMulti(ctx context.Context, keys []*Key) error {
	err := s.client.DeleteMulti(ctx, keys)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if cache := GetCache(k.Kind); cache != nil {
			cache.Delete(k)
		}
	}
	return nil
}

func (s *CloudStore) Delete(ctx context.Context, key *Key) error {
	err := s.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	if cache := GetCache(key.Kind); cache != nil {
		cache.Delete(key)
	}
	return nil
}

// CloudQuery implements Query for the Google Cloud Datastore.
type CloudQuery struct {
	query *datastore.Query
}

func (q *CloudQuery) Filter(filterStr string, value interface{}) error {
	if value == nil {
		return nil
	}
	q.query = q.query.Filter(filterStr, value)
	return nil
}

// FilterField filters a query.
func (q *CloudQuery) FilterField(fieldName string, operator string, value interface{}) error {
	if value == nil {
		return nil
	}
	q.query = q.query.FilterField(fieldName, operator, value)
	return nil
}

func (q *CloudQuery) Order(fieldName string) {
	q.query = q.query.Order(fieldName)
}

// Limit limits the number of results returned.
func (q *CloudQuery) Limit(limit int) {
	q.query = q.query.Limit(limit)
}

// Offset sets the number of keys to skip before returning results.
func (q *CloudQuery) Offset(offset int) {
	q.query = q.query.Offset(offset)
}

// FileStore implements Store for file storage. Each entity is
// represented as a file named <key> under the directory
// <dir>/<id>/<kind>, where <kind> and <key> are the entity kind and
// key respectively.
//
// FileStore implements a simple form of indexing based on the
// period-separated parts of key names. For example, a User has the
// key structure <Skey>.<Email>. Filestore queries that utilize
// indexes must specify the optional key parts when constructing the
// query.
//
//	q = store.NewQuery(ctx, "User", "Skey", "Email")
//
// All but the last part of the key must not contain periods. In this
// example, only Email may contain periods.
//
// The key "671314941988.test@ausocean.org" would match 67131494198 or
// "test@ausocean.org" or both.
//
// To match both Skey and Email:
//
//	q.Filter("Skey =", 671314941988)
//	q.Filter("Email =", "test@ausocean.org)
//
// To match just Skey, i.e., return all entities for a given site,
// simply omit the Email part as follows.
//
//	q.Filter("Skey =", 671314941988)
//
// Alternatively, specify nil for the Email which is a wild card.
//
//	q.Filter("Skey =", 671314941988)
//	q.Filter("Email =", nil)
//
// Using nil is most useful however when ignoring a key part which has
// key parts before and after that much be matched.
//
// To match just Email, i.e., return all entities for a given user.
//
//	q.Filter("Email =", "test@ausocean.org)
//
// The following query would however fail due to the wrong order:
//
//	q.Filter("Email =", "test@ausocean.org)
//	q.Filter("Skey =", 671314941988)
//
// FileStore represents all keys as strings. To faciliate substring
// comparisons, 64-bit ID keys are represented as two 32-bit integers
// separated by a period and the keyParts in NewQuery must reflect
// this.
type FileStore struct {
	mu  sync.Mutex
	id  string
	dir string
}

// newFileStore returns a new FileStore, using the given id and dir to
// create the storage directory, <dir>/<id>. If dir is empty it
// defaults to the current directory.
func newFileStore(ctx context.Context, id, dir string) (*FileStore, error) {
	if dir == "" {
		dir = "."
	}
	store := FileStore{id: id, dir: dir}
	dir = filepath.Join(dir, id)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766)
		if err != nil {
			return &store, err
		}
	}
	return &store, nil
}

// IDKey returns a FileStore ID key for the given kind, setting Name
// to the file name. A 64-bit ID is represented as two 32-bit unsigned
// integers separated by a period.
func (s *FileStore) IDKey(kind string, id int64) *Key {
	dir := filepath.Join(s.dir, s.id, kind)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0766)
	}
	var name string
	id_ := uint64(id)
	if id_ < 1<<32 {
		name = strconv.FormatUint(id_, 10)
	} else {
		name = strconv.FormatUint(id_>>32, 10) + "." + strconv.FormatUint(id_&0xffffffff, 10)
	}
	return &Key{Kind: kind, ID: id, Name: name}
}

// NameKey returns a FileStore name key for the given kind, setting Name to the file name.
func (s *FileStore) NameKey(kind, name string) *Key {
	dir := filepath.Join(s.dir, s.id, kind)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0766)
	}
	return &Key{Kind: kind, Name: name}
}

// IncompleteKey returns an incomplete key given a kind.
func (s *FileStore) IncompleteKey(kind string) *Key {
	// Continue selecting an ID until we find one not used.
	for {
		var name string
		ID := rand.Int63()
		if ID < 1<<32 {
			name = strconv.FormatInt(ID, 10)
		} else {
			name = strconv.FormatInt(ID>>32, 10) + "." + strconv.FormatInt(ID&0xffffffff, 10)
		}
		path := filepath.Join(s.dir, s.id, kind, name)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return &Key{Kind: kind, ID: ID, Name: name}
		}
	}
}

// NewQuery creates and returns a new FileQuery.
func (s *FileStore) NewQuery(kind string, keysOnly bool, keyParts ...string) Query {
	query := FileQuery{kind: kind, keysOnly: keysOnly, keyParts: keyParts, limit: MaxKeys}
	return &query
}

// Get returns a single entity by its key.
func (s *FileStore) Get(ctx context.Context, key *Key, dst Entity) error {
	bytes, err := ioutil.ReadFile(filepath.Join(s.dir, s.id, key.Kind, key.Name))
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoSuchEntity
		}
		return err
	}
	dst.Decode(bytes)
	return nil
}

// GetAll implements Store.GetAll for FileStore. Unlike CloudStore
// queries are only valid against key parts. If the entity has a
// property named Key, it is automatically populated with its Key
// value.
func (s *FileStore) GetAll(ctx context.Context, query Query, dst interface{}) ([]*Key, error) {
	q, ok := query.(*FileQuery)
	if !ok {
		return nil, errors.New("expected *FileQuery type")
	}

	// Construct keys from file names.
	var keys []*Key
	dir := filepath.Join(s.dir, s.id, q.kind)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0766)
		return keys, nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		keys = append(keys, &Key{Kind: q.kind, ID: extractID(f.Name()), Name: f.Name()})
	}

	// Apply filters, if any.
	if q.filter {
		var filtered []*Key
		for _, k := range keys {
			var matches bool
			// Split the key into its parts.
			parts := strings.SplitN(k.Name, ".", len(q.keyParts))
			if len(parts) != len(q.keyParts) {
				continue // No match.
			}
			// Apply filter to each key part.
			matches = true
			for i, part := range parts {
				if q.cmp[i] == nil {
					continue
				}
				// All comparisons must pass in order for the filter to match as a whole.
				for j := range q.cmp[i] {
					if !q.cmp[i][j](part, q.value[i][j]) {
						matches = false
						break
					}
				}
			}
			if matches {
				filtered = append(filtered, k)
			}
		}
		keys = filtered
	}

	// Sort the keys if ordering.
	if q.order {
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].ID != 0 {
				return keys[i].ID < keys[j].ID
			} else {
				return keys[i].Name < keys[j].Name
			}
		})
	}

	// Apply query's limit and offset to the keys.
	if q.offset+q.limit > len(keys) {
		keys = keys[q.offset:]
	} else {
		keys = keys[q.offset:(q.offset + q.limit)]
	}

	if q.keysOnly {
		return keys, nil
	}

	// Get the constructor for this kind of entity.
	entity := newEntity[q.kind]
	if entity == nil {
		return nil, ErrUnimplemented
	}

	// Get each entity and append it to dst.
	dv := reflect.ValueOf(dst).Elem()
	for _, k := range keys {
		e := entity()
		err := s.Get(ctx, k, e)
		if err != nil {
			return nil, err
		}
		// The underlying entity type is a pointer, so we need its indirect type.
		ev := reflect.Indirect(reflect.ValueOf(e))
		// As with the Cloud datastore, if the Key field is present we populate it.
		fld := ev.FieldByName("Key")
		if fld.IsValid() {
			fld.Set(reflect.ValueOf(k))
		}
		dv.Set(reflect.Append(dv, ev))
	}

	return keys, nil
}

// extractID attempts to extract an integer ID from a key name,
// returning zero otherwise. IDs are names comprising either a single
// unsigned integer or two dot-separated 32-bit unsigned integers. The
// latter are recombined into a 64-bit integer.
func extractID(name string) int64 {
	sep := strings.Index(name, ".")
	if sep < 0 {
		n, err := strconv.ParseUint(name, 10, 64)
		if err != nil {
			return 0
		}
		return int64(n)
	}

	msb, err := strconv.ParseUint(name[:sep], 10, 64)
	if err != nil {
		return 0
	}
	lsb, err := strconv.ParseUint(name[sep+1:], 10, 64)
	if err != nil {
		return 0
	}
	return int64(msb<<32 | lsb&0xffffffff)
}

func (s *FileStore) Create(ctx context.Context, key *Key, src Entity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.Get(ctx, key, src)
	if err == nil {
		return ErrEntityExists
	}
	if err != ErrNoSuchEntity {
		return err
	}
	_, err = s.Put(ctx, key, src)
	return err
}

func (s *FileStore) Put(ctx context.Context, key *Key, src Entity) (*Key, error) {
	bytes := src.Encode()
	err := ioutil.WriteFile(filepath.Join(s.dir, s.id, key.Kind, key.Name), bytes, 0777)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *FileStore) Update(ctx context.Context, key *Key, fn func(Entity), dst Entity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.Get(ctx, key, dst)
	if err != nil {
		return err
	}
	fn(dst)
	_, err = s.Put(ctx, key, dst)
	return err
}

func (s *FileStore) Delete(ctx context.Context, key *Key) error {
	return os.Remove(filepath.Join(s.dir, s.id, key.Kind, key.Name))
}

func (s *FileStore) DeleteMulti(ctx context.Context, keys []*Key) error {
	for _, k := range keys {
		err := os.Remove(filepath.Join(s.dir, s.id, k.Kind, k.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

// FileQuery implements Query for FileStore.
type FileQuery struct {
	kind     string                        // Our datastore type.
	keysOnly bool                          // True if this is a keys only query.
	order    bool                          // True if results are ordered.
	filter   bool                          // True if a filter is defined.
	keyParts []string                      // Defines the parts of the key.
	value    [][]string                    // Filter value(s).
	cmp      [][]func(string, string) bool // Filter comparison function(s).
	limit    int                           // Limits the number of results returned.
	offset   int                           // How many keys to skip before returning results.
}

// Filter implements FileQuery matching against key parts.
// This function transforms certain properties commonly used in queries:
//
//   - Properties named ID or MID are reduced to their least-significant 32 bits.
//   - Properties named Timestamp are converted from the Unix epoch to the AusOcean epoch.
func (q *FileQuery) Filter(filterStr string, value interface{}) error {
	if !q.filter {
		q.filter = true
		q.value = make([][]string, len(q.keyParts))
		q.cmp = make([][]func(string, string) bool, len(q.keyParts))
	}

	if value == nil {
		return nil
	}

	// The filter field must match one of the key parts, otherwise it is invalid.
	pos := strings.IndexAny(filterStr, " =<>")
	if pos == -1 {
		return ErrInvalidFilter
	}
	fld := filterStr[:pos]
	idx := -1
	for i, part := range q.keyParts {
		if fld == part {
			idx = i
			break
		}
	}
	if idx == -1 {
		return ErrInvalidField
	}

	// Now extract the comparison operation.
	i := strings.IndexAny(filterStr, "=<>")
	if i == -1 {
		return ErrOperatorMissing
	}

	// We support string values and numeric (int64) values, but both are
	// represented as strings to simplify key/value matching.
	s, ok := value.(string)
	if ok {
		// String value.
		q.value[idx] = append(q.value[idx], s)
		switch filterStr[i:] {
		case "=":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return a == b })
		case "<":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return a > b })
		case ">":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return a > b })
		case "<=":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return a <= b })
		case ">=":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return a >= b })
		default:
			return ErrInvalidOperator
		}
		return nil
	}

	n, ok := value.(int64)
	if ok {
		// Numeric value.
		switch fld {
		case "ID", "MID":
			n &= 0xffffffff
		case "Timestamp":
			n = (n - EpochStart) << SubTimeBits
		}
		q.value[idx] = append(q.value[idx], strconv.FormatInt(n, 10))
		switch filterStr[i:] {
		case "=":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return toInt64(a) == toInt64(b) })
		case "<":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return toInt64(a) < toInt64(b) })
		case ">":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return toInt64(a) > toInt64(b) })
		case "<=":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return toInt64(a) <= toInt64(b) })
		case ">=":
			q.cmp[idx] = append(q.cmp[idx], func(a, b string) bool { return toInt64(a) >= toInt64(b) })
		default:
			return ErrInvalidOperator
		}
		return nil
	}

	return ErrInvalidValue
}

// FilterField filters a query.
func (q *FileQuery) FilterField(fieldName string, operator string, value interface{}) error {
	return q.Filter(fieldName+" "+operator, value)
}

// toInt64 converts a string to an int64, or returns 0.
func toInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func (q *FileQuery) Order(fieldName string) {
	q.order = true
}

// Limit limits the number of results returned.
func (q *FileQuery) Limit(limit int) {
	q.limit = limit
}

// Offset sets the number of keys to skip before returning results.
func (q *FileQuery) Offset(offset int) {
	q.offset = offset
}

// IDKey makes a datastore ID key by combining an ID, a Unix timestamp
// and an optional subtime (st). The latter is used to
// disambiguate otherwise-identical timestamps.
//
//   - The least significant 32 bits of the ID.
//   - The timestamp from the start of the "AusOcean epoch", 2017-01-01 00:00:00+Z (29 bits).
//   - The subtime (3 bits).
func IDKey(id, ts, st int64) int64 {
	ts = ts - EpochStart
	if ts < 0 {
		ts = 0
	}
	return id<<32 | ts<<SubTimeBits | st&((1<<SubTimeBits)-1)
}

// SplitIDKey is the inverse of IDKey and splits an ID key into
// its parts. See IDKey.
func SplitIDKey(id int64) (int64, int64, int64) {
	return int64(uint64(id) >> 32), ((id & 0xffffffff) >> SubTimeBits) + EpochStart, id & ((1 << SubTimeBits) - 1)
}

// DeleteMulti is a wrapper for store.DeleteMulti which returns the
// number of deletions.
func DeleteMulti(ctx context.Context, store Store, keys []*Key) (int, error) {
	n := 0
	for sz := len(keys); sz > 0; sz = len(keys) {
		if sz > MaxKeys {
			sz = MaxKeys
		}
		err := store.DeleteMulti(ctx, keys[:sz])
		if err != nil {
			return n, err
		}
		n += sz
		keys = keys[sz:]
	}
	return n, nil
}

// GetCache returns a cache corresponding to a given kind, or nil otherwise.
func GetCache(kind string) Cache {
	entity := newEntity[kind]
	if entity == nil {
		return nil
	}
	return entity().GetCache()
}
