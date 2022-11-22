package main

import (
	h "codes-kitchen/pkg/handler"
	r "codes-kitchen/pkg/route"
	ss "codes-kitchen/pkg/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Static("/", "./public")

	mainSwitchService := ss.NewMainSwitchService()
	stoveService := ss.NewStoveService(mainSwitchService)
	toasterService := ss.NewToasterService(mainSwitchService)

	mainSwitchHandler := h.NewMainSwitchHandler(mainSwitchService)
	stoveHandler := h.NewStoveHandler(stoveService)
	toasterHandler := h.NewToasterHandler(toasterService)

	stoveApi := app.Group("/stove")
	toasterApi := app.Group("/toaster")
	mainSwitchApi := app.Group("/main-switch")

	r.StoveRoute(stoveApi, stoveHandler)
	r.ToasterRoute(toasterApi, toasterHandler)
	r.MainSwitchRoute(mainSwitchApi, mainSwitchHandler)

	log.Fatal(app.Listen(":3002"))
}
