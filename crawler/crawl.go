package crawler

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

// Crawl is the public method that initializes the recursive crawling.
func (c *Crawler) Crawl(pageURL ...string) (map[string]string, error) {
	c.mode = "crawl"
	c.wg = sync.WaitGroup{}
	c.errors = []error{} // Initialize errors slice
	for _, url := range c.Selectors.ExcludedUrls {
		c.pagesContent[url] = ""
	}
	for _, url := range pageURL {
		c.wg.Add(1) // Add to the wait group before starting the recursive crawl
		go func(url string) {
			defer c.wg.Done()
			err := c.recursiveCrawl(url, 1)
			if err != nil {
				c.mutex.Lock()
				c.errors = append(c.errors, err)
				c.mutex.Unlock()
				if !c.Silent {
					fmt.Printf("Error crawling %s: %v\n", url, err)
				}
			}
		}(url)
	}
	c.wg.Wait() // Wait for all goroutines to finish

	for url := range c.pagesContent {
		if slices.Contains(c.Selectors.ExcludedUrls, url) {
			delete(c.pagesContent, url)
		}
	}

	// Return the first error if any occurred (but still return the results)
	if len(c.errors) > 0 {
		return c.pagesContent, c.errors[0]
	}
	return c.pagesContent, nil
}

// recursiveCrawl is a private method that performs the recursive crawling, respecting MaxDepth.
func (c *Crawler) recursiveCrawl(pageURL string, currentDepth int) error {
	if currentDepth > c.MaxDepth {
		return nil
	}
	useJavascript := c.JsDepth >= currentDepth

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

	c.mutex.Lock()
	if len(c.SearchAny) == 0 {
		c.pagesContent[pageURL] = htmlContent
	} else {
		c.pagesContent[pageURL] = ""
	}
	c.mutex.Unlock()

	links, err := c.extractLinks(htmlContent)
	if err != nil {
		// HTML parsing errors are common with malformed HTML - log but continue
		if !c.Silent {
			fmt.Printf("Warning: failed to extract links from %s: %v\n", pageURL, err)
		}
		return nil // Continue with other pages
	}

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
				err := c.recursiveCrawl(url, currentDepth+1)
				if err != nil {
					c.mutex.Lock()
					c.errors = append(c.errors, err)
					c.mutex.Unlock()
					if !c.Silent {
						fmt.Printf("Error crawling %s: %v\n", url, err)
					}
				}
			}(fullURL)
		}
	}

	return nil
}
