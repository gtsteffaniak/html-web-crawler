# html-web-crawler

Created as a Golang library, this web crawler was initially conceived in Python -- a language I deemed suitable for these type of tasks. However, upon realizing I would need multithreaded processing to make it fast enough, I quickly realized it would more aptly benefit from go's native concurrency.

In stark contrast to the Python implementation, the Go counterpart -- even without leveraging concurrency -- astoundingly outperformed its predecessor. Processing the identical task in under 4 seconds, Go showcased an 8-fold acceleration over Python's 32-second execution time, while consuming considerably fewer resources.

The decision to opt for Python over Go is an interesting topic which I intend to delve into extensively on my blog. In the meantime, I have this as a library, ready to integrate it seamlessly into another of my projects.

# How to use

```
Crawler := NewCrawler()
crawledData, _ := Crawler.Crawl("https://apnews.com/hub/earthquakes")
fmt.Println("Total: ", len(crawledData))
```

then set my options:
```
Crawler.SelectorClass = "PageList-items-item" # entire html document by default
Crawler.Threads = 50                          # single threaded by default
```

Lastly, run results which return `map[string]string` with data
```
crawledData, _ := Crawler.Crawl("https://apnews.com/hub/earthquakes")
fmt.Println("Total: ", len(crawledData))
```