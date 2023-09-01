package spotify

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
)

var (
	playlistID = spotify.ID("6Wf5xOSUPT0CU4IZPpTTs2")
	// playlistID = spotify.ID("06iUfe2cSQJwdXvk77W2Me")

)

func getPlaylistData(client *spotify.Client, playlistData *PlaylistData) (err error) {
	playlistOptions := "name,description,id"
	playlistRequest, err := client.GetPlaylist(context.Background(), playlistID, spotify.Fields(playlistOptions))
	if err != nil {
		return err
	}

	// Apparently GetPlaylistTracks is soon to be deprecated and replaced with GetPlayListItems.
	// GetPlaylistItems does not work with the fields argument so cannot be used
	playlistOptions = "limit,offset,total,items(track(id))"
	playlistItems, err := client.GetPlaylistTracks(context.Background(), playlistID, spotify.Limit(50), spotify.Fields(playlistOptions))
	if err != nil {
		return fmt.Errorf("Error:%w", err)
	}
	if len(playlistItems.Tracks) == 0 {
		return fmt.Errorf("No tracks in playlist")
	}

	playlistData.ID = string(playlistRequest.ID)
	playlistData.Name = playlistRequest.Name
	playlistData.Description = playlistRequest.Description
	playlistData.TrackCount = playlistItems.Total

	totalDownloaded := 0
	playlistData.TrackIDs = make([]string, playlistData.TrackCount)

	for totalDownloaded < playlistData.TrackCount {
		playlistItems, err := client.GetPlaylistTracks(context.Background(), playlistID, spotify.Limit(50), spotify.Fields(playlistOptions), spotify.Offset(totalDownloaded))
		if err != nil {
			return fmt.Errorf("Error:%w", err)
		}

		length := len(playlistItems.Tracks)

		for i := 0; i < length; i++ {
			playlistData.TrackIDs[totalDownloaded+i] = string(playlistItems.Tracks[i].Track.ID)
		}
		totalDownloaded += length
	}

	return nil
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

func requestAndSavePlaylist(client *spotify.Client, fileName string, playlistData *PlaylistData) (err error) {
	err = getPlaylistData(client, playlistData)
	if err != nil {
		return err
	}
	fileName = fmt.Sprintf("%v.json", playlistData.ID)
	SaveStruct(fileName, playlistData)

	return nil
}
