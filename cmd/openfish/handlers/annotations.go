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
	"strings"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/timespan"

	"github.com/gofiber/fiber/v2"
)

// AnnotationResult describes the JSON format for annotations in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type AnnotationResult struct {
	ID            *int64                `json:"id,omitempty"`
	VideoStreamID *int64                `json:"videostreamId,omitempty"`
	TimeSpan      *timespan.TimeSpan    `json:"timespan,omitempty"`
	BoundingBox   *entities.BoundingBox `json:"boundingBox,omitempty"`
	Observer      *string               `json:"observer,omitempty"`
	Observation   map[string]string     `json:"observation,omitempty"`
}

// FromAnnotation creates an AnnotationResult from a entities.Annotation and key, formatting it according to the requested format.
func FromAnnotation(annotation *entities.Annotation, id int64, format *api.Format) AnnotationResult {
	var result AnnotationResult
	if format.Requires("id") {
		result.ID = &id
	}
	if format.Requires("videostream_id") {
		result.VideoStreamID = &annotation.VideoStreamID
	}
	if format.Requires("timespan") {
		result.TimeSpan = &annotation.TimeSpan
	}
	if format.Requires("bounding_box") {
		result.BoundingBox = annotation.BoundingBox
	}
	if format.Requires("observer") {
		result.Observer = &annotation.Observer
	}
	if format.Requires("observation") {
		result.Observation = annotation.Observation
	}

	return result
}

// GetAnnotationsQuery describes the URL query parameters required for the GetAnnotations endpoint.
type GetAnnotationsQuery struct {
	TimeSpan      *string           `query:"timespan"`       // Optional. TODO: choose more appropriate type.
	CaptureSource *int64            `query:"capture_source"` // Optional.
	Observer      *string           `query:"observer"`       // Optional.
	Observation   map[string]string `query:"observation"`    // Optional.
	api.LimitAndOffset
	api.Format
}

// CreateAnnotationBody describes the JSON format required for the CreateAnnotation endpoint.
//
// ID is omitted because it is chosen automatically.
// BoundingBox is optional because some annotations might not be described by a rectangular area.
type CreateAnnotationBody struct {
	VideoStreamID int64                 `json:"videostreamId"`
	TimeSpan      timespan.TimeSpan     `json:"timespan"`
	BoundingBox   *entities.BoundingBox `json:"boundingBox"` // Optional.
	Observer      string                `json:"observer"`
	Observation   map[string]string     `json:"observation"`
}

// GetAnnotationByID gets an annotation when provided with an ID.
func GetAnnotationByID(ctx *fiber.Ctx) error {
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
	annotation, err := services.GetAnnotationByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromAnnotation(annotation, id, format)
	return ctx.JSON(result)
}

// GetAnnotations gets a list of annotations, filtering by timespan, capturesource, observer & observation if specified.
func GetAnnotations(ctx *fiber.Ctx) error {
	qry := new(GetAnnotationsQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	format := new(api.Format)
	if err := ctx.QueryParser(format); err != nil {
		return api.InvalidRequestURL(err)
	}

	// NOTE: fiber's QueryParser does not handle map[string]string so we need to parse the query manually.
	// This can be revisited if PR https://github.com/gofiber/fiber/issues/2524 is merged.
	qry.Observation = make(map[string]string)
	for k, v := range ctx.Queries() {
		if strings.HasPrefix(k, "observation[") && strings.HasSuffix(k, "]") {
			k = strings.TrimPrefix(k, "observation[")
			k = strings.TrimSuffix(k, "]")
			qry.Observation[k] = v
		}
	}

	// Fetch data from the datastore.
	annotations, ids, err := services.GetAnnotations(qry.Limit, qry.Offset, qry.Observer, qry.Observation)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]AnnotationResult, len(annotations))
	for i := range annotations {
		results[i] = FromAnnotation(&annotations[i], ids[i], format)
	}

	return ctx.JSON(api.Result[AnnotationResult]{
		Results: results,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(results),
	})
}

// CreateAnnotation creates a new annotation.
func CreateAnnotation(ctx *fiber.Ctx) error {
	// Parse URL.
	var body CreateAnnotationBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Get logged in user.
	observer := ctx.Locals("email").(string)

	// Check logged in user is in annotator_list.
	videostream, err := services.GetVideoStreamByID(body.VideoStreamID)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}
	if len(videostream.AnnotatorList) != 0 && !slices.Contains(videostream.AnnotatorList, observer) {
		return api.Forbidden(fmt.Errorf("logged in user is not within annotator list for this videostream (%d)", body.VideoStreamID))
	}

	// Write data to the datastore.
	id, err := services.CreateAnnotation(body.VideoStreamID,
		body.TimeSpan, body.BoundingBox,
		observer, body.Observation)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(AnnotationResult{
		ID: &id,
	})
}

// TODO: Implement UpdateAnnotation.

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
