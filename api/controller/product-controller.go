package controller

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"pnm-todo-be/internal/service"
	"pnm-todo-be/pkg"
	"strconv"
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

	switch c.Method() {
	case fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete:
		result, err := controller.Product.CRUDProduct(input, c)
		if err != nil {
			return response.Respose(fiber.StatusBadRequest, err.Error(), true, nil)
		}
		return response.Respose(fiber.StatusAccepted, "success", false, result)
	default:
		return response.Respose(fiber.StatusMethodNotAllowed, "Method not allowed", true, nil)
	}

}

func (controller *ProductController) GetProduct(c *fiber.Ctx) error {
	response := pkg.InitialResponse{Ctx: c}

	id, _ := strconv.Atoi(c.Query("id", "0"))
	result, err := controller.Product.GetProduct(id, c)
	if err != nil {
		return response.Respose(fiber.StatusBadRequest, err.Error(), true, nil)
	}
	return response.Respose(fiber.StatusOK, "success", false, result)
}
