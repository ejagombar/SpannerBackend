package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/ejagombar/SpannerBackend/internal/spotify"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func (s *SpannerController) retrieveToken(c *fiber.Ctx) (token *oauth2.Token, err error) {
	sess, err := s.session.Get(c)
	if err != nil {
		return nil, err
	}

	timeOut, err := time.Parse(time.RFC1123Z, fmt.Sprintf("%v", sess.Get("tokenExpiry")))
	if err != nil {
		return nil, err
	}

	token = new(oauth2.Token)
	token.AccessToken = fmt.Sprintf("%v", sess.Get("accessToken"))
	token.RefreshToken = fmt.Sprintf("%v", sess.Get("refreshToken"))
	token.Expiry = timeOut

	return token, nil
}

func (s *SpannerController) DisplayName(c *fiber.Ctx) error {
	token, err := s.retrieveToken(c)
	if err != nil {
		return err
	}

	str, err := spotify.GetUserName(spotify.Client(token, c.Context()))
	if err != nil {
		return err
	}

	return c.SendString(str)
}

func (s *SpannerController) TopPlaylistSongs(c *fiber.Ctx) error {
	token, err := s.retrieveToken(c)
	if err != nil {
		return err
	}

	client, ctx := spotify.Client(token, c.Context())
	idSubset, err := spotify.GetTopPlaylistSongs(client, ctx, playlistID, 30)
	if err != nil {
		return err
	}

	out := strings.Join(idSubset, "\n")
	return c.SendString(out)
}

func (s *SpannerController) UserPlaylists(c *fiber.Ctx) error {
	token, err := s.retrieveToken(c)
	if err != nil {
		return err
	}

	client, ctx := spotify.Client(token, c.Context())

	userID, err := spotify.GetUserID(client, ctx)
	if err != nil {
		return err
	}

	playlistIDs, err := spotify.GetAllUserPlaylists(client, ctx, userID)
	if err != nil {
		return err
	}

	out := ""
	for i := 0; i < len(playlistIDs); i++ {
		out = out + playlistIDs[i].Name + "\n"
	}

	return c.SendString(out)
}
