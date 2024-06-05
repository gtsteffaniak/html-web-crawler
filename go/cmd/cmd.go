package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gtsteffaniak/html-web-crawler/crawler"
)

// flags for command-line arguments
var (
	threads  = flag.Int("threads", 1, "Number of concurrent urls to check when crawling")
	timeout  = flag.Int("timeout", 10, "Timeout in seconds for each HTTP request")
	maxDepth = flag.Int("maxDepth", 1, "Maximum depth for pages to crawl, 1 will only return links from the given URLs")
	urls     = flag.String("urls", "", "Comma separated URLs to crawl (required)")
	classes  = flag.String("classes", "", "Comma separated list of classes inside the html that links need to be inside to crawl")
	ids      = flag.String("ids", "", "Comma separated list of ids inside the html that links need to be inside to crawl")
	domains  = flag.String("domains", "", "Comma separated list of domains to crawl")

	help = flag.Bool("help", false, "Show help message")
)

func usage() {
	fmt.Println("Usage: crawler [options] <urls>")
	fmt.Println("Options:")
	flag.PrintDefaults() // Print flag details with defaults
}

func Execute() (map[string]string, error) {
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	if *threads <= 0 {
		fmt.Println("Error: threads must be a positive integer")
		usage()
		os.Exit(1)
	}

	if *timeout <= 0 {
		fmt.Println("Error: timeout must be a positive integer")
		usage()
		os.Exit(1)
	}

	if *maxDepth < 1 {
		fmt.Println("Error: maxDepth cannot be less than 1")
		usage()
		os.Exit(1)
	}

	if *urls == "" {
		usage()
		os.Exit(1)
	}

	// Create a new crawler instance with flag values
	crawler := crawler.NewCrawler()
	crawler.Threads = *threads
	crawler.Timeout = *timeout
	crawler.MaxDepth = *maxDepth
	if *ids != "" {
		crawler.Selectors.Ids = strings.Split(*ids, ",")
	}
	if *classes != "" {
		crawler.Selectors.Classes = strings.Split(*classes, ",")
	}
	if *domains != "" {
		crawler.Selectors.Domains = strings.Split(*domains, ",")
	}
	// Split the URLs by comma
	urls := strings.Split(*urls, ",")
	return crawler.Crawl(urls...)
}
