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
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

const typeNameValue = "NameValue" // NameValue datastore type.
const typeMixed = "Mixed"         // Mixed datastore type.

// NameValue represents a key/value pair.
type NameValue struct {
	Name  string
	Value string
	cache Cache
}

// Encode serializes a NameValue into tab-separated values.
func (v *NameValue) Encode() []byte {
	return []byte(fmt.Sprintf("%s\t%s", v.Name, v.Value))
}

// Decode deserializes a NameValue from tab-separated values.
func (v *NameValue) Decode(b []byte) error {
	p := strings.Split(string(b), "\t")
	if len(p) != 2 {
		return ErrDecoding
	}
	v.Name = p[0]
	v.Value = p[1]
	return nil
}

// Copy copies a NameValue to dst, or returns a copy of the NameValue when dst is nil.
func (v *NameValue) Copy(dst Entity) (Entity, error) {
	var kv *NameValue
	if dst == nil {
		kv = new(NameValue)
	} else {
		var ok bool
		kv, ok = dst.(*NameValue)
		if !ok {
			return nil, ErrWrongType
		}
	}
	*kv = *v
	return kv, nil
}

// GetCache returns the NameValue cache.
func (v *NameValue) GetCache() Cache {
	return v.cache
}

// Mixed represents a mixed entity for testing various comparisons.
type Mixed struct {
	ID      string    // primary key
	Str     string    // test string comparisons
	Int     int64     // test integer comparisons
	Float   float64   // test float comparisons
	Created time.Time // test time comparisons if needed
	cache   Cache
}

// Copy copies a Mixed to dst, or returns a copy of the Mixed when dst is nil.
func (m *Mixed) Copy(dst Entity) (Entity, error) {
	var mix *Mixed
	if dst == nil {
		mix = new(Mixed)
	} else {
		var ok bool
		mix, ok = dst.(*Mixed)
		if !ok {
			return nil, ErrWrongType
		}
	}
	*mix = *m
	return mix, nil
}

// GetCache returns the Mixed cache.
func (m *Mixed) GetCache() Cache {
	return m.cache
}

// CreateNameValue creates a NameValue.
func CreateNameValue(ctx context.Context, store Store, key, value string) error {
	k := store.NameKey(typeNameValue, key)
	v := &NameValue{Name: key, Value: value}
	return store.Create(ctx, k, v)
}

// PutNameValue creates or updates a NameValue.
func PutNameValue(ctx context.Context, store Store, key, value string) error {
	k := store.NameKey(typeNameValue, key)
	v := &NameValue{Name: key, Value: value}
	_, err := store.Put(ctx, k, v)
	return err
}

// GetNameValue gets a NameValue.
func GetNameValue(ctx context.Context, store Store, key string) (*NameValue, error) {
	k := store.NameKey(typeNameValue, key)
	v := new(NameValue)
	return v, store.Get(ctx, k, v)
}

// UpdateNameValue updates a NameValue by applying the given function.
func UpdateNameValue(ctx context.Context, store Store, key string, fn func(Entity)) (*NameValue, error) {
	k := store.NameKey(typeNameValue, key)
	v := new(NameValue)
	return v, store.Update(ctx, k, fn, v)
}

// DeleteNameValue deletes a NameValue.
func DeleteNameValue(ctx context.Context, store Store, key string) error {
	k := store.NameKey(typeNameValue, key)
	return store.DeleteMulti(ctx, []*Key{k})
}

// init registers the NameValue entitity.
func init() {
	RegisterEntity(typeNameValue, func() Entity { return new(NameValue) })
	RegisterEntity(typeMixed, func() Entity { return new(Mixed) })
}

// TestFile tests the file store.
func TestFile(t *testing.T) {
	testNameValue(t, "file", nil)
}

// TestCloud tests the cloud store without caching.
// OPENFISH_CREDENTIALS must be supplied.
func TestCloud(t *testing.T) {
	if os.Getenv("OPENFISH_CREDENTIALS") == "" {
		t.Skip("OPENFISH_CREDENTIALS")
	}
	testNameValue(t, "cloud", nil)
}

