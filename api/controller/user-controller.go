package controller

import "pnm-todo-be/internal/service"

type UserController struct {
	FaceService service.UserService
}

func HandlerUserController() *UserController {

	return &UserController{
		FaceService: service.UserService{},
	}
}
