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
	"fmt"
	"strconv"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"

	"github.com/gofiber/fiber/v2"
)

// GetSpeciesQuery describes the URL query parameters required for the GetRecommendedSpecies endpoint.
type GetSpeciesQuery struct {
	VideoStream   *int64  `query:"videostream"`   // Optional.
	CaptureSource *int64  `query:"capturesource"` // Optional.
	Search        *string `query:"search"`        // Optional.
	api.LimitAndOffset
}

type ImportFromINaturalistQuery struct {
	DescendantsOf []string `query:"descendants_of"`
}

// GetSpeciesByID gets a species when provided with an ID.
//
//	@Summary		Get species by ID
//	@Description	Gets a species when provided with an ID.
//	@Tags			Species
//	@Produce		json
//	@Param			id	path		int	true	"Species ID"	example(1234567890)
//	@Success		200	{object}	services.Species
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/species/{id} [get]
func GetSpeciesByID(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	species, err := services.GetSpeciesByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(species)
}

// GetSpecies gets a list of species.
//
//	@Summary		Get species
//	@Description	Get paginated species, with options to filter by name and location.
//	@Tags			Species
//	@Produce		json
//	@Param			limit	query		int		false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset	query		int		false	"Number of results to skip."	minimum(0)
//	@Param			search	query		string		false	"Search Query"
//	@Success		200		{object}	api.Result[services.Species]
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/species [get]
func GetSpecies(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetSpeciesQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	species, err := services.GetSpecies(qry.Limit, qry.Offset, qry.VideoStream, qry.CaptureSource, qry.Search)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(api.Result[services.Species]{
		Results: species,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(species),
	})
}

// CreateSpecies creates a new species.
//
//	@Summary		Create species
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Creates a new species from provided JSON body.
//	@Tags			Species
//	@Accept			json
//	@Produce		json
//	@Param			body	body		services.SpeciesContents	true	"New Species"
//	@Success		201		{object}	services.Species
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/species [post]
func CreateSpecies(ctx *fiber.Ctx) error {
	// Parse body.
	var body services.SpeciesContents
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	created, err := services.CreateSpecies(body)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(created)
}

// ImportFromINaturalist imports species from INaturalist's taxa API.
//
//	@Summary		Import from iNaturalist
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Imports all species that are descendants of a Phylum/Class/Order/etc from iNaturalist's taxa API.
//	@Tags			Species
//	@Param			descendants_of	query	string	true	"Phylum/Class/Order/etc to import"	example(Infraorder Cetacea)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Router			/api/v1/species/inaturalist-import [post]
func ImportFromINaturalist(ctx *fiber.Ctx) error {

	qry := new(ImportFromINaturalistQuery)

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	for _, parentName := range qry.DescendantsOf {

		// Get parent ID.
		parentTaxa, err := services.GetTaxonByName(parentName)
		if err != nil {
			return fmt.Errorf("could not get taxon by name %s, error: %w", parentName, err)
		}

		// Get descendants.
		species, err := services.GetSpeciesByDescendant(parentTaxa.ID)
		if err != nil {
			return fmt.Errorf("could not get species as descendant of %s, error: %w", parentTaxa.Name, err)
		}

		// Insert species into datastore or update existing entry.
		for _, s := range species {
			// Skip species without a photo.
			if s.DefaultPhoto == nil {
				continue
			}

			species, _ := services.GetSpeciesByINaturalistID(s.ID)
			if species == nil {
				services.CreateSpecies(services.SpeciesContents{
					ScientificName: s.Name,
					CommonName:     s.PreferredCommonName,
					Images: []services.SpeciesImage{
						{
							Src:         s.DefaultPhoto.MediumURL,
							Attribution: s.DefaultPhoto.Attribution,
						},
					},
					INaturalistTaxonID: &s.ID,
				})
			} else {
				services.UpdateSpecies(species.ID, services.PartialSpeciesContents{
					ScientificName: &s.Name,
					CommonName:     &s.PreferredCommonName,
					Images: &[]services.SpeciesImage{
						{
							Src:         s.DefaultPhoto.MediumURL,
							Attribution: s.DefaultPhoto.Attribution,
						},
					},
					INaturalistTaxonID: &s.ID,
				})
			}
		}
	}

	return nil
}

// DeleteSpecies deletes a species.
//
//	@Summary		Delete species
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete a species by providing the species ID.
//	@Tags			Species
//	@Param			id	path	int	true	"Species ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/species/{id} [delete]
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
