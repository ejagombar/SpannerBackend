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
	if s.client == nil {

		auth := processing.CreateAuthRequest(s.config.CLIENT_ID, s.config.CLIENT_SECRET)

		storedToken, err := s.storage.GetToken()
		if err != nil {
			return nil, err
		}

		timeOut, err := time.Parse(time.RFC1123Z, storedToken.Expiry)
		if err != nil {
			return nil, err
		}

		token := &oauth2.Token{
			AccessToken:  storedToken.Access,
			RefreshToken: storedToken.Refresh,
			Expiry:       timeOut,
		}

		s.client = spotify.New(auth.Client(context.Background(), token))

		token, err = s.client.Token()
		if err != nil {
			return nil, err
		}

		s.storage.SaveToken(ConvertTokenType(token))
	}

	return s.client, nil
}

func ConvertTokenType(token *oauth2.Token) Token {
	return Token{
		Access:  token.AccessToken,
		Refresh: token.RefreshToken,
		Expiry:  token.Expiry.Format(time.RFC1123Z),
	}
}
