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

	"cloud.google.com/go/datastore"
	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/model"

	"github.com/gofiber/fiber/v2"
)

// CaptureSourceResult describes the JSON format for capture sources in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type CaptureSourceResult struct {
	ID             *int    `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Location       *string `json:"location,omitempty"`
	CameraHardware *string `json:"camera_hardware,omitempty"`
	SiteID         *int64  `json:"site_id,omitempty"`
}

// FromCaptureSource creates a CaptureSourceResult from a model.CaptureSource and key, formatting it according to the requested format.
func FromCaptureSource(captureSource *model.CaptureSource, id int, format *api.Format) CaptureSourceResult {
	var result CaptureSourceResult
	if format.Requires("id") {
		result.ID = &id
	}
	if format.Requires("name") {
		result.Name = &captureSource.Name
	}
	if format.Requires("location") {
		location := fmt.Sprintf("%f,%f", captureSource.Location.Lat, captureSource.Location.Lng)
		result.Location = &location
	}
	if format.Requires("camera_hardware") {
		result.CameraHardware = &captureSource.CameraHardware
	}
	if format.Requires("site_id") {
		result.SiteID = captureSource.SiteID
	}
	return result
}

// GetCaptureSourcesQuery describes the URL query parameters required for the GetCaptureSources endpoint.
type GetCaptureSourcesQuery struct {
	Name     *string `query:"name"`     // Optional.
	Location *string `query:"location"` // Optional.
	api.LimitAndOffset
	// TODO: Code could be simplified if api.Format could be embedded here. Testing shows it cannot, reason unknown.
}

// CreateCaptureSourceBody describes the JSON format required for the CreateCaptureSource endpoint.
// ID is omitted because it is chosen automatically. All other fields are required.
type CreateCaptureSourceBody struct {
	Name           string `json:"name"`
	Location       string `json:"location"`
	CameraHardware string `json:"camera_hardware"`
	SiteID         *int64 `json:"site_id"` // Optional.
}

// GetCaptureSourceByID gets a capture source when provided with an ID.
func GetCaptureSourceByID(ctx *fiber.Ctx) error {
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
	key := store.IDKey("CaptureSource", id)
	var captureSource model.CaptureSource
	if store.Get(context.Background(), key, &captureSource) != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromCaptureSource(&captureSource, int(id), format)

	return ctx.JSON(result)
}

// GetCaptureSources gets a list of capture sources, filtering by name, location if specified.
func GetCaptureSources(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetCaptureSourcesQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	format := new(api.Format)
	if err := ctx.QueryParser(format); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery("CaptureSource", false)

	if qry.Name != nil {
		query.FilterField("Name", "=", qry.Name)
	}

	// TODO: implement filtering based on location

	query.Limit(qry.Limit)
	query.Offset(qry.Offset)

	var captureSources []model.CaptureSource
	keys, err := store.GetAll(context.Background(), query, &captureSources)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]CaptureSourceResult, len(captureSources))
	for i := range captureSources {
		results[i] = FromCaptureSource(&captureSources[i], int(keys[i].ID), format)
	}

	return ctx.JSON(api.Result[CaptureSourceResult]{
		Results: results,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(results),
	})
}

// CreateCaptureSource creates a new capture source.
func CreateCaptureSource(ctx *fiber.Ctx) error {
	var body CreateCaptureSourceBody

	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Parse location.
	parts := strings.Split(body.Location, ",")
	lat, _ := strconv.ParseFloat(parts[0], 64)
	long, _ := strconv.ParseFloat(parts[1], 64)

	// Get a unique ID for the new capturesource.
	store := ds_client.Get()
	key := store.IncompleteKey("CaptureSource")

	// Create capture source entity and add to the datastore.
	cs := model.CaptureSource{
		Name:           body.Name,
		Location:       datastore.GeoPoint{Lat: lat, Lng: long},
		CameraHardware: body.CameraHardware,
		SiteID:         body.SiteID,
	}

	// Add to datastore.
	key, err = store.Put(context.Background(), key, &cs)
	if err != nil {
		print(err.Error())
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created capture source.
	id := int(key.ID)
	return ctx.JSON(CaptureSourceResult{
		ID: &id,
	})
}
