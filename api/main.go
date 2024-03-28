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

package main

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/ausocean/openfish/api/api"
	"github.com/ausocean/openfish/api/ds_client"
	"github.com/ausocean/openfish/api/handlers"

	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/api/idtoken"
)

// registerAPIRoutes registers all handler functions to their routes.
func registerAPIRoutes(app *fiber.App) {

	v1 := app.Group("/api/v1")

	// Capture sources.
	v1.Get("/capturesources/:id", handlers.GetCaptureSourceByID)
	v1.Get("/capturesources", handlers.GetCaptureSources)
	v1.Post("/capturesources", handlers.CreateCaptureSource)
	v1.Patch("/capturesources/:id", handlers.UpdateCaptureSource)
	v1.Delete("/capturesources/:id", handlers.DeleteCaptureSource)

	// Video streams.
	v1.Get("/videostreams/:id", handlers.GetVideoStreamByID)
	v1.Get("/videostreams", handlers.GetVideoStreams)
	v1.Post("/videostreams/live", handlers.StartVideoStream)
	v1.Patch("/videostreams/:id/live", handlers.EndVideoStream)
	v1.Post("/videostreams", handlers.CreateVideoStream)
	v1.Patch("/videostreams/:id", handlers.UpdateVideoStream)
	v1.Delete("/videostreams/:id", handlers.DeleteVideoStream)

	// Annotations.
	v1.Get("/annotations/:id", handlers.GetAnnotationByID)
	v1.Get("/annotations", handlers.GetAnnotations)
	v1.Post("/annotations", handlers.CreateAnnotation)
	v1.Delete("/annotations/:id", handlers.DeleteAnnotation)

	// Species.
	v1.Get("/species/recommended", handlers.GetRecommendedSpecies)
	v1.Get("/species/:id", handlers.GetSpeciesByID)
	v1.Post("/species", handlers.CreateSpecies)
	v1.Delete("/species/:id", handlers.DeleteSpecies)

	// Auth.
	v1.Get("/auth/me", handlers.GetSelf)
}

// errorHandler creates a HTTP response with the given status code or 500 by default.
// The response body is JSON: {"message": "<error message here>"}
func errorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500.
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error.
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send JSON response.
	ctx.Status(code).JSON(api.Failure{Message: err.Error()})
	return nil
}

// validateJWT creates a validator middleware that validate JWT tokens returned from Google IAP.
// Otherwise, it returns a 401 Unauthorized http error.
// See more: https://cloud.google.com/iap/docs/signed-headers-howto#iap_validate_jwt-go
func validateJWT(aud string) func(*fiber.Ctx) error {

	fmt.Println("jwt audience: ", aud)

	return func(ctx *fiber.Ctx) error {

		// Get JWT from header.
		iapJWT := ctx.Get("X-Goog-IAP-JWT-Assertion")

		// Validate JWT token.
		payload, err := idtoken.Validate(context.Background(), iapJWT, aud)
		if err != nil {
			return api.Unauthorized(err)
		}

		ctx.Locals("subject", payload.Subject)
		ctx.Locals("email", payload.Claims["email"])

		return ctx.Next()
	}
}

func envOrFlag[T any](flagName string, envName string, description string, defaultVal T, parse func(string) (T, error), flag func(string, T, string) *T) *T {
	v, err := parse(os.Getenv(envName))

	if err == nil {
		description = fmt.Sprintf("%s [set %s=%s]", description, envName, os.Getenv(envName))
		defaultVal = v
	}

	return flag(flagName, defaultVal, description)
}

// parseString does nothing, but is required for using envOrFlag for strings.
func parseString(s string) (string, error) {
	return s, nil
}

func main() {

	// Get app configuration. Configurations can be set using environment variables or command line arguments.
	// Command line arguments take priority.

	port := envOrFlag("port", "PORT", "Port to listen on", 8080, strconv.Atoi, flag.Int)
	useFilestore := envOrFlag("filestore", "FILESTORE", "Use local datastore", false, strconv.ParseBool, flag.Bool)
	useIAP := envOrFlag("iap", "IAP", "Use Google's Identity Aware Proxy for authentication", false, strconv.ParseBool, flag.Bool)
	jwtAudience := envOrFlag("jwt-audience", "JWT_AUDIENCE", "Audience to use to validate JWT token", "", parseString, flag.String)

	flag.Parse()

	// Datastore setup.
	if *useFilestore {
		fmt.Println("using filestore")
	} else {
		fmt.Println("using cloud datastore")
	}
	ds_client.Init(*useFilestore)

	// Create app.
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})

	// CORS middleware.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// IAP middleware.
	if *useIAP {
		fmt.Println("using IAP for authentication")
		app.Use(validateJWT(*jwtAudience))
	}

	// Register routes.
	registerAPIRoutes(app)

	// Start web server.
	listenOn := fmt.Sprintf(":%d", *port)
	fmt.Printf("starting web server on %s\n", listenOn)
	app.Listen(listenOn)
}
