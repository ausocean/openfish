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

package middleware

import (
	"context"
	"fmt"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/services"
	"github.com/ausocean/openfish/api/types/user"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

// ValidateJWT creates a validator middleware that validate JWT tokens returned from Google IAP.
// Otherwise, it returns a 401 Unauthorized http error.
// See more: https://cloud.google.com/iap/docs/signed-headers-howto#iap_validate_jwt-go
func ValidateJWT(aud string) func(*fiber.Ctx) error {

	fmt.Println("jwt audience: ", aud)

	return func(ctx *fiber.Ctx) error {

		// Get JWT from header.
		iapJWT := ctx.Get("X-Goog-IAP-JWT-Assertion")

		// Validate JWT token.
		payload, err := idtoken.Validate(context.Background(), iapJWT, aud)
		if err != nil {
			return api.Unauthorized(err)
		}

		// Extract subject and email.
		email := payload.Claims["email"].(string)
		subject := payload.Subject
		role := user.DefaultRole

		// Fetch user from datastore if they exist, else
		// create new user.
		user, err := services.GetUserByEmail(email)
		if err != nil {
			services.CreateUser(email, role)
		} else {
			role = user.Role
		}

		// Add subject, email, and user role to ctx.Locals.
		ctx.Locals("subject", subject)
		ctx.Locals("email", email)
		ctx.Locals("role", role)

		return ctx.Next()
	}
}

func Guard(requiredRole user.Role) func(*fiber.Ctx) error {

	return func(ctx *fiber.Ctx) error {
		// Skip if IAP authentication is disabled.
		if ctx.Locals("role") == nil {
			return ctx.Next()
		}

		userRole := ctx.Locals("role").(user.Role)
		if userRole >= requiredRole {
			return ctx.Next()
		} else {
			return api.Forbidden(fmt.Errorf("this operation requires %s role, requesting user has %s role", requiredRole.String(), userRole.String()))
		}

	}
}
