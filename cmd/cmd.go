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
	maxLinks = flag.Int("maxLinks", 0, "Maximum number of links to crawl, 0 will crawl all links found.")

	urls             = flag.String("urls", "", "Comma separated URLs to crawl (required). ie \"https://example.com,https://example2.com\"")
	classes          = flag.String("classes", "", "Comma separated list of classes inside the html that links need to be inside to crawl (inclusive with ids)")
	ids              = flag.String("ids", "", "Comma separated list of ids inside the html that links need to be inside to crawl  (inclusive with classes)")
	domains          = flag.String("domains", "", "Comma separated list of exact match domains to crawl, e.g. 'ap.com,reuters.com'")
	linkTextPatterns = flag.String("linkTextPatterns", "", "Comma separated list of link text to crawl (inclusive with urlPatterns)")
	urlPatterns      = flag.String("urlPatterns", "", "Comma separated list of URL patterns to crawl (inclusive with linkTextPatterns)")
	contentPatterns  = flag.String("contentPatterns", "", "Comma separated list terms that must exist in page contents")
	IgnoredUrls      = flag.String("ignoredUrls", "", "Comma separated list of URLs to ignore")
	help             = flag.Bool("help", false, "Show help message")
)

func usage() {
	fmt.Println("Usage: ./html-web-crawler [options] --urls <urls>")
	fmt.Println("Options:")
	flag.PrintDefaults() // Print flag details with defaults
	fmt.Println("\nNote: \"(inclusive with ___)\" means program will crawl link if either property matches (OR condition)")
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
	if *linkTextPatterns != "" {
		crawler.Selectors.LinkTextPatterns = strings.Split(*linkTextPatterns, ",")
	}
	if *urlPatterns != "" {
		crawler.Selectors.UrlPatterns = strings.Split(*urlPatterns, ",")
	}
	if *contentPatterns != "" {
		crawler.Selectors.ContentPatterns = strings.Split(*contentPatterns, ",")
	}
	if *maxLinks > 0 {
		crawler.MaxLinks = *maxLinks
	}
	// Split the URLs by comma
	urls := strings.Split(*urls, ",")
	return crawler.Crawl(urls...)
}
