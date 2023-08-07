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

// EntityCache represents a cache for holding datastore entities indexed by key.
type EntityCache struct {
	data  map[Key]Entity
	mutex sync.RWMutex
}

// ErrCacheMiss is the error returned when a value is not found in the cache.
var ErrCacheMiss = errors.New("cache miss")

// Set adds or updates a value to the cache.
func (c *EntityCache) Set(key *Key, value Entity) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data == nil {
		c.data = map[Key]Entity{}
	}
	c.data[*key] = value
}

// Get retrieves a value from the cache, or returns ErrcacheMiss.
func (c *EntityCache) Get(key *Key) (Entity, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data == nil {
		c.data = map[Key]Entity{}
	}
	value, ok := c.data[*key]
	if !ok {
		return nil, ErrCacheMiss
	}
	return value, nil
}

// Delete removes a value from the cache.
func (c *EntityCache) Delete(key *Key) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data == nil {
		c.data = map[Key]Entity{}
	}
	delete(c.data, *key)
}

// Reset resets (clears) the cache.
func (c *EntityCache) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = map[Key]Entity{}
}
