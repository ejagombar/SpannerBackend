package api

import (
	"fmt"
	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (s *SpannerController) TopPlaylistTracks(c *fiber.Ctx) error {
	tokenData, err := s.getTokenData(c)
	if err != nil {
		return err
	}

	client, err := spotify.GetClient(c.Context(), tokenData)
	if err != nil {
		return err
	}

	playlistID := fmt.Sprintf("%v", c.Params("id"))
	fmt.Println("playlistID", playlistID)
	maxItemCount, err := strconv.Atoi(c.Params("maxcount"))
	if err != nil {
		return err
	}

	topTracks, err := spotify.GetPlaylistTopTracks(client, playlistID, maxItemCount)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(topTracks)
}

func (s *SpannerController) TopPlaylistTracks(c *fiber.Ctx) error {
	tokenData, err := s.getTokenData(c)
	if err != nil {
		return err
	}

	client, err := spotify.GetClient(c.Context(), tokenData)
	if err != nil {
		return err
	}

}
