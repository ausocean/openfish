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
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ausocean/openfish/cmd/openfish/api"
	"github.com/ausocean/openfish/cmd/openfish/features"
	"github.com/ausocean/openfish/cmd/openfish/globals"
	"github.com/ausocean/openfish/cmd/openfish/handlers"
	"github.com/ausocean/openfish/cmd/openfish/middleware"
	"github.com/ausocean/openfish/cmd/openfish/types/role"

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
		Post("/", middleware.Guard(role.Admin), handlers.CreateCaptureSource).
		Patch("/:id", middleware.Guard(role.Admin), handlers.UpdateCaptureSource).
		Delete("/:id", middleware.Guard(role.Admin), handlers.DeleteCaptureSource)

	// Video streams.
	v1.Group("/videostreams").
		Get("/:id", handlers.GetVideoStreamByID).
		Get("/:id/media/:type/:subtype", middleware.Guard(role.Admin), handlers.GetVideoStreamMedia).
		Delete("/:id/media/:type/:subtype", middleware.Guard(role.Admin), handlers.DeleteVideoStreamMedia).
		Get("/", handlers.GetVideoStreams).
		Post("/live", middleware.Guard(role.Curator), handlers.StartVideoStream).
		Patch("/:id/live", middleware.Guard(role.Curator), handlers.EndVideoStream).
		Post("/", middleware.Guard(role.Curator), handlers.CreateVideoStream).
		Patch("/:id", middleware.Guard(role.Curator), handlers.UpdateVideoStream).
		Delete("/:id", middleware.Guard(role.Admin), handlers.DeleteVideoStream)

	// Annotations.
	v1.Group("/annotations").
		Get("/:id", handlers.GetAnnotationByID).
		Get("/", handlers.GetAnnotations).
		Post("/", middleware.Guard(role.Annotator), handlers.CreateAnnotation).
		Post("/:id/identifications/:species_id", middleware.Guard(role.Annotator), handlers.AddIdentification).
		Delete("/:id/identifications/:species_id", middleware.Guard(role.Annotator), handlers.DeleteIdentification).
		Delete("/:id", middleware.Guard(role.Admin), handlers.DeleteAnnotation)

	// Species.
	species := v1.Group("/species")
	features.RegisterINaturalistImport(species)
	species.
		Get("/", handlers.GetSpecies).
		Get("/:id", handlers.GetSpeciesByID).
		Post("/", middleware.Guard(role.Admin), handlers.CreateSpecies).
		Delete("/:id", middleware.Guard(role.Admin), handlers.DeleteSpecies)

	// Users.
	v1.Group("/users", middleware.Guard(role.Admin)).
		Get("/:id", handlers.GetUserByID).
		Get("/", handlers.GetUsers).
		Patch("/:id", handlers.UpdateUser).
		Delete("/:id", handlers.DeleteUser)

	// Auth.
	v1.Group("/auth/me").
		Get("/", handlers.GetSelf).
		Post("/", handlers.CreateSelf)

	// Tasks.
	v1.Get("/tasks/:id/status", handlers.PollTask)

}

// envOrFlag configures a setting using either an environment variable or a command-line flag.
// The environment variable value is parsed and used as the default if valid, but a flag value takes priority.
// Returns a pointer to the final value managed by the flag package.
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
// The following comments use swaggo notation, which is a tool that generates OpenAPI/Swagger documentation
// from Go source code annotations. These annotations define API metadata.
//
//	@title				OpenFish API
//	@version			1.0
//	@description		OpenFish API
//	@termsOfService		http://swagger.io/terms/
//	@contact.name		Scott Barnard
//	@contact.email		scott@ausocean.org
//	@license.name		BSD 3
//	@license.url		https://github.com/ausocean/openfish/blob/master/LICENSE
//	@host				https://openfish.appspot.com
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
//	@tag.name			Tasks
//	@tag.description	Some operations are long-running and execute asynchronously. These APIs return immediately with a task ID. You track task progress by polling the task API endpoint.
//	@tag.name			Media
//	@tag.description	Media is video or images that can be downloaded to be used as training data from annotated video streams.
//	@title				OpenFish API
//	@version			1.0
//	@description		OpenFish API
//	@termsOfService		http://swagger.io/terms/
//	@contact.name		Scott Barnard
//	@contact.email		scott@ausocean.org
//	@license.name		BSD 3
//	@license.url		https://github.com/ausocean/openfish/blob/master/LICENSE
//	@host				https://openfish.appspot.com
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
//	@tag.name			Tasks
//	@tag.description	Some operations are long-running and execute asynchronously. These APIs return immediately with a task ID. You track task progress by polling the task API endpoint.
func main() {

	// Get app configuration. Configurations can be set using environment variables or command line arguments.
	// Command line arguments take priority.

	port := envOrFlag("port", "PORT", "Port to listen on", 8080, strconv.Atoi, flag.Int)
	useFilestore := envOrFlag("filestore", "FILESTORE", "Use local datastore", false, strconv.ParseBool, flag.Bool)
	useCloudStorage := envOrFlag("cloud-storage", "CLOUD-STORAGE", "Use Cloud Buckets for storage of large data", false, strconv.ParseBool, flag.Bool)
	useIAP := envOrFlag("iap", "IAP", "Use Google's Identity Aware Proxy for authentication", false, strconv.ParseBool, flag.Bool)
	jwtAudience := envOrFlag("jwt-audience", "JWT_AUDIENCE", "Audience to use to validate JWT token", "", parseString, flag.String)

	flag.Parse()

	// Datastore setup.
	globals.InitStore(*useFilestore)

	// Storage setup.
	err := globals.InitStorage(!*useCloudStorage)
	if err != nil {
		panic(err.Error())
	}

	// Create app.
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

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
	} else {
		app.Use(middleware.NoAuth())
	}

	// Register routes.
	registerAPIRoutes(app)

	// Start web server.
	listenOn := fmt.Sprintf(":%d", *port)
	fmt.Printf("starting web server on %s\n", listenOn)
	app.Listen(listenOn)
}
