package spotify

import (
	"github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/api/account/callback"

// Creates a authentication request with all the nessecary scopes needed for the CLI tool
func CreateAuthRequest(spotify_id string, spotify_client string) *spotifyauth.Authenticator {
	auth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
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
