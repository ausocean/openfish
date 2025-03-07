//go:build inat

package features

import (
	"fmt"

	"github.com/ausocean/openfish/cmd/openfish/handlers"
	"github.com/ausocean/openfish/cmd/openfish/middleware"
	"github.com/ausocean/openfish/cmd/openfish/types/role"
	"github.com/gofiber/fiber/v2"
)

func RegisterINaturalistImport(group fiber.Router) {
	fmt.Println("INaturalist species import enabled")
	group.Post("/inaturalist-import", middleware.Guard(role.Admin), handlers.ImportFromINaturalist)
}
