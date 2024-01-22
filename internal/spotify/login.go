package spotify

import (
	"context"
	"time"

	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/api/account/callback"

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	tokCh = make(chan *oauth2.Token)
)

// Creates a authentication request with all the nessecary scopes needed for the CLI tool
func CreateAuthRequest(spotify_id string, spotify_client string) *spotifyauth.Authenticator {
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithClientID(spotify_id),
		spotifyauth.WithClientSecret(spotify_client),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeStreaming,
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopeImageUpload,
			spotifyauth.ScopeUserTopRead,
			spotifyauth.ScopeUserReadCurrentlyPlaying))

	return auth
}

func GetClient(accessTok, refreshTok, TokExpiry string) (*spotify.Client, error) {
	timeOut, err := time.Parse(time.RFC1123Z, TokExpiry)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
		Expiry:       timeOut,
	}

	return spotify.New(auth.Client(context.Background(), token)), nil
}
