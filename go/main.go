package main

import (
	"fmt"

	"github.com/gtsteffaniak/html-web-crawler/cmd"
)

func main() {
	crawledData, _ := cmd.Execute()
	fmt.Println("Total: ", len(crawledData))
}
