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

// handlers package handles HTTP requests.
package handlers

import (
	"strconv"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/services"
	"github.com/ausocean/openfish/api/types/species"

	"github.com/gofiber/fiber/v2"
)

// SpeciesResult describes the JSON format for species in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type SpeciesResult struct {
	ID         *int64           `json:"id,omitempty"`
	Species    *string          `json:"species,omitempty"`
	CommonName *string          `json:"common_name,omitempty"`
	Images     *[]species.Image `json:"images,omitempty"`
}

// FromSpecies creates a SpeciesResult from a species.Species and key, formatting it according to the requested format.
func FromSpecies(species *species.Species, id int64, format *api.Format) SpeciesResult {
	var result SpeciesResult
	if format.Requires("id") {
		result.ID = &id
	}
	if format.Requires("species") {
		result.Species = &species.Species
	}
	if format.Requires("common_name") {
		result.CommonName = &species.CommonName
	}
	if format.Requires("images") {
		result.Images = &species.Images
	}
	return result
}

// GetRecommendedSpeciesQuery describes the URL query parameters required for the GetRecommendedSpecies endpoint.
type GetRecommendedSpeciesQuery struct {
	VideoStream   *int64 `query:"videostream"`   // Optional.
	CaptureSource *int64 `query:"capturesource"` // Optional.
	api.LimitAndOffset
}

type ImportFromINaturalistQuery struct {
	DescendantsOf []string `query:"descendants_of"`
}

// CreateSpeciesBody describes the JSON format required for the CreateSpecies endpoint.
//
// ID is omitted because it is chosen automatically.
type CreateSpeciesBody struct {
	Species    string          `json:"species"`
	CommonName string          `json:"common_name"`
	Images     []species.Image `json:"images"`
}

// GetSpeciesByID gets a species when provided with an ID.
func GetSpeciesByID(ctx *fiber.Ctx) error {
	// Parse URL.
	format := new(api.Format)

	if err := ctx.QueryParser(format); err != nil {
		return api.InvalidRequestURL(err)
	}

	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	species, err := services.GetSpeciesByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromSpecies(species, id, format)
	return ctx.JSON(result)
}

// GetRecommendedSpecies gets a list of species, most relevant for the specified stream and capture source.
func GetRecommendedSpecies(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetRecommendedSpeciesQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	format := new(api.Format)
	if err := ctx.QueryParser(format); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	species, ids, err := services.GetRecommendedSpecies(qry.Limit, qry.Offset, qry.VideoStream, qry.CaptureSource)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]SpeciesResult, len(species))
	for i := range species {
		results[i] = FromSpecies(&species[i], int64(ids[i]), format)
	}

	return ctx.JSON(api.Result[SpeciesResult]{
		Results: results,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(results),
	})
}

// CreateSpecies creates a new species.
func CreateSpecies(ctx *fiber.Ctx) error {
	// Parse body.
	var body CreateSpeciesBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	id, err := services.CreateSpecies(body.Species, body.CommonName, body.Images, nil)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(VideoStreamResult{
		ID: &id,
	})
}

// ImportFromINaturalist imports species from INaturalist's taxa API.
func ImportFromINaturalist(ctx *fiber.Ctx) error {

	qry := new(ImportFromINaturalistQuery)

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	for _, parentName := range qry.DescendantsOf {

		// Get parent ID.
		parentTaxa, err := services.GetTaxonByName(parentName)
		if err != nil {
			return err // TODO: add more descriptive error.
		}

		// Get descendants.
		descendants, err := services.GetSpeciesByDescendant(parentTaxa.ID)
		if err != nil {
			return err // TODO: add more descriptive error.
		}

		// Insert species into datastore or update existing entry.
		for _, s := range descendants {

			img := species.Image{Src: s.DefaultPhoto.MediumURL, Attribution: s.DefaultPhoto.Attribution}

			spec, id, _ := services.GetSpeciesByINaturalistID(s.ID)
			if spec == nil {
				services.CreateSpecies(s.Name, s.PreferredCommonName, []species.Image{img}, &s.ID)
			}

			services.UpdateSpecies(id, &s.Name, &s.PreferredCommonName, &[]species.Image{img}, nil)
		}
	}

	return nil
}

// DeleteSpecies deletes a species.
func DeleteSpecies(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Delete entity.
	err = services.DeleteSpecies(id)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
