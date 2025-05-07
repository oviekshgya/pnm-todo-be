package routes

import (
	"pnm-todo-be/api/controller"
	"pnm-todo-be/db"
)

var (
	UserController    *controller.UserController
	ProductController *controller.ProductController
)

func InitialRoute() {
	UserController = controller.HandlerUserController(db.ConnDB)
	ProductController = controller.HandlerProductController(db.ConnDB)
}
