package api

import (
	"fmt"

	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2/auth"
)

var (
	auth *spotifyauth.Authenticator
)

const redirectURI = "http://localhost:8080/callback"

func Login(c *fiber.Ctx) error {
	fmt.Fprintf(c, "test")
	address := spotify.GetLoginURL(spotify_id, spotify_client, "ed")
	return c.SendString(address)
}

// Handler function that is used to retrieve the token from the spotify authentication webpage
// This toek is used to create a client.
func CompleteAuth(c *fiber.Ctx) error {

	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithClientID(spotify_id),
		spotifyauth.WithClientSecret(spotify_client),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeStreaming,
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopeUserTopRead,
			spotifyauth.ScopeUserReadCurrentlyPlaying))

	code := c.Query("code")

	tok, err := auth.Exchange(c.Context(), code)
	if err != nil {
		fmt.Print(err)
	}

	if st := c.FormValue("state"); st != "d" {
		fmt.Printf("state mismatch")
	}

	// use the token to get an authenticated client
	// client := spotify.New(auth.Client(r.Context(), tok))
	// w.Header().Set("Content-Type", "text/html; charset=utf8")
	// fmt.Fprintf(w, form)

	fmt.Println("token: ", tok)
	// ch <- client
	// return c.SendString(tok.AccessToken)
	return nil
}
