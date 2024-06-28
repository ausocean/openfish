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
	"errors"
	"fmt"
	"time"

	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/types/timespan"
	"github.com/ausocean/openfish/datastore"
)

// GetAnnotationByID gets an annotation from datastore when provided with an ID.
func GetAnnotationByID(id int64) (*entities.Annotation, error) {
	store := ds_client.Get()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	var annotation entities.Annotation
	err := store.Get(context.Background(), key, &annotation)
	if err != nil {
		return nil, err
	}

	return &annotation, nil
}

func AnnotationExists(id int64) bool {
	store := ds_client.Get()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	var annotation entities.Annotation
	err := store.Get(context.Background(), key, &annotation)
	return err == nil
}

// GetAnnotations gets a list of annotations, filtering by timespan, capturesource, observer & observation if specified.
func GetAnnotations(limit int, offset int, observer *string, observation map[string]string) ([]entities.Annotation, []int64, error) {
	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery(entities.ANNOTATION_KIND, false)

	// Filter by observer.
	if observer != nil {
		query.FilterField("Observer", "=", *observer)
	}

	// Filter by observation records.
	for k, v := range observation {
		if v == "*" {
			query.FilterField("ObservationKeys", "=", k)
		} else {
			query.FilterField("ObservationPairs", "=", fmt.Sprintf("%s:%s", k, v))
		}
	}

	query.Limit(limit)
	query.Offset(offset)

	var annotations []entities.Annotation
	keys, err := store.GetAll(context.Background(), query, &annotations)
	if err != nil {
		return []entities.Annotation{}, []int64{}, err
	}
	ids := make([]int64, len(annotations))
	for i, k := range keys {
		ids[i] = k.ID
	}

	return annotations, ids, nil
}

// CreateAnnotation creates a new annotation.
func CreateAnnotation(videoStreamID int64, timeSpan timespan.TimeSpan, boundingBox *entities.BoundingBox, observer string, observation map[string]string) (int64, error) {
	// Convert observation map into a format the datastore can take.
	obsKeys := make([]string, 0, len(observation))
	obsPairs := make([]string, 0, len(observation))

	for k, v := range observation {
		obsKeys = append(obsKeys, k)
		obsPairs = append(obsPairs, fmt.Sprintf("%s:%s", k, v))
	}

	// Create annotation entity and add to the datastore.
	an := entities.Annotation{
		VideoStreamID:    videoStreamID,
		TimeSpan:         timeSpan,
		BoundingBox:      boundingBox,
		Observer:         observer,
		ObservationPairs: obsPairs,
		ObservationKeys:  obsKeys,
	}

	// Verify VideoStream exists.
	if !VideoStreamExists(int64(videoStreamID)) {
		return 0, errors.New("VideoStream does not exist")
	}

	// Get a unique ID for the new annotation.
	store := ds_client.Get()
	key := store.IncompleteKey(entities.ANNOTATION_KIND)
	key, err := store.Put(context.Background(), key, &an)
	if err != nil {
		return 0, err
	}

	// Return ID of created video stream.
	return key.ID, nil
}

// UpdateAnnotation updates an Annotation.
func UpdateAnnotation(id int64, streamURL *string, captureSource *int64, startTime *time.Time, endTime *time.Time) error {

	// Update data in the datastore.
	store := ds_client.Get()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	var videoStream entities.VideoStream

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*entities.VideoStream)
		if ok {
			if streamURL != nil {
				v.StreamUrl = *streamURL
			}
			if captureSource != nil {
				// TODO: Check that captureSource exists.
				v.CaptureSource = *captureSource
			}
			if startTime != nil {
				v.StartTime = *startTime
			}
			if endTime != nil {
				v.EndTime = endTime
			}
		}
	}, &videoStream)
}

// DeleteAnnotation deletes an annotation.
func DeleteAnnotation(id int64) error {
	// TODO: Check that annotation exists.

	// Delete entity.
	store := ds_client.Get()
	key := store.IDKey(entities.ANNOTATION_KIND, id)
	return store.Delete(context.Background(), key)
}
