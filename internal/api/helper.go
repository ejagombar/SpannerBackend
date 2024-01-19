package api

import (
	"fmt"

	"github.com/ejagombar/SpannerBackend/config"
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
)

func (s *SpannerController) getTokenData2(c *fiber.Ctx) (spotify.TokenData, error) {
	tokenData := spotify.TokenData{}

	env := c.Locals("env").(*config.EnvVars)

	tokenData.AccessToken = env.ACCESS_TOKEN
	tokenData.RefreshToken = env.REFRESH_TOKEN
	tokenData.Expiry = env.TOKEN_TIMEOUT

	return tokenData, nil
}

func (s *SpannerController) getTokenData(c *fiber.Ctx) (spotify.TokenData, error) {
	tokenData := spotify.TokenData{}

	sess, err := s.session.Get(c)
	if err != nil {
		return tokenData, err
	}

	tokenData.AccessToken = fmt.Sprintf("%v", sess.Get("accessToken"))
	tokenData.RefreshToken = fmt.Sprintf("%v", sess.Get("refreshToken"))
	tokenData.Expiry = fmt.Sprintf("%v", sess.Get("tokenExpiry"))

	return tokenData, nil
}
