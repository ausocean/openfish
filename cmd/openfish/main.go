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
		Post(":import-from-inaturalist", middleware.Guard(entities.AdminRole), handlers.ImportFromINaturalist).
		Delete("/:id", middleware.Guard(entities.AdminRole), handlers.DeleteSpecies)

	// Users.
	v1.Group("/users", middleware.Guard(entities.AdminRole)).
		Get("/:email", handlers.GetUserByEmail).
		Get("/", handlers.GetUsers).
		Patch("/:email", handlers.UpdateUser).
		Delete("/:email", handlers.DeleteUser)

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

// main creates and starts the web server.
//
//	@title			OpenFish API
//	@version		1.0
//	@description	### Overview
//	@description	OpenFish provides an API to access stored marine footage, and video
//	@description	annotations / labels, allowing clients to retrieve and filter the data.
//	@description	Clients can download segments of footage or video annotations by querying
//	@description	by location, time, and other parameters.
//	@description
//	@description	### Authentication
//	@description	OpenFish has optional support for requiring user authentication.
//	@description
//	@description	User authentication is provided using Google Cloud's Identity Aware Proxy (IAP). By default it is disabled, to use it you need to pass the command line flag `--iap` or set the environmental variable `IAP="true"` to enable it.
//	@description
//	@description	### Roles and permissions
//	@description	If user authentication is enabled, the following roles and permissions apply:
//	@description	| Role               | Permissions                                                                       |
//	@description	| ------------------ | --------------------------------------------------------------------------------- |
//	@description	| Admin              | Can add and remove annotations, videostreams, capturesources, users, and species. |
//	@description	| Curator            | Can select streams for classification.                                            |
//	@description	| Annotator          | Can add annotations, and delete their own annotations                             |
//	@description	| Readonly (default) | A readonly user is only be able to look at annotations, not make any              |
//	@description
//	@termsOfService		http://swagger.io/terms/
//	@contact.name		Scott Barnard
//	@contact.email		scott@ausocean.org
//	@license.name		BSD 3
//	@license.url		https://github.com/ausocean/openfish/blob/master/LICENSE
//	@host				http://localhost:8080
//	@tag.name			Annotations
//	@tag.description	Annotations are used for labeling interesting things in videos. They store a linked video stream, a bounding box, a start and end time, the observer's name, and the observations themselves. Observations have a flexible format, they use key-value pairs so you can add all sorts of different information. Most commonly used is species=<species name>.
//	@tag.description	OpenFish provides APIs to create, retrieve, update (under development) and delete (under development) annotations, and features to query annotations by the person who made the observation, what kind of observations were made (presence of a key), what was observed (presence of a key and given value), and by the location (under development), video stream or capture source (under development).
//	@tag.name			Capture Sources
//	@tag.description	Capture sources are cameras that produces video streams. AusOcean may have one or many capture sources at a rig or jetty,  depending on how many cameras are set up.
//	@tag.name			Video Streams
//	@tag.description	The video stream API provides the metadata for video streams. A video stream has a start time, end time, stream URL and linked capture source. Using the video stream API we can register our video streams with OpenFish so it can annotate and play back that stream. The stream URL specifies where the video data is stored. Examples: http://vidgrind.ausocean.org/get?id=1, https://www.youtube.com/watch?v=abcdefghijk
//	@tag.name			Video Streams (Live)
//	@tag.description	Live streams are different to registering an existing video. This is because we don't know the end time when we start it. To register a stream when it starts use POST. It takes the current time as the start time. To finish a stream use PATCH. It uses the current time as the end time. See also: Video Streams
//	@tag.name			Species
//	@tag.description	Species are used for providing suggestions to our users when annotating videos. They have the scientific and common name, and an images or images. Images have a source and attribution - we use this to give the author credit and to abide by the rules of the license.
//	@tag.name			Users
//	@tag.description	A user is identified by their email and has a role that gives them permissions. A user is created when they first login to OpenFish. There are APIs for updating user's role, listing users and deleting a user account.
//	@tag.name			Authentication
//	@tag.description	Operations related to user authentication & authorization.
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
