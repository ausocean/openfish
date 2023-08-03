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
	"errors"
	"sync"
)

// Type Cache defines the (optional) caching interface used by Entity.
type Cache interface {
	Set(key *Key, value Entity)   // Set adds or updates a value to the cache.
	Get(key *Key) (Entity, error) // Get retrieves a value from the cache, or returns ErrcacheMiss.
	Delete(key *Key)              // Delete removes a value from the cache.
	Reset()                       // Reset resets (clears) the cache.
}

// Type genericCache represents a generic cache for holding datastore entities.
// The key K is either an int64 or a string.
type genericCache[K comparable] struct {
	data  map[K]Entity
	mutex sync.RWMutex
}

// ErrCacheMiss is the error returned when a value is not found in the cache.
var ErrCacheMiss = errors.New("cache miss")

// Set adds or updates a value to the cache.
func (c *genericCache[K]) Set(key K, value Entity) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data == nil {
		c.data = map[K]Entity{}
	}
	c.data[key] = value
}

// Get retrieves a value from the cache, or returns ErrcacheMiss.
func (c *genericCache[K]) Get(key K) (Entity, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data == nil {
		c.data = map[K]Entity{}
	}
	value, ok := c.data[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	return value, nil
}

// Delete removes a value from the cache.
func (c *genericCache[K]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data == nil {
		c.data = map[K]Entity{}
	}
	delete(c.data, key)
}

// Reset resets (clears) the cache.
func (c *genericCache[K]) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = map[K]Entity{}
}

// Below we define two implementations of the Cache interface, which
// are thin wrappers of genericCache.

// Type IDcache caches entities by (int64) ID keys.
type IDCache struct {
	cache genericCache[int64]
}

// Set adds or updates a value to the cache.
func (c *IDCache) Set(key *Key, value Entity) {
	c.cache.Set(key.ID, value)
}

// Get retrieves a value from the cache, or returns ErrcacheMiss.
func (c *IDCache) Get(key *Key) (Entity, error) {
	return c.cache.Get(key.ID)
}

// Delete removes a value from the cache.
func (c *IDCache) Delete(key *Key) {
	c.cache.Delete(key.ID)
}

// Reset resets (clears) the cache.
func (c *IDCache) Reset() {
	c.cache.Reset()
}

// Type Namecache caches entities by (string) name keys.
type NameCache struct {
	cache genericCache[string]
}

// Set adds or updates a value to the cache.
func (c *NameCache) Set(key *Key, value Entity) {
	c.cache.Set(key.Name, value)
}

// Get retrieves a value from the cache, or returns ErrcacheMiss.
func (c *NameCache) Get(key *Key) (Entity, error) {
	return c.cache.Get(key.Name)
}

// Delete removes a value from the cache.
func (c *NameCache) Delete(key *Key) {
	c.cache.Delete(key.Name)
}

// Reset resets (clears) the cache.
func (c *NameCache) Reset() {
	c.cache.Reset()
}
