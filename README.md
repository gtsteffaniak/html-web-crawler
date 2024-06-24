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
// Make sure your go bin is in your path if installing via go
go install github.com/gtsteffaniak/html-web-crawler@latest
```

```
html-web-crawler crawl --urls https://apnews.com/
```

Use `--help` to see more options:

```
usage: ./html-web-crawler <command> [options] --urls <urls>
  commands:
    collect  Collect data from URLs
    crawl    Crawl URLs and collect data
    install  Install chrome browser for javascript enabled scraping.
               Note: Consider instead to install via native package manager,
                     then set "CHROME_EXECUTABLE" in the environment
```

Available flags will very by command given.


## Example CMD commands and purpose

To get all links on a given page, but not crawl any further:
```
html-web-crawler crawl --urls https://apnews.com/ --max-depth 1
```
To query duck duck go search with javascript enabled:

Note: Javascript requires chrome be installed and `CHROME_EXECUTABLE` path set for js enabled searching.
```
$ ./html-web-crawler collect --urls https://duckduckgo.com/?t=h_&q=puppies&iax=images&ia=images \
--js-depth 1 --filetypes images
Collect function returned data:
https://duckduckgo.com/assets/icons/meta/DDG-iOS-icon_60x60.png
https://duckduckgo.com/assets/icons/meta/DDG-iOS-icon_76x76.png
https://duckduckgo.com/assets/icons/meta/DDG-iOS-icon_120x120.png
https://duckduckgo.com/assets/icons/meta/DDG-iOS-icon_152x152.png
https://duckduckgo.com/assets/icons/meta/DDG-icon_256x256.png
https://duckduckgo.com/i/a49fa21e.jpg
```

To collect pages that include text:
```

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
