package spotify

import (
	"math/rand"
	"errors"
)

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
