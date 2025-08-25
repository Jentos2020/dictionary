package handler

import (
	"fmt"
	"leetgo/internal/app/controller"
	"leetgo/internal/gen"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	c *controller.Controller
}

func New(ctr *controller.Controller) *fiber.App {
	app := fiber.New()

	h := &Handler{
		ctr,
	}

	api := app.Group("/api")

	gen.RegisterHandlers(api, h)

	app.Get("/ws/search", websocket.New(WSHandler(h)))

	routes := app.Stack()
	for _, handlers := range routes {
		for _, route := range handlers {
			ctr.Logger.Debug(fmt.Sprintf("Route %s %s", route.Method, route.Path))
		}
	}

	return app
}
