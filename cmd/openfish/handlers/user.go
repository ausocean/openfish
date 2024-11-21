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
	"strconv"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/role"

	"github.com/gofiber/fiber/v2"
)

// GetUsersQuery describes the URL query parameters required for the GetUsers endpoint.
type GetUsersQuery struct {
	api.LimitAndOffset
}

// GetUserByID gets a user when provided with an ID.
//
//	@Summary		Get user by ID
//	@Description	Gets a user when provided with an ID. When invoked by an admin, the full user object is returned, otherwise the public user only. When the user is yourself, the full user object is returned regardless of your permissions.
//	@Tags			Users
//	@Produce		json
//	@Param			id	path		int	true	"ID"	example(1234567890)
//	@Success		200	{object}	services.User
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/users/{id} [get]
func GetUserByID(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Fetch data from the datastore.
	user, err := services.GetUserByID(id)
	if err != nil {
		return api.DatastoreReadFailure(err)
	}

	// Return result depending on user role.
	loggedInUser := ctx.Locals("user").(*services.User)
	if loggedInUser != nil && loggedInUser.Role == role.Admin {
		return ctx.JSON(user)
	} else {
		return ctx.JSON(user.ToPublicUser())
	}
}

// GetUsers gets a list of users.
// TODO: support filtering by role.
//
//	@Summary		Get users
//	@Description	Gets paginated users. When invoked by an admin, the full user object is returned, otherwise the public user only.
//	@Tags			Users
//	@Produce		json
//	@Param			limit	query		int	false	"Number of results to return."	minimum(1)	default(20)
//	@Param			offset	query		int	false	"Number of results to skip."	minimum(0)
//	@Success		200		{object}	api.Result[services.User]
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

	// Return results depending on user role and identity.
	loggedInUser := ctx.Locals("user").(*services.User)
	if loggedInUser != nil && loggedInUser.Role == role.Admin {
		return ctx.JSON(api.Result[services.User]{
			Results: users,
			Offset:  qry.Offset,
			Limit:   qry.Limit,
			Total:   len(users),
		})
	} else {
		results := make([]services.PublicUser, len(users))
		for i := range users {
			results[i] = users[i].ToPublicUser()
		}

		return ctx.JSON(api.Result[services.PublicUser]{
			Results: results,
			Offset:  qry.Offset,
			Limit:   qry.Limit,
			Total:   len(results),
		})
	}

}

// UpdateUser updates a user.
//
//	@Summary		Update role
//	@Description	Roles required: <role-tag>Admin</role-tag>
//	@Description
//	@Description	Update a user's role.
//	@Tags			Users
//	@Accept			json
//	@Param			id		path	string							true	"ID"	example(1234567890)
//	@Param			body	body	services.PartialUserContents	true	"Update User"
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Router			/api/v1/users/{id} [patch]
func UpdateUser(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Parse body.
	var body services.PartialUserContents
	err = ctx.BodyParser(&body)
	if err != nil {
		return api.InvalidRequestJSON(err)
	}

	// Update data in the datastore.
	err = services.UpdateUser(id, body)
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
//	@Description	Delete a user by providing the user's ID.
//	@Tags			Users
//	@Param			id	path	string	true	"ID"	example(1234567890)
//	@Success		200
//	@Failure		400	{object}	api.Failure
//	@Failure		404	{object}	api.Failure
//	@Router			/api/v1/users/{id} [delete]
func DeleteUser(ctx *fiber.Ctx) error {
	// Parse URL.
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return api.InvalidRequestURL(err)
	}

	// Delete user.
	err = services.DeleteUser(id)
	if err != nil {
		return api.DatastoreWriteFailure(err)
	}

	return nil
}
