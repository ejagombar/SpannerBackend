package api

import (
	"context"
	"time"

	"github.com/ejagombar/SpannerBackend/internal/processing"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/account/callback"

func (s *SpannerController) GetClient(c *fiber.Ctx) (*spotify.Client, error) {
	auth := processing.CreateAuthRequest(s.config.CLIENT_ID, s.config.CLIENT_SECRET)

	accessTok, refreshTok, TokExpiry, err := s.storage.GetToken()
	if err != nil {
		return nil, err
	}
	print("\n access token: ", accessTok, "\n refresh token: ", refreshTok, "\n expiry: ", TokExpiry, "\n")

	timeOut, err := time.Parse(time.RFC1123Z, TokExpiry)

	token := &oauth2.Token{
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
		Expiry:       timeOut,
	}

	x := auth.Client(context.Background(), token)
	client := spotify.New(x)
	return client, nil
}
