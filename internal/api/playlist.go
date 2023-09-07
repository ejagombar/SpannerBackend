package api

//
// import (
// 	"strings"
//
// 	"github.com/ejagombar/SpannerBackend/internal/spotify"
// 	"github.com/gofiber/fiber/v2"
// )
//
// func (s *SpannerController) TopPlaylistSongs(c *fiber.Ctx) error {
// 	token, err := s.retrieveToken(c)
// 	if err != nil {
// 		return err
// 	}
//
// 	client, ctx := spotify.Client(token, c.Context())
// 	idSubset, err := spotify.GetTopPlaylistSongIDs(client, ctx, playlistID, 30)
// 	if err != nil {
// 		return err
// 	}
//
// 	out := strings.Join(idSubset, "\n")
// 	return c.SendString(out)
// }
