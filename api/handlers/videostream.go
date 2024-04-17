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
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/entities"
	"github.com/ausocean/openfish/api/services"

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

// FromVideoStream creates a VideoStreamResult from a entities.VideoStream and key, formatting it according to the requested format.
func FromVideoStream(videoStream *entities.VideoStream, id int64, format *api.Format) VideoStreamResult {
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
	CaptureSource *int64             `query:"capturesource"` // Optional.
	TimeSpan      *entities.TimeSpan `query:"timespan"`      // Optional.
	api.LimitAndOffset
}

// GetMediaQuery describes the URL query parameters required for the GetVideoStreamMedia endpoint.
type GetMediaQuery struct {
	Time string `query:"time"` // Optional.
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
	videoStream, err := services.GetVideoStreamByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromVideoStream(videoStream, id, format)
	return ctx.JSON(result)
}

func GetVideoStreamMedia(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetMediaQuery)

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	timespan, err := entities.TimeSpanFromString(qry.Time)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	r, filename, err := services.GetVideoStreamMedia(id, *timespan)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Set mime type and content-disposition headers.
	ctx.Type(".mkv")
	ctx.Attachment(filename)

	// Write file to response.
	if _, err := io.Copy(ctx, r); err != nil {
		return err
	}

	return nil
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

	// Validate timespan
	if qry.TimeSpan != nil {
		if qry.TimeSpan.Start.After(qry.TimeSpan.End) {
			return api.InvalidRequestURL(errors.New("start time not before end time"))
		}
	}

	// Fetch data from the datastore.
	videoStreams, ids, err := services.GetVideoStreams(qry.Limit, qry.Offset, qry.TimeSpan, qry.CaptureSource)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]VideoStreamResult, len(videoStreams))
	for i := range videoStreams {
		results[i] = FromVideoStream(&videoStreams[i], int64(ids[i]), format)
	}

	return ctx.JSON(api.Result[VideoStreamResult]{
		Results: results,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(results),
	})
}

// CreateVideoStream creates a new video stream.
// BUG: start and end time are required but this is not being enforced.
// https://github.com/ausocean/openfish/issues/18
func CreateVideoStream(ctx *fiber.Ctx) error {
	// Parse body.
	var body CreateVideoStreamBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	id, err := services.CreateVideoStream(body.StreamUrl, body.CaptureSource, body.StartTime, &body.EndTime)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(VideoStreamResult{
		ID: &id,
	})
}

// StartVideoStream creates a new video stream at the current time.
func StartVideoStream(ctx *fiber.Ctx) error {
	// Parse body.
	var body StartVideoStreamBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	id, err := services.CreateVideoStream(body.StreamUrl, body.CaptureSource, time.Now(), nil)
	if err != nil {
		return api.DatastoreWriteFailure(err)
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
	err = services.UpdateVideoStream(id, nil, nil, nil, &now)
	if err != nil {
		return api.DatastoreWriteFailure(err)
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
	err = services.UpdateVideoStream(id, body.StreamUrl, body.CaptureSource, body.StartTime, body.EndTime)
	if err != nil {
		return api.DatastoreWriteFailure(err)
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
	err = services.DeleteVideoStream(id)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
