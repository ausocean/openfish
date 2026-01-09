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
	"fmt"
	"slices"
	"strconv"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/keypoint"

	"github.com/gofiber/fiber/v2"
)

// GetAnnotationsQuery describes the URL query parameters required for the GetAnnotations endpoint.
type GetAnnotationsQuery struct {
	// TimeSpan      *string           `query:"timespan"`      // Optional. TODO: choose more appropriate type.
	// CaptureSource *int64            `query:"capturesource"` // Optional.
	VideoStream *int64 `query:"videostream"` // Optional.
	api.LimitAndOffset
	api.Sort
}

// GetAnnotationByID gets an annotation when provided with an ID.
//
//	@Summary		Get annotation by ID
//	@Description	Gets an annotation when provided with an ID.
//	@Tags			Annotations
//	@Produce		json
//	@Param			id	path		int	true	"Annotation ID"	example(1234567890)
//	@Success		200	{object}	services.AnnotationWithJoins
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/annotations/{id} [get]
func GetAnnotationByID(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	annotation, err := services.GetAnnotationByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	joined, err := annotation.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
}

// GetAnnotations gets a list of annotations, filtering by videostream if specified.
//
//	@Summary		Get annotations
//	@Description	Get paginated annotations, with options to filter by video stream.
//	@Tags			Annotations
//	@Produce		json
//	@Param			limit		query		int		false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset		query		int		false	"Number of results to skip."	minimum(0)
//	@Param			name		query		string	false	"Name to filter by."
//	@Param			videostream	query		int		false	"Video stream to filter by."
//	@Success		200			{object}	api.Result[services.AnnotationWithJoins]
//	@Failure		400			{object}	api.Failure
//	@Failure		401			{object}	api.Failure
//	@Failure		403			{object}	api.Failure
//	@Router			/api/v1/annotations [get]
func GetAnnotations(ctx *fiber.Ctx) error {
	qry := new(GetAnnotationsQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	annotations, err := services.GetAnnotations(qry.Limit, qry.Offset, qry.Order, qry.VideoStream)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Apply Joins.
	joined := make([]services.AnnotationWithJoins, len(annotations))
	for i, annotation := range annotations {
		j, err := annotation.JoinFields()
		if err != nil {
			return api.DatastoreReadFailure(err)
		}
		joined[i] = *j
	}

	return ctx.JSON(api.Result[services.AnnotationWithJoins]{
		Results: joined,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(joined),
	})
}

// NewAnnotationBody describes the JSON body required for the CreateAnnotation endpoint.
type NewAnnotationBody struct {
	KeyPoints      []keypoint.KeyPoint `json:"keypoints"`
	Identification *int64              `json:"identification" example:"1234567890" validate:"optional"`
	VideostreamID  int64               `json:"videostream_id" example:"1234567890"`
}

// CreateAnnotation creates a new annotation.
//
//	@Summary		Create annotation
//	@Description	Roles required: <role-tag>Annotator</role-tag>, <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Creates a new annotation from provided JSON body.
//	@Tags			Annotations
//	@Accept			json
//	@Produce		json
//	@Param			body	body		NewAnnotationBody	true	"New Annotation"
//	@Success		201		{object}	services.AnnotationWithJoins
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/annotations [post]
func CreateAnnotation(ctx *fiber.Ctx) error {
	// Parse URL.
	var body NewAnnotationBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Get logged in user.
	annotator, ok := ctx.Locals("user").(*services.User)
	if !ok {
		return fmt.Errorf("failed to assert type: expected *services.User but got %T", ctx.Locals("user"))
	}
	if annotator == nil {
		return api.Unauthorized(fmt.Errorf("user not logged in"))
	}

	// Check logged in user is in annotator_list.
	videostream, err := services.GetVideoStreamByID(body.VideostreamID)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}
	if len(videostream.AnnotatorList) != 0 && !slices.Contains(videostream.AnnotatorList, annotator.ID) {
		return api.Forbidden(fmt.Errorf("logged in user is not within annotator list for this videostream (%d)", body.VideostreamID))
	}

	// Check if we have any identifications.
	var ids map[int64][]int64
	if body.Identification != nil {
		ids = map[int64][]int64{
			*body.Identification: {annotator.ID},
		}
	} else {
		ids = nil
	}

	// Write data to the datastore.
	annotation := services.AnnotationContents{
		KeyPoints:       body.KeyPoints,
		VideostreamID:   body.VideostreamID,
		CreatedByID:     annotator.ID,
		Identifications: ids,
	}

	created, err := services.CreateAnnotation(annotation)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return joined form.
	joined, err := created.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
}

// AddIdentification adds a new identification to an annotation.
//
//	@Summary		Add Identification
//	@Description	Roles required: <role-tag>Annotator</role-tag>, <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Adds a new identification to an existing annotation.
//	@Tags			Annotations
//	@Produce		json
//	@Param			id		path		int	true	"Annotation ID"	example(1234567890)
//	@Param			species	path		int	true	"Species ID"	example(1234567890)
//	@Success		201		{object}	services.AnnotationWithJoins
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Failure		404		{object}	api.Failure
//	@Router			/api/v1/annotations/{id}/identifications/{species_id} [post]
func AddIdentification(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	speciesID, err := strconv.ParseInt(ctx.Params("species_id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Get logged in user.
	creator, ok := ctx.Locals("user").(*services.User)
	if !ok {
		return fmt.Errorf("failed to assert type: expected *services.User but got %T", ctx.Locals("user"))
	}
	if creator == nil {
		return api.Unauthorized(fmt.Errorf("user not logged in"))
	}

	// Write data to the datastore.
	err = services.AddIdentification(id, creator.ID, speciesID)
	if err != nil {
		return err
	}

	// Get updated annotation.
	modified, err := services.GetAnnotationByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	joined, err := modified.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
}

// DeleteIdentification removes an identification from an annotation.
//
//	@Summary		Remove Identification
//	@Description	Roles required: <role-tag>Annotator</role-tag>, <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Removes an identification from an annotation.
//	@Tags			Annotations
//	@Produce		json
//	@Param			id		path		int	true	"Annotation ID"	example(1234567890)
//	@Param			species	path		int	true	"Species ID"	example(1234567890)
//	@Success		201		{object}	services.AnnotationWithJoins
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Failure		404		{object}	api.Failure
//	@Router			/api/v1/annotations/{id}/identifications/{species_id} [delete]
func DeleteIdentification(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	speciesID, err := strconv.ParseInt(ctx.Params("species_id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Get logged in user.
	creator, ok := ctx.Locals("user").(*services.User)
	if !ok {
		return fmt.Errorf("failed to assert type: expected *services.User but got %T", ctx.Locals("user"))
	}
	if creator == nil {
		return api.Unauthorized(fmt.Errorf("user not logged in"))
	}

	// Write data to the datastore.
	err = services.DeleteIdentification(id, creator.ID, speciesID)
	if err != nil {
		return err
	}

	// Get updated annotation.
	modified, err := services.GetAnnotationByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	joined, err := modified.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
}

// DeleteAnnotation deletes an annotation.
//
//	@Summary		Delete annotation
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete an annotation by providing the annotation ID.
//	@Tags			Annotations
//	@Param			id	path	int	true	"Annotation ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/annotations/{id} [delete]
func DeleteAnnotation(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Delete entity.
	err = services.DeleteAnnotation(id)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
