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
	"github.com/ausocean/openfish/cmd/openfish/types/keypoint"
	"github.com/ausocean/openfish/datastore"
)

// Kind of entity to store / fetch from the datastore.
const ANNOTATION_KIND = "Annotation"

// An Annotation holds information about observations at a particular moment and region within a video stream.
type Annotation struct {
	VideoStreamID int64
	Start         int64 `datastore:"StartTime"` // for indexing purposes.
	Keypoints     []struct {
		Time string
		keypoint.BoundingBox
	}
	CreatedBy int64

	// User ID and Species ID are stored separately, so we can query for annotations made
	// by a particular user or with a particular species. NOTE: a species ID may appear
	// multiple times if it is identified by many users.
	IdentificationUserID    []int64
	IdentificationSpeciesID []int64
	datastore.NoCache
}

// Implements Copy from the Entity interface.
func (an *Annotation) Copy(dst datastore.Entity) (datastore.Entity, error) {
	return datastore.CopyEntity(an, dst)
}

// NewAnnotation returns a new Annotation entity.
func NewAnnotation() datastore.Entity {
	return &Annotation{}
}
