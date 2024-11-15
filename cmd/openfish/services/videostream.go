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
	"fmt"
	"time"

	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/types/timespan"
	"github.com/ausocean/openfish/datastore"
)

// GetVideoStreamByID gets a video stream when provided with an ID.
func GetVideoStreamByID(id int64) (*entities.VideoStream, error) {
	store := ds_client.Get()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	var videoStream entities.VideoStream
	err := store.Get(context.Background(), key, &videoStream)
	if err != nil {
		return nil, err
	}

	return &videoStream, nil
}

func VideoStreamExists(id int64) bool {
	store := ds_client.Get()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	var videoStream entities.VideoStream
	err := store.Get(context.Background(), key, &videoStream)
	return err == nil
}

// GetVideoStreams gets a list of video streams, filtering by timespan, capturesource if specified.
func GetVideoStreams(limit int, offset int, timespan *timespan.TimeSpan, captureSource *int64) ([]entities.VideoStream, []int64, error) {
	// Fetch data from the datastore.
	store := ds_client.Get()
	query := store.NewQuery(entities.VIDEOSTREAM_KIND, false)

	if captureSource != nil {
		query.FilterField("CaptureSource", "=", captureSource)
	}

	if timespan != nil {
		// BUG: Because of a limitation of google cloud's datastore, we can only use inequality filters on one
		// field at a time. This query needs to filter out videostreams that started after the specified range.
		// https://github.com/ausocean/openfish/issues/23
		query.FilterField("EndTime", ">=", timespan.Start)
	}

	// TODO: implement filtering based on location

	query.Limit(limit)
	query.Offset(offset)

	var videoStreams []entities.VideoStream
	keys, err := store.GetAll(context.Background(), query, &videoStreams)
	if err != nil {
		return []entities.VideoStream{}, []int64{}, err
	}
	ids := make([]int64, len(videoStreams))
	for i, k := range keys {
		ids[i] = k.ID
	}

	return videoStreams, ids, nil
}

// CreateVideoStream puts a video stream in the datastore, checking if the capture source exists.
func CreateVideoStream(streamURL string, captureSource int64, startTime time.Time, endTime *time.Time, annotatorList []string) (int64, error) {

	// Verify CaptureSource exists.
	if !CaptureSourceExists(captureSource) {
		return 0, fmt.Errorf("CaptureSource does not exist")
	}

	// Create VideoStream entity.
	store := ds_client.Get()
	key := store.IncompleteKey(entities.VIDEOSTREAM_KIND)

	vs := entities.VideoStream{
		StreamUrl:     streamURL,
		CaptureSource: captureSource,
		StartTime:     startTime,
		EndTime:       endTime,
		AnnotatorList: annotatorList,
	}
	key, err := store.Put(context.Background(), key, &vs)
	if err != nil {
		return 0, err
	}

	// Return ID of created video stream.
	return key.ID, nil
}

// UpdateCaptureSource updates a capture source.
func UpdateVideoStream(id int64, streamURL *string, captureSource *int64, startTime *time.Time, endTime *time.Time, annotatorList *[]string) error {

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
			if annotatorList != nil {
				v.AnnotatorList = *annotatorList
			}
		}
	}, &videoStream)
}

// DeleteCaptureSource deletes a capture source.
func DeleteVideoStream(id int64) error {
	// TODO: Check that video stream has no annotations associated with it.

	// Delete entity.
	store := ds_client.Get()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	return store.Delete(context.Background(), key)
}
