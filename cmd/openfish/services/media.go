/*
AUTHORS
  Scott Barnard <scott@ausocean.org>
LICENSE
  Copyright (c) 2023-2024, The OpenFish Contributors.
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

// services contains the main logic for the OpenFish API.
package services

import (
	"context"
	"fmt"
	"io"

	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/types/mediatype"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"
)

// MediaKey is a composite key with all the attributes that define media.
// All media are uniquely identified by the combination of these three parameters.
type MediaKey struct {
	Type          mediatype.MediaType
	VideoStreamID int64
	StartTime     videotime.VideoTime
	EndTime       *videotime.VideoTime
}

// Valid tests a MediaKey has a valid VideoStreamID and either has an end time for the case
// of videos, or does not have an end time in the case of images.
func (q *MediaKey) Valid() bool {
	return q.Type.IsVideo() == (q.EndTime != nil) && VideoStreamExists(q.VideoStreamID)
}

// ToStorageName returns the name used to store the media in a bucket.
func (q *MediaKey) ToStorageName() string {
	if q.EndTime == nil {
		return fmt.Sprintf("images/%d[%s].%s", q.VideoStreamID, q.StartTime.String(), q.Type.FileExtension())
	} else {
		return fmt.Sprintf("videos/%d[%s-%s].%s", q.VideoStreamID, q.StartTime.String(), q.EndTime.String(), q.Type.FileExtension())
	}
}

// GetMedia gets the media with the specified type, source video stream and time.
func GetMedia(q MediaKey) ([]byte, error) {

	if !q.Valid() {
		return nil, fmt.Errorf("video media types must be provided with an end time and image media types must not")
	}

	// Get file from storage.
	storage := globals.GetStorage()
	name := q.ToStorageName()
	handle := storage.Object(name)
	r, err := handle.NewReader(context.Background())
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// MediaExists checks if the media exists in the datastore.
func MediaExists(q MediaKey) bool {
	name := q.ToStorageName()
	storage := globals.GetStorage()
	handle := storage.Object(name)
	exists, _ := handle.Exists(context.Background())
	return exists
}

// CreateMedia puts an image or video in the datastore.
func CreateMedia(q MediaKey, data []byte) (string, error) {

	if !q.Valid() {
		return "", fmt.Errorf("video media types must be provided with an end time and image media types must not")
	}

	// Put binary file into storage.
	name := q.ToStorageName()
	storage := globals.GetStorage()
	handle := storage.Object(name)
	w, err := handle.NewWriter(context.Background())
	if err != nil {
		return "", err
	}
	defer w.Close()

	_, err = w.Write(data)
	if err != nil {
		return "", err
	}

	// Return name of created media.
	return name, nil
}

// DeleteMedia deletes an image/video.
func DeleteMedia(q MediaKey) error {
	name := q.ToStorageName()
	storage := globals.GetStorage()
	handle := storage.Object(name)
	return handle.Delete(context.Background())
}
