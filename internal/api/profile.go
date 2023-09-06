package api

import (
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
)

func (s *SpannerController) TopTracks(c *fiber.Ctx) error {
	tokenData, err := s.getTokenData(c)
	if err != nil {
		return err
	}

	client, err := spotify.GetClient(c.Context(), tokenData)
	if err != nil {
		return err
	}

	topTracks, err := spotify.RequestTopTracks(client, c.Context(), "short_term")
	if err != nil {
		return err
	}

	// return c.Status(fiber.StatusOK).JSON(topTracks) will add this along with better error handling later
	return c.JSON(topTracks)
}

func (s *SpannerController) GetName(c *fiber.Ctx) error {
	tokenData, err := s.getTokenData(c)
	if err != nil {
		return err
	}

	client, err := spotify.GetClient(c.Context(), tokenData)
	if err != nil {
		return err
	}

	str, err := spotify.GetUserName(client, c.Context())
	if err != nil {
		return err
	}

	return c.SendString(str)
}

func (s *SpannerController) GetAllUserPlaylist(c *fiber.Ctx) error {
	tokenData, err := s.getTokenData(c)
	if err != nil {
		return err
	}

	client, err := spotify.GetClient(c.Context(), tokenData)
	if err != nil {
		return err
	}

	userID, err := spotify.GetUserID(client, c.Context())
	if err != nil {
		return err
	}

	playlistIDs, err := spotify.GetAllUserPlaylists(client, c.Context(), userID)
	if err != nil {
		return err
	}

	out := ""
	for i := 0; i < len(playlistIDs); i++ {
		out = out + playlistIDs[i].Name + "\n"
	}

	return c.SendString(out)
}
