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

package api

import "github.com/gofiber/fiber/v2"

// JSON format of the response body for a list of results.
type Result[T any] struct {
	Results []T `json:"results"`
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
}

// JSON format of the response body for a failure.
type Failure struct {
	Message string `json:"message"`
}

// HTTP status code and JSON for datastore read failure.
func DatastoreReadFailure(ctx *fiber.Ctx) error {
	return ctx.
		Status(fiber.StatusInternalServerError).
		JSON(Failure{Message: "could not read from datastore"})
}

// HTTP status code and JSON for datastore write failure.
func DatastoreWriteFailure(ctx *fiber.Ctx) error {
	return ctx.
		Status(fiber.StatusInternalServerError).
		JSON(Failure{Message: "could not write to datastore"})
}

// HTTP status code and JSON for invalid request JSON.
func InvalidRequestJSON(ctx *fiber.Ctx) error {
	return ctx.
		Status(fiber.StatusBadRequest).
		JSON(Failure{Message: "invalid json in request"})
}

// HTTP status code and JSON for invalid request URL.
func InvalidRequestURL(ctx *fiber.Ctx) error {
	return ctx.
		Status(fiber.StatusBadRequest).
		JSON(Failure{Message: "invalid URL in request"})
}