// TestCloudCaching tests the cloud store with caching.
// OPENFISH_CREDENTIALS must be supplied.
func TestCloudCaching(t *testing.T) {
	if os.Getenv("OPENFISH_CREDENTIALS") == "" {
		t.Skip("OPENFISH_CREDENTIALS")
	}
	testNameValue(t, "cloud", NewEntityCache())
}

// testNameValue tests various NameValue methods.
func testNameValue(t *testing.T, kind string, cache Cache) {
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
		err = PutNameValue(ctx, store, test.key, test.value)
		if err != nil {
			t.Errorf("PutNameValue %d failed with error: %v", i, err)
		}
		err = CreateNameValue(ctx, store, test.key, test.value)
		if err != ErrEntityExists {
			t.Errorf("CreateNameValue %d failed with unexpected error: %v", i, err)
		}
		v, err := GetNameValue(ctx, store, test.key)
		if err != nil {
			t.Errorf("GetNameValue %d failed with error: %v", i, err)
		}
		if v.Value != test.value {
			t.Errorf("GetNameValue %d returned wrong value; expected %s, got %s", i, test.value, v.Value)
		}
		v, err = UpdateNameValue(ctx, store, test.key, clearValue)
		if err != nil {
			t.Errorf("UpdateNameValue %d failed with error: %v", i, err)
		}
		if v.Value != "" {
			t.Errorf("GetNameValue %d returned wrong value; expected empty string, got %s", i, v.Value)
		}
		err = DeleteNameValue(ctx, store, test.key)
		if err != nil {
			t.Errorf("DeleteNameValue %d failed with error: %v", i, err)
		}
		if test.cache != nil {
			// Check that the value was cleared from the cache.
			k := store.NameKey(typeNameValue, test.key)
			var v NameValue
			err := cache.Get(k, &v)
			if err == nil {
				t.Errorf("cache.Get %d returned no error", i)
			}
			var errCacheMiss ErrCacheMiss
			if !errors.As(err, &errCacheMiss) {
				t.Errorf("cache.Get %d returned wrong error: %v", i, err)
			}
		}
	}
}

// clearValue clears the value of a NameValue.
func clearValue(e Entity) {
	v, ok := e.(*NameValue)
	if ok {
		v.Value = ""
	}
}

