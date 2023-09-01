package spotify

import (
	"context"
	_ "embed"
	// "fmt"
	// "log"
	// "net/http"

	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/callback"

var (
	//go:embed callback.html
	form  string
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	tokCh = make(chan *oauth2.Token)
	state = "1234567IshouldProbablyChangeThis"
)

type TokenStore struct {
	token *oauth2.Token
}

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

// Starts the callback server, generates a link for the user to login with spotify, and waits
// until a client is recieved which is then returned from the function.
func GetLoginURL(spotify_id string, spotify_client string, state string) string {
	CreateAuthRequest(spotify_id, spotify_client)
	url := auth.AuthURL(state, spotifyauth.ShowDialog)
	return url
}

func Client(tok *oauth2.Token, ctx context.Context) (*spotify.Client, context.Context) {
	return spotify.New(auth.Client(ctx, tok)), ctx
}
