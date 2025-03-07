//go:build !inat

package features

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterINaturalistImport(group fiber.Router) {
	// No-op
}
