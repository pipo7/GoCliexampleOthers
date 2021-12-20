package main

// Directory name is perftest and it has package cmd
import (
	"perftest/cmd"
	"log" 
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
