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
	"github.com/ausocean/openfish/api/utils"

	"github.com/gofiber/fiber/v2"
)

type CaptureSourceResult struct {
	ID             *int    `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Location       *string `json:"location,omitempty"`
	CameraHardware *string `json:"camera_hardware,omitempty"`
}

type CaptureSourcePost struct {
	Name           string `json:"name"`
	Location       string `json:"location"`
	CameraHardware string `json:"camera_hardware"`
}

func GetCaptureSourceByID(ctx *fiber.Ctx) error {
	// Parse URL.
	format := utils.GetFormat(ctx)
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(ctx)
	}

	// Fetch data from the datastore.
	store := ds_client.Get()
	key := store.IDKey("CaptureSource", id)
	var captureSource model.CaptureSource
	if store.Get(context.Background(), key, &captureSource) != nil {
		return api.DatastoreReadFailure(ctx)
	}

	// Format result.
	var result CaptureSourceResult

	if format.Requires("id") {
		id := int(id)
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

	return ctx.JSON(result)
}

func GetCaptureSources(ctx *fiber.Ctx) error {
	// Parse URL.
	name := ctx.Query("name")
	location := ctx.Query("location")
	format := utils.GetFormat(ctx)
	limit, offset := utils.GetLimitAndOffset(ctx, 20)

	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery("CaptureSource", false)

	if name != "" {
		query.FilterField("Name", "=", name)
	}

	// TODO: implement filtering based on location
	_ = location

	query.Limit(limit)
	query.Offset(offset)

	var captureSources []model.CaptureSource
	keys, err := store.GetAll(context.Background(), query, &captureSources)
	if err != nil {
		return api.DatastoreReadFailure(ctx)
	}

	// Format results.
	results := make([]CaptureSourceResult, len(captureSources))
	for i := range captureSources {
		if format.Requires("id") {
			id := int(keys[i].ID)
			results[i].ID = &id
		}
		if format.Requires("name") {
			results[i].Name = &captureSources[i].Name
		}
		if format.Requires("location") {
			location := fmt.Sprintf("%f,%f", captureSources[i].Location.Lat, captureSources[i].Location.Lng)
			results[i].Location = &location
		}
		if format.Requires("camera_hardware") {
			results[i].CameraHardware = &captureSources[i].CameraHardware
		}
	}

	return ctx.JSON(api.Result[CaptureSourceResult]{
		Results: results,
		Offset:  offset,
		Limit:   limit,
		Total:   len(results),
	})
}

func CreateCaptureSource(ctx *fiber.Ctx) error {
	var body CaptureSourcePost

	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(ctx)
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
	}

	// Add to datastore.
	key, err = store.Put(context.Background(), key, &cs)
	if err != nil {
		print(err.Error())
		return api.DatastoreWriteFailure(ctx)
	}

	// Return ID of created capture source.
	id := int(key.ID)
	return ctx.JSON(CaptureSourceResult{
		ID: &id,
	})
}
