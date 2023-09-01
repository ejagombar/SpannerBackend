package api

import (
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/gofiber/fiber/v2"
)

func AddTodoRoutes(app *fiber.App, env config.EnvVars, spannerStorage *SpannerController) {

	// add middlewares here
	app.Use(AppConfigMiddleware(&env))

	// add routes here
	app.Get("/login", spannerStorage.Login)
	app.Get("/logout", spannerStorage.Logout)
	app.Get("/callback", spannerStorage.CompleteAuth)
	app.Get("/user", spannerStorage.DisplayName)
	app.Get("/check", spannerStorage.GetLogged)
	app.Get("/top", spannerStorage.TopPlaylistSongs)
	app.Get("/playlists", spannerStorage.UserPlaylists)
}
