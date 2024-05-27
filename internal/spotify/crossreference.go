package spotify

import (
	"errors"
	"math/rand"
)

// Finds elements that are found in both string array and returns them
func findCommonElements(slice1, slice2 []string) []string {
	elementsMap := make(map[string]bool)

	for _, elem := range slice1 {
		elementsMap[elem] = true
	}

	var commonElements []string
	for _, elem := range slice2 {
		if elementsMap[elem] {
			commonElements = append(commonElements, elem)
		}
	}

	return commonElements
}

// Returns a slice of strings that contains the elements of allElements that are not
// found in illegalElements
func addToSliceIfNotPresent(allElements, illegalElements []string) (out []string) {
	existingElements := make(map[string]bool)

	for _, element := range illegalElements {
		existingElements[element] = true
	}

	for _, element := range allElements {
		if !existingElements[element] {
			out = append(out, element)
		}
	}

	return out
}

// Fisher-Yates shuffle algorithm
func shuffleStringSlice(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// select a subset of ids from commonIDs and and playlistIDs for a given length. It ensures that the
// result contains a certain number of elements from commonIds and the remaining from playlistIDs with
// no duplicates.
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
