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

package services

import (
	"context"

	googlestore "cloud.google.com/go/datastore"

	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/entities"
	"github.com/ausocean/openfish/datastore"
)

// GetCaptureSourceByID gets a capture source when provided with an ID.
func GetCaptureSourceByID(id int64) (*entities.CaptureSource, error) {
	store := ds_client.Get()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	var captureSource entities.CaptureSource
	err := store.Get(context.Background(), key, &captureSource)
	if err != nil {
		return nil, err
	}
	return &captureSource, nil
}

func CaptureSourceExists(id int64) bool {
	store := ds_client.Get()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	var captureSource entities.CaptureSource
	err := store.Get(context.Background(), key, &captureSource)
	return err == nil
}

// GetCaptureSources gets a list of capture sources, filtering by name, location if specified.
func GetCaptureSources(limit int, offset int, name *string) ([]entities.CaptureSource, []int64, error) {
	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery(entities.CAPTURESOURCE_KIND, false)

	if name != nil {
		query.FilterField("Name", "=", name)
	}

	// TODO: implement filtering based on location

	query.Limit(limit)
	query.Offset(offset)

	var captureSources []entities.CaptureSource
	keys, err := store.GetAll(context.Background(), query, &captureSources)
	if err != nil {
		return []entities.CaptureSource{}, []int64{}, err
	}
	ids := make([]int64, len(captureSources))
	for i, k := range keys {
		ids[i] = k.ID
	}

	return captureSources, ids, nil
}

// CreateCaptureSource creates a new capture source.
func CreateCaptureSource(name string, lat float64, long float64, cameraHardware string, siteID *int64) (int64, error) {

	// Get a unique ID for the new capturesource.
	store := ds_client.Get()
	key := store.IncompleteKey(entities.CAPTURESOURCE_KIND)

	cs := entities.CaptureSource{
		Name:           name,
		Location:       googlestore.GeoPoint{Lat: lat, Lng: long},
		CameraHardware: cameraHardware,
		SiteID:         siteID,
	}

	// Add to datastore.
	key, err := store.Put(context.Background(), key, &cs)
	if err != nil {
		return 0, err
	}
	return key.ID, nil
}

// UpdateCaptureSource updates a capture source.
func UpdateCaptureSource(id int64, name *string, lat *float64, long *float64, cameraHardware *string, siteID *int64) error {
	// TODO: Check that capture source has no video streams associated with it.

	// Update data in the datastore.
	store := ds_client.Get()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	var captureSource entities.CaptureSource

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*entities.CaptureSource)
		if ok {
			if name != nil {
				v.Name = *name
			}
			if lat != nil && long != nil {
				v.Location = googlestore.GeoPoint{Lat: *lat, Lng: *long}
			}
			if cameraHardware != nil {
				v.CameraHardware = *cameraHardware
			}
			if siteID != nil {
				v.SiteID = siteID
			}
		}
	}, &captureSource)
}

// DeleteCaptureSource deletes a capture source.
func DeleteCaptureSource(id int64) error {
	// TODO: Check that capture source has no video streams associated with it.

	// Delete entity.
	store := ds_client.Get()
	key := store.IDKey(entities.CAPTURESOURCE_KIND, id)
	return store.Delete(context.Background(), key)
}
