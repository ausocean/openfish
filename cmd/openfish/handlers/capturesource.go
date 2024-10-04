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

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/services"

	"github.com/gofiber/fiber/v2"
)

// CaptureSourceResult describes the JSON format for capture sources in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
//
//	@Description	contains information about something that produces a video stream
type CaptureSourceResult struct {
	ID             *int64  `json:"id,omitempty" example:"1234567890"`                               // Unique ID of the capture source.
	Name           *string `json:"name,omitempty" example:"Stony Point Cuttle Cam"`                 // Name of rig or camera.
	Location       *string `json:"location,omitempty" example:"-32.12345,139.12345"`                // Where the rig or camera is located.
	CameraHardware *string `json:"camera_hardware,omitempty" example:"pi cam v2 (wide angle lens)"` // Short description of the camera hardware.
	SiteID         *int64  `json:"site_id,omitempty" example:"246813579"`                           // Site ID is used to reference sites in OceanBench. Optional.
}

// EntityIDResult contains the ID of a newly created entity.
//
//	@Description	ID of newly created entity.
type EntityIDResult struct {
	ID *int64 `json:"id,omitempty" example:"1234567890"` // Unique ID of the entity.
}

// FromCaptureSource creates a CaptureSourceResult from a entities.CaptureSource and key, formatting it according to the requested format.
func FromCaptureSource(captureSource *entities.CaptureSource, id int64, format *api.Format) CaptureSourceResult {
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
	Name           string `json:"name" example:"Stony Point Cuttle Cam" validate:"required"`                 // Name of rig or camera.
	Location       string `json:"location" example:"-32.12345,139.12345" validate:"required"`                // Location of the rig or camera.
	CameraHardware string `json:"camera_hardware" example:"pi cam v2 (wide angle lens)" validate:"required"` // Short description of the camera hardware.
	SiteID         *int64 `json:"site_id" example:"246813579" validate:"optional"`                           // ID used to reference sites in OceanBench.
}

// UpdateCaptureSourceBody describes the JSON format required for the UpdateCaptureSource endpoint.
type UpdateCaptureSourceBody struct {
	Name           *string `json:"name" example:"Stony Point Cuttle Cam" validate:"optional"`                 // Name of rig or camera.
	Location       *string `json:"location" example:"-32.12345,139.12345" validate:"optional"`                // Location of the rig or camera.
	CameraHardware *string `json:"camera_hardware" example:"pi cam v2 (wide angle lens)" validate:"optional"` // Short description of the camera hardware.
	SiteID         *int64  `json:"site_id" example:"246813579" validate:"optional"`                           // ID used to reference sites in OceanBench.
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
//
//	@Summary		Get capture source by ID
//	@Description	Gets a capture source when provided with an ID.
//	@Tags			Capture Sources
//	@Produce		json
//	@Param			id	path		int	true	"Capture Source ID"	example(1234567890)
//	@Success		200	{object}	CaptureSourceResult
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/capturesources/{id} [get]
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
//
//	@Summary		Get capture sources
//	@Description	Get paginated capture sources, with options to filter by name and location.
//	@Tags			Capture Sources
//	@Produce		json
//	@Param			limit	query		int		false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset	query		int		false	"Number of results to skip."	minimum(0)
//	@Param			name	query		string	false	"Name to filter by."
//	@Success		200		{object}	api.Result[CaptureSourceResult]
//	@Failure		400		{object}	api.Failure
//	@Router			/api/v1/capturesources [get]
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
//
//	@Summary		Create capture source
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Creates a new capture source from provided JSON body.
//	@Tags			Capture Sources
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateCaptureSourceBody	true	"New Capture Source"
//	@Success		201		{object}	EntityIDResult
//	@Failure		400		{object}	api.Failure
//	@Router			/api/v1/capturesources [post]
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
	return ctx.JSON(EntityIDResult{
		ID: &id,
	})
}

// UpdateCaptureSource updates a capture source.
//
//	@Summary		Update capture source
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Partially update a capture source by specifying the properties to update.
//	@Tags			Capture Sources
//	@Accept			json
//	@Param			id		path	int						true	"Capture Source ID"	example(1234567890)
//	@Param			body	body	UpdateCaptureSourceBody	true	"Update Capture Source"
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Router			/api/v1/capturesources/{id} [patch]
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
//
//	@Summary		Delete capture source
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete a capture source by providing the capture source ID.
//	@Tags			Capture Sources
//	@Param			id	path	int	true	"Capture Source ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
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
