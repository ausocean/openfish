package plugins

import (
	"fmt"
	"net/url"

	"github.com/ausocean/openfish/api/model"
	"github.com/gofiber/fiber/v2"
)

type VidgrindProvider struct{}

func (vg VidgrindProvider) Get(ctx *fiber.Ctx, streamURL *url.URL, timespan model.TimeSpan) error {

	// Validate Stream URL.
	qry := streamURL.Query()
	_, ok := qry["id"]
	if !ok {
		panic("TODO")
	}

	// Add timestamp and duration to stream URL.
	// TODO: remove timestamp placeholders.
	qry.Add("ts", fmt.Sprintf("%d,%d", 1594547570, 10))
	streamURL.RawQuery = qry.Encode()

	// TODO: Add authorisation

	// Redirect to vidgrind's API for serving video.
	return ctx.Redirect(streamURL.String())
}
