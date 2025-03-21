/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023-2025, The OpenFish Contributors.

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
	"strings"

	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/sliceutils"
	"github.com/ausocean/openfish/datastore"
)

// Species describes a species that can be chosen in identifications on a stream.
type Species struct {
	ID int64 `json:"id" example:"1234567890"` // Unique ID of the species.
	SpeciesContents
}

// SpeciesContents is the contents of a Species.
type SpeciesContents struct {
	ScientificName     string         `json:"scientific_name" example:"Rhincodon typus"` // Scientific name of the species.
	CommonName         string         `json:"common_name" example:"Whale Shark"`         // Common name (in English) of the species.
	Images             []SpeciesImage `json:"images"`                                    // Image or images of the species.
	INaturalistTaxonID *int           `json:"inaturalist_taxon_id" example:"1234567890"`
}

// PartialSpeciesContents is for updating a species with a partial update (such as a PATCH request).
type PartialSpeciesContents struct {
	ScientificName     *string         `json:"scientific_name,omitempty" example:"Rhincodon typus"` // Scientific name of the species.
	CommonName         *string         `json:"common_name,omitempty" example:"Whale Shark"`         // Common name (in English) of the species.
	Images             *[]SpeciesImage `json:"images,omitempty"`                                    // Image or images of the species.
	INaturalistTaxonID *int            `json:"inaturalist_taxon_id" example:"1234567890"`
}

// SpeciesImage represents an image URL and attribution pair.
type SpeciesImage struct {
	Src         string `json:"src" example:"https://inaturalist-open-data.s3.amazonaws.com/photos/340064435/medium.jpg"`
	Attribution string `json:"attribution" example:"Tiffany Kosch, CC BY-NC-SA 4.0"`
}

// SpeciesSummary is a summary of a species.
type SpeciesSummary struct {
	ID             int64  `json:"id" example:"1234567890"`
	CommonName     string `json:"common_name" example:"Whale Shark"`
	ScientificName string `json:"scientific_name" example:"Rhincodon typus"`
}

// SpeciesContentsFromEntity converts an entities.Species to a SpeciesContents.
func SpeciesContentsFromEntity(e entities.Species) SpeciesContents {
	images := make([]SpeciesImage, len(e.ImageSources))
	for i, _ := range e.ImageSources {
		images[i].Src = e.ImageSources[i]
		images[i].Attribution = e.ImageAttributions[i]
	}

	return SpeciesContents{
		ScientificName:     e.ScientificName,
		CommonName:         e.CommonName,
		Images:             images,
		INaturalistTaxonID: e.INaturalistTaxonID,
	}
}

// ToEntity converts a SpeciesContents to an entities.Species for storage in the datastore.
func (s *SpeciesContents) ToEntity() entities.Species {
	sources := make([]string, len(s.Images))
	attributions := make([]string, len(s.Images))
	for i, img := range s.Images {
		sources[i] = img.Src
		attributions[i] = img.Attribution
	}

	e := entities.Species{
		ScientificName:     s.ScientificName,
		CommonName:         s.CommonName,
		ImageSources:       sources,
		ImageAttributions:  attributions,
		INaturalistTaxonID: s.INaturalistTaxonID,
		SearchIndex:        makeSearchIndex(s.ScientificName, s.CommonName),
	}

	return e
}

// ToSummary converts a Species to a SpeciesSummary.
func (s *Species) ToSummary() SpeciesSummary {
	return SpeciesSummary{
		ID:             s.ID,
		CommonName:     s.CommonName,
		ScientificName: s.ScientificName,
	}
}

// makeSearchIndex derives the search index field from the scientific name and common name.
func makeSearchIndex(scientificName string, commonName string) []string {
	searchableStrings := make([]string, 0, 10)

	for subslice := range sliceutils.WindowPermutations(strings.Split(scientificName, " ")) {
		str := strings.ToLower(strings.Join(subslice, " "))
		searchableStrings = append(searchableStrings, str)
	}

	for subslice := range sliceutils.WindowPermutations(strings.Split(commonName, " ")) {
		str := strings.ToLower(strings.Join(subslice, " "))
		searchableStrings = append(searchableStrings, str)
	}
	return searchableStrings
}

