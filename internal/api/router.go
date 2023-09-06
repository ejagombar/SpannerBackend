package api

import (
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/gofiber/fiber/v2"
)

func AddTodoRoutes(app *fiber.App, env config.EnvVars, spannerStorage *SpannerController) {

	api := app.Group("/api", AppConfigMiddleware(&env))

	// Anything related to spotify authentication and Spanner related account data
	account := api.Group("/account")

	account.Get("/login", spannerStorage.Login)
	account.Get("/logout", spannerStorage.Logout)
	account.Get("/callback", spannerStorage.CompleteAuth)
	account.Get("/check", spannerStorage.GetLogged)

	// Anything available on the user's spotify profile or any data related directly to them.
	profile := api.Group("/profile")

	profile.Get("/toptracks")
	profile.Get("/topartists")
	profile.Get("/playlists", spannerStorage.UserPlaylists)
	profile.Get("/name", spannerStorage.DisplayName)

	// playlist := api.Group("/playlist")

}
