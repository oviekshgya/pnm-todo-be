package routes

import (
	"github.com/gofiber/fiber/v2"
	"pnm-todo-be/api/middleware"
	"time"
)

var (
	Router *fiber.App
)

func Route() {
	Router.Use(middleware.CORSMiddleware())
	Router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World <UNK>!")
	})

	v1 := Router.Group("/v.01")
	v1.Use(middleware.BasicAuthMiddleware(), middleware.APIKeyMiddleware(), middleware.RateLimitMiddleware(5, 10*time.Second))
	{
		v1.Post("/create", UserController.RegisterUser)
		v1.Post("/login", UserController.Login)
	}

	product := v1.Group("/product")
	product.Use(middleware.AuthBearer())
	{
		product.All("/:jenis", ProductController.CRUD)
	}
}
