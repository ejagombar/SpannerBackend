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
