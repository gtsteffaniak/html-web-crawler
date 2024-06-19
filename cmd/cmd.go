package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gtsteffaniak/html-web-crawler/browser"
	"github.com/gtsteffaniak/html-web-crawler/crawler"
)

// flags for command-line arguments
var (
	threads          = flag.Int("threads", 1, "Number of concurrent urls to check when crawling")
	timeout          = flag.Int("timeout", 10, "Timeout in seconds for each HTTP request")
	maxDepth         = flag.Int("maxDepth", 1, "Maximum depth for pages to crawl, 1 will only return links from the given URLs")
	maxLinks         = flag.Int("maxLinks", 0, "Maximum number of links to crawl, 0 will crawl all links found.")
	searchAny        = flag.String("searchAny", "", "search string")
	urls             = flag.String("urls", "", "Comma separated URLs to crawl (required). ie \"https://example.com,https://example2.com\"")
	classes          = flag.String("classes", "", "Comma separated list of classes inside the html that links need to be inside to crawl (inclusive with ids)")
	ids              = flag.String("ids", "", "Comma separated list of ids inside the html that links need to be inside to crawl  (inclusive with classes)")
	domains          = flag.String("domains", "", "Comma separated list of exact match domains to crawl, e.g. 'ap.com,reuters.com'")
	excludeDomains   = flag.String("excludeDomains", "", "Comma separated list of exact match domains NOT to crawl, e.g. 'ap.com,reuters.com'")
	linkTextPatterns = flag.String("linkTextPatterns", "", "Comma separated list of link text to crawl (inclusive with urlPatterns)")
	urlPatterns      = flag.String("urlPatterns", "", "Comma separated list of URL patterns to crawl (inclusive with linkTextPatterns)")
	contentPatterns  = flag.String("contentPatterns", "", "Comma separated list terms that must exist in page contents")
	IgnoredUrls      = flag.String("ignoredUrls", "", "Comma separated list of URLs to ignore")
	jsDepth          = flag.Int("jsDepth", 0, "Depth to use javascript rendering")
	Images           = flag.Bool("images", false, "Include images in the search")
	filetypes        = flag.String("filetypes", "", "Comma separated list of filetypes for collection (e.g. 'pdf,docx,doc'), also supports by group name 'images','video','audio','pdf','doc','archive','code','shell','text','json','yaml','font'")
	help             = flag.Bool("help", false, "Show help message")
)

func usage() {
	fmt.Println("Usage: ./html-web-crawler <command> [options] --urls <urls>")
	fmt.Println("Commands:")
	fmt.Println("  collect  Collect data from URLs")
	fmt.Println("  crawl    Crawl URLs and collect data")
	fmt.Println("Options:")
	flag.PrintDefaults() // Print flag details with defaults
	fmt.Println("\nNote: \"(inclusive with ___)\" means program will crawl link if either property matches (OR condition)")
}

func Execute() (interface{}, error) {
	if len(os.Args) < 2 {
		usage()
		return nil, errors.New("no command provided")
	}

	command := os.Args[1]
	err := flag.CommandLine.Parse(os.Args[2:])
	if *help || err != nil {
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
		if command == "install" {
			browser.Install()
			os.Exit(0)
		} else {
			usage()
			os.Exit(1)
		}
	}

	// Create a new crawler instance with flag values
	c := crawler.NewCrawler()
	c.Threads = *threads
	c.Timeout = *timeout
	c.MaxDepth = *maxDepth
	c.MaxLinks = *maxLinks
	c.JsDepth = *jsDepth
	if *ids != "" {
		c.Selectors.Ids = strings.Split(*ids, ",")
	}
	if *classes != "" {
		c.Selectors.Classes = strings.Split(*classes, ",")
	}
	if *domains != "" {
		c.Selectors.Domains = strings.Split(*domains, ",")
	}
	if *excludeDomains != "" {
		c.Selectors.ExcludeDomains = strings.Split(*excludeDomains, ",")
	}
	if *linkTextPatterns != "" {
		c.Selectors.LinkTextPatterns = strings.Split(*linkTextPatterns, ",")
	}
	if *urlPatterns != "" {
		c.Selectors.UrlPatterns = strings.Split(*urlPatterns, ",")
	}
	if *contentPatterns != "" {
		c.Selectors.ContentPatterns = strings.Split(*contentPatterns, ",")
	}
	if *searchAny != "" {
		c.SearchAny = strings.Split(*searchAny, ",")
	}
	// Split the URLs by comma
	urls := strings.Split(*urls, ",")

	switch command {
	case "install":
		return nil, nil
	case "collect":
		if *filetypes != "" {
			c.Selectors.Collections = strings.Split(*filetypes, ",")
		}
		fmt.Printf("collecting %v...\n", c.Selectors.Collections)
		return c.Collect(urls...)
	case "crawl":
		fmt.Println("crawling...")
		return c.Crawl(urls...)
	default:
		usage()
		return nil, errors.New("unknown command")
	}
}
