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
	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/services"

	"github.com/gofiber/fiber/v2"
)

// UserResult describes the JSON format for users in API responses.
type UserResult struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// FromUserEntity creates a UserResult from a entities.User.
func FromUserEntity(user *entities.User) UserResult {
	return UserResult{Email: user.Email, Role: user.Role.String()}
}

// GetUsersQuery describes the URL query parameters required for the GetUsers endpoint.
type GetUsersQuery struct {
	api.LimitAndOffset
}

// UpdateUserBody describes the JSON format required for the UpdateUser endpoint.
type UpdateUserBody struct {
	Role string `json:"role"` // Required.
}

// GetUserByEmail gets a user when provided with an email.
func GetUserByEmail(ctx *fiber.Ctx) error {
	// Parse URL.
	email := ctx.Params("email")

	// Fetch data from the datastore.
	user, err := services.GetUserByEmail(email)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format result.
	result := FromUserEntity(user)
	return ctx.JSON(result)
}

// GetUsers gets a list of users.
// TODO: support filtering by role.
func GetUsers(ctx *fiber.Ctx) error {
	// Parse URL.
	qry := new(GetUsersQuery)
	qry.SetLimit()

	if err := ctx.QueryParser(qry); err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	users, err := services.GetUsers(qry.Limit, qry.Offset)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Format results.
	results := make([]UserResult, len(users))
	for i := range users {
		results[i] = FromUserEntity(&users[i])
	}

	return ctx.JSON(api.Result[UserResult]{
		Results: results,
		Offset:  qry.Offset,
		Limit:   qry.Limit,
		Total:   len(results),
	})
}

// UpdateUser updates a user.
func UpdateUser(ctx *fiber.Ctx) error {
	// Parse URL.
	email := ctx.Params("email")

	// Parse body.
	var body UpdateUserBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	role, err := entities.ParseRole(body.Role)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Update data in the datastore.
	err = services.UpdateUser(email, role)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}

// DeleteUser deletes a user.
func DeleteUser(ctx *fiber.Ctx) error {
	// Parse URL.
	email := ctx.Params("email")

	// Delete capture source.
	err := services.DeleteUser(email)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
