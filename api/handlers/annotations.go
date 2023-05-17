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

package handlers

import (
	"fmt"
	"time"

	"github.com/ausocean/openfish/api/utils"

	"github.com/gofiber/fiber/v2"
)

type TimeSpan struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type BoundingBox struct {
	x1, x2, y1, y2 int
}

type AnnotationResult struct {
	ID            *int              `json:"id,omitempty"`
	VideoStreamID *int              `json:"videostreamId,omitempty"`
	TimeSpan      *TimeSpan         `json:"timespan,omitempty"`
	BoundingBox   *BoundingBox      `json:"boundingBox,omitempty"`
	Observer      *string           `json:"observer,omitempty"`
	Observation   map[string]string `json:"observation,omitempty"`
}

func GetAnnotationByID(ctx *fiber.Ctx) error {
	// TODO: implement handler

	// id, _ := ctx.ParamsInt("id", 1)
	return ctx.JSON("TODO")
}

func GetAnnotations(ctx *fiber.Ctx) error {
	timespan := ctx.Query("timespan")
	captureSource := ctx.Query("capturesource")
	species := ctx.Query("observation[species]")
	format := utils.GetFormat(ctx)
	limit, offset := utils.GetLimitAndOffset(ctx, 20)

	// Debugging info.
	fmt.Println(species, timespan, captureSource, format, limit, offset)

	// Placeholder code: returns an empty result.
	// TODO: implement fetching from datastore.
	result := utils.Result[AnnotationResult]{
		Results: []AnnotationResult{},
		Offset:  offset,
		Limit:   limit,
		Total:   0,
	}
	return ctx.JSON(result)
}

func CreateAnnotation(ctx *fiber.Ctx) error {
	// TODO: implement handler
	return ctx.JSON("TODO")
}
