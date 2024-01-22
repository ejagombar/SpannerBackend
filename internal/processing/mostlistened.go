package processing

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

type Tracks struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	ImageURL string `json:"imageUrl"`
}

type Artists struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

func getAllTopTrackIDs(client *spotify.Client) (topTrackIDs []string, err error) {
	totalDownloaded := 0
	timeRanges := [3]string{"short_term", "medium_term", "long_term"}
	topTrackIDs = make([]string, 150)

	for _, timeRange := range timeRanges {
		topTracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Limit(50), spotify.Timerange(spotify.Range(timeRange)))
		if err != nil {
			return nil, err
		}

		length := len(topTracks.Tracks)

		for i := 0; i < length; i++ {
			topTrackIDs[totalDownloaded+i] = string(topTracks.Tracks[i].ID)
		}
		totalDownloaded += length
	}
	return topTrackIDs, nil
}

func GetTopTracks(client *spotify.Client, ctx context.Context, timeRange string) (tracks []Tracks, err error) {
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

func GetTopArtists(client *spotify.Client, ctx context.Context, timeRange string) (artists []Artists, err error) {
	topArtists, err := client.CurrentUsersTopArtists(ctx, spotify.Limit(50), spotify.Timerange(spotify.Range(timeRange)))
	if err != nil {
		return nil, err
	}

	length := len(topArtists.Artists)

	artists = make([]Artists, length)

	for i := 0; i < length; i++ {
		artists[i].ID = string(topArtists.Artists[i].ID)
		artists[i].Name = topArtists.Artists[i].Name
		artists[i].ImageURL = topArtists.Artists[i].Images[0].URL
	}

	return artists, nil
}
