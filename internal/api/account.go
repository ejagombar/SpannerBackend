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

var ch = make(chan *spotify.Client)

func (s *SpannerController) GetClient(c *fiber.Ctx) (*spotify.Client, error) {
	if s.client == nil {
		storedToken, err := s.storage.GetToken()
		if err != nil {
			return nil, err
		}

		auth := processing.CreateAuthRequest(s.config.CLIENT_ID, s.config.CLIENT_SECRET)

		timeOut, err := time.Parse(time.RFC1123Z, storedToken.Expiry)
		if err != nil {
			address := processing.GetLoginURL(s.config.CLIENT_ID, s.config.CLIENT_SECRET, auth)
			print(address)
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

func (s *SpannerController) CompleteAuth(c *fiber.Ctx) error {
	auth := processing.CreateAuthRequest(s.config.CLIENT_ID, s.config.CLIENT_SECRET)

	token, err := auth.Exchange(c.Context(), c.Query("code"))
	if err != nil {
		return err
	}

	s.client = spotify.New(auth.Client(context.Background(), token))

	token, err = s.client.Token()
	if err != nil {
		return err
	}

	s.storage.SaveToken(ConvertTokenType(token))
	return nil
}
