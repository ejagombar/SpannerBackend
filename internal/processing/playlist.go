package processing

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

type PlaylistData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ImageLink   string   `json:"imagelink"`
	TrackCount  int      `json:"trackcount"`
	Followers   int      `json:"followers"`
	TrackIDs    []string `json:"trackids"`
}

type PlaylistAnalysisData struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	ImageLink         string         `json:"imagelink"`
	Followers         string         `json:"followers"`
	TrackCount        string         `json:"trackcount"`
	TopPlaylistTracks []Tracks       `json:"topplaylisttracks"`
	AudioFeatures     []AudioFeature `json:"audiofeatures"`
}

type AudioFeatures struct {
	Acousticness     float32 `json:"acousticness"`
	Danceability     float32 `json:"danceability"`
	Energy           float32 `json:"energy"`
	Instrumentalness float32 `json:"instrumental"`
	Valence          float32 `json:"valence"`
}

type AudioFeature struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
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

func getPlaylistData(client *spotify.Client, playlistID string) (playlistData PlaylistData, err error) {
	playlistData = PlaylistData{}
	id := spotify.ID(playlistID)

	playlistOptions := "name,description,id,images,tracks(total),followers"
	playlistRequest, err := client.GetPlaylist(context.Background(), id, spotify.Fields(playlistOptions))
	if err != nil {
		return playlistData, err
	}

	// Apparently GetPlaylistTracks is soon to be deprecated and replaced with GetPlayListItems.
	// GetPlaylistItems does not work with the fields argument so cannot be used
	playlistOptions = "limit,offset,total,items(track(id))"
	// playlistItems, err := client.GetPlaylistTracks(context.Background(), id, spotify.Limit(50), spotify.Fields(playlistOptions))
	// if err != nil {
	// 	return playlistData, fmt.Errorf("Error:%w", err)
	// }

	playlistData.ID = string(playlistRequest.ID)
	playlistData.Name = playlistRequest.Name
	playlistData.Description = playlistRequest.Description
	playlistData.Followers = int(playlistRequest.Followers.Count)

	if len(playlistRequest.Images) > 0 {
		playlistData.ImageLink = playlistRequest.Images[0].URL
	}
	playlistData.TrackCount = playlistRequest.Tracks.Total

	totalDownloaded := 0
	playlistData.TrackIDs = make([]string, playlistData.TrackCount)

	for totalDownloaded < playlistData.TrackCount {
		playlistItems, err := client.GetPlaylistTracks(context.Background(), id, spotify.Limit(50), spotify.Offset(totalDownloaded), spotify.Fields(playlistOptions))
		if err != nil {
			return playlistData, fmt.Errorf("Error:%w", err)
		}

		length := len(playlistItems.Tracks)

		for i := 0; i < length; i++ {
			playlistData.TrackIDs[totalDownloaded+i] = string(playlistItems.Tracks[i].Track.ID)
		}
		totalDownloaded += length

		fmt.Println("totalDownloaded: ", totalDownloaded)
	}

	return playlistData, nil
}

func GetPlaylistInfo(client *spotify.Client, ctx context.Context, playlistID string) (playlistInfo PlaylistAnalysisData, err error) {
	playlistInfo = PlaylistAnalysisData{}
	topTracks, err := getAllTopTrackIDs(client)
	if err != nil {
		return playlistInfo, err
	}

	playlistData, err := getPlaylistData(client, playlistID)
	if err != nil {
		return playlistInfo, err
	}

	playlistInfo.ID = playlistData.ID
	playlistInfo.Name = playlistData.Name
	playlistInfo.Description = playlistData.Description
	playlistInfo.ImageLink = playlistData.ImageLink
	playlistInfo.TrackCount = fmt.Sprint(playlistData.TrackCount)
	playlistInfo.Followers = fmt.Sprint(playlistData.Followers)

	topTrackIDs := findCommonElements(topTracks, playlistData.TrackIDs)

	length := min(100, len(playlistData.TrackIDs))
	randomSelectedIDs, err := selectIDSubset(topTrackIDs, playlistData.TrackIDs, length)

	selectedTrackAudioFeatures, err := GetTrackAudioFeatures(client, ctx, randomSelectedIDs)
	if err != nil {
		return playlistInfo, err
	}

	AudioFeatures := calculateAverageFeatures(selectedTrackAudioFeatures)
	initializeAudioFeatures(&playlistInfo, &AudioFeatures)

	playlistInfo.TopPlaylistTracks, err = GetTracks(client, ctx, topTrackIDs)
	if err != nil {
		return playlistInfo, fmt.Errorf("Error getting top playlist tracks: %w", err)
	}

	return playlistInfo, nil
}

