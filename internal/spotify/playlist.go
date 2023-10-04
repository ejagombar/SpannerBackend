package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
	// "golang.org/x/oauth2"
)

type PlaylistData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ImageLink   string   `json:"imagelink"`
	TrackCount  int      `json:"trackcount"`
	TrackIDs    []string `json:"trackids"`
}

type PlaylistInfo struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	ImageLink         string   `json:"imagelink"`
	TopPlaylistTracks []Tracks `json:"topplaylisttracks"`
	AudioFeatures     AudioFeatures
}

type AudioFeatures struct {
	acousticness     int
	danceability     int
	energy           int
	instrumentalness int
	valence          int
	tempo            int
}

func GetPlaylistTopTracks(client *spotify.Client, playlistID string, idCount int) (idSubset []string, err error) {
	topTracks, err := getAllTopTrackIDs(client)
	if err != nil {
		return nil, err
	}

	playlistData, err := getPlaylistData(client, playlistID)
	if err != nil {
		return nil, err
	}

	commonElements := findCommonElements(topTracks, playlistData.TrackIDs)

	length := min(idCount, len(playlistData.TrackIDs))
	return selectIDSubset(commonElements, playlistData.TrackIDs, length)
}

func getPlaylistData(client *spotify.Client, playlistID string) (playlistData *PlaylistData, err error) {
	playlistOptions := "name,description,id"
	id := spotify.ID(playlistID)

	playlistRequest, err := client.GetPlaylist(context.Background(), id, spotify.Fields(playlistOptions))
	if err != nil {
		return nil, err
	}

	// Apparently GetPlaylistTracks is soon to be deprecated and replaced with GetPlayListItems.
	// GetPlaylistItems does not work with the fields argument so cannot be used
	playlistOptions = "limit,offset,total,items(track(id))"
	playlistItems, err := client.GetPlaylistTracks(context.Background(), id, spotify.Limit(50), spotify.Fields(playlistOptions))
	if err != nil {
		return nil, fmt.Errorf("Error:%w", err)
	}
	if len(playlistItems.Tracks) == 0 {
		return nil, fmt.Errorf("No tracks in playlist")
	}

	playlistData.ID = string(playlistRequest.ID)
	playlistData.Name = playlistRequest.Name
	playlistData.Description = playlistRequest.Description
	playlistData.TrackCount = playlistItems.Total

	totalDownloaded := 0
	playlistData.TrackIDs = make([]string, playlistData.TrackCount)

	for totalDownloaded < playlistData.TrackCount {
		playlistItems, err := client.GetPlaylistTracks(context.Background(), id, spotify.Limit(50), spotify.Fields(playlistOptions), spotify.Offset(totalDownloaded))
		if err != nil {
			return nil, fmt.Errorf("Error:%w", err)
		}

		length := len(playlistItems.Tracks)

		for i := 0; i < length; i++ {
			playlistData.TrackIDs[totalDownloaded+i] = string(playlistItems.Tracks[i].Track.ID)
		}
		totalDownloaded += length
	}

	return playlistData, nil
}

func GetPlaylistInfo(client *spotify.Client, playlistID string) (playlistInfo PlaylistInfo, err error) {
	playlistInfo = PlaylistInfo{}
	topTracks, err := getAllTopTrackIDs(client)
	if err != nil {
		return playlistInfo, err
	}

	playlistData, err := getPlaylistData(client, playlistID)
	if err != nil {
		return playlistInfo, err
	}

	topTrackIDs := findCommonElements(topTracks, playlistData.TrackIDs)

	length := min(100, len(playlistData.TrackIDs))
	randomSelectedIDs, err := selectIDSubset(topTrackIDs, playlistData.TrackIDs, length)

	fmt.Print(randomSelectedIDs)
	return playlistInfo, err

}

func GetTrackAudioFeatures(client *spotify.Client, ctx context.Context, trackIDs []string) (trackAudioFeatures []AudioFeatures, err error) {

	const maxTrackIDs = 100
	arrayLength := min(maxTrackIDs, len(trackIDs))
	var idArray = make([]spotify.ID, arrayLength)

	for i := 0; i < arrayLength; i++ {
		idArray[i] = spotify.ID(trackIDs[i])
	}

	playlists, err := client.GetAudioFeatures(ctx, idArray...)
	if err != nil {
		return nil, err
	}

	totalLength := len(trackIDs)
	trackAudioFeatures = make([]AudioFeatures, totalLength)
	totalDownloaded := 0

	for totalDownloaded < totalLength {
		length := len(playlists.Playlists)

		for i := 0; i < length; i++ {
			userPlaylists[totalDownloaded+i].TrackCount = playlists.Playlists[i].Tracks.Total
		}
		totalDownloaded += length

		playlists, err = client.GetPlaylistsForUser(ctx, userID, spotify.Limit(50), spotify.Offset(totalDownloaded))
		if err != nil {
			return nil, fmt.Errorf("Error getting user playlists: %w", err)
		}
	}

	return userPlaylists, nil
}
