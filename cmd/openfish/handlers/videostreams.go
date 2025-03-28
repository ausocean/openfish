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

// handlers package handles HTTP requests.
package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/mediatype"
	"github.com/ausocean/openfish/cmd/openfish/types/timespan"
	"github.com/ausocean/openfish/cmd/openfish/types/videotime"

	"github.com/gofiber/fiber/v2"
)

// GetVideoStreamsQuery describes the URL query parameters required for the GetVideoStreams endpoint.
type GetVideoStreamsQuery struct {
	CaptureSource *int64             `query:"capturesource"` // Optional.
	TimeSpan      *timespan.TimeSpan `query:"timespan"`      // Optional.
	api.LimitAndOffset
}

// GetMediaVideoQuery describes the URL query parameters required for the GetVideoStreamMedia endpoint, for video mime types.
type GetMediaVideoQuery struct {
	TimeSpan timespan.TimeSpan `query:"time"`
}

// GetMediaImageQuery describes the URL query parameters required for the GetVideoStreamMedia endpoint, for image mime types.
type GetMediaImageQuery struct {
	Time videotime.VideoTime `query:"time"`
}

// CreateVideoStreamBody describes the JSON format required for the CreateVideoStream endpoint.
//
// ID is omitted because it is chosen automatically.
// StartVideoStreamBody describes the JSON format required for the StartVideoStream endpoint.
//
// ID is omitted because it is chosen automatically.
// Datetime is omitted because it uses the current time.
// Duration is omitted because it will be set once the stream concludes.
type StartVideoStreamBody struct {
	StreamUrl     string  `json:"stream_url" example:"https://www.youtube.com/watch?v=abcdefghijk" validate:"required"` // URL of live video stream.
	CaptureSource int64   `json:"capturesource" example:"1234567890" validate:"required"`                               // ID of the capture source that produced the stream.
	AnnotatorList []int64 `json:"annotator_list" example:"1234567890" validate:"optional"`                              // Users that are permitted to add annotations.
}

// GetVideoStreamByID gets a video stream when provided with an ID.
//
//	@Summary		Get video stream by ID
//	@Description	Gets a video stream when provided with an ID.
//	@Tags			Video Streams
//	@Produce		json
//	@Param			id	path		int	true	"Video Stream ID"	example(1234567890)
//	@Success		200	{object}	services.VideoStreamWithJoins
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id} [get]
func GetVideoStreamByID(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	videoStream, err := services.GetVideoStreamByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	joined, err := videoStream.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
}

// GetVideoStreamMedia gets the image/video snippet from this video stream at the given time.
//
//	@Summary		Get video stream media
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Gets the image or video snippet from this video stream at the given time.
//	@Tags			Media
//	@Param			id		path	int		true	"Video Stream ID"	example(1234567890)
//	@Param			type	path	string	true	"Type"				example(image)
//	@Param			subtype	path	string	true	"Subtype"			example(jpeg)
//	@Param			time	query	string	true	"Time"				example(00:00:01.000-00:00:05.500)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id}/media/{type}/{subtype} [get]
func GetVideoStreamMedia(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	mtype, err := mediatype.ParseMimeType(fmt.Sprintf("%s/%s", ctx.Params("type"), ctx.Params("subtype")))
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse query params.
	var start videotime.VideoTime
	var end *videotime.VideoTime
	if mtype.IsVideo() {
		qry := new(GetMediaVideoQuery)
		if err := ctx.QueryParser(qry); err != nil {
			return api.InvalidRequestURL(err)
		}
		if !qry.TimeSpan.Valid() {
			return api.InvalidRequestURL(fmt.Errorf("invalid time span, start time must occur before end time"))
		}
		start = qry.TimeSpan.Start
		end = &qry.TimeSpan.End
	} else {
		qry := new(GetMediaImageQuery)
		if err := ctx.QueryParser(qry); err != nil {
			return api.InvalidRequestURL(err)
		}
		start = qry.Time
	}

	// Fetch data from storage
	bytes, err := services.GetMedia(services.MediaKey{
		Type:          mtype,
		VideoStreamID: id,
		StartTime:     start,
		EndTime:       end,
	})
	if err != nil {
		return err
	}

	ctx.Type(mtype.FileExtension())
	ctx.Attachment(fmt.Sprintf("%d.%s", id, mtype.FileExtension()))
	ctx.Write(bytes)

	return nil
}

