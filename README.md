# html-web-crawler

Simple html web crawler

created as a golang library

Before creating this I preferred to create it in python, considering it a more apt language for this kind of task.

However, within a few minutes I realized I wanted to do multithreaded processing to reduce the time. While I could do this in python with threading, I thought concurrency was go's strongsuit. So, I created another version in go.

I have both versions in this repo and have some interesting findings. Without even using concurrency, the go version completes the same task in a fraction of the time, with a fraction of the resources.

Go takes 0.5 seconds to process the same results as python takes 32 second, go showing an 64x speed improvement without even considering concurrency.

So, why bother with python? Well, a more in-depth comparison may go on my blog. For now I will create this library so I can use it in another project.

How to use? see the main.go main() function, you'll see I instantiate:

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