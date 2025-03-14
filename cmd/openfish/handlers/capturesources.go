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

// handlers package handles HTTP requests.
package handlers

import (
	"strconv"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/latlong"

	"github.com/gofiber/fiber/v2"
)

// EntityIDResult contains the ID of a newly created entity.
//
//	@Description	ID of newly created entity.
type EntityIDResult struct {
	ID int64 `json:"id" example:"1234567890"` // Unique ID of the entity.
}

// GetCaptureSourcesQuery describes the URL query parameters required for the GetCaptureSources endpoint.
type GetCaptureSourcesQuery struct {
	Name     *string          `query:"name"`     // Optional.
	Location *latlong.LatLong `query:"location"` // Optional.
	api.LimitAndOffset
	// TODO: Code could be simplified if api.Format could be embedded here. Testing shows it cannot, reason unknown.
}

// GetCaptureSourceByID gets a capture source when provided with an ID.
//
//	@Summary		Get capture source by ID
//	@Description	Gets a capture source when provided with an ID.
//	@Tags			Capture Sources
//	@Produce		json
//	@Param			id	path		int	true	"Capture Source ID"	example(1234567890)
//	@Success		200	{object}	services.CaptureSource
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/capturesources/{id} [get]
func GetCaptureSourceByID(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	src, err := services.GetCaptureSourceByID(id)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	return ctx.JSON(src)
}

// GetCaptureSources gets a list of capture sources, filtering by name, location if specified.
//
//	@Summary		Get capture sources
//	@Description	Get paginated capture sources, with options to filter by name and location.
//	@Tags			Capture Sources
//	@Produce		json
//	@Param			limit	query		int		false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset	query		int		false	"Number of results to skip."	minimum(0)
//	@Param			name	query		string	false	"Name to filter by."
//	@Success		200		{object}	api.Result[services.CaptureSource]
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/capturesources [get]
func GetCaptureSources(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetCaptureSourcesQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	srcs, err := services.GetCaptureSources(qry.Limit, qry.Offset, qry.Name)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	return ctx.JSON(api.Result[services.CaptureSource]{
		Results: srcs,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(srcs),
	})
}

// CreateCaptureSource creates a new capture source.
//
//	@Summary		Create capture source
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Creates a new capture source from provided JSON body.
//	@Tags			Capture Sources
//	@Accept			json
//	@Produce		json
//	@Param			body	body		services.CaptureSourceContents	true	"New Capture Source"
//	@Success		201		{object}	services.CaptureSource
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/capturesources [post]
func CreateCaptureSource(ctx *fiber.Ctx) error {
	// Parse body.
	var body services.CaptureSourceContents

	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create capture source entity and add to the datastore.
	created, err := services.CreateCaptureSource(body)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created capture source.
	return ctx.JSON(created)
}

// UpdateCaptureSource updates a capture source.
//
//	@Summary		Update capture source
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Partially update a capture source by specifying the properties to update.
//	@Tags			Capture Sources
//	@Accept			json
//	@Param			id		path	int										true	"Capture Source ID"	example(1234567890)
//	@Param			body	body	services.PartialCaptureSourceContents	true	"Update Capture Source"
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/capturesources/{id} [patch]
func UpdateCaptureSource(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse body.
	var body services.PartialCaptureSourceContents
	if ctx.BodyParser(&body) != nil {
		return api.InvalidRequestJSON(err)
	}

	// Update data in the datastore.
	err = services.UpdateCaptureSource(id, body)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}

// DeleteCaptureSource deletes a capture source.
//
//	@Summary		Delete capture source
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete a capture source by providing the capture source ID.
//	@Tags			Capture Sources
//	@Param			id	path	int	true	"Capture Source ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/capturesources/{id} [delete]
func DeleteCaptureSource(ctx *fiber.Ctx) error {

	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Delete capture source.
	err = services.DeleteCaptureSource(id)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
