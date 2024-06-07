package main

import (
	"fmt"

	"github.com/gtsteffaniak/html-web-crawler/cmd"
)

func main() {
	crawledData, err := cmd.Execute()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Total: ", len(crawledData))
}
