package controller

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"pnm-todo-be/internal/service"
	"pnm-todo-be/pkg"
)

type ProductController struct {
	Product service.ProductService
}

func HandlerProductController(db *gorm.DB) *ProductController {

	return &ProductController{
		Product: service.ProductService{
			DB: db,
		},
	}
}

func (controller *ProductController) CRUD(c *fiber.Ctx) error {
	response := pkg.InitialResponse{Ctx: c}
	var input pkg.CRUDProduct
	if err := c.BodyParser(&input); err != nil {
		return response.Respose(fiber.StatusUnprocessableEntity, err.Error(), true, nil)
	}

	result, err := controller.Product.CRUDProduct(input)
	if err != nil {
		return response.Respose(fiber.StatusBadRequest, err.Error(), true, nil)
	}
	return response.Respose(fiber.StatusAccepted, "success", false, result)
}
