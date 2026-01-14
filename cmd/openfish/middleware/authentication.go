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
	"errors"
	"fmt"
	"strings"

	"github.com/ausocean/cloud/datastore"
	"github.com/ausocean/cloud/gauth"
	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/services"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

// NoAuth skips authentication, for when we are running the OpenFish API locally.
func NoAuth() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		ctx.Locals("email", "no-user@localhost")
		user, err := services.GetUserByEmail("no-user@localhost")
		if err != nil {
			return err
		}
		ctx.Locals("user", user)

		return ctx.Next()
	}
}

// ValidateIAPJWT creates a validator middleware that validate JWT tokens returned from Google IAP.
// Otherwise, it returns a 401 Unauthorized http error.
// See more: https://cloud.google.com/iap/docs/signed-headers-howto#iap_validate_jwt-go
func ValidateIAPJWT(aud string) func(*fiber.Ctx) error {

	fmt.Println("jwt audience: ", aud)

	return func(ctx *fiber.Ctx) error {

		// Get JWT from header.
		iapJWT := ctx.Get("X-Goog-IAP-JWT-Assertion")

		// Validate JWT token.
		payload, err := idtoken.Validate(context.Background(), iapJWT, aud)
		if err != nil {
			return api.Unauthorized(err)
		}

		// Extract email.
		email := payload.Claims["email"].(string)
		ctx.Locals("email", email)

		// Fetch user from datastore if they exist, otherwise create the user.
		user, err := services.GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, datastore.ErrNoSuchEntity) {
				userContents := services.UserContents{DisplayName: strings.Split(email, "@")[0], Email: email, Role: role.Readonly}
				userId, err := services.CreateUser(userContents)
				if err != nil {
					return err
				}
				user, err = services.GetUserByID(userId)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		ctx.Locals("user", user)

		return ctx.Next()
	}
}

// ValidateJWT creates a validator middleware that validates JWT tokens returned from another service.
// Otherwise, it returns a 401 Unauthorized http error.
// The token must have an audience, subject (which should be the email), display name, and issuer.
func ValidateJWT(aud, validIssuer string, secret []byte) func(*fiber.Ctx) error {

	fmt.Println("jwt audience: ", aud)

	return func(ctx *fiber.Ctx) error {

		// Get JWT from header.
		token := ctx.Get("Authorization")

		// Validate JWT token.
		claims, err := gauth.GetClaims(token, secret)
		if err != nil {
			return api.Unauthorized(err)
		}

		if claims["aud"] != aud {
			return api.Unauthorized(fmt.Errorf("Invalid Audience"))
		}

		if claims["iss"] != validIssuer {
			return api.Unauthorized(fmt.Errorf("Invalid Issuer"))
		}

		_, expExists := claims["exp"]
		if !expExists {
			return api.Unauthorized(fmt.Errorf("Token does not contain an expiry time"))
		}

		// Extract email.
		email := claims["sub"].(string)
		ctx.Locals("email", email)

		// Fetch user from datastore if they exist, otherwise create the user.
		user, err := services.GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, datastore.ErrNoSuchEntity) {
				userContents := services.UserContents{DisplayName: claims["name"].(string), Email: email, Role: role.Annotator}
				userId, err := services.CreateUser(userContents)
				if err != nil {
					return err
				}
				user, err = services.GetUserByID(userId)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		ctx.Locals("user", user)

		return ctx.Next()
	}
}

func Guard(requiredRole role.Role) func(*fiber.Ctx) error {

	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*services.User)
		if !ok {
			return fmt.Errorf("failed to assert type: expected *services.User but got %T", ctx.Locals("user"))
		}
		if user != nil && user.Role >= requiredRole {
			return ctx.Next()
		} else {
			return api.Forbidden(fmt.Errorf("this operation requires %s role", requiredRole.String()))
		}

	}
}
