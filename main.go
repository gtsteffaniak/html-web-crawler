package main

import (
	"fmt"

	"github.com/gtsteffaniak/html-web-crawler/cmd"
)

func main() {
	crawledData, err := cmd.Execute()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Use a type switch to determine the type of crawledData
	switch data := crawledData.(type) {
	case []string:
		fmt.Println("Collect function returned data:")
		for _, item := range data {
			fmt.Println(item)
		}
	case map[string]string:
		// nothing
	default:
		fmt.Println("Unknown data type returned")
	}
}
