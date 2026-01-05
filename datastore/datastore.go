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
	"encoding/json"
	"errors"
	"math"

	"cloud.google.com/go/datastore"
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
	ErrNoSuchEntity    = datastore.ErrNoSuchEntity
	ErrEntityExists    = errors.New("entity exists")
	ErrWrongType       = errors.New("wrong type")
	ErrInvalidType     = datastore.ErrInvalidEntityType
	ErrInvalidStoreID  = errors.New("invalid datastore id")
)

// We reuse some Google Datastore types.
type (
	Key      = datastore.Key      // Datastore key.
	Property = datastore.Property // Datastore property.
)

// Entity defines the common interface for our datastore entities.
type Entity interface {
	Copy(dst Entity) (Entity, error) // Copy an entity to dst, or return a copy of the entity when dst is nil.
	GetCache() Cache                 // Returns a cache, or nil for no caching.
}

// EntityEncoder defines an Entity with a custom encoding function.
type EntityEncoder interface {
	Encode() []byte // Encode an entity into bytes.
}

// EntityDecoder defines an Entity with a custom decoding function.
type EntityDecoder interface {
	Decode([]byte) error // Decode bytes into an entity.
}

// CopyEntity copies src into dst (if non-nil) or allocates a new *T.
// It also enforces that dst is of the correct concrete type.
//
// It is a generic helper function to reduce boilerplate code when implementing the Entity interface.
// T is the concrete struct type (e.g. Plan), and PT is a pointer to T (e.g. *Plan).
// PT must implement the Entity interface.
func CopyEntity[T any, PT interface {
	*T
	Entity
}](src PT, dst Entity) (Entity, error) {
	if dst == nil {
		dst = PT(new(T))
	}

	v, ok := dst.(PT)
	if !ok {
		return nil, ErrWrongType
	}

	*v = *src
	return v, nil
}

// encode encodes an entity into bytes, by default using json.Marshal.
func encode(e Entity) []byte {
	encodable, ok := e.(EntityEncoder)
	if ok {
		return encodable.Encode()
	}
	bytes, _ := json.Marshal(e)
	return bytes
}

// decode decodes an entity from bytes, by default using json.Unmarshal.
func decode(e Entity, b []byte) error {
	decodable, ok := e.(EntityDecoder)
	if ok {
		return decodable.Decode(b)
	}

	// Default implementation.
	return json.Unmarshal(b, e)
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
