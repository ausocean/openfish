/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023-2025, The OpenFish Contributors.

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

	"github.com/ausocean/cloud/datastore"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/types/timespan"
	"github.com/ausocean/openfish/cmd/openfish/types/timezone"
)

// VideoStream is a video stream registered to OpenFish, for annotation and playback.
type VideoStream struct {
	ID int64 `json:"id,omitempty" example:"1234567890"`
	VideoStreamContents
}

// VideoStreamContents is the contents of a video stream.
type VideoStreamContents struct {
	StartTime     time.Time  `json:"startTime" example:"2023-05-25T08:00:00Z"`
	EndTime       *time.Time `json:"endTime" example:"2023-05-25T16:30:00Z"`
	AnnotatorList []int64    `json:"annotator_list" example:"1234567890"`
	BaseVideoStreamFields
}

// BaseVideoStreamFields contains shared fields.
type BaseVideoStreamFields struct {
	TimeZone      timezone.TimeZone `json:"timezone" swaggertype:"string" example:"Australia/Adelaide"`
	StreamURL     string            `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk"`
	CaptureSource int64             `json:"capturesource" example:"1234567890"`
}

// PartialVideoStreamContents represents optional fields that can be updated for a video stream.
type PartialVideoStreamContents struct {
	StartTime     *time.Time         `json:"startTime,omitempty" example:"2023-05-25T08:00:00Z"`
	EndTime       *time.Time         `json:"endTime,omitempty" example:"2023-05-25T16:30:00Z"`
	TimeZone      *timezone.TimeZone `json:"timezone" swaggertype:"string" example:"Australia/Adelaide"`
	StreamURL     *string            `json:"stream_url,omitempty" example:"https://www.youtube.com/watch?v=abcdefghijk"`
	CaptureSource *int64             `json:"capturesource,omitempty" example:"1234567890"`
	AnnotatorList *[]int64           `json:"annotator_list,omitempty" example:"1234567890"`
}

// VideoStreamWithJoins represents a video stream with its related entities joined.
type VideoStreamWithJoins struct {
	ID            int64                `json:"id,omitempty" example:"1234567890"`
	StartTime     time.Time            `json:"startTime" example:"2023-05-25T08:00:00Z"`
	EndTime       *time.Time           `json:"endTime" example:"2023-05-25T16:30:00Z"`
	TimeZone      timezone.TimeZone    `json:"timezone" swaggertype:"string" example:"Australia/Adelaide"`
	StreamURL     string               `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk"`
	CaptureSource CaptureSourceSummary `json:"capturesource"`
	AnnotatorList []PublicUser         `json:"annotator_list"`
}

// VideoStreamSummary is a summary of a video stream.
type VideoStreamSummary struct {
	ID        int64  `json:"id" example:"1234567890"`
	StreamURL string `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk"`
}

// ToSummary converts a VideoStream to a VideoStreamSummary.
func (v *VideoStream) ToSummary() VideoStreamSummary {
	return VideoStreamSummary{
		ID:        v.ID,
		StreamURL: v.StreamURL,
	}
}

// JoinFields joins the foreign key fields of a videostream with their respective entities.
func (v *VideoStream) JoinFields() (*VideoStreamWithJoins, error) {

	// Get video stream details.
	captureSource, err := GetCaptureSourceByID(v.CaptureSource)
	if err != nil {
		return nil, err // TODO: Return a more informative message.
	}

	// Get user details.
	annotatorList := make([]PublicUser, 0, len(v.AnnotatorList))
	for _, uid := range v.AnnotatorList {

		user, err := GetUserByID(uid)
		if err != nil {
			return nil, err // TODO: Return a more informative message.
		}

		annotatorList = append(annotatorList, user.ToPublicUser())
	}

	return &VideoStreamWithJoins{
		ID:            v.ID,
		StartTime:     v.StartTime,
		EndTime:       v.EndTime,
		TimeZone:      v.TimeZone,
		StreamURL:     v.StreamURL,
		CaptureSource: captureSource.ToSummary(),
		AnnotatorList: annotatorList,
	}, nil
}

