/*
AUTHORS
  Alan Noble <alan@ausocean.org>
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2026, The OpenFish Contributors.

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

package datastore

import (
	"context"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FileStore implements Store for file storage. Each entity is
// represented as a file named <key> under the directory
// <dir>/<id>/<kind>, where <kind> and <key> are the entity kind and
// key respectively.
//
// FileStore implements a simple form of indexing based on the
// period-separated parts of key names. For example, a User has the
// key structure <Skey>.<Email>. Filestore queries that utilize
// key-based indexes must specify the optional key parts when
// constructing the query.
//
//	q = store.NewQuery(ctx, "User", "Skey", "Email")
//
// All but the last part of the key must not contain periods. In this
// example, only Email may contain periods.
//
// The key "671314941988.test@ausocean.org" would match 671314941988 or
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
// key parts before and after that must be matched.
//
// To match just Email, i.e., return all entities for a given user:
//
//	q.Filter("Email =", "test@ausocean.org")
//
// The following query would however fail due to the wrong order:
//
//	q.Filter("Email =", "test@ausocean.org")
//	q.Filter("Skey =", 671314941988)
//
// FileStore represents all keys as strings. To facilitate substring
// comparisons, 64-bit ID keys are represented as two 32-bit integers
// separated by a period and the keyParts in NewQuery must reflect
// this.
//
// In addition to key-based filtering, FileStore supports content-based
// filtering using FilterField. If the field name matches a key part,
// the query is evaluated efficiently using the file name. Otherwise,
// the entity file is read, parsed, and the field is evaluated using
// reflection. This allows for flexible filtering even on fields not
// encoded in the key.
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
	decode(dst, bytes)
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

		// Apply field filters if needed.
		if len(q.fieldFilters) > 0 && !matchesFieldFilters(e, q.fieldFilters) {
			continue
		}

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

// matchesFieldFilters returns true if the given entity matches all the specified fieldFilters.
// It uses reflection to extract struct fields by name and compares their values based on the
// provided operator (e.g., "=", "<=", ">"). If any filter does not match, it returns false.
func matchesFieldFilters(e Entity, filters []fieldFilter) bool {
	v := reflect.Indirect(reflect.ValueOf(e))
	for _, f := range filters {
		field := v.FieldByName(f.Field)
		if !field.IsValid() {
			return false
		}

		switch f.Operator {
		case "=":
			if field.Interface() != f.Value {
				return false
			}
		case "<", ">", "<=", ">=":
			if !compare(field.Interface(), f.Value, f.Operator) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func compare(a, b interface{}, op string) bool {
	// Try string comparison first.
	as, aok := a.(string)
	bs, bok := b.(string)
	if aok && bok {
		switch op {
		case "=":
			return as == bs
		case "<":
			return as < bs
		case ">":
			return as > bs
		case "<=":
			return as <= bs
		case ">=":
			return as >= bs
		}
		return false
	}

	// Try time.Time comparison.
	at, aok := a.(time.Time)
	bt, bok := b.(time.Time)
	if aok && bok {
		switch op {
		case "=":
			return at.Equal(bt)
		case "<":
			return at.Before(bt)
		case ">":
			return at.After(bt)
		case "<=":
			return at.Before(bt) || at.Equal(bt)
		case ">=":
			return at.After(bt) || at.Equal(bt)
		}
		return false
	}

	// Compare as integers if both values are int-compatible.
	ai, aok := toInt64Strict(a)
	bi, bok := toInt64Strict(b)
	if aok && bok {
		switch op {
		case "=":
			return ai == bi
		case "<":
			return ai < bi
		case ">":
			return ai > bi
		case "<=":
			return ai <= bi
		case ">=":
			return ai >= bi
		}
		return false
	}

	// Fallback to float64 comparison.
	af, aok := toFloat64Safe(a)
	bf, bok := toFloat64Safe(b)
	if !aok || !bok {
		return false
	}
	switch op {
	case "=":
		return af == bf
	case "<":
		return af < bf
	case ">":
		return af > bf
	case "<=":
		return af <= bf
	case ">=":
		return af >= bf
	default:
		return false
	}
}

func toFloat64Safe(x interface{}) (float64, bool) {
	switch v := x.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func toInt64Strict(x interface{}) (int64, bool) {
	switch v := x.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case uint64:
		if v <= math.MaxInt64 {
			return int64(v), true
		}
		return 0, false // Value too large to fit in int64.
	default:
		return 0, false
	}
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
	bytes := encode(src)
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

type fieldFilter struct {
	Field    string
	Operator string
	Value    interface{}
}

// FileQuery implements Query for FileStore.
type FileQuery struct {
	kind         string                        // Our datastore type.
	keysOnly     bool                          // True if this is a keys only query.
	order        bool                          // True if results are ordered.
	filter       bool                          // True if a filter is defined.
	keyParts     []string                      // Defines the parts of the key.
	value        [][]string                    // Filter value(s).
	cmp          [][]func(string, string) bool // Filter comparison function(s).
	limit        int                           // Limits the number of results returned.
	offset       int                           // How many keys to skip before returning results.
	fieldFilters []fieldFilter                 // Filters on entity fields in the file.
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

// FilterField filters a query by field name and value using the specified operator.
// If the field name matches one of the declared keyParts, the query is handled efficiently
// using filename-based key filtering via Filter. Otherwise, the query falls back to
// reflection-based entity field filtering, which reads and inspects each entity file.
func (q *FileQuery) FilterField(fieldName string, operator string, value interface{}) error {
	for _, part := range q.keyParts {
		if part == fieldName {
			return q.Filter(fieldName+" "+operator, value)
		}
	}
	// Fallback to reflection-based filtering.
	q.fieldFilters = append(q.fieldFilters, fieldFilter{
		Field:    fieldName,
		Operator: operator,
		Value:    value,
	})
	return nil
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
