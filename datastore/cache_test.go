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
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		action, key, value, want string
		ok                       bool // true if action returns an error and is expected to succeed.
	}{
		{
			action: "get",
			key:    "a",
			ok:     false,
		},
		{
			action: "set",
			key:    "a",
			value:  "aa",
		},
		{
			action: "get",
			key:    "a",
			want:   "aa",
			ok:     true,
		},
		{
			action: "set",
			key:    "b",
			value:  "bb",
		},
		{
			action: "delete",
			key:    "a",
		},
		{
			action: "get",
			key:    "a",
			ok:     false,
		},
		{
			action: "get",
			key:    "b",
			want:   "bb",
			ok:     true,
		},
		{
			action: "reset",
		},
		{
			action: "get",
			key:    "b",
			ok:     false,
		},
	}

	var cache Cache = NewEntityCache()

	for i, test := range tests {
		var k Key = Key{Name: test.key}

		switch test.action {
		case "get":
			var kv KeyValue
			err := cache.Get(&k, &kv)
			if err != nil {
				if test.ok {
					t.Errorf("Test %d: Get(%s) returned unexpected error: %v", i, test.key, err)
				}
				var errCacheMiss ErrCacheMiss
				if !errors.As(err, &errCacheMiss) {
					t.Errorf("Test %d: Get(%s) returned wrong error: %v", i, test.key, err)
				}
				continue // Got expected type of error.
			}
			if !test.ok {
				t.Errorf("Test %d: Get(%s) did not return error", i, test.key)
			}
			if test.want != kv.Value {
				t.Errorf("Test %d: Get(%s) returned wrong value: %s", i, test.key, kv.Value)
			}

		case "set":
			kv := KeyValue{Key: test.key, Value: test.value}
			err := cache.Set(&k, &kv)
			if err != nil {
				t.Errorf("Test %d: Set(%s,%s) returned unexpected error: %v", i, test.key, test.value, err)
			}

		case "delete":
			cache.Delete(&k)

		case "reset":
			cache.Reset()

		default:
			panic("unexpected test action: " + test.action)
		}

	}
}
