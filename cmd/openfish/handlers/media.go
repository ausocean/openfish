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

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/datastore"
	"github.com/gofiber/fiber/v2"
)

// GetMediaByID gets an image/video when provided with an ID.
//
//	@Summary		Get media by ID
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Gets an image or video file when provided with an ID.
//	@Tags			Media
//	@Produce		video/mp4, image/jpeg
//	@Param			id	path	int	true	"Media ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/media/{id} [get]
func GetMediaByID(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	media, err := services.GetMediaByID(id)
	if err == datastore.ErrNoSuchEntity {
		return api.NotFound()
	}
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	ctx.Type(fmt.Sprintf("%d.%s", id, media.Type.FileExtension()))
	ctx.Attachment(fmt.Sprintf("%d.%s", id, media.Type.FileExtension()))
	ctx.Write(media.Bytes)

	return nil
}

// DeleteMedia deletes an image/video.
//
//	@Summary		Delete media
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete an image or video, by providing a media id.
//	@Tags			Media
//	@Param			id	path	int	true	"Media ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/media/{id} [delete]
func DeleteMedia(ctx *fiber.Ctx) error {

	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Delete media.
	err = services.DeleteMedia(id)
	if err == datastore.ErrNoSuchEntity {
		return api.NotFound()
	}
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
