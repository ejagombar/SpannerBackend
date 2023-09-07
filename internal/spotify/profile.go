package spotify

import (
	"context"
	"github.com/zmb3/spotify/v2"
)

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
