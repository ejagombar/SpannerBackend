package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

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

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *SpannerStorage) Login(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)

	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	state, err := generateState(16)
	if err != nil {
		return err
	}

	sess.Set("authed", false)
	sess.Set("state", state)

	fmt.Println(state)
	fmt.Println(sess.Get("state"))

	if err := sess.Save(); err != nil {
		return err
	}

	address := spotify.GetLoginURL(env.CLIENT_ID, env.CLIENT_SECRET, state)
	return c.SendString(address)
}

// Handler function that is used to retrieve the token from the spotify authentication webpage
// This toek is used to create a client.
func (s *SpannerStorage) CompleteAuth(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)
	auth = spotify.CreateAuthRequest(env.CLIENT_ID, env.CLIENT_SECRET)

	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}
	fmt.Println(sess.Get("state"))
	fmt.Println(c.FormValue("state"))

	if state := c.FormValue("state"); state != sess.Get("state") {
		fmt.Println("state mismatch")
		return fmt.Errorf("state mismatch")
	}

	code := c.Query("code")
	tok, err := auth.Exchange(c.Context(), code)
	if err != nil {
		fmt.Print(err)
	}

	sess.Set("authed", true)

	if err := sess.Save(); err != nil {
		return err
	}

	fmt.Println("token: ", tok)
	return c.SendString("WE ARE THERE")
}

func (s *SpannerStorage) GetLogged(c *fiber.Ctx) error {

	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	name := sess.Get("authed")
	str := fmt.Sprintf("%v", name)
	fmt.Println(str)

	return c.SendString(str)
}
