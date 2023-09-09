package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
	// "golang.org/x/oauth2"
)

var ()

type PlaylistData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ImageLink   string   `json:"imagelink"`
	TrackCount  int      `json:"trackcount"`
	TrackIDs    []string `json:"trackids"`
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

func GetAllUserPlaylists(client *spotify.Client, ctx context.Context, userID string) (userPlaylists []PlaylistData, err error) {
	playlists, err := client.GetPlaylistsForUser(ctx, userID, spotify.Limit(50))
	if err != nil {
		return nil, err
	}

	totalLength := playlists.Total
	userPlaylists = make([]PlaylistData, totalLength)
	totalDownloaded := 0

	for totalDownloaded < playlists.Total {
		length := len(playlists.Playlists)

		for i := 0; i < length; i++ {
			userPlaylists[totalDownloaded+i].ID = string(playlists.Playlists[i].ID)
			userPlaylists[totalDownloaded+i].Name = playlists.Playlists[i].Name
		}
		totalDownloaded += length

		playlists, err = client.GetPlaylistsForUser(ctx, userID, spotify.Limit(50), spotify.Offset(totalDownloaded))
		if err != nil {
			return nil, fmt.Errorf("Error: %w", err)
		}
	}

	return userPlaylists, nil
}
