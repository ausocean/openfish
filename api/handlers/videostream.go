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
	"errors"
	"strconv"
	"time"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/model"
	"github.com/ausocean/openfish/datastore"

	"github.com/gofiber/fiber/v2"
)

// VideoStreamResult  describes the JSON format for video streams in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type VideoStreamResult struct {
	ID            *int64     `json:"id,omitempty"`
	StartTime     *time.Time `json:"startTime,omitempty"`
	EndTime       *time.Time `json:"endTime,omitempty"`
	StreamUrl     *string    `json:"stream_url,omitempty"`
	CaptureSource *int64     `json:"capturesource,omitempty"`
}

// FromVideoStream creates a VideoStreamResult from a model.VideoStream and key, formatting it according to the requested format.
func FromVideoStream(videoStream *model.VideoStream, id int64, format *api.Format) VideoStreamResult {
	var result VideoStreamResult
	if format.Requires("id") {
		result.ID = &id
	}
	if format.Requires("start_time") {
		result.StartTime = &videoStream.StartTime
	}
	if format.Requires("end_time") {
		result.EndTime = videoStream.EndTime
	}
	if format.Requires("stream_url") {
		result.StreamUrl = &videoStream.StreamUrl
	}
	if format.Requires("capturesource") {
		result.CaptureSource = &videoStream.CaptureSource
	}
	return result
}

// GetVideoStreamsQuery describes the URL query parameters required for the GetVideoStreams endpoint.
type GetVideoStreamsQuery struct {
	CaptureSource *int64          `query:"capturesource"` // Optional.
	TimeSpan      *model.TimeSpan `query:"timespan"`      // Optional.
	api.LimitAndOffset
}

// CreateVideoStreamBody describes the JSON format required for the CreateVideoStream endpoint.
//
// ID is omitted because it is chosen automatically.
type CreateVideoStreamBody struct {
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	StreamUrl     string    `json:"stream_url"`
	CaptureSource int64     `json:"capturesource"`
}

// StartVideoStreamBody describes the JSON format required for the StartVideoStream endpoint.
//
// ID is omitted because it is chosen automatically.
// Datetime is omitted because it uses the current time.
// Duration is omitted because it will be set once the stream concludes.
type StartVideoStreamBody struct {
	StreamUrl     string `json:"stream_url"`
	CaptureSource int64  `json:"capturesource"`
}

// UpdateVideoStreamBody describes the JSON format required for the UpdateVideoStream endpoint.
type UpdateVideoStreamBody struct {
	StartTime     *time.Time `json:"startTime"`     // Optional.
	EndTime       *time.Time `json:"endTime"`       // Optional.
	StreamUrl     *string    `json:"stream_url"`    // Optional.
	CaptureSource *int64     `json:"capturesource"` // Optional.
}

// GetVideoStreamByID gets a video stream when provided with an ID.
func GetVideoStreamByID(ctx *fiber.Ctx) error {
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
	key := store.IDKey("VideoStream", id)
	var videoStream model.VideoStream
	if store.Get(context.Background(), key, &videoStream) != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromVideoStream(&videoStream, id, format)

	return ctx.JSON(result)
}

// GetVideoStreams gets a list of video streams, filtering by timespan, capture source if specified.
func GetVideoStreams(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetVideoStreamsQuery)
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
	query := store.NewQuery("VideoStream", false)

	// Filter by timespan (database side).
	if qry.TimeSpan != nil {
		// Validate timespan
		if qry.TimeSpan.Start.After(qry.TimeSpan.End) {
			return api.InvalidRequestURL(errors.New("start time not before end time"))
		}

		// BUG: Because of a limitation of google cloud's datastore, we can only use inequality filters on one
		// field at a time. This query needs to filter out videostreams that started after the specified range.
		// https://github.com/ausocean/openfish/issues/23
		query.FilterField("EndTime", ">=", qry.TimeSpan.Start)
	}

	// Filter by capture source.
	if qry.CaptureSource != nil {
		query.FilterField("CaptureSource", "=", *qry.CaptureSource)
	}

	query.Limit(qry.Limit)
	query.Offset(qry.Offset)

	var videoStreams []model.VideoStream
	keys, err := store.GetAll(context.Background(), query, &videoStreams)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]VideoStreamResult, len(videoStreams))
	for i := range videoStreams {
		results[i] = FromVideoStream(&videoStreams[i], keys[i].ID, format)
	}

	return ctx.JSON(api.Result[VideoStreamResult]{
		Results: results,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(results),
	})
}

