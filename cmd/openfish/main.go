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
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/entities"
	"github.com/ausocean/openfish/cmd/openfish/handlers"
	"github.com/ausocean/openfish/cmd/openfish/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// registerAPIRoutes registers all handler functions to their routes.
func registerAPIRoutes(app *fiber.App) {

	v1 := app.Group("/api/v1")

	// Capture sources.
	v1.Group("/capturesources").
		Get("/:id", handlers.GetCaptureSourceByID).
		Get("/", handlers.GetCaptureSources).
		Post("/", middleware.Guard(entities.AdminRole), handlers.CreateCaptureSource).
		Patch("/:id", middleware.Guard(entities.AdminRole), handlers.UpdateCaptureSource).
		Delete("/:id", middleware.Guard(entities.AdminRole), handlers.DeleteCaptureSource)

	// Video streams.
	v1.Group("/videostreams").
		Get("/:id", handlers.GetVideoStreamByID).
		Get("/", handlers.GetVideoStreams).
		Post("/live", middleware.Guard(entities.CuratorRole), handlers.StartVideoStream).
		Patch("/:id/live", middleware.Guard(entities.CuratorRole), handlers.EndVideoStream).
		Post("/", middleware.Guard(entities.CuratorRole), handlers.CreateVideoStream).
		Patch("/:id", middleware.Guard(entities.CuratorRole), handlers.UpdateVideoStream).
		Delete("/:id", middleware.Guard(entities.AdminRole), handlers.DeleteVideoStream)

	// Annotations.
	v1.Group("/annotations").
		Get("/:id", handlers.GetAnnotationByID).
		Get("/", handlers.GetAnnotations).
		Post("/", middleware.Guard(entities.AnnotatorRole), handlers.CreateAnnotation).
		Delete("/:id", middleware.Guard(entities.AdminRole), handlers.DeleteAnnotation)

	// Species.
	v1.Group("/species").
		Get("/recommended", handlers.GetRecommendedSpecies).
		Get("/:id", handlers.GetSpeciesByID).
		Post("/", middleware.Guard(entities.AdminRole), handlers.CreateSpecies).
		Post("/import-from-inaturalist", middleware.Guard(entities.AdminRole), handlers.ImportFromINaturalist).
		Delete("/:id", middleware.Guard(entities.AdminRole), handlers.DeleteSpecies)

	// Users.
	v1.Group("/users", middleware.Guard(entities.AdminRole)).
		Get("/:email", handlers.GetUserByEmail).
		Get("/", handlers.GetUsers).
		Patch("/:email", middleware.Guard(entities.AdminRole), handlers.UpdateUser).
		Delete("/:email", middleware.Guard(entities.AdminRole), handlers.DeleteUser)

	// Auth.
	v1.Get("/auth/me", handlers.GetSelf)

	// Tasks.
	v1.Get("/tasks/:id/status", handlers.PollTask)

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

	// Recover from panics.
	app.Use(recover.New())

	// CORS middleware.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// IAP middleware.
	if *useIAP {
		fmt.Println("using IAP for authentication")
		app.Use(middleware.ValidateJWT(*jwtAudience))
	}

	// Register routes.
	registerAPIRoutes(app)

	// Start web server.
	listenOn := fmt.Sprintf(":%d", *port)
	fmt.Printf("starting web server on %s\n", listenOn)
	app.Listen(listenOn)
}
