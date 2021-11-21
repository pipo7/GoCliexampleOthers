package main

import (
	"encoding/json"
	"io/ioutil"
)

func returnErr(e error) {
	if e != nil {
		panic(e)
	}
}

type video struct {
	Id          string
	Title       string
	Description string
	ImageURL    string
	URL         string
}

func getVideos() (videos []video) {
	fileBytes, err := ioutil.ReadFile("./videos.json")
	err = json.Unmarshal(fileBytes, &videos)
	returnErr(err)
	return videos
}

func saveVideos(videos []video) {
	videoBytes, err := json.Marshal(videos)
	returnErr(err)
	//Save back in bytes format
	returnErr(ioutil.WriteFile("./videos.json", videoBytes, 0644))
}
