/*
AUTHORS
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

package model

import (
	"encoding/json"
	"time"
)

// Kind of entity to store / fetch from the datastore.
const VIDEOSTREAM_KIND = "VideoStream"

// VideoStream holds the information about a single video stream.
// VideoStream contains the url for a live or completed stream off of youtube, the start time,
// the end time (unless it is still ongoing), and the ID of its capture source.
// In the future, it will also contain data about how to retrieve the video data from the cloud.
type VideoStream struct {
	StartTime     time.Time
	EndTime       *time.Time // Optional.
	StreamUrl     string
	CaptureSource int64
	// TODO: Add cloud storage location.
}

// Encode serializes VideoStream. Implements Entity interface. Used for FileStore datastore.
func (vs *VideoStream) Encode() []byte {
	bytes, _ := json.Marshal(vs)
	return bytes
}

// Encode deserializes VideoStream. Implements Entity interface. Used for FileStore datastore.
func (vs *VideoStream) Decode(b []byte) error {
	return json.Unmarshal(b, vs)
}