// GetSpeciesByID gets a species when provided with an ID.
func GetSpeciesByID(id int64) (*Species, error) {
	store := globals.GetStore()
	key := store.IDKey(entities.SPECIES_KIND, id)
	var e entities.Species
	err := store.Get(context.Background(), key, &e)
	if err != nil {
		return nil, err
	}

	species := Species{
		ID:              id,
		SpeciesContents: SpeciesContentsFromEntity(e),
	}

	return &species, nil
}

// GetSpeciesByINaturalist gets a species when provided with an iNaturalist ID.
func GetSpeciesByINaturalistID(id int) (*Species, error) {
	store := globals.GetStore()
	query := store.NewQuery(entities.SPECIES_KIND, false)

	query.FilterField("INaturalistTaxonID", "=", id)
	query.Limit(1)

	var ents []entities.Species
	keys, err := store.GetAll(context.Background(), query, &ents)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}

	species := Species{
		ID:              keys[0].ID,
		SpeciesContents: SpeciesContentsFromEntity(ents[0]),
	}

	return &species, nil
}

// SpeciesExists checks if a species exists with the given ID.
func SpeciesExists(id int64) bool {
	store := globals.GetStore()
	key := store.IDKey(entities.SPECIES_KIND, id)
	var species entities.Species
	err := store.Get(context.Background(), key, &species)
	return err == nil
}

// GetSpecies gets a list of species, most relevant for the specified stream and capture source.
func GetSpecies(limit int, offset int, videostream *int64, captureSource *int64, search *string) ([]Species, error) {
	// Fetch data from the datastore.
	store := globals.GetStore()
	query := store.NewQuery(entities.SPECIES_KIND, false)

	if search != nil {
		trimmed := strings.TrimSpace(*search)
		lower := strings.ToLower(trimmed)

		// Datastore does not support starts with or contains queries so we do two inequalities.
		query.FilterField("SearchIndex", ">=", lower)
		lastChar := (lower)[len(lower)-1]
		bytes := []byte(lower)
		bytes[len(bytes)-1] = lastChar + 1
		query.FilterField("SearchIndex", "<=", string(bytes))
	}

	// TODO: implement ordering by most relevant species.

	query.Limit(limit)
	query.Offset(offset)

	var ents []entities.Species
	keys, err := store.GetAll(context.Background(), query, &ents)
	if err != nil {
		return []Species{}, err
	}
	species := make([]Species, len(ents))
	for i, key := range keys {
		species[i] = Species{
			ID:              key.ID,
			SpeciesContents: SpeciesContentsFromEntity(ents[i]),
		}
	}

	return species, nil
}

// CreateSpecies puts a species in the datastore.
func CreateSpecies(contents SpeciesContents) (*Species, error) {

	// Create Species entity.
	store := globals.GetStore()
	key := store.IncompleteKey(entities.SPECIES_KIND)
	ent := contents.ToEntity()
	key, err := store.Put(context.Background(), key, &ent)
	if err != nil {
		return nil, err
	}

	species := Species{
		ID:              key.ID,
		SpeciesContents: contents,
	}

	return &species, nil
}

// UpdateSpecies updates existing species with partial species data
func UpdateSpecies(id int64, updates PartialSpeciesContents) error {
	store := globals.GetStore()
	key := store.IDKey(entities.SPECIES_KIND, id)

	var sp entities.Species

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		s, ok := e.(*entities.Species)
		if ok {
			if updates.ScientificName != nil {
				s.ScientificName = *updates.ScientificName
			}
			if updates.CommonName != nil {
				s.CommonName = *updates.CommonName
			}
			s.SearchIndex = makeSearchIndex(s.ScientificName, s.CommonName)
		}
	}, &sp)

}

// DeleteSpecies deletes a species.
func DeleteSpecies(id int64) error {
	// Delete entity.
	store := globals.GetStore()
	key := store.IDKey(entities.SPECIES_KIND, id)
	return store.Delete(context.Background(), key)
}