func initializeAudioFeatures(playlistInfo *PlaylistAnalysisData, audioFeatures *AudioFeatures) {
	playlistInfo.AudioFeatures = []AudioFeature{
		{"Acousticness", audioFeatures.Acousticness},
		{"Danceability", audioFeatures.Danceability},
		{"Energy", audioFeatures.Energy},
		{"Instrumentalness", audioFeatures.Instrumentalness},
		{"Valence", audioFeatures.Valence},
	}
}

func calculateAverageFeatures(features []AudioFeatures) AudioFeatures {
	if len(features) == 0 {
		return AudioFeatures{}
	}

	var totalAcousticness, totalDanceability, totalEnergy, totalInstrumentalness, totalValence float32

	for _, f := range features {
		totalAcousticness += f.Acousticness
		totalDanceability += f.Danceability
		totalEnergy += f.Energy
		totalInstrumentalness += f.Instrumentalness
		totalValence += f.Valence
	}

	averageFeatures := AudioFeatures{
		Acousticness:     totalAcousticness / float32(len(features)),
		Danceability:     totalDanceability / float32(len(features)),
		Energy:           totalEnergy / float32(len(features)),
		Instrumentalness: totalInstrumentalness / float32(len(features)),
		Valence:          totalValence / float32(len(features)),
	}

	return averageFeatures
}

func GetTrackAudioFeatures(client *spotify.Client, ctx context.Context, trackIDs []string) (trackAudioFeatures []AudioFeatures, err error) {
	const maxTrackIDs = 100
	arrayLength := min(maxTrackIDs, len(trackIDs))
	var idArray = make([]spotify.ID, arrayLength)

	for i := 0; i < arrayLength; i++ {
		idArray[i] = spotify.ID(trackIDs[i])
	}

	track, err := client.GetAudioFeatures(ctx, idArray...)
	if err != nil {
		return nil, err
	}

	totalLength := len(trackIDs)
	trackAudioFeatures = make([]AudioFeatures, totalLength)

	for i := 0; i < arrayLength; i++ {
		trackAudioFeatures[i].Energy = track[i].Energy
		trackAudioFeatures[i].Valence = track[i].Valence
		trackAudioFeatures[i].Acousticness = track[i].Acousticness
		trackAudioFeatures[i].Danceability = track[i].Danceability
		trackAudioFeatures[i].Instrumentalness = track[i].Instrumentalness
	}

	return trackAudioFeatures, nil
}

func GetTracks(client *spotify.Client, ctx context.Context, trackIDs []string) (tracks []Tracks, err error) {
	const maxSpotifyTracks = 50
	arrayLength := len(trackIDs)
	var idArray = make([]spotify.ID, arrayLength)

	for i := 0; i < arrayLength; i++ {
		idArray[i] = spotify.ID(trackIDs[i])
	}

	tracks = make([]Tracks, arrayLength)
	totalDownloaded := 0

	for totalDownloaded < arrayLength {

		passedTracksLength := min(maxSpotifyTracks, (arrayLength - totalDownloaded))
		passedTracks := idArray[totalDownloaded:(totalDownloaded + passedTracksLength)]

		spotifyTracks, err := client.GetTracks(ctx, passedTracks)
		if err != nil {
			return nil, fmt.Errorf("Error getting user playlists: %w", err)
		}

		length := len(spotifyTracks)

		for i := 0; i < length; i++ {
			tracks[totalDownloaded+i].ID = string(spotifyTracks[i].ID)
			tracks[totalDownloaded+i].Name = spotifyTracks[i].Name
			tracks[totalDownloaded+i].Artist = spotifyTracks[i].Artists[0].Name
			if len(spotifyTracks[i].Album.Images) > 0 {
				tracks[totalDownloaded+i].ImageURL = spotifyTracks[i].Album.Images[0].URL
			}

		}

		totalDownloaded += length
	}

	return tracks, nil
}
