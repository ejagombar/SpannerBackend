package api

import (
	"fmt"

	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
)

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
