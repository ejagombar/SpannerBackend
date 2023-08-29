package api

import (
	"fmt"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/zmb3/spotify/v2/auth"
)

var (
	auth *spotifyauth.Authenticator
)

const redirectURI = "http://localhost:8080/callback"

type SpannerStorage struct {
	session *session.Store
}

func NewSpannerStorage(session *session.Store) *SpannerStorage {
	return &SpannerStorage{session: session}
}

func AppConfigMiddleware(env *config.EnvVars) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("env", env)
		return c.Next()
	}
}

func (s *SpannerStorage) Login(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)

	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	sess.Set("authed", "false")

	if err := sess.Save(); err != nil {
		return err
	}

	address := spotify.GetLoginURL(env.CLIENT_ID, env.CLIENT_SECRET, "ed")
	return c.SendString(address)
}

// Handler function that is used to retrieve the token from the spotify authentication webpage
// This toek is used to create a client.
func (S *SpannerStorage) CompleteAuth(c *fiber.Ctx) error {
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
