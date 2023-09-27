package spotify

import (
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
