package spotify 

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"errors"

	"github.com/zmb3/spotify/v2"
)

func GetTopTracks(client *spotify.Client, playlistID string, idCount int) (idSubset []string, err error) {
	var playlistData PlaylistData

	fileName := fmt.Sprintf("%v.json", playlistID)
	err = LoadStruct(fileName, &playlistData)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = requestAndSavePlaylist(client, fileName, &playlistData)
		}
		log.Fatal(err)
	}

	topTracks := make([]string, 150)
	fileName = fmt.Sprintf("%v.json", "userTopTracks")
	err = LoadStruct(fileName, &topTracks)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = requestAndSaveTopTracks(client, topTracks)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	commonElements := findCommonElements(topTracks, playlistData.TrackIDs)

	length := min(idCount, len(playlistData.TrackIDs))
	return selectIDSubset(commonElements, playlistData.TrackIDs, length)
}

func findCommonElements(slice1, slice2 []string) []string {
	elementsMap := make(map[string]bool)

	for _, elem := range slice1 {
		elementsMap[elem] = true
	}

	var commonElements []string
	for _, elem := range slice2 {
		// Need to investigate how this map search is actually implemented.
		if elementsMap[elem] {
			commonElements = append(commonElements, elem)
		}
	}

	return commonElements
}

func addToSliceIfNotPresent(illegalElements, allElements []string) (out []string) {
	existingElements := make(map[string]bool)

	for _, element := range allElements {
		existingElements[element] = true
	}

	for _, element := range illegalElements {
		if !existingElements[element] {
			out = append(out, element)
		}
	}

	return out
}

func shuffleStringSlice(slice []string) {
	rand.Seed(time.Now().UnixNano())

	// Fisher-Yates shuffle algorithm
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func selectIDSubset(commonIDs []string, playlistIDs []string, length int) (idsOut []string, err error) {
	if length > len(playlistIDs) {
		return nil, errors.New("Request length is larger than playlist length")
	}
	idsOut = make([]string, length)

	minLength := min(length, len(commonIDs))

	shuffleStringSlice(commonIDs)
	copy(idsOut, commonIDs)

	uniquePlaylistIDs := addToSliceIfNotPresent(playlistIDs, commonIDs)
	shuffleStringSlice(uniquePlaylistIDs)

	copy(idsOut[minLength:], uniquePlaylistIDs)

	return idsOut, err
}

func min(num1, num2 int) (out int) {
	out = num1
	if num2 < out {
		out = num2
	}
	return out
}
