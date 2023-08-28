package api

import (
	"fmt"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2/auth"
)

var (
	auth *spotifyauth.Authenticator
)

const redirectURI = "http://localhost:8080/callback"

func AppConfigMiddleware(env *config.EnvVars) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("env", env)
		return c.Next()
	}
}

func Login(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)

	address := spotify.GetLoginURL(env.CLIENT_ID, env.CLIENT_SECRET, "ed")
	return c.SendString(address)
}

// Handler function that is used to retrieve the token from the spotify authentication webpage
// This toek is used to create a client.
func CompleteAuth(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)
	auth = spotify.CreateAuthRequest(env.CLIENT_ID, env.CLIENT_SECRET)

	if state := c.FormValue("state"); state != "ed" {
		fmt.Printf("state mismatch")
	}

	code := c.Query("code")
	tok, err := auth.Exchange(c.Context(), code)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("token: ", tok)
	return nil
}
