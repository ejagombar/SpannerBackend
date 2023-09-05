package api

import (
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
)

func (s *SpannerController) TopTracks(c *fiber.Ctx) error {
	token, err := s.retrieveToken(c)
	if err != nil {
		return err
	}

	client, _ := spotify.Client(token, c.Context())
	topTracks, err := spotify.RequestTopTracks(client, "short_term")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(topTracks)
}
