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
	"strconv"
	"strings"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/services"
	"github.com/ausocean/openfish/api/types/capturesource"

	"github.com/gofiber/fiber/v2"
)

// CaptureSourceResult describes the JSON format for capture sources in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type CaptureSourceResult struct {
	ID             *int64  `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Location       *string `json:"location,omitempty"`
	CameraHardware *string `json:"camera_hardware,omitempty"`
	SiteID         *int64  `json:"site_id,omitempty"`
}

// FromCaptureSource creates a CaptureSourceResult from a capturesource.CaptureSource and key, formatting it according to the requested format.
func FromCaptureSource(captureSource *capturesource.CaptureSource, id int64, format *api.Format) CaptureSourceResult {
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

// UpdateCaptureSourceBody describes the JSON format required for the UpdateCaptureSource endpoint.
type UpdateCaptureSourceBody struct {
	Name           *string `json:"name"`            // Optional.
	Location       *string `json:"location"`        // Optional.
	CameraHardware *string `json:"camera_hardware"` // Optional.
	SiteID         *int64  `json:"site_id"`         // Optional.
}

// parseGeoPoint converts a string containing two comma-separated values into a GeoPoint.
func parseGeoPoint(location string) (float64, float64, error) {
	errMsg := "invalid location string: %w"

	parts := strings.Split(location, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf(errMsg, "string split failed")
	}
	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf(errMsg, err)
	}
	long, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf(errMsg, err)
	}
	return lat, long, nil
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
	captureSource, err := services.GetCaptureSourceByID(id)

	// Format result.
	result := FromCaptureSource(captureSource, id, format)
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
	captureSources, ids, err := services.GetCaptureSources(qry.Limit, qry.Offset, qry.Name)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]CaptureSourceResult, len(captureSources))
	for i := range captureSources {
		results[i] = FromCaptureSource(&captureSources[i], ids[i], format)
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
	// Parse body.
	var body CreateCaptureSourceBody

	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	lat, long, err := parseGeoPoint(body.Location)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create capture source entity and add to the datastore.
	id, err := services.CreateCaptureSource(body.Name, lat, long, body.CameraHardware, body.SiteID)

	// Return ID of created capture source.
	return ctx.JSON(CaptureSourceResult{
		ID: &id,
	})
}

// UpdateCaptureSource updates a capture source.
func UpdateCaptureSource(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse body.
	var body UpdateCaptureSourceBody
	if ctx.BodyParser(&body) != nil {
		return api.InvalidRequestJSON(err)
	}

	var lat, long *float64
	if body.Location != nil {
		*lat, *long, err = parseGeoPoint(*body.Location)
		if err != nil {
			return api.InvalidRequestJSON(err)
		}
	}

	// Update data in the datastore.
	err = services.UpdateCaptureSource(id, body.Name, lat, long, body.CameraHardware, body.SiteID)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}

// DeleteCaptureSource deletes a capture source.
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
