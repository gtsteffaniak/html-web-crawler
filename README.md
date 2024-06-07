# HTML Web Crawler

A Golang library to crawl the web for links and information.

## About

This web crawler was initially conceived in Python -- a language I deemed suitable for these type of tasks. However, upon realizing I would need multithreaded processing to make it fast enough, I quickly realized it would more aptly benefit from go's native concurrency.

In stark contrast to the Python implementation, the Go counterpart -- even without leveraging concurrency -- astoundingly outperformed its predecessor. Processing the identical task in under 4 seconds, Go showcased an 8-fold acceleration over Python's 32-second execution time, while consuming considerably fewer resources.

The decision to opt for Python over Go is an interesting topic which I intend to delve into extensively on my blog. In the meantime, I have this as a library, ready to integrate it seamlessly into another of my projects.

# How to use

## CLI

First, install or download the program
```
go install github.com/gtsteffaniak/html-web-crawler@latest
```

Make sure your go bin is in your path. Then, run with the commands
```
html-web-crawler --urls https://apnews.com/
```

Use `--help` to see more options:

```
Usage: ./html-web-crawler [options] --urls <urls>
Options:
  -classes string
        Comma separated list of classes inside the html that links need to be inside to crawl (inclusive with ids)
  -contentPatterns string
        Comma separated list terms that must exist in page contents
  -domains string
        Comma separated list of exact match domains to crawl, e.g. 'ap.com,reuters.com'
  -help
        Show help message
  -ids string
        Comma separated list of ids inside the html that links need to be inside to crawl  (inclusive with classes)
  -ignoredUrls string
        Comma separated list of URLs to ignore
  -linkTextPatterns string
        Comma separated list of link text to crawl (inclusive with urlPatterns)
  -maxDepth int
        Maximum depth for pages to crawl, 1 will only return links from the given URLs (default 1)
  -maxLinks int
        Maximum number of links to crawl, 0 will crawl all links found.
  -threads int
        Number of concurrent urls to check when crawling (default 1)
  -timeout int
        Timeout in seconds for each HTTP request (default 10)
  -urlPatterns string
        Comma separated list of URL patterns to crawl (inclusive with linkTextPatterns)
  -urls string
        Comma separated URLs to crawl (required). ie "https://example.com,https://example2.com"

Note: "(inclusive with ___)" means program will crawl link if either property matches (OR condition)
```

## Include as a module in your go program

Note: you can also see [ai-earthquake-tracker](https://github.com/gtsteffaniak/ai-earthquake-tracker) as an example.

```
package main

import (
	"fmt"

	"github.com/gtsteffaniak/html-web-crawler/crawler"
)

func main() {
	Crawler := crawler.NewCrawler()
	// add crawling html selector classes
	Crawler.Selectors.Classes = []string{"PageList-items-item"}
	// Allow 50 consecutive pages to crawl at a time
	Crawler.Threads = 50
	// Crawl starting with a given url
	crawledData, _ := Crawler.Crawl("https://apnews.com/hub/earthquakes")
	// Print/Process the crawled data
	fmt.Println("Total: ", len(crawledData))
}
```