// TestFileDirect tests direct file store operations.
func TestFileDirect(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(ctx, "file", "test", "store")
	if err != nil {
		t.Fatalf("could not create file store: %v", err)
	}

	const (
		name1 = "1"
		name2 = "2"
		value = "localuser@localhost"
	)

	// Put two NameValue entities with the same value.
	_, err = store.Put(ctx, store.NameKey(typeNameValue, name1+"."+value), &NameValue{Name: name1, Value: value})
	if err != nil {
		t.Errorf("Put name1 failed: %v", err)
	}
	_, err = store.Put(ctx, store.NameKey(typeNameValue, name2+"."+value), &NameValue{Name: name2, Value: value})
	if err != nil {
		t.Errorf("Put name2 failed: %v", err)
	}

	// GetAll by name, returning key only.
	q := store.NewQuery(typeNameValue, true, "Name", "Value")
	q.Filter("Name =", name1)
	keys, err := store.GetAll(ctx, q, nil)
	if err != nil {
		t.Errorf("GetAll by name, keys only failed: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("GetAll by name returned %d keys, expected 1", len(keys))
	}

	// GetAll by name, returning 1 entity.
	q = store.NewQuery(typeNameValue, false, "Name", "Value")
	q.Filter("Name =", name1)
	var entities []NameValue
	_, err = store.GetAll(ctx, q, &entities)
	if err != nil {
		t.Errorf("GetAll by name, returning entities failed: %v", err)
	}
	if len(entities) != 1 {
		t.Errorf("GetAll by name returned %d entities, expected 1", len(entities))
	}

	// GetAll by value, returning 2 entities.
	q = store.NewQuery(typeNameValue, false, "Name", "Value")
	q.Filter("Value =", value)
	entities = nil
	_, err = store.GetAll(ctx, q, &entities)
	if err != nil {
		t.Errorf("GetAll by value failed: %v", err)
	}
	if len(entities) != 2 {
		t.Errorf("GetAll by value returned %d entities, expected 2", len(entities))
	}

	// GetAll by value, returning entities, limited to 1.
	q = store.NewQuery(typeNameValue, false, "Name", "Value")
	q.Filter("Value =", value)
	q.Limit(1)
	entities = nil
	_, err = store.GetAll(ctx, q, &entities)
	if err != nil {
		t.Errorf("GetAll by value failed: %v", err)
	}
	if len(entities) != 1 {
		t.Errorf("GetAll by value returned %d entities, expected 1", len(entities))
	}
}

// TestFileFilterField tests filestore filtering by a field that is not part of the key.
func TestFileFilterField(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(ctx, "file", "test", "store")
	if err != nil {
		t.Fatalf("could not create file store: %v", err)
	}

	// Clear store dir for clean test (optional in your setup).
	_ = os.RemoveAll("store/test/NameValue")

	// Insert test data.
	values := []NameValue{
		{Name: "A", Value: "abalone"},
		{Name: "B", Value: "bullseye"},
		{Name: "C", Value: "abalone"},
	}
	for _, nv := range values {
		_, err := store.Put(ctx, store.NameKey(typeNameValue, nv.Name), &nv)
		if err != nil {
			t.Fatalf("failed to insert test entity %s: %v", nv.Name, err)
		}
	}

	// Filter by Value using FilterField (not part of key name).
	q := store.NewQuery(typeNameValue, false, "Name") // keyParts only include Name.
	q.FilterField("Value", "=", "abalone")

	var results []NameValue
	_, err = store.GetAll(ctx, q, &results)
	if err != nil {
		t.Fatalf("FilterField query failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results for Value = 'abalone', got %d", len(results))
	}

	// Filter by Name using FilterField (should fall back to Filter for efficiency).
	q = store.NewQuery(typeNameValue, false, "Name")
	q.FilterField("Name", "=", "B")

	results = nil
	_, err = store.GetAll(ctx, q, &results)
	if err != nil {
		t.Fatalf("FilterField fallback-to-key query failed: %v", err)
	}
	if len(results) != 1 || results[0].Name != "B" {
		t.Errorf("expected 1 result for Name = 'B', got %+v", results)
	}
}

func TestFileFilterFieldVariants(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(ctx, "file", "test", "store")
	if err != nil {
		t.Fatalf("could not create file store: %v", err)
	}

	// Reset the directory for a clean run.
	_ = os.RemoveAll("store/test/NameValue")

	// Insert test data.
	data := []NameValue{
		{Name: "a", Value: "apple"},
		{Name: "b", Value: "banana"},
		{Name: "c", Value: "carrot"},
		{Name: "d", Value: "apple"},
		{Name: "e", Value: "eggplant"},
		{Name: "f", Value: ""},
		{Name: "g", Value: "ápple"}, // Unicode variant
		{Name: "h", Value: "appl"},  // Prefix of "apple"
		{Name: "i", Value: "Apple"}, // Capitalized
		{Name: "j", Value: "zucchini"},
	}
	for _, nv := range data {
		_, err := store.Put(ctx, store.NameKey(typeNameValue, nv.Name), &nv)
		if err != nil {
			t.Fatalf("failed to insert test entity %s: %v", nv.Name, err)
		}
	}

	tests := []struct {
		name        string
		field       string
		operator    string
		value       interface{}
		keyParts    []string
		limit       int
		offset      int
		expectNames []string
	}{
		{
			name:        "Value = apple",
			field:       "Value",
			operator:    "=",
			value:       "apple",
			keyParts:    []string{"Name"},
			expectNames: []string{"a", "d"},
		},
		{
			name:        "Value < banana",
			field:       "Value",
			operator:    "<",
			value:       "banana",
			keyParts:    []string{"Name"},
			expectNames: []string{"a", "d", "f", "h", "i"},
		},
		{
			name:        "Value >= carrot",
			field:       "Value",
			operator:    ">=",
			value:       "carrot",
			keyParts:    []string{"Name"},
			expectNames: []string{"c", "e", "g", "j"},
		},
		{
			name:        "Value = '' (empty string)",
			field:       "Value",
			operator:    "=",
			value:       "",
			keyParts:    []string{"Name"},
			expectNames: []string{"f"},
		},
		{
			name:        "Unicode comparison (Value > apple)",
			field:       "Value",
			operator:    ">",
			value:       "apple",
			keyParts:    []string{"Name"},
			expectNames: []string{"b", "c", "e", "g", "j"},
		},
		{
			name:        "Name = h (uses key path)",
			field:       "Name",
			operator:    "=",
			value:       "h",
			keyParts:    []string{"Name"},
			expectNames: []string{"h"},
		},
		{
			name:        "Limit and offset combo",
			field:       "Value",
			operator:    ">",
			value:       "a",
			keyParts:    []string{"Name"},
			limit:       2,
			offset:      1,
			expectNames: []string{"b", "c"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := store.NewQuery(typeNameValue, false, tc.keyParts...)
			q.FilterField(tc.field, tc.operator, tc.value)
			q.Order(tc.field)
			if tc.limit > 0 {
				q.Limit(tc.limit)
			}
			if tc.offset > 0 {
				q.Offset(tc.offset)
			}

			var results []NameValue
			_, err := store.GetAll(ctx, q, &results)
			if err != nil {
				t.Fatalf("GetAll failed: %v", err)
			}

			var gotNames []string
			for _, r := range results {
				gotNames = append(gotNames, r.Name)
			}

			if !equalUnordered(gotNames, tc.expectNames) {
				t.Errorf("unexpected result names: got %v, want %v", gotNames, tc.expectNames)
			}
		})
	}
}

func TestFileFilterFieldMixedTypes(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(ctx, "file", "test", "store")
	if err != nil {
		t.Fatalf("could not create file store: %v", err)
	}
	_ = os.RemoveAll("store/test/Mixed")

	now := time.Now().UTC()
	entities := []Mixed{
		{ID: "1", Str: "alpha", Int: 10, Float: 1.1, Created: now.Add(-10 * time.Hour)},
		{ID: "2", Str: "bravo", Int: 20, Float: 2.2, Created: now.Add(-5 * time.Hour)},
		{ID: "3", Str: "charlie", Int: 30, Float: 3.3, Created: now.Add(-1 * time.Hour)},
		{ID: "4", Str: "delta", Int: 40, Float: 4.4, Created: now},
	}

	for _, e := range entities {
		_, err := store.Put(ctx, store.NameKey(typeMixed, e.ID), &e)
		if err != nil {
			t.Fatalf("failed to insert test Mixed entity %s: %v", e.ID, err)
		}
	}

	tests := []struct {
		name      string
		field     string
		operator  string
		value     interface{}
		keyParts  []string
		expectIDs []string
	}{
		{
			name:      "Int >= 30",
			field:     "Int",
			operator:  ">=",
			value:     int64(30),
			keyParts:  []string{"ID"},
			expectIDs: []string{"3", "4"},
		},
		{
			name:      "Float < 3.0",
			field:     "Float",
			operator:  "<",
			value:     3.0,
			keyParts:  []string{"ID"},
			expectIDs: []string{"1", "2"},
		},
		{
			name:      "Str > 'bravo'",
			field:     "Str",
			operator:  ">",
			value:     "bravo",
			keyParts:  []string{"ID"},
			expectIDs: []string{"3", "4"},
		},
		{
			name:      "Created < now",
			field:     "Created",
			operator:  "<",
			value:     now,
			keyParts:  []string{"ID"},
			expectIDs: []string{"1", "2", "3"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := store.NewQuery(typeMixed, false, tc.keyParts...)
			q.FilterField(tc.field, tc.operator, tc.value)
			q.Order(tc.field)

			var results []Mixed
			_, err := store.GetAll(ctx, q, &results)
			if err != nil {
				t.Fatalf("GetAll failed: %v", err)
			}

			var gotIDs []string
			for _, r := range results {
				gotIDs = append(gotIDs, r.ID)
			}
			if !equalUnordered(gotIDs, tc.expectIDs) {
				t.Errorf("unexpected result IDs: got %v, want %v", gotIDs, tc.expectIDs)
			}
		})
	}
}

// equalUnordered checks that two string slices have the same elements, regardless of order.
func equalUnordered(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[string]int)
	for _, x := range a {
		counts[x]++
	}
	for _, x := range b {
		if counts[x] == 0 {
			return false
		}
		counts[x]--
	}
	return true
}
