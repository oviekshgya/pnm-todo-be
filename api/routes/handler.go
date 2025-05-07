package routes

import (
	"pnm-todo-be/api/controller"
	"pnm-todo-be/db"
)

var (
	UserController *controller.UserController
)

func InitialRoute() {
	UserController = controller.HandlerUserController(db.ConnDB)
}
