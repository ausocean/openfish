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

// services contains the main logic for the OpenFish API.
package services

import (
	"context"

	"github.com/ausocean/cloud/datastore"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/types/latlong"
)

// CaptureSource stores information about the camera that produces video streams.
type CaptureSource struct {
	ID int64 `json:"id" example:"1234567890"`
	CaptureSourceContents
}

// CaptureSourceContents is the contents of a CaptureSource.
type CaptureSourceContents struct {
	Name           string          `json:"name" example:"Stony Point Cuttle Cam"`                       // Name of rig or camera.
	Location       latlong.LatLong `json:"location" swaggertype:"string" example:"-32.12345,139.12345"` // Where the rig or camera is located.
	CameraHardware string          `json:"camera_hardware" example:"pi cam v2 (wide angle lens)"`       // Short description of the camera hardware.
	SiteID         *int64          `json:"site_id" example:"246813579"`
}

// PartialCaptureSourceContents is for updating a capture source with a partial update (such as a PATCH request).
type PartialCaptureSourceContents struct {
	Name           *string          `json:"name,omitempty" example:"Stony Point Cuttle Cam"`                       // Name of rig or camera.
	Location       *latlong.LatLong `json:"location,omitempty" swaggertype:"string" example:"-32.12345,139.12345"` // Where the rig or camera is located.
	CameraHardware *string          `json:"camera_hardware,omitempty" example:"pi cam v2 (wide angle lens)"`       // Short description of the camera hardware.
	SiteID         *int64           `json:"site_id,omitempty" example:"246813579"`
}

// CaptureSourceSummary is a summary of a capture source.
type CaptureSourceSummary struct {
	ID   int64  `json:"id" example:"1234567890"`
	Name string `json:"name" example:"Stony Point Cuttle Cam"` // Name of rig or camera.
}

// ToSummary converts a CaptureSource to a CaptureSourceSummary.
func (c *CaptureSource) ToSummary() CaptureSourceSummary {
	return CaptureSourceSummary{
		ID:   c.ID,
		Name: c.Name,
	}
}

// CaptureSourceContentsFromEntity converts an entities.CaptureSource to a CaptureSourceContents.
func CaptureSourceContentsFromEntity(c entities.CaptureSource) CaptureSourceContents {
	l, _ := latlong.New(c.Location.Lat, c.Location.Lng)
	return CaptureSourceContents{
		Name:           c.Name,
		Location:       l,
		CameraHardware: c.CameraHardware,
		SiteID:         c.SiteID,
	}
}

// ToEntity converts a CaptureSourceContents to an entities.CaptureSource for storage in the datastore.
func (c *CaptureSourceContents) ToEntity() entities.CaptureSource {
	return entities.CaptureSource{
		Name:           c.Name,
		Location:       c.Location.GeoPoint,
		CameraHardware: c.CameraHardware,
		SiteID:         c.SiteID,
	}
}

// GetCaptureSourceByID gets a capture source when provided with an ID.
func GetCaptureSourceByID(id int64) (*CaptureSource, error) {
	store := globals.GetStore()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	var e entities.CaptureSource
	err := store.Get(context.Background(), key, &e)
	if err != nil {
		return nil, err
	}

	return &CaptureSource{
		ID:                    id,
		CaptureSourceContents: CaptureSourceContentsFromEntity(e),
	}, nil
}

// CaptureSourceExists checks if a capture source exists with the given ID.
func CaptureSourceExists(id int64) bool {
	store := globals.GetStore()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	var captureSource entities.CaptureSource
	err := store.Get(context.Background(), key, &captureSource)
	return err == nil
}

// GetCaptureSources gets a list of capture sources, filtering by name, location if specified.
func GetCaptureSources(limit int, offset int, name *string) ([]CaptureSource, error) {
	// Fetch data from the datastore.
	store := globals.GetStore()
	query := store.NewQuery(entities.CAPTURESOURCE_KIND, false)

	if name != nil {
		query.FilterField("Name", "=", name)
	}

	// TODO: implement filtering based on location

	query.Limit(limit)
	query.Offset(offset)

	var srcs []entities.CaptureSource
	ids, err := store.GetAll(context.Background(), query, &srcs)
	if err != nil {
		return []CaptureSource{}, err
	}

	results := make([]CaptureSource, len(srcs))
	for i := range srcs {
		results[i] = CaptureSource{
			ID:                    ids[i].ID,
			CaptureSourceContents: CaptureSourceContentsFromEntity(srcs[i]),
		}
	}

	return results, nil

}

// CreateCaptureSource creates a new capture source.
func CreateCaptureSource(contents CaptureSourceContents) (*CaptureSource, error) {

	// Get a unique ID for the new capturesource.
	store := globals.GetStore()
	key := store.IncompleteKey(entities.CAPTURESOURCE_KIND)

	c := contents.ToEntity()

	// Add to datastore.
	key, err := store.Put(context.Background(), key, &c)
	if err != nil {
		return nil, err
	}

	// Return newly created capture source.
	created := CaptureSource{
		ID:                    key.ID,
		CaptureSourceContents: contents,
	}
	return &created, nil
}

// UpdateCaptureSource updates a capture source.
func UpdateCaptureSource(id int64, updates PartialCaptureSourceContents) error {
	// TODO: Check that capture source has no video streams associated with it.

	// Update data in the datastore.
	store := globals.GetStore()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	var captureSource entities.CaptureSource

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		c, ok := e.(*entities.CaptureSource)
		if ok {
			if updates.Name != nil {
				c.Name = *updates.Name
			}
			if updates.Location != nil {
				c.Location.Lat = updates.Location.Lat
				c.Location.Lng = updates.Location.Lng
			}
			if updates.CameraHardware != nil {
				c.CameraHardware = *updates.CameraHardware
			}
			if updates.SiteID != nil {
				c.SiteID = updates.SiteID
			}
		}
	}, &captureSource)
}

// DeleteCaptureSource deletes a capture source.
func DeleteCaptureSource(id int64) error {
	// TODO: Check that capture source has no video streams associated with it.

	// Delete entity.
	store := globals.GetStore()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	return store.Delete(context.Background(), key)
}
