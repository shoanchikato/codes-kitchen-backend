package route

import (
	h "codes-kitchen/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func MainSwitchRoute(app fiber.Router, route h.MainSwitchHandler) {
	app.Post("/switch", route.Switch())
	app.Get("/is-on", route.IsOn())
}
