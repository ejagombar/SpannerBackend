package api

import (
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/gofiber/fiber/v2"
)

func AddTodoRoutes(app *fiber.App, env config.EnvVars) {

	// add middlewares here
	app.Use(AppConfigMiddleware(&env))

	// add routes here
	app.Get("/login", Login)
	app.Get("/callback", CompleteAuth)

}
