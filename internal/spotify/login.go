package spotify 

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	//go:embed callback.html
	form  string
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "1234567IshouldProbablyChangeThis"
)

// Creates a authentication request with all the nessecary scopes needed for the CLI tool
func createAuthRequest(spotify_id string, spotify_client string) {
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
			spotifyauth.ScopeUserTopRead,
			spotifyauth.ScopeUserReadCurrentlyPlaying))
}

// Handler function that is used to retrieve the token from the spotify authentication webpage
// This toek is used to create a client.
func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	fmt.Fprintf(w, form)

	ch <- client
}

// Starts the callback server, generates a link for the user to login with spotify, and waits
// until a client is recieved which is then returned from the function.
func getNewClient(spotify_id string, spotify_client string) (client *spotify.Client) {
	createAuthRequest(spotify_id, spotify_client)
	http.HandleFunc("/callback", completeAuth)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client = <-ch

	return client
}
