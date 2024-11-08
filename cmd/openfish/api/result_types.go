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

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Result is the JSON format to use in response bodies for returning a list of results.
type Result[T any] struct {
	Results []T `json:"results"`
	Offset  int `json:"offset" example:"0"`
	Limit   int `json:"limit" example:"20"`
	Total   int `json:"total" example:"1"`
}

// Failure is the JSON format to use in response bodies for returning errors.
type Failure struct {
	Message string `json:"message" example:"error message here"`
}

// DatastoreReadFailure returns an error for datastore read failures.
func DatastoreReadFailure(err error) error {
	return fiber.NewError(500, fmt.Errorf("could not read from datastore: %w", err).Error())
}

// DatastoreWriteFailure returns an error for datastore write failures.
func DatastoreWriteFailure(err error) error {
	return fiber.NewError(500, fmt.Errorf("could not write to datastore: %w", err).Error())
}

// InvalidRequestJSON returns an error for requests with invalid JSON.
func InvalidRequestJSON(err error) error {
	return fiber.NewError(400, fmt.Errorf("invalid JSON in request: %w", err).Error())
}

// InvalidRequestURL returns an error for requests with invalid URLs.
func InvalidRequestURL(err error) error {
	return fiber.NewError(400, fmt.Errorf("invalid URL in request: %w", err).Error())
}

func Unauthorized(err error) error {
	return fiber.NewError(401, fmt.Errorf("Unauthorized: %w", err).Error())
}

func Forbidden(err error) error {
	return fiber.NewError(403, fmt.Errorf("Forbidden: %w", err).Error())
}

func NotAcceptable() error {
	return fiber.NewError(406)
}

func NotFound() error {
	return fiber.NewError(404)
}
