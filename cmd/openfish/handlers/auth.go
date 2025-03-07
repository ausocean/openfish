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

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
	"github.com/gofiber/fiber/v2"
)

// CreateSelfBody is the body of a request to create a new user.
type CreateSelfBody struct {
	DisplayName string `json:"display_name" example:"Coral Fischer"`
}

// GetSelf gets information about the current user.
//
//	@Summary		Get current user
//	@Description	Gets information about the current user.
//	@Tags			Authentication
//	@Produce		json
//	@Success		200	{object}	services.User
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/auth/me [get]
func GetSelf(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(*services.User)
	if !ok {
		return fmt.Errorf("failed to assert type: expected *services.User but got %T", ctx.Locals("user"))
	}
	if user == nil {
		return api.NotFound(fmt.Errorf("user not found"))
	}
	return ctx.JSON(user)
}

// CreateSelf creates a new user.
//
//	@Summary		Create user
//	@Description	Creates a new user.
//	@Tags			Authentication
//	@Accept			json
//	@Param			body	body	CreateSelfBody	true	"New User"
//	@Produce		json
//	@Success		201	{object}	services.User
//	@Failure		400	{object}	api.Failure
//	@Failure		409	{object}	api.Failure
//	@Router			/api/v1/auth/me [post]
func CreateSelf(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(*services.User)
	if !ok {
		return fmt.Errorf("failed to assert type: expected *services.User but got %T", ctx.Locals("user"))
	}
	if user != nil {
		return api.Conflict(fmt.Errorf("user already exists"))
	}

	// Parse body.
	var body CreateSelfBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Create user.
	email := ctx.Locals("email").(string)
	id, err := services.CreateUser(services.UserContents{
		Email:       email,
		Role:        role.Default,
		DisplayName: body.DisplayName,
	})
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return ctx.JSON(EntityIDResult{
		ID: id,
	})
}
