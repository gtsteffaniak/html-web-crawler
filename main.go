package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gtsteffaniak/html-web-crawler/cmd"
)

func main() {
	crawledData, err := cmd.Execute()
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Use a type switch to determine the type of crawledData
	switch data := crawledData.(type) {
	case []string:
		fmt.Println("Collect function returned data:")
		for _, item := range data {
			fmt.Println(item)
		}
	case map[string]string:
		// Crawl returns map, but we don't print it by default
		// User can pipe to file if needed
	default:
		if data != nil {
			fmt.Println("Unknown data type returned")
		}
	}
}
