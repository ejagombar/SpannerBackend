package api

import (
	"fmt"
	"strconv"

	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
)

func (s *SpannerController) TopPlaylistTracks(c *fiber.Ctx) error {
	client, err := s.GetClient(c)
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

func (s *SpannerController) PlaylistAnalysis(c *fiber.Ctx) error {
	client, err := s.GetClient(c)
	if err != nil {
		return err
	}

	playlistID := fmt.Sprintf("%v", c.Params("id"))
	fmt.Println("playlistID", playlistID)

	playlistAnalysis, err := spotify.GetPlaylistInfo(client, c.Context(), playlistID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(playlistAnalysis)
}
