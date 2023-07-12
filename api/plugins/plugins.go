package plugins

import (
	"net/url"

	"github.com/ausocean/openfish/api/model"
	"github.com/gofiber/fiber/v2"
)

type VideoProviderPlugin interface {
	Get(ctx *fiber.Ctx, streamURL *url.URL, timespan model.TimeSpan) error
}

var VideoProviderPlugins map[string]VideoProviderPlugin

func Init() {
	VideoProviderPlugins = make(map[string]VideoProviderPlugin)
	VideoProviderPlugins["vidgrind.ausocean.org"] = VidgrindProvider{}
}
