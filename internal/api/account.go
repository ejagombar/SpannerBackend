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
	"github.com/gofiber/fiber/v2/middleware/session"
)

const redirectURI = "http://localhost:8080/account/callback"

type SpannerController struct {
	session *session.Store
}

func NewSpannerStorage(session *session.Store) *SpannerController {
	return &SpannerController{session: session}
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

func (s *SpannerController) CompleteAuth(c *fiber.Ctx) error {
	env := c.Locals("env").(*config.EnvVars)
	auth := spotify.CreateAuthRequest(env.CLIENT_ID, env.CLIENT_SECRET)

	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	// IMPORTANT: STATE IS CURRENTLY NOT BEING CHECKED
	// if state := c.FormValue("state"); state != sess.Get("state") {
	// 	return fmt.Errorf("state mismatch")
	// }

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

	name := sess.Get("authed")
	str := fmt.Sprintf("%v", name)
	fmt.Println(str)

	return c.SendString("User" + str)
}

func (s *SpannerController) Logout(c *fiber.Ctx) error {
	sess, err := s.session.Get(c)
	if err != nil {
		return err
	}

	if err := sess.Destroy(); err != nil {
		return err
	}

	return c.SendString("Logged out")
}
