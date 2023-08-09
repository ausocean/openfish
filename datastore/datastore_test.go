/*
AUTHORS
  Alan Noble <alan@ausocean.org>

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

package datastore

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

const typeKeyValue = "KeyValue" // KeyValue datastore type.

// KeyValue represents a key/value pair.
type KeyValue struct {
	Key   string
	Value string
	cache Cache
}

// Encode serializes a KeyValue into tab-separated values.
func (v *KeyValue) Encode() []byte {
	return []byte(fmt.Sprintf("%s\t%s", v.Key, v.Value))
}

// Decode deserializes a KeyValue from tab-separated values.
func (v *KeyValue) Decode(b []byte) error {
	p := strings.Split(string(b), "\t")
	if len(p) != 2 {
		return ErrDecoding
	}
	v.Key = p[0]
	v.Value = p[1]
	return nil
}

// Copy copies another KeyValue entity onto self, or clones self when other is nil.
func (v *KeyValue) Copy(other Entity) (Entity, error) {
	if other == nil {
		kv := new(KeyValue)
		*kv = *v
		return kv, nil
	}
	if kv, ok := other.(*KeyValue); ok {
		*v = *kv
		return v, nil
	}
	return nil, ErrWrongType
}

// GetCache returns the KeyValue cache.
func (v *KeyValue) GetCache() Cache {
	return v.cache
}

// CreateKeyValue creates a KeyValue.
func CreateKeyValue(ctx context.Context, store Store, key, value string) error {
	k := store.NameKey(typeKeyValue, key)
	v := &KeyValue{Key: key, Value: value}
	return store.Create(ctx, k, v)
}

// PutKeyValue creates or updates a KeyValue.
func PutKeyValue(ctx context.Context, store Store, key, value string) error {
	k := store.NameKey(typeKeyValue, key)
	v := &KeyValue{Key: key, Value: value}
	_, err := store.Put(ctx, k, v)
	return err
}

// GetKeyValue gets a KeyValue.
func GetKeyValue(ctx context.Context, store Store, key string) (*KeyValue, error) {
	k := store.NameKey(typeKeyValue, key)
	v := new(KeyValue)
	return v, store.Get(ctx, k, v)
}

// UpdateKeyValue updates a KeyValue by applying the given function.
func UpdateKeyValue(ctx context.Context, store Store, key string, fn func(Entity)) (*KeyValue, error) {
	k := store.NameKey(typeKeyValue, key)
	v := new(KeyValue)
	return v, store.Update(ctx, k, fn, v)
}

// DeleteKeyValue deletes a KeyValue.
func DeleteKeyValue(ctx context.Context, store Store, key string) error {
	k := store.NameKey(typeKeyValue, key)
	return store.DeleteMulti(ctx, []*Key{k})
}

// init registers the KeyValue entitity.
func init() {
	RegisterEntity(typeKeyValue, func() Entity { return new(KeyValue) })
}

// TestFile tests the file store.
func TestFile(t *testing.T) {
	testKeyValue(t, "file", nil)
}

// TestCloud tests the cloud store without caching.
// OPENFISH_CREDENTIALS must be supplied.
func TestCloud(t *testing.T) {
	if os.Getenv("OPENFISH_CREDENTIALS") == "" {
		t.Skip("OPENFISH_CREDENTIALS")
	}
	testKeyValue(t, "cloud", nil)
}

// TestCloudCaching tests the cloud store with caching.
// OPENFISH_CREDENTIALS must be supplied.
func TestCloudCaching(t *testing.T) {
	if os.Getenv("OPENFISH_CREDENTIALS") == "" {
		t.Skip("OPENFISH_CREDENTIALS")
	}
	testKeyValue(t, "cloud", NewEntityCache())
}

// testKeyValue tests various KeyValue methods.
func testKeyValue(t *testing.T, kind string, cache Cache) {
	ctx := context.Background()

	store, err := NewStore(ctx, kind, "openfish", "")
	if err != nil {
		t.Errorf("NewStore(%s, openfish) failed with error: %v", kind, err)
	}

	tests := []struct {
		key   string
		value string
		cache Cache
	}{
		{
			key:   "foo",
			value: "bar1",
			cache: cache,
		},
		{
			key:   "_foo",
			value: "bar2",
			cache: cache,
		},
		{
			key:   "dev.foo",
			value: "bar3",
			cache: cache,
		},
	}

	for i, test := range tests {
		err = PutKeyValue(ctx, store, test.key, test.value)
		if err != nil {
			t.Errorf("PutKeyValue %d failed with error: %v", i, err)
		}
		err = CreateKeyValue(ctx, store, test.key, test.value)
		if err != ErrEntityExists {
			t.Errorf("CreateKeyValue %d failed with unexpected error: %v", i, err)
		}
		v, err := GetKeyValue(ctx, store, test.key)
		if err != nil {
			t.Errorf("GetKeyValue %d failed with error: %v", i, err)
		}
		if v.Value != test.value {
			t.Errorf("GetKeyValue %d returned wrong value; expected %s, got %s", i, test.value, v.Value)
		}
		v, err = UpdateKeyValue(ctx, store, test.key, clearValue)
		if err != nil {
			t.Errorf("UpdateKeyValue %d failed with error: %v", i, err)
		}
		if v.Value != "" {
			t.Errorf("GetKeyValue %d returned wrong value; expected empty string, got %s", i, v.Value)
		}
		err = DeleteKeyValue(ctx, store, test.key)
		if err != nil {
			t.Errorf("DeleteKeyValue %d failed with error: %v", i, err)
		}
	}
}

// clearValue clears the value of a KeyValue.
func clearValue(e Entity) {
	v, ok := e.(*KeyValue)
	if ok {
		v.Value = ""
	}
}
