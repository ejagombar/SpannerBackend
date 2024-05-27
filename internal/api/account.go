package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
)

const redirectURI = "http://localhost:8080/account/callback"

func AppConfigMiddleware(env *config.EnvVars) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("env", env)
		return c.Next()
	}
}

// Generate random state string
func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// Returns a login URL to authenticate with Spotify.
// A random state is created and stored in the session cookie.
func (s *SpannerController) Login(c *fiber.Ctx) error {
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

	if err := sess.Save(); err != nil {
		fmt.Println(err)
		return err
	}

	address := spotify.GetLoginURL(env.CLIENT_ID, env.CLIENT_SECRET, state)
	return c.SendString(address)
}

// This method completes the OAuth flow by exchanging an authorisation code
// for an access token which are then stored in the session cookie. THis allows future API calls
// to retreive the API key from the user to perform API requests against the spotify API
func (s *SpannerController) CompleteAuth(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)
	auth := spotify.CreateAuthRequest(env.CLIENT_ID, env.CLIENT_SECRET)

	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	if state := c.FormValue("state"); state != sess.Get("state") {
		return fmt.Errorf("state mismatch")
	}

	tok, err := auth.Exchange(c.Context(), c.Query("code"))
	if err != nil {
		fmt.Print(err)
	}

	sess.Set("authed", true)
	sess.Set("accessToken", tok.AccessToken)
	sess.Set("refreshToken", tok.RefreshToken)
	sess.Set("tokenExpiry", tok.Expiry.Format(time.RFC1123Z))

	if err := sess.Save(); err != nil {
		return err
	}

	c.Set("Content-Type", "text/html")

	js := `
        <script>
            window.close();
        </script>
    `
	return c.SendString(js)
}

func (s *SpannerController) LoggedStatus(c *fiber.Ctx) error {
	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	authed := sess.Get("authed")
	str := fmt.Sprintf("%v", authed)

	return c.SendString(str)
}

func (s *SpannerController) Logout(c *fiber.Ctx) error {
	sess, err := s.session.Get(c)
	if err != nil {
		return fmt.Errorf("Here we are %w", err)
	}

	if err := sess.Destroy(); err != nil {
		return fmt.Errorf("Here we are2 %w", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
