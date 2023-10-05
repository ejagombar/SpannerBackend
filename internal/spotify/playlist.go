package spotify

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
	Acousticness     float32
	Danceability     float32
	Energy           float32
	Instrumentalness float32
	Valence          float32
	Tempo            float32
	Loudness         float32
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
	playlistOptions := "name,description,id"
	id := spotify.ID(playlistID)

	fmt.Println("getting playlistRequest")
	playlistRequest, err := client.GetPlaylist(context.Background(), id, spotify.Fields(playlistOptions))
	if err != nil {
		return playlistData, err
	}

	// Apparently GetPlaylistTracks is soon to be deprecated and replaced with GetPlayListItems.
	// GetPlaylistItems does not work with the fields argument so cannot be used
	playlistOptions = "limit,offset,total,items(track(id))"
	playlistItems, err := client.GetPlaylistTracks(context.Background(), id, spotify.Limit(50), spotify.Fields(playlistOptions))
	if err != nil {
		return playlistData, fmt.Errorf("Error:%w", err)
	}

	playlistData.ID = string(playlistRequest.ID)
	playlistData.Name = playlistRequest.Name
	playlistData.Description = playlistRequest.Description
	playlistData.TrackCount = playlistRequest.Tracks.Total
	playlistData.TrackCount = 951
	if playlistRequest.Tracks.Total == 0 {
		fmt.Println("ERROR NO TRAKCS")
	}

	totalDownloaded := 0
	playlistData.TrackIDs = make([]string, playlistData.TrackCount)

	fmt.Println("starting loop")
	for totalDownloaded < playlistData.TrackCount {
		playlistItems, err = client.GetPlaylistTracks(context.Background(), id, spotify.Limit(50), spotify.Fields(playlistOptions), spotify.Offset(totalDownloaded))
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

func GetPlaylistInfo(client *spotify.Client, ctx context.Context, playlistID string) (playlistInfo PlaylistInfo, err error) {
	playlistInfo = PlaylistInfo{}
	topTracks, err := getAllTopTrackIDs(client)
	if err != nil {
		return playlistInfo, err
	}
	fmt.Println("topTracks length: ", len(topTracks))

	playlistData, err := getPlaylistData(client, playlistID)
	if err != nil {
		return playlistInfo, err
	}

	fmt.Println("playlist data: ", playlistData)

	topTrackIDs := findCommonElements(topTracks, playlistData.TrackIDs)

	fmt.Println("top track IDs", topTrackIDs)

	length := min(100, len(playlistData.TrackIDs))
	randomSelectedIDs, err := selectIDSubset(topTrackIDs, playlistData.TrackIDs, length)

	fmt.Println("randomSelectedIDs: ", randomSelectedIDs)

	selectedTrackAudioFeatures, err := GetTrackAudioFeatures(client, ctx, randomSelectedIDs)
	if err != nil {
		return playlistInfo, err
	}

	playlistInfo.AudioFeatures = calculateAverageFeatures(selectedTrackAudioFeatures)

	playlistInfo.TopPlaylistTracks, err = GetTracks(client, ctx, topTrackIDs)
	if err != nil {
		return playlistInfo, fmt.Errorf("Error getting top playlist tracks: %w", err)
	}

	return playlistInfo, nil
}
func calculateAverageFeatures(features []AudioFeatures) AudioFeatures {
	if len(features) == 0 {
		return AudioFeatures{}
	}

	var totalAcousticness, totalDanceability, totalEnergy, totalInstrumentalness, totalValence, totalTempo, totalLoudness float32

	for _, f := range features {
		totalAcousticness += f.Acousticness
		totalDanceability += f.Danceability
		totalEnergy += f.Energy
		totalInstrumentalness += f.Instrumentalness
		totalValence += f.Valence
		totalTempo += f.Tempo
		totalLoudness += f.Loudness
	}

	averageFeatures := AudioFeatures{
		Acousticness:     totalAcousticness / float32(len(features)),
		Danceability:     totalDanceability / float32(len(features)),
		Energy:           totalEnergy / float32(len(features)),
		Instrumentalness: totalInstrumentalness / float32(len(features)),
		Valence:          totalValence / float32(len(features)),
		Tempo:            totalTempo / float32(len(features)),
		Loudness:         totalLoudness / float32(len(features)),
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
		trackAudioFeatures[i].Tempo = track[i].Tempo
		trackAudioFeatures[i].Energy = track[i].Energy
		trackAudioFeatures[i].Valence = track[i].Valence
		trackAudioFeatures[i].Acousticness = track[i].Acousticness
		trackAudioFeatures[i].Danceability = track[i].Danceability
		trackAudioFeatures[i].Instrumentalness = track[i].Instrumentalness
		trackAudioFeatures[i].Loudness = track[i].Loudness

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

	passedTracks := idArray[:(maxSpotifyTracks - 1)]
	spotifyTracks, err := client.GetTracks(ctx, passedTracks, spotify.Limit(maxSpotifyTracks))
	if err != nil {
		return nil, err
	}

	tracks = make([]Tracks, arrayLength)
	totalDownloaded := 0

	for totalDownloaded < arrayLength {
		length := len(spotifyTracks)

		for i := 0; i < length; i++ {
			tracks[totalDownloaded+i].ID = string(spotifyTracks[i].ID)
			tracks[totalDownloaded+i].Name = spotifyTracks[i].Name
			tracks[totalDownloaded+i].Artist = spotifyTracks[i].Artists[0].Name
			tracks[totalDownloaded+i].ImageURL = spotifyTracks[i].Album.Images[0].URL

		}
		totalDownloaded += length

		passedTracks = idArray[totalDownloaded:(totalDownloaded + maxSpotifyTracks)]
		spotifyTracks, err = client.GetTracks(ctx, passedTracks)
		if err != nil {
			return nil, fmt.Errorf("Error getting user playlists: %w", err)
		}
	}

	return tracks, nil
}
