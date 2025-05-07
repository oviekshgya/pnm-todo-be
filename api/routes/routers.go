package routes

import (
	"github.com/gofiber/fiber/v2"
	"pnm-todo-be/api/middleware"
)

var (
	Router *fiber.App
)

func Route() {
	Router.Use(middleware.CORSMiddleware())
	Router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World <UNK>!")
	})
}
