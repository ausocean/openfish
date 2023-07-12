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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/model"
	"github.com/ausocean/openfish/datastore"

	"github.com/gofiber/fiber/v2"
)

// AnnotationResult describes the JSON format for annotations in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type AnnotationResult struct {
	ID            *int64             `json:"id,omitempty"`
	VideoStreamID *int64             `json:"videostreamId,omitempty"`
	TimeSpan      *model.TimeSpan    `json:"timespan,omitempty"`
	BoundingBox   *model.BoundingBox `json:"boundingBox,omitempty"`
	Observer      *string            `json:"observer,omitempty"`
	Observation   map[string]string  `json:"observation,omitempty"`
}

// FromAnnotation creates an AnnotationResult from a model.Annotation and key, formatting it according to the requested format.
func FromAnnotation(annotation *model.Annotation, key *datastore.Key, format *api.Format) AnnotationResult {
	var result AnnotationResult
	if format.Requires("id") {
		result.ID = &key.ID
	}
	if format.Requires("videostream_id") {
		// result.VideoStreamID = &annotation.VideoStreamID
		result.VideoStreamID = &key.Parent.ID
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
		observation := make(map[string]string)
		for _, o := range annotation.ObservationPairs {
			parts := strings.Split(o, ":")
			observation[parts[0]] = parts[1]
		}
		result.Observation = observation
	}

	return result
}

// GetAnnotationsQuery describes the URL query parameters required for the GetAnnotations endpoint.
type GetAnnotationsQuery struct {
	TimeSpan      *string           `query:"timespan"`       // Optional. TODO: choose more appropriate type.
	CaptureSource *int64            `query:"capture_source"` // Optional.
	VideoStream   *int64            `query:"video_stream"`   // Optional.
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
	VideoStreamID int64              `json:"videostreamId"`
	TimeSpan      model.TimeSpan     `json:"timespan"`
	BoundingBox   *model.BoundingBox `json:"boundingBox"` // Optional.
	Observer      string             `json:"observer"`
	Observation   map[string]string  `json:"observation"`
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
	store := ds_client.Get()
	key := store.IDKey("Annotation", id, nil)
	var annotation model.Annotation
	err = store.Get(context.Background(), key, &annotation)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromAnnotation(&annotation, key, format)

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
	store := ds_client.Get()
	query := store.NewQuery("Annotation", false)

	// Filter by observer.
	if qry.Observer != nil {
		query.FilterField("Observer", "=", *qry.Observer)
	}

	// Filter by observation records.
	for k, v := range qry.Observation {
		if v == "*" {
			query.FilterField("ObservationKeys", "=", k)
		} else {
			query.FilterField("ObservationPairs", "=", fmt.Sprintf("%s:%s", k, v))
		}
	}

	// Filter by videostream.
	if qry.VideoStream != nil {
		parent := store.IDKey("VideoStream", *qry.VideoStream, nil)
		query.Ancestor(parent)
	}

	// TODO: support filtering by timespan, capturesource, locations.

	query.Limit(qry.Limit)
	query.Offset(qry.Offset)

	var annotations []model.Annotation
	keys, err := store.GetAll(context.Background(), query, &annotations)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]AnnotationResult, len(annotations))
	for i := range annotations {
		results[i] = FromAnnotation(&annotations[i], keys[i], format)
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
	var body CreateAnnotationBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Convert observation map into a format the datastore can take.
	obsKeys := make([]string, 0, len(body.Observation))
	obsPairs := make([]string, 0, len(body.Observation))

	for k, v := range body.Observation {
		obsKeys = append(obsKeys, k)
		obsPairs = append(obsPairs, fmt.Sprintf("%s:%s", k, v))
	}

	// Create annotation entity and add to the datastore.
	an := model.Annotation{
		TimeSpan:         body.TimeSpan,
		BoundingBox:      body.BoundingBox,
		Observer:         body.Observer,
		ObservationPairs: obsPairs,
		ObservationKeys:  obsKeys,
	}

	// Store annotation entity in the datastore.
	store := ds_client.Get()
	parent := store.IDKey("VideoStream", body.VideoStreamID, nil)
	key := store.IncompleteKey("Annotation", parent)
	key, err = store.Put(context.Background(), key, &an)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(AnnotationResult{
		ID: &key.ID,
	})
}
