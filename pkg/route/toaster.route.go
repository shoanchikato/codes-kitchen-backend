package route

import (
	h "codes-kitchen/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func ToasterRoute(app fiber.Router, route h.ToasterHandler) {
	app.Post("/switch", route.Switch())
	app.Get("/is-on", route.IsOn())
}
