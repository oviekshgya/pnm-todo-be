package routes

import (
	"pnm-todo-be/api/controller"
)

var (
	UserController *controller.UserController
)

func InitialRoute() {
	UserController = controller.HandlerUserController()
}