// VideoStreamContentsFromEntity converts an entities.VideoStream to a VideoStreamContents.
func VideoStreamContentsFromEntity(v entities.VideoStream) VideoStreamContents {
	return VideoStreamContents{
		StartTime:     v.StartTime,
		EndTime:       v.EndTime,
		AnnotatorList: v.AnnotatorList,
		BaseVideoStreamFields: BaseVideoStreamFields{
			TimeZone:      timezone.UncheckedParse(v.TimeZone),
			StreamURL:     v.StreamURL,
			CaptureSource: v.CaptureSource,
		},
	}
}

// ToEntity converts a VideoStreamContents to an entities.VideoStream for storage in the datastore.
func (v *VideoStreamContents) ToEntity() entities.VideoStream {
	return entities.VideoStream{
		StartTime:     v.StartTime,
		EndTime:       v.EndTime,
		TimeZone:      v.TimeZone.String(),
		StreamURL:     v.StreamURL,
		CaptureSource: v.CaptureSource,
		AnnotatorList: v.AnnotatorList,
	}
}

// GetVideoStreamByID gets a video stream when provided with an ID.
func GetVideoStreamByID(id int64) (*VideoStream, error) {
	store := globals.GetStore()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	var e entities.VideoStream
	err := store.Get(context.Background(), key, &e)
	if err != nil {
		return nil, err
	}

	return &VideoStream{
		ID:                  id,
		VideoStreamContents: VideoStreamContentsFromEntity(e),
	}, nil
}

// VideoStreamExists checks if a video stream exists with the given ID.
func VideoStreamExists(id int64) bool {
	store := globals.GetStore()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	var videoStream entities.VideoStream
	err := store.Get(context.Background(), key, &videoStream)
	return err == nil
}

// GetVideoStreams gets a list of video streams, filtering by timespan, capturesource if specified.
func GetVideoStreams(limit int, offset int, timespan *timespan.TimeSpan, captureSource *int64) ([]VideoStream, error) {
	// Fetch data from the datastore.
	store := globals.GetStore()
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

	var ents []entities.VideoStream
	keys, err := store.GetAll(context.Background(), query, &ents)
	if err != nil {
		return []VideoStream{}, err
	}

	// Convert entities.
	videoStreams := make([]VideoStream, len(ents))
	for i := range ents {
		videoStreams[i] = VideoStream{
			ID:                  keys[i].ID,
			VideoStreamContents: VideoStreamContentsFromEntity(ents[i]),
		}
	}

	return videoStreams, nil
}

// CreateVideoStream puts a video stream in the datastore, checking if the capture source exists.
func CreateVideoStream(contents VideoStreamContents) (*VideoStream, error) {

	// Verify CaptureSource exists.
	if !CaptureSourceExists(contents.CaptureSource) {
		return nil, fmt.Errorf("CaptureSource does not exist")
	}

	// Create VideoStream entity.
	store := globals.GetStore()
	key := store.IncompleteKey(entities.VIDEOSTREAM_KIND)
	ent := contents.ToEntity()
	key, err := store.Put(context.Background(), key, &ent)
	if err != nil {
		return nil, err
	}

	// Return newly created videostream.
	created := VideoStream{
		ID:                  key.ID,
		VideoStreamContents: contents,
	}
	return &created, nil
}

// UpdateVideoStream updates an existing video stream with the provided partial contents.
func UpdateVideoStream(id int64, updates PartialVideoStreamContents) error {

	// Update data in the datastore.
	store := globals.GetStore()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	var videoStream entities.VideoStream

	return store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*entities.VideoStream)
		if ok {
			if updates.StreamURL != nil {
				v.StreamURL = *updates.StreamURL
			}
			if updates.CaptureSource != nil {
				// TODO: Check that captureSource exists.
				v.CaptureSource = *updates.CaptureSource
			}
			if updates.StartTime != nil {
				v.StartTime = *updates.StartTime
			}
			if updates.EndTime != nil {
				v.EndTime = updates.EndTime
			}
			if updates.AnnotatorList != nil {
				v.AnnotatorList = *updates.AnnotatorList
			}
		}
	}, &videoStream)
}

// DeleteVideoStream deletes a video stream.
func DeleteVideoStream(id int64) error {
	// TODO: Check that video stream has no annotations associated with it.

	// Delete entity.
	store := globals.GetStore()
	key := store.IDKey(entities.VIDEOSTREAM_KIND, id)
	return store.Delete(context.Background(), key)
}
