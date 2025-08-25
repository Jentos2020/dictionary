package handler

import (
	"fmt"
	"leetgo/internal/app/controller"
	"leetgo/internal/gen"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(func(c *fiber.Ctx) error {
		ctr.Logger.Debug(fmt.Sprintf("Request: %s %s", c.Method(), c.Path()))
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://127.0.0.1:5173", // Оба варианта для WSL
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
	}))

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
