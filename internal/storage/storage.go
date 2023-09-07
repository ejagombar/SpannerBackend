
// package spotify
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )
//
// // type Artist struct {
// // 	ID     string   `json:"id"`
// // 	Name   string   `json:"name"`
// // 	Genres []string `json:"genres"`
// // }
// //
// // type Track struct {
// // 	ID      string   `json:"id"`
// // 	Name    string   `json:"name"`
// // 	Artists []Artist `json:"artists"`
// // }
//
// type PlaylistData struct {
// 	ID          string   `json:"id"`
// 	Name        string   `json:"name"`
// 	Description string   `json:"description"`
// 	ImageLink   string   `json:"imagelink"`
// 	TrackCount  int      `json:"trackcount"`
// 	TrackIDs    []string `json:"trackids"`
// }
//
// func SaveStruct(filename string, data any) (err error) {
// 	marshalData, err := json.MarshalIndent(data, "", "\t")
// 	if err != nil {
// 		return fmt.Errorf("error marshaling AuthStore to JSON: %v", err)
// 	}
//
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("error creating file: %v", err)
// 	}
// 	defer file.Close()
//
// 	_, err = file.Write(marshalData)
// 	if err != nil {
// 		return fmt.Errorf("error writing to file: %v", err)
// 	}
//
// 	return nil
// }
//
// func LoadStruct(filename string, data any) (err error) {
// 	marshalData, err := os.ReadFile(filename)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = json.Unmarshal(marshalData, &data)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

	// fileName := fmt.Sprintf("%v.json", playlistID)
	// err = LoadStruct(fileName, &playlistData)
	// if err != nil {
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		err = requestAndSavePlaylist(client, fileName, &playlistData)
	// 	}
	// 	log.Fatal(err)
	// }
	//
	// topTracks := make([]string, 150)
	// fileName = fmt.Sprintf("%v.json", "userTopTracks")
	// err = LoadStruct(fileName, &topTracks)
	// if err != nil {
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		err = requestAndSaveTopTracks(client, topTracks)
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
