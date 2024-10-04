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
	Email string `json:"email" example:"user@example.com"`
	Role  string `json:"role" example:"annotator"`
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
	Role string `json:"role" example:"annotator" validate:"required" enums:"readonly,annotator,curator,admin"` // User role.
}

// GetUserByEmail gets a user when provided with an email.
//
//	@Summary		Get user by email
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Gets a user when provided with an email.
//	@Tags			Users
//	@Produce		json
//	@Param			email	path		string	true	"Email"	example(user@example.com)
//	@Success		200		{object}	UserResult
//	@Failure		400		{object}	api.Failure
//	@Failure		404		{object}	api.Failure
//	@Router			/api/v1/users/{email} [get]
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
//
//	@Summary		Get users
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Get paginated users.
//	@Tags			Users
//	@Produce		json
//	@Param			limit	query		int	false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset	query		int	false	"Number of results to skip."	minimum(0)
//	@Success		200		{object}	api.Result[UserResult]
//	@Failure		400		{object}	api.Failure
//	@Router			/api/v1/users [get]
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
//
//	@Summary		Update role
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Update a user's role.
//	@Tags			Users
//	@Accept			json
//	@Param			email	path	string			true	"Email"	example(user@example.com)
//	@Param			body	body	UpdateUserBody	true	"Update User"
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Router			/api/v1/users/{email} [patch]
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
//
//	@Summary		Delete user
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Delete a user by providing the user's email.
//	@Tags			Users
//	@Param			email	path	string	true	"Email"	example(user@example.com)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/users/{email} [delete]
func DeleteUser(ctx *fiber.Ctx) error {
	// Parse URL.
	email := ctx.Params("email")

	// Delete user.
	err := services.DeleteUser(email)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
