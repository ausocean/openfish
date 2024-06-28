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

package services

import (
	"context"

	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/types/species"
	"github.com/ausocean/openfish/datastore"
)

// GetSpeciesByID gets a species when provided with an ID.
func GetSpeciesByID(id int64) (*species.Species, error) {
	store := ds_client.Get()
	key := store.IDKey(species.KIND, id)
	var s species.Species
	err := store.Get(context.Background(), key, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetSpeciesByINaturalist gets a species when provided with an iNaturalist ID.
func GetSpeciesByINaturalistID(id int) (*species.Species, int64, error) {
	store := ds_client.Get()
	query := store.NewQuery(species.KIND, false)

	query.FilterField("INaturalistTaxonID", "=", id)
	query.Limit(1)

	var s []species.Species
	keys, err := store.GetAll(context.Background(), query, &s)
	if err != nil {
		return nil, 0, err
	}

	if len(keys) == 0 {
		return nil, 0, nil
	}

	return &s[0], keys[0].ID, nil
}

func SpeciesExists(id int64) bool {
	store := ds_client.Get()
	key := store.IDKey(species.KIND, id)
	var s species.Species
	err := store.Get(context.Background(), key, &s)
	return err == nil
}

// GetRecommendedSpecies gets a list of species, most relevant for the specified stream and capture source.
func GetRecommendedSpecies(limit int, offset int, videostream *int64, captureSource *int64) ([]species.Species, []int64, error) {
	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery(species.KIND, false)

	// TODO: implement returning most relevant species.

	query.Limit(limit)
	query.Offset(offset)

	var s []species.Species
	keys, err := store.GetAll(context.Background(), query, &s)
	if err != nil {
		return []species.Species{}, []int64{}, err
	}
	ids := make([]int64, len(s))
	for i, k := range keys {
		ids[i] = k.ID
	}

	return s, ids, nil
}

// CreateSpecies puts a species in the datastore.
func CreateSpecies(latin string, commonName string, images []species.Image, iNaturalistTaxonID *int) (int64, error) {

	// Create Species entity.
	store := ds_client.Get()
	key := store.IncompleteKey(species.KIND)

	sp := species.Species{
		Species:            latin,
		CommonName:         commonName,
		Images:             images,
		INaturalistTaxonID: iNaturalistTaxonID,
	}
	key, err := store.Put(context.Background(), key, &sp)
	if err != nil {
		return 0, err
	}

	// Return ID of created species.
	return key.ID, nil
}

// UpdateSpecies finds the species with a given
func UpdateSpecies(id int64, latin *string, commonName *string, images *[]species.Image, iNaturalistTaxonID *int) error {
	store := ds_client.Get()
	key := store.IDKey(species.KIND, id)

	var sp species.Species

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		s, ok := e.(*species.Species)
		if ok {
			if latin != nil {
				s.Species = *latin
			}
			if commonName != nil {
				s.CommonName = *commonName
			}
			if images != nil {
				s.Images = *images
			}
			if iNaturalistTaxonID != nil {
				s.INaturalistTaxonID = iNaturalistTaxonID
			}
		}
	}, &sp)

}

// DeleteSpecies deletes a species.
func DeleteSpecies(id int64) error {
	// Delete entity.
	store := ds_client.Get()
	key := store.IDKey(species.KIND, id)
	return store.Delete(context.Background(), key)
}
