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
	"log"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// CloudStore implements Store for the Google Cloud Datastore.
type CloudStore struct {
	client *datastore.Client
}

// newCloudStore returns a new CloudStore, using the given URL to
// retrieve credentials and authenticate.
// The ID can be passed with an optional database name in the format
// <ID>/<Database_Name>, if there is no database name given, the default
// database will be used.
// To obtain credentials from a Google storage bucket, URL takes the
// form gs://bucket_name/creds. A URL without a scheme is interpreted
// as a file. If the environment variable <ID>_CREDENTIALS is defined
// it overrides the supplied URL.
func newCloudStore(ctx context.Context, id, url string) (*CloudStore, error) {
	s := new(CloudStore)

	var db string
	parts := strings.Split(id, "/")
	if len(parts) == 2 {
		db = parts[1]
	} else if len(parts) != 1 {
		return nil, ErrInvalidStoreID
	}

	id = parts[0]

	ev := strings.ToUpper(id) + "_CREDENTIALS"
	if os.Getenv(ev) != "" {
		url = os.Getenv(ev)
	}

	var err error
	if url == "" {
		// Attempt authentication using the default credentials.
		s.client, err = datastore.NewClientWithDatabase(ctx, id, db)
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

	s.client, err = datastore.NewClientWithDatabase(ctx, id, db, option.WithCredentialsJSON(creds))
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
		if err != ErrNoSuchEntity {
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
