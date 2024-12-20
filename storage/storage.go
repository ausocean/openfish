/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2024, The OpenFish Contributors.

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

// Package storage offers a storage service for large binary data with
// multiple implementations.
//
//   - CloudStorage is a Google Cloud Storage implementation.
//   - FileStorage is a file based implementation.
package storage

import (
	"context"
	"io"
	"os"
	"path"

	"cloud.google.com/go/storage"
)

// Storage defines the storage interface. It lets us store large binary data in a bucket or file.
type Storage interface {
	Object(name string) ObjectHandle
}

// ObjectHandle defines the object handle interface. It provides a ReadCloser and WriteCloser for reading
// and writing to the object.
type ObjectHandle interface {
	NewWriter(ctx context.Context) (io.WriteCloser, error)
	NewReader(ctx context.Context) (io.ReadCloser, error)
	Delete(ctx context.Context) error
	Exists(ctx context.Context) (bool, error)
}

// CloudStorage implements Storage using Google Cloud Buckets.
type CloudStorage struct {
	bkt *storage.BucketHandle
}

// CloudObjectHandle implements ObjectHandle using Google Cloud Buckets.
type CloudObjectHandle struct {
	objHandle *storage.ObjectHandle
}

// NewCloudStorage creates a new CloudStorage, using name as the bucket URI.
func NewCloudStorage(name string) (Storage, error) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &CloudStorage{
		bkt: client.Bucket(name),
	}, nil
}

// Object returns an ObjectHandle, which provides operations on the named object.
func (s *CloudStorage) Object(name string) ObjectHandle {
	return &CloudObjectHandle{
		objHandle: s.bkt.Object(name),
	}
}

// NewWriter returns a WriteCloser that writes to the storage object.
func (h *CloudObjectHandle) NewWriter(ctx context.Context) (io.WriteCloser, error) {
	return h.objHandle.NewWriter(ctx), nil
}

// NewReader returns a ReadCloser that reads from the storage object.
func (h *CloudObjectHandle) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return h.objHandle.NewReader(ctx)
}

func (h *CloudObjectHandle) Delete(ctx context.Context) error {
	return h.objHandle.Delete(ctx)
}

func (h *CloudObjectHandle) Exists(ctx context.Context) (bool, error) {
	_, err := h.objHandle.Attrs(ctx)
	if err == nil {
		return true, nil
	}
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	return false, err
}

// FileStorage implements Storage using files on your local machine.
type FileStorage struct {
	base string
}

// FileObjectHandle implements ObjectHandle using files on your local machine.
type FileObjectHandle struct {
	filepath string
}

// NewFileStorage creates a new FileStorage, using base as the path of the directory to store objects in.
func NewFileStorage(base string) Storage {
	return &FileStorage{
		base: base,
	}
}

// Object returns an ObjectHandle, which provides operations on the named object.
func (s *FileStorage) Object(name string) ObjectHandle {
	return &FileObjectHandle{
		filepath: path.Join(s.base, name),
	}
}

// NewWriter returns a WriteCloser that writes to the storage object.
func (h *FileObjectHandle) NewWriter(ctx context.Context) (io.WriteCloser, error) {
	return os.Create(h.filepath)
}

// NewReader returns a ReadCloser that reads from the storage object.
func (h *FileObjectHandle) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return os.Open(h.filepath)
}

func (h *FileObjectHandle) Delete(ctx context.Context) error {
	return os.Remove(h.filepath)
}

func (h *FileObjectHandle) Exists(ctx context.Context) (bool, error) {
	_, err := os.Stat(h.filepath)
	if err == nil {
		return true, nil
	}
	if err == os.ErrNotExist {
		return false, nil
	}
	return false, err
}
