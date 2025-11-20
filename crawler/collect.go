package crawler

import (
	"fmt"
	"regexp"
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
	if err := c.compileCollections(); err != nil {
		return nil, fmt.Errorf("failed to compile collection patterns: %w", err)
	}
	c.wg = sync.WaitGroup{}
	c.errors = []error{} // Initialize errors slice
	// Initialize shared semaphore for concurrency control
	if c.Threads > 0 {
		c.semaphore = make(chan struct{}, c.Threads)
	} else {
		c.semaphore = make(chan struct{}, 1) // Default to 1 if not set
	}
	for _, url := range c.Selectors.ExcludedUrls {
		c.pagesContent[url] = ""
	}
	for _, url := range pageURL {
		url := url // Capture loop variable
		c.wg.Go(func() {
			err := c.recursiveCollect(url, 1)
			if err != nil {
				c.mutex.Lock()
				c.errors = append(c.errors, err)
				c.mutex.Unlock()
				if !c.Silent {
					fmt.Printf("Error crawling %s: %v\n", url, err)
				}
			}
		})
	}
	c.wg.Wait() // Wait for all goroutines to finish
	slices.Sort(c.collectedItems)
	// Return the first error if any occurred
	if len(c.errors) > 0 {
		return slices.Compact(c.collectedItems), c.errors[0]
	}
	return slices.Compact(c.collectedItems), nil
}

func (c *Crawler) compileCollections() error {
	for _, collectionType := range c.Selectors.Collections {
		pattern, exists := collectionTypes[collectionType]
		if !exists {
			pattern = fmt.Sprintf(`([https?:]|\/)[^\s()'"]+\.(?:%v)`, collectionType)
		}
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("error compiling regex pattern for collection type '%s': %w", collectionType, err)
		}
		c.regexPatterns = append(c.regexPatterns, *regex)
	}
	return nil
}

// recursiveCrawl is a private method that performs the recursive crawling, respecting MaxDepth.
func (c *Crawler) recursiveCollect(pageURL string, currentDepth int) error {
	useJavascript := c.JsDepth >= currentDepth
	if currentDepth > c.MaxDepth {
		return nil
	}
	// Use single lock to check and set atomically to prevent race conditions
	c.mutex.Lock()
	if _, ok := c.pagesContent[pageURL]; ok {
		c.mutex.Unlock()
		return nil
	}
	if c.MaxLinks > 0 && len(c.pagesContent) >= c.MaxLinks {
		c.mutex.Unlock()
		return nil
	}
	// Mark as processing before releasing lock
	c.pagesContent[pageURL] = ""
	c.mutex.Unlock()
	htmlContent, err := c.FetchHTML(pageURL, useJavascript)
	if err != nil {
		// Log transient HTTP errors but don't fail the entire crawl
		// These are expected when scraping (403, 404, network issues, etc.)
		if !c.Silent {
			fmt.Printf("Warning: failed to fetch %s: %v\n", pageURL, err)
		}
		return nil // Continue crawling other pages
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
	// No need to set pagesContent again - already set above
	links, err := c.extractLinks(htmlContent)
	if err != nil {
		// HTML parsing errors are common with malformed HTML - log but continue
		if !c.Silent {
			fmt.Printf("Warning: failed to extract links from %s: %v\n", pageURL, err)
		}
		return nil // Continue with other pages
	}
	items, err := c.extractItems(htmlContent, pageURL)
	if err != nil {
		// HTML parsing errors are common - log but continue
		if !c.Silent {
			fmt.Printf("Warning: failed to extract items from %s: %v\n", pageURL, err)
		}
		return nil // Continue with other pages
	}
	// Batch mutex operations for better performance
	c.mutex.Lock()
	c.collectedItems = append(c.collectedItems, items...)
	// If "html" is in Collections, also collect the page URL itself
	if slices.Contains(c.Selectors.Collections, "html") {
		c.collectedItems = append(c.collectedItems, pageURL)
	}
	c.mutex.Unlock()

	// Process links with shared semaphore for concurrency control
	for link, linkText := range links {
		// Check if already processed (with lock)
		c.mutex.Lock()
		_, alreadyProcessed := c.pagesContent[link]
		c.mutex.Unlock()
		if alreadyProcessed {
			continue
		}
		if !c.linkTextCheck(link, linkText) {
			continue
		}

		fullURL := toAbsoluteURL(pageURL, link)
		if !c.validDomainCheck(fullURL) {
			continue
		}

		// Acquire semaphore slot before starting goroutine
		c.semaphore <- struct{}{}
		urlToProcess := fullURL // Capture for goroutine
		c.wg.Go(func() {
			defer func() {
				<-c.semaphore // Release the slot
			}()
			err := c.recursiveCollect(urlToProcess, currentDepth+1)
			if err != nil {
				c.mutex.Lock()
				c.errors = append(c.errors, err)
				c.mutex.Unlock()
				if !c.Silent {
					fmt.Printf("Error collecting %s: %v\n", urlToProcess, err)
				}
			}
		})
	}
	return nil
}
