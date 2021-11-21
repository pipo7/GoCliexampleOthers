package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCmd.Bool("all", false, "get all videos")
	getID := getCmd.String("id", "", "get a particular video")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addID := addCmd.String("id", "", "id of video")
	addTitle := addCmd.String("title", "", "title of video")
	addUrl := addCmd.String("url", "", "url of video")
	addDesc := addCmd.String("desc", "", "description of video")

	if len(os.Args) < 2 {
		fmt.Println("Use either get or add subcommand")
		os.Exit(1)
	} else {
		switch os.Args[1] {
		case "get":
			HandleGet(getCmd, getAll, getID)
		case "add":
			HandleAdd(addCmd, addID, addTitle, addUrl, addDesc)
		default:
			fmt.Println("Unable to understand the subcommand use --get or --add ")
		}
	}
}
func HandleGet(getCmd *flag.FlagSet, all *bool, id *string) {
	getCmd.Parse(os.Args[2:])
	if *all == false && *id == "" {
		getCmd.PrintDefaults() // prints defaults
		os.Exit(1)
	}
	if *all {
		videos := getVideos()
		for _, video := range videos {
			fmt.Printf(" %s %s %s %s \n", video.Id, video.Title, video.URL, video.Description)
		}
	}

	if *id != "" {
		videos := getVideos()
		for _, video := range videos {
			if *id == video.Id {
				fmt.Printf(" %s %s %s %s \n", video.Id, video.Title, video.URL, video.Description)
			}

		}
	}
}
func ValidateAdd(addCmd *flag.FlagSet, id *string, title *string, url *string, desc *string) {
	addCmd.Parse(os.Args[2:])
	if *id == "" || *title == "" || *url == "" || *desc == "" {
		fmt.Println("All fields are requried using add")
		fmt.Println("Values mising for id or title or url or desc")
		addCmd.PrintDefaults()
		os.Exit(1)
	}
}
func HandleAdd(addCmd *flag.FlagSet, id *string, title *string, url *string, desc *string) {
	ValidateAdd(addCmd, id, title, url, desc)
	video := video{
		Id:          *id,
		Title:       *title,
		URL:         *url,
		Description: *desc,
	}
	videos := getVideos()          // get all videos
	videos = append(videos, video) // append newly added video to videos struct
	saveVideos(videos)             // save all videos again
	fmt.Println("Values saved :", video, " Use --get to retirve this value")
}
