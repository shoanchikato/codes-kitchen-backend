package route

import (
	h "codes-kitchen/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func StoveRoute(app fiber.Router, route h.StoveHandler) {
	app.Post("/switch", route.Switch())
	app.Get("/temp", route.Temperature())
}
