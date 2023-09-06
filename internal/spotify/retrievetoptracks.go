package spotify

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/zmb3/spotify/v2"
)

func requestTopTrackIDs(client *spotify.Client, topTrackIDs []string) (err error) {
	totalDownloaded := 0
	timeRanges := [3]string{"short_term", "medium_term", "long_term"}
	for _, timeRange := range timeRanges {
		topTracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Limit(50), spotify.Timerange(spotify.Range(timeRange)))
		if err != nil {
			return err
		}

		length := len(topTracks.Tracks)

		for i := 0; i < length; i++ {
			topTrackIDs[totalDownloaded+i] = string(topTracks.Tracks[i].ID)
		}
		totalDownloaded += length
	}
	return nil
}

func requestAndSaveTopTracks(client *spotify.Client, topTrackIDs []string) (err error) {
	err = requestTopTrackIDs(client, topTrackIDs)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%v.json", "userTopTracks")
	SaveStruct(fileName, topTrackIDs)
	return nil
}

type Tracks struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	ImageURL string `json:"imageUrl"`
}

func RequestTopTracks(client *spotify.Client, ctx context.Context, timeRange string) (tracks []Tracks, err error) {

	topTracks, err := client.CurrentUsersTopTracks(ctx, spotify.Limit(50), spotify.Timerange(spotify.Range(timeRange)))
	if err != nil {
		return nil, err
	}

	length := len(topTracks.Tracks)

	tracks = make([]Tracks, length)

	for i := 0; i < length; i++ {
		tracks[i].ID = string(topTracks.Tracks[i].ID)
		tracks[i].Name = topTracks.Tracks[i].Name
		tracks[i].Artist = topTracks.Tracks[i].Artists[0].Name
		tracks[i].ImageURL = topTracks.Tracks[i].Album.Images[0].URL
	}

	return tracks, nil
}
