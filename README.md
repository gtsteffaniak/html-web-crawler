# html-web-crawler

Created as a Golang library, this web crawler was initially conceived in Python -- a language I deemed suitable for these type of tasks. However, upon realizing I would need multithreaded processing to make it fast enough, I quickly realized it would more aptly benefit from go's native concurrency.

In stark contrast to the Python implementation, the Go counterpart -- even without leveraging concurrency -- astoundingly outperformed its predecessor. Processing the identical task in under 4 seconds, Go showcased an 8-fold acceleration over Python's 32-second execution time, while consuming considerably fewer resources.

The decision to opt for Python over Go is an interesting topic which I intend to delve into extensively on my blog. In the meantime, I have this as a library, ready to integrate it seamlessly into another of my projects.

# How to use (go)

## cli 

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
Usage: crawler [options] <urls>
Options:
  -classes string
        Comma separated list of classes inside the html that links need to be inside to crawl
  -domains string
        Comma separated list of domains to crawl
  -help
        Show help message
  -ids string
        Comma separated list of ids inside the html that links need to be inside to crawl
  -maxDepth int
        Maximum depth for pages to crawl, 1 will only return links from the given URLs (default 1)
  -threads int
        Number of concurrent urls to check when crawling (default 1)
  -timeout int
        Timeout in seconds for each HTTP request (default 10)
  -urls string
        Comma separated URLs to crawl (required)
```

## include as a module in your go program

Note: you can also see [cmd.go](cmd/cmd.go) as an example.

First, add to imports
```
import (
    "github.com/gtsteffaniak/html-web-crawler/crawler"
)
```
Then initialize with whatever options you want
```
Crawler := NewCrawler()
Crawler.Selector.Classes = []string{"PageList-items-item"} # entire html document by default
Crawler.Threads = 50                          # single threaded by default

```

Lastly, run results which return `map[string]string` with data
```
crawledData, _ := Crawler.Crawl("https://apnews.com/hub/earthquakes") # can  have multipe urls
fmt.Println("Total: ", len(crawledData))
```