// DeleteVideoStreamMedia deletes the cached image/video snippet from this video stream at the given time.
//
//	@Summary		Delete video stream media
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Deletes the cached image or video snippet from this video stream at the given time.
//	@Tags			Media
//	@Param			id		path	int		true	"Video Stream ID"	example(1234567890)
//	@Param			type	path	string	true	"Type"				example(image)
//	@Param			subtype	path	string	true	"Subtype"			example(jpeg)
//	@Param			time	query	string	true	"Time"				example(00:00:01.000-00:00:05.500)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id}/media/{type}/{subtype} [delete]
func DeleteVideoStreamMedia(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	mtype, err := mediatype.ParseMimeType(fmt.Sprintf("%s/%s", ctx.Params("type"), ctx.Params("subtype")))
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse query params.
	var start videotime.VideoTime
	var end *videotime.VideoTime
	if mtype.IsVideo() {
		qry := new(GetMediaVideoQuery)
		if err := ctx.QueryParser(qry); err != nil {
			return api.InvalidRequestURL(err)
		}
		if !qry.TimeSpan.Valid() {
			return api.InvalidRequestURL(fmt.Errorf("invalid time span, start time must occur before end time"))
		}
		start = qry.TimeSpan.Start
		end = &qry.TimeSpan.End
	} else {
		qry := new(GetMediaImageQuery)
		if err := ctx.QueryParser(qry); err != nil {
			return api.InvalidRequestURL(err)
		}
		start = qry.Time
	}

	// Delete media from storage.
	return services.DeleteMedia(services.MediaKey{
		Type:          mtype,
		VideoStreamID: id,
		StartTime:     start,
		EndTime:       end,
	})
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
//	@Success		200				{object}	api.Result[services.VideoStreamWithJoins]
//	@Failure		400				{object}	api.Failure
//	@Failure		401				{object}	api.Failure
//	@Failure		403				{object}	api.Failure
//	@Router			/api/v1/videostreams [get]
func GetVideoStreams(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetVideoStreamsQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Validate timespan.
	if qry.TimeSpan != nil {
		if !qry.TimeSpan.Valid() {
			return api.InvalidRequestURL(fmt.Errorf("invalid time span, start time must occur before end time"))
		}
	}

	// Fetch data from the datastore.
	videoStreams, err := services.GetVideoStreams(qry.Limit, qry.Offset, qry.TimeSpan, qry.CaptureSource)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Apply Joins.
	joined := make([]services.VideoStreamWithJoins, len(videoStreams))
	for i, annotation := range videoStreams {
		j, err := annotation.JoinFields()
		if err != nil {
			return api.DatastoreReadFailure(err)
		}
		joined[i] = *j
	}

	return ctx.JSON(api.Result[services.VideoStreamWithJoins]{
		Results: joined,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(joined),
	})
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
//	@Param			body	body	services.PartialVideoStreamContents	true	"Update Video Stream"
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		401	{object}	api.Failure
//	@Failure		403	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/videostreams/{id} [patch]
func UpdateVideoStream(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse body.
	var body services.PartialVideoStreamContents
	if ctx.BodyParser(&body) != nil {
		return api.InvalidRequestJSON(err)
	}

	// Update data in the datastore.
	err = services.UpdateVideoStream(id, body)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
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
//	@Param			body	body		services.VideoStreamContents	true	"New Video Stream"
//	@Success		201		{object}	services.VideoStream
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/videostreams [post]
func CreateVideoStream(ctx *fiber.Ctx) error {
	// Parse body.
	var body services.VideoStreamContents
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create video stream entity and add to the datastore.
	created, err := services.CreateVideoStream(body)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return joined form.
	joined, err := created.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
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
//	@Param			body	body		services.BaseVideoStreamFields	true	"New Video Stream"
//	@Success		201		{object}	services.VideoStream
//	@Failure		400		{object}	api.Failure
//	@Failure		401		{object}	api.Failure
//	@Failure		403		{object}	api.Failure
//	@Router			/api/v1/videostreams/live [post]
func StartVideoStream(ctx *fiber.Ctx) error {
	// Parse body.
	var body services.BaseVideoStreamFields
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	contents := services.VideoStreamContents{
		StartTime:             time.Now(),
		BaseVideoStreamFields: body,
	}

	// Create video stream entity and add to the datastore.
	created, err := services.CreateVideoStream(contents)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	// Return joined form.
	joined, err := created.JoinFields()
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	return ctx.JSON(joined)
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
	err = services.UpdateVideoStream(id, services.PartialVideoStreamContents{EndTime: &now})
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
