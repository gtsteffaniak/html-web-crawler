package crawler

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

var collectionTypes = map[string]string{
	"images":  `([https?:]|\/)[^\s()'"]+\.(?:jpg|jpeg|png|gif|bmp|svg|webp|tiff)`,
	"video":   `([https?:]|\/)[^\s()'"]+\.(?:mp4|webm|ogg|flv|avi|mov|wmv|3gp)`,
	"audio":   `([https?:]|\/)[^\s()'"]+\.(?:mp3|wav|ogg|flac|wma|aac|alac|aiff)`,
	"pdf":     `([https?:]|\/)[^\s()'"]+\.(?:pdf)`,
	"doc":     `([https?:]|\/)[^\s()'"]+\.(?:doc|docx)`,
	"xls":     `([https?:]|\/)[^\s()'"]+\.(?:xls|xlsx)`,
	"ppt":     `([https?:]|\/)[^\s()'"]+\.(?:ppt|pptx)`,
	"archive": `([https?:]|\/)[^\s()'"]+\.(?:zip|rar|7z|tar|gz|bz2|tgz|tbz2|txz)`,
	"code":    `([https?:]|\/)[^\s()'"]+\.(?:py|rb|java|c|cpp|cs|go|swift|kt`,
	"shell":   `([https?:]|\/)[^\s()'"]+\.(?:sh|bat|ps1|bash`,
	"text":    `([https?:]|\/)[^\s()'"]+\.(?:txt|md|csv|log|toml|ini|cfg|conf|txt|text|rtf)`,
	"json":    `([https?:]|\/)[^\s()'"]+\.(?:json)`,
	"yaml":    `([https?:]|\/)[^\s()'"]+\.(?:yml|yaml)`,
	"font":    `([https?:]|\/)[^\s()'"]+\.(?:ttf|otf|woff|woff2|eot|svg)`,
}

// Crawl is the public method that initializes the recursive crawling.
func (c *Crawler) Collect(pageURL ...string) ([]string, error) {
	c.mode = "collect"
	c.wg = sync.WaitGroup{}
	for _, url := range c.Selectors.ExcludedUrls {
		c.pagesContent[url] = ""
	}
	for _, url := range pageURL {
		c.wg.Add(1) // Add to the wait group before starting the recursive crawl
		go func(url string) {
			defer c.wg.Done()
			err := c.recursiveCollect(url, 1)
			if err != nil {
				fmt.Printf("Error crawling %s: %v\n", url, err)
			}
		}(url)
	}
	c.wg.Wait() // Wait for all goroutines to finish
	return slices.Compact(c.collectedItems), nil
}

// recursiveCrawl is a private method that performs the recursive crawling, respecting MaxDepth.
func (c *Crawler) recursiveCollect(pageURL string, currentDepth int) error {
	useJavascript := c.JsDepth >= currentDepth
	if currentDepth > c.MaxDepth {
		return nil
	}
	c.mutex.Lock()
	if _, ok := c.pagesContent[pageURL]; ok {
		c.mutex.Unlock()
		return nil
	}
	if len(c.pagesContent) >= c.MaxLinks && c.MaxLinks != 0 {
		c.mutex.Unlock()
		return nil
	}
	// Update crawledData before recursive calls
	c.pagesContent[pageURL] = ""
	c.mutex.Unlock()
	htmlContent, err := c.FetchHTML(pageURL, useJavascript)

	if err != nil {
		return nil // return nil on page load error because the site could be down
	}
	if currentDepth > 0 && len(c.Selectors.ContentPatterns) > 0 {
		matchContentPattern := false
		for _, pattern := range c.Selectors.ContentPatterns {
			if strings.Contains(htmlContent, pattern) {
				matchContentPattern = true
			}
		}
		if !matchContentPattern {
			return nil
		}
	}
	c.mutex.Lock()
	c.pagesContent[pageURL] = ""
	c.mutex.Unlock()
	links, err := c.extractLinks(htmlContent)
	if err != nil {
		return err
	}
	items, err := c.extractItems(htmlContent, pageURL)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	c.collectedItems = append(c.collectedItems, items...)
	c.mutex.Unlock()
	// Limit the number of concurrent goroutines based on Threads
	semaphore := make(chan struct{}, c.Threads)
	for link, linkText := range links {
		c.mutex.Lock()
		_, ok := c.pagesContent[link]
		c.mutex.Unlock()
		if ok {
			continue
		}
		if !c.linkTextCheck(link, linkText) {
			continue
		}

		fullURL := toAbsoluteURL(pageURL, link)
		if c.validDomainCheck(fullURL) {
			c.wg.Add(1)
			semaphore <- struct{}{}
			go func(url string) {
				defer func() {
					<-semaphore // Release the slot
					c.wg.Done() // Decrement counter after goroutine finishes
				}()
				err := c.recursiveCollect(url, currentDepth+1)
				if err != nil {
					fmt.Printf("Error collecting %s: %v\n", url, err)
				}
			}(fullURL)
		}
	}

	return nil
}
