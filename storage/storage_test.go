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
	"testing"
)

const baseDir = "./test-storage"

// TestWriteObject verifies that we can create and write to a file.
func TestWriteObject(t *testing.T) {

	os.MkdirAll(baseDir, os.ModePerm)
	defer os.RemoveAll(baseDir)

	// File path to be used for the test
	storage := NewFileStorage(baseDir)
	handle := storage.Object("my-object")

	// Create a writer and write data
	writer, err := handle.NewWriter(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = writer.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("expected no error while writing, got %v", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		t.Fatalf("expected no error while closing writer, got %v", err)
	}

	// Verify the file has been created
	_, err = os.Stat(path.Join(baseDir, "my-object"))
	if err != nil {
		t.Fatalf("expected file to be created, but got error: %v", err)
	}
}

// TestReadObject verifies that we can read data from a file.
func TestReadObject(t *testing.T) {

	os.MkdirAll(baseDir, 0755)
	defer os.RemoveAll(baseDir)

	// File path to be used for the test
	storage := NewFileStorage(baseDir)
	handle := storage.Object("my-object")

	// Create a file with some test data
	filePath := path.Join(baseDir, "my-object")
	err := os.WriteFile(filePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("expected no error while writing file, got %v", err)
	}

	// Create a reader and read the content
	reader, err := handle.NewReader(context.Background())
	if err != nil {
		t.Fatalf("expected no error while creating reader, got %v", err)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("expected no error while reading file, got %v", err)
	}

	// Verify the content is correct
	expectedContent := "test data"
	if string(content) != expectedContent {
		t.Fatalf("expected content %v, got %v", expectedContent, string(content))
	}
}
