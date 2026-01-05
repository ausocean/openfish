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

// entities package has the data types of data we keep in the datastore.
package entities

import (
	googlestore "cloud.google.com/go/datastore"

	"github.com/ausocean/openfish/datastore"
)

// Kind of entity to store / fetch from the datastore.
const CAPTURESOURCE_KIND = "CaptureSource"

// A CaptureSource holds the information about where a video stream came from.
// A single capture source will produce multiple video streams (typically one per day).
type CaptureSource struct {
	Name           string
	Location       googlestore.GeoPoint
	CameraHardware string
	SiteID         *int64 // Optional.
	datastore.NoCache
}

// Implements Copy from the Entity interface.
func (cs *CaptureSource) Copy(dst datastore.Entity) (datastore.Entity, error) {
	return datastore.CopyEntity(cs, dst)
}

// NewCaptureSource returns a new CaptureSource entity.
func NewCaptureSource() datastore.Entity {
	return &CaptureSource{}
}
