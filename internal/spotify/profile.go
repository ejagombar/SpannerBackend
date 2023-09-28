package spotify

import (
    "fmt"
	"context"
	"github.com/zmb3/spotify/v2"
)

type User struct {
	DisplayName   string `json:"displayname"`
	FollowerCount uint   `json:"followercount"`
	ImageURL      string `json:"imageurl"`
}

func GetUserName(client *spotify.Client, ctx context.Context) (string, error) {
	usr, err := client.CurrentUser(ctx)
	if err != nil {
		return "", err
	}
	return usr.DisplayName, nil
}

func GetUserID(client *spotify.Client, ctx context.Context) (string, error) {
	usr, err := client.CurrentUser(ctx)
	if err != nil {
		return "", err
	}
	return usr.ID, nil
}

func GetUserProfileInfo(client *spotify.Client, ctx context.Context) (User, error) {
	userInfo := User{}

	usr, err := client.CurrentUser(ctx)
	if err != nil {
		return userInfo, err
	}

	userInfo.DisplayName = usr.DisplayName
	userInfo.FollowerCount = usr.Followers.Count
	userInfo.ImageURL = usr.Images[0].URL

	return userInfo, nil
}

func AllUserPlaylistIds(client *spotify.Client, ctx context.Context, userID string) (userPlaylists []PlaylistData, err error) {
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
