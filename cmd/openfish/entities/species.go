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

package entities

import (
	"github.com/ausocean/openfish/datastore"
)

// Kind of entity to store / fetch from the datastore.
const SPECIES_KIND = "Species"

// Species is used for our guide.
type Species struct {
	Species            string
	CommonName         string
	Images             []Image
	INaturalistTaxonID *int // Optional.
	SearchIndex        []string
}

type Image struct {
	Src         string `json:"src" example:"https://inaturalist-open-data.s3.amazonaws.com/photos/340064435/medium.jpg"`
	Attribution string `json:"attribution" example:"Tiffany Kosch, CC BY-NC-SA 4.0"`
}

// Implements Copy from the Entity interface.
func (vs *Species) Copy(dst datastore.Entity) (datastore.Entity, error) {
	var v *Species
	if dst == nil {
		v = new(Species)
	} else {
		var ok bool
		v, ok = dst.(*Species)
		if !ok {
			return nil, datastore.ErrWrongType
		}
	}
	*v = *vs
	return v, nil
}

// No caching is used.
func (vs *Species) GetCache() datastore.Cache {
	return nil
}

func NewSpecies() datastore.Entity {
	return &Species{}
}