// putVideoStream puts a video stream in the datastore, checking if the capture source exists.
func putVideoStream(ctx *fiber.Ctx, vs model.VideoStream) (int64, error) {
	// Verify CaptureSource exists.
	store := ds_client.Get()
	key := store.IDKey("CaptureSource", vs.CaptureSource)
	var captureSource model.CaptureSource

	err := store.Get(context.Background(), key, &captureSource)
	if err != nil {
		return 0, api.DatastoreReadFailure(err)
	}

	// Get a unique ID for the new video stream.
	key = store.IncompleteKey("VideoStream")

	key, err = store.Put(context.Background(), key, &vs)
	if err != nil {
		print(err.Error())
		return 0, api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return key.ID, nil

}

// CreateVideoStream creates a new video stream.
// BUG: start and end time are required but this is not being enforced.
// https://github.com/ausocean/openfish/issues/18
func CreateVideoStream(ctx *fiber.Ctx) error {
	var body CreateVideoStreamBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	vs := model.VideoStream{
		StartTime:     body.StartTime,
		EndTime:       &body.EndTime,
		StreamUrl:     body.StreamUrl,
		CaptureSource: body.CaptureSource,
	}
	id, err := putVideoStream(ctx, vs)
	if err != nil {
		return err
	}

	// Return ID of created video stream.
	return ctx.JSON(VideoStreamResult{
		ID: &id,
	})
}

// StartVideoStream creates a new video stream at the current time.
func StartVideoStream(ctx *fiber.Ctx) error {
	now := time.Now()

	var body StartVideoStreamBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	vs := model.VideoStream{
		StartTime:     now,
		StreamUrl:     body.StreamUrl,
		CaptureSource: body.CaptureSource,
	}
	id, err := putVideoStream(ctx, vs)
	if err != nil {
		return err
	}

	// Return ID of created video stream.
	return ctx.JSON(VideoStreamResult{
		ID: &id,
	})
}

// EndVideoStream updates the video stream's duration.
func EndVideoStream(ctx *fiber.Ctx) error {
	now := time.Now()

	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Update data in the datastore.
	store := ds_client.Get()
	key := store.IDKey("VideoStream", id)
	var videoStream model.VideoStream

	err = store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*model.VideoStream)
		if ok {
			v.EndTime = &now
		}
	}, &videoStream)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return nil
}

// UpdateVideoStream updates a video stream.
func UpdateVideoStream(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse body.
	var body UpdateVideoStreamBody
	if ctx.BodyParser(&body) != nil {
		return api.InvalidRequestJSON(err)
	}

	// Update data in the datastore.
	store := ds_client.Get()
	key := store.IDKey("VideoStream", id)
	var videoStream model.VideoStream

	err = store.Update(context.Background(), key, func(e datastore.Entity) {
		v, ok := e.(*model.VideoStream)
		if ok {
			if body.StartTime != nil {
				v.StartTime = *body.StartTime
			}
			if body.EndTime != nil {
				v.EndTime = body.EndTime
			}
			if body.CaptureSource != nil {
				v.CaptureSource = *body.CaptureSource
			}
			if body.StreamUrl != nil {
				v.StreamUrl = *body.StreamUrl
			}
		}
	}, &videoStream)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return nil
}

// DeleteVideoStream deletes a video stream.
// BUG: endpoint returns 200 ok for nonexistent IDs.
// https://github.com/ausocean/openfish/issues/17
func DeleteVideoStream(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Delete entity.
	store := ds_client.Get()
	key := store.IDKey("VideoStream", id)

	if store.Delete(context.Background(), key) != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil

}
