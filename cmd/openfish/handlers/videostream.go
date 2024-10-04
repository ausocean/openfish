/*
AUTHORS
  Scott Barnard <scott@ausocean.org>

LICENSE
  Copyright (c) 2023-2024, The OpenFish Contributors.

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
	"time"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/timespan"

	"github.com/gofiber/fiber/v2"
)

// VideoStreamResult  describes the JSON format for video streams in API responses.
// Fields use pointers because they are optional (this is what the format URL param is for).
type VideoStreamResult struct {
	ID            *int64     `json:"id,omitempty" example:"1234567890"`
	StartTime     *time.Time `json:"startTime,omitempty" example:"2023-05-25T08:00:00Z"`
	EndTime       *time.Time `json:"endTime,omitempty" example:"2023-05-25T16:30:00Z"`
	StreamUrl     *string    `json:"stream_url,omitempty" example:"https://www.youtube.com/watch?v=abcdefghijk"`
	CaptureSource *int64     `json:"capturesource,omitempty" example:"1234567890"`
	AnnotatorList *[]string  `json:"annotator_list,omitempty" example:"user@example.com"`
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
	if format.Requires("annotator_list") {
		result.AnnotatorList = &videoStream.AnnotatorList
	}
	return result
}

// GetVideoStreamsQuery describes the URL query parameters required for the GetVideoStreams endpoint.
type GetVideoStreamsQuery struct {
	CaptureSource *int64             `query:"capturesource"` // Optional.
	TimeSpan      *timespan.TimeSpan `query:"timespan"`      // Optional.
	api.LimitAndOffset
}

// CreateVideoStreamBody describes the JSON format required for the CreateVideoStream endpoint.
//
// ID is omitted because it is chosen automatically.
type CreateVideoStreamBody struct {
	StartTime     time.Time `json:"startTime" example:"2023-05-25T08:00:00Z" validate:"required"`                         // Start time of stream.
	EndTime       time.Time `json:"endTime" example:"2023-05-25T16:30:00Z" validate:"required"`                           // End time of stream.
	StreamUrl     string    `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk" validate:"required"` // URL of video stream.
	CaptureSource int64     `json:"capturesource" example:"1234567890" validate:"required"`                               // ID of the capture source that produced the stream.
	AnnotatorList []string  `json:"annotator_list" example:"user@example.com" validate:"optional"`                        // Users that are permitted to add annotations.
}

// StartVideoStreamBody describes the JSON format required for the StartVideoStream endpoint.
//
// ID is omitted because it is chosen automatically.
// Datetime is omitted because it uses the current time.
// Duration is omitted because it will be set once the stream concludes.
type StartVideoStreamBody struct {
	StreamUrl     string   `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk" validate:"required"` // URL of live video stream.
	CaptureSource int64    `json:"capturesource" example:"1234567890" validate:"required"`                               // ID of the capture source that produced the stream.
	AnnotatorList []string `json:"annotator_list" example:"user@example.com" validate:"optional"`                        // Users that are permitted to add annotations.
}

// UpdateVideoStreamBody describes the JSON format required for the UpdateVideoStream endpoint.
type UpdateVideoStreamBody struct {
	StartTime     *time.Time `json:"startTime" example:"2023-05-25T08:00:00Z" validate:"optional"`                         // Start time of stream.
	EndTime       *time.Time `json:"endTime" example:"2023-05-25T16:30:00Z" validate:"optional"`                           // End time of stream.
	StreamUrl     *string    `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk" validate:"optional"` // URL of video stream.
	CaptureSource *int64     `json:"capturesource" example:"1234567890" validate:"optional"`                               // ID of the capture source that produced the stream.
	AnnotatorList *[]string  `json:"annotator_list" example:"user@example.com" validate:"optional"`                        // Users that are permitted to add annotations.
}

// GetVideoStreamByID gets a video stream when provided with an ID.
//
//	@Summary		Get video stream by ID
//	@Description	Gets a video stream when provided with an ID.
//	@Tags			Video Streams
//	@Produce		json
//	@Param			id	path		int	true	"Video Stream ID"	example(1234567890)
//	@Success		200	{object}	VideoStreamResult
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id} [get]
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

// GetVideoStreams gets a list of video streams, filtering by timespan, capture source if specified.
//
//	@Summary		Get video streams
//	@Description	Get paginated video streams, with options to filter by timespan and capturesource.
//	@Tags			Video Streams
//	@Produce		json
//	@Param			limit			query		int		false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset			query		int		false	"Number of results to skip."	minimum(0)
//	@Param			capturesource	query		int		false	"Capture source ID to filter by."
//	@Param			timespan[start]	query		string	false	"Start time to filter by."
//	@Param			timespan[end]	query		string	false	"End time to filter by."
//	@Success		200				{object}	api.Result[VideoStreamResult]
//	@Failure		400				{object}	api.Failure
//	@Router			/api/v1/videostreams [get]
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

	// Validate timespan.
	if qry.TimeSpan != nil {
		if !qry.TimeSpan.Valid() {
			return api.InvalidRequestURL(fmt.Errorf("invalid time span, start time must occur before end time"))
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
//
//	@Summary		Register video stream
//	@Description	Roles required: <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Registers a new video stream with OpenFish.
//	@Tags			Video Streams
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateVideoStreamBody	true	"New Video Stream"
//	@Success		201		{object}	EntityIDResult
//	@Failure		400		{object}	api.Failure
//	@Router			/api/v1/videostreams [post]
func CreateVideoStream(ctx *fiber.Ctx) error {
	// Parse body.
	var body CreateVideoStreamBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	id, err := services.CreateVideoStream(body.StreamUrl, body.CaptureSource, body.StartTime, &body.EndTime, body.AnnotatorList)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(EntityIDResult{
		ID: &id,
	})
}

// StartVideoStream creates a new video stream at the current time.
//
//	@Summary		Register live stream
//	@Description	Roles required: <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Registers a new live video stream with OpenFish. The API takes the current time as the start time of the video stream.
//	@Tags			Video Streams (Live)
//	@Accept			json
//	@Produce		json
//	@Param			body	body		StartVideoStreamBody	true	"New Video Stream"
//	@Success		201		{object}	EntityIDResult
//	@Failure		400		{object}	api.Failure
//	@Router			/api/v1/videostreams/live [post]
func StartVideoStream(ctx *fiber.Ctx) error {
	// Parse body.
	var body StartVideoStreamBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	id, err := services.CreateVideoStream(body.StreamUrl, body.CaptureSource, time.Now(), nil, body.AnnotatorList)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return ID of created video stream.
	return ctx.JSON(EntityIDResult{
		ID: &id,
	})
}

// EndVideoStream updates the video stream's duration.
//
//	@Summary		Finish live stream
//	@Description	Roles required: <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Notify OpenFish that a live video stream has finished. The API takes the current time as the end time.
//	@Tags			Video Streams (Live)
//	@Param			id	path	int	true	"Video Stream ID"	Example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id}/live [patch]
func EndVideoStream(ctx *fiber.Ctx) error {
	now := time.Now()

	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Update data in the datastore.
	err = services.UpdateVideoStream(id, nil, nil, nil, &now, nil)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}

// UpdateVideoStream updates a video stream.
//
//	@Summary		Update video stream
//	@Description	Roles required: <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Partially update a video stream by specifying the properties to update.
//	@Tags			Video Streams
//	@Accept			json
//	@Param			id		path	int						true	"Video Stream ID"	example(1234567890)
//	@Param			body	body	UpdateVideoStreamBody	true	"Update Video Stream"
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id} [patch]
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
	err = services.UpdateVideoStream(id, body.StreamUrl, body.CaptureSource, body.StartTime, body.EndTime, body.AnnotatorList)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}

// DeleteVideoStream deletes a video stream.
// BUG: endpoint returns 200 ok for nonexistent IDs.
// https://github.com/ausocean/openfish/issues/17
//
//	@Summary		Delete video stream
//	@Description	Roles required: <role-tag>Curator</role-tag> or <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete a video stream by providing the video stream ID.
//	@Tags			Video Streams
//	@Param			id	path	int	true	"Video Stream ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id} [delete]
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
