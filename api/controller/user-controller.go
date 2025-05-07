package controller

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"pnm-todo-be/internal/service"
	"pnm-todo-be/pkg"
)

type UserController struct {
	FaceService service.UserService
}

func HandlerUserController(db *gorm.DB) *UserController {

	return &UserController{
		FaceService: service.UserService{
			DB: db,
		},
	}
}

func (controller *UserController) RegisterUser(c *fiber.Ctx) error {
	response := pkg.InitialResponse{Ctx: c}
	var input pkg.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return response.Respose(fiber.StatusUnprocessableEntity, err.Error(), true, nil)
	}

	result, err := controller.FaceService.RegisterUser(input)
	if err != nil {
		return response.Respose(fiber.StatusBadRequest, err.Error(), true, nil)
	}
	return response.Respose(fiber.StatusCreated, "success", false, result)
}
