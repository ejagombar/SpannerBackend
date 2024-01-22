package api

import (
	"github.com/gofiber/fiber/v2"
)

func AddTodoRoutes(app *fiber.App, spannerController *SpannerController) {
	api := app.Group("/api")

	// account := api.Group("/account")
	profile := api.Group("/profile")
	playlist := api.Group("/playlist")

	// Anything related to spotify authentication and Spanner related account data
	// account.Get("/login", Login)
	// account.Post("/logout", Logout)
	// account.Get("/callback", CompleteAuth)
	// account.Get("/authenticated", LoggedStatus)

	// Anything available on the user's spotify profile or any data related directly to them.
	profile.Get("/toptracks/:timerange", spannerController.TopTracks)
	profile.Get("/topartists/:timerange", spannerController.TopArtists)
	profile.Get("/userplaylists", spannerController.UserPlaylists)
	profile.Get("/info", spannerController.ProfileInfo)

	// Anything related to playlist analysis
	playlist.Get("/:id/toptracks/maxcount=:maxcount", spannerController.TopPlaylistTracks)
	playlist.Get("/:id/analysis", spannerController.PlaylistAnalysis)
}
