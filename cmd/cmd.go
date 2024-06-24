package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gtsteffaniak/html-web-crawler/crawler"
)

func generalUsage() {
	fmt.Printf(`usage: ./html-web-crawler <command> [options] --urls <urls>
  commands:
    collect  Collect data from URLs
    crawl    Crawl URLs and collect data
    install  Install chrome browser for javascript enabled scraping.
               Note: Consider instead to install via native package manager,
                     then set "CHROME_EXECUTABLE" in the environment
	` + "\n")
}

func commandHelp(flagset *flag.FlagSet) {
	fmt.Println("Options:")
	flagset.PrintDefaults()
	os.Exit(1)
}

func Execute() (interface{}, error) {
	if len(os.Args) < 2 {
		generalUsage()
		return nil, errors.New("no command provided")
	}
	var crawlCmd = flag.NewFlagSet(os.Args[1], flag.ExitOnError) // Flags specific to "crawl" command
	//var collectCmd = flag.NewFlagSet("collect", flag.ExitOnError) // Flags specific to "collect" command
	// general flags
	help := crawlCmd.Bool("help", false, "Show help message")
	threads := crawlCmd.Int("threads", 1, "Number of concurrent urls to check when crawling")
	timeout := crawlCmd.Int("timeout", 10, "Timeout in seconds for each HTTP request")
	maxDepth := crawlCmd.Int("max-depth", 2, "Maximum depth for pages to crawl, 1 will only return links from the given URLs")
	maxLinks := crawlCmd.Int("max-links", 0, "Will limit crawling to a number of pages given")
	urls := crawlCmd.String("urls", "",
		`Comma separated URLs to crawl (required).
example: --urls "https://example.com,https://example2.com"`)
	classes := crawlCmd.String("class-selectors", "",
		`Comma separated list of classes inside the html that links need to be inside to crawl
Note: combined using OR condition with ids
example: --classes "button,center_col"`)
	ids := crawlCmd.String("id-selectors", "", `Comma separated list of ids inside the html that links need to be inside to crawl.
Note: combined using OR condition with classes
example: --ids "main,content"`)
	domains := crawlCmd.String("domains", "", "Comma separated list of exact match domains to crawl, e.g. 'ap.com,reuters.com'")
	excludeDomains := crawlCmd.String("domains-excluded", "", "Comma separated list of exact match domains NOT to crawl, e.g. 'ap.com,reuters.com'")
	linkTextPatterns := crawlCmd.String("link-text-selectors", "", `Comma separated list of link text to crawl
Note: combined using OR condition with urlPatterns`)
	urlPatterns := crawlCmd.String("url-selectors", "", `Comma separated list of URL patterns to crawl
Note: combined using OR condition with linkTextPatterns`)
	contentPatterns := crawlCmd.String("content-selectors", "", "Comma separated list terms that must exist in page contents")
	excludedUrls := crawlCmd.String("exclude-urls", "", "Comma separated list of URLs to ignore")
	jsDepth := crawlCmd.Int("js-depth", 0, "Depth to use javascript rendering")
	searchAny := crawlCmd.String("search-any", "", "search for any pattern (with OR operator)")
	searchAll := crawlCmd.String("search-all", "", "search must contain all patterns (with AND operator)")
	filetypes := crawlCmd.String("filetypes", "",
		`Comma separated list of filetypes for collection.
Example: --filetypes "pdf,docx,doc"
Also supports by group name such as:
images, video, audio, pdf, doc, archive, code, shell, text, json, yaml, font`)
	command := os.Args[1]
	err := crawlCmd.Parse(os.Args[2:])
	if *help || err != nil {
		generalUsage()
		os.Exit(0)
	}
	if *threads <= 0 {
		fmt.Println("Error: threads must be a positive integer")
		generalUsage()
		os.Exit(1)
	}

	if *timeout <= 0 {
		fmt.Println("Error: timeout must be a positive integer")
		generalUsage()
		os.Exit(1)
	}

	if *maxDepth < 1 {
		fmt.Println("Error: maxDepth cannot be less than 1")
		generalUsage()
		os.Exit(1)
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
	if *excludedUrls != "" {
		c.Selectors.ExcludedUrls = strings.Split(*excludedUrls, ",")
	}
	if *searchAny != "" {
		c.SearchAny = strings.Split(*searchAny, ",")
	}
	if *searchAll != "" {
		c.SearchAll = strings.Split(*searchAll, ",")
	}
	if *filetypes != "" {
		c.Selectors.Collections = strings.Split(*filetypes, ",")
	}
	// Split the URLs by comma
	searchUrls := strings.Split(*urls, ",")
	switch command {
	case "install":
		return nil, nil
	case "collect":
		if *urls == "" {
			commandHelp(crawlCmd)
		}
		return c.Collect(searchUrls...)
	case "crawl":
		if *urls == "" {
			commandHelp(crawlCmd)
		}
		return c.Crawl(searchUrls...)
	default:
		generalUsage()
		return nil, errors.New("unknown command")
	}
}
