package api

import (
	"github.com/gofiber/fiber/v2"
)

func AddTodoRoutes(app *fiber.App) {

	// add middlewares here

	// add routes here
	app.Get("/login", Login)
	app.Get("/callback", CompleteAuth)

}
