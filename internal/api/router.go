package api

import (
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/gofiber/fiber/v2"
)

func AddTodoRoutes(app *fiber.App, env config.EnvVars, spannerController *SpannerController) {

	app.Use(AppConfigMiddleware(&env))
	api := app.Group("/api")

	// Anything related to spotify authentication and Spanner related account data
	account := api.Group("/account")

	// Anything available on the user's spotify profile or any data related directly to them.
	profile := api.Group("/profile")

	// playlist := api.Group("/playlist")

	account.Get("/login", spannerController.Login)
	account.Get("/logout", spannerController.Logout)
	account.Get("/callback", spannerController.CompleteAuth)
	account.Get("/check", spannerController.GetLoggedStatus)

	profile.Get("/toptracks/:timerange", spannerController.TopTracks)
	// profile.Get("/topartists/:timerange")
	// profile.Get("/playlists", spannerController.GetAllUserPlaylistIds)
	profile.Get("/name", spannerController.GetName)

}
