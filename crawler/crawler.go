package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Crawler struct {
	Threads   int
	Timeout   int
	MaxDepth  int
	Selectors Selectors
	// private fields
	pagesContent map[string]string
	mutex        sync.Mutex
	wg           sync.WaitGroup
}

type Selectors struct {
	Classes []string
	Ids     []string
	Domains []string
}

func NewCrawler() *Crawler {
	return &Crawler{
		pagesContent: make(map[string]string),
		Threads:      1,
		Timeout:      10,
		MaxDepth:     1,
		Selectors: Selectors{
			Classes: []string{},
			Ids:     []string{},
			Domains: []string{},
		},
	}
}

// FetchHTML retrieves the HTML content of the given URL.
func FetchHTML(pageURL string) (string, error) {
	fmt.Println("fetching", pageURL)
	resp, err := http.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// Crawl is the public method that initializes the recursive crawling.
func (c *Crawler) Crawl(pageURL ...string) (map[string]string, error) {
	for _, url := range pageURL {
		err := c.recursiveCrawl(url, 0)
		if err != nil {
			return nil, err
		}
	}
	c.wg.Wait() // Wait for all goroutines to finish
	return c.pagesContent, nil
}

// recursiveCrawl is a private method that performs the recursive crawling, respecting MaxDepth.
func (c *Crawler) recursiveCrawl(pageURL string, currentDepth int) error {
	if currentDepth > c.MaxDepth {
		return nil
	}

	c.mutex.Lock()
	if _, ok := c.pagesContent[pageURL]; ok {
		c.mutex.Unlock()
		return nil
	}
	// Update crawledData before recursive calls
	c.pagesContent[pageURL] = ""
	c.mutex.Unlock()

	htmlContent, err := FetchHTML(pageURL)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	c.pagesContent[pageURL] = htmlContent
	c.mutex.Unlock()

	links, err := extractLinks(htmlContent, c.Selectors.Classes, c.Selectors.Ids)
	if err != nil {
		return err
	}

	// Limit the number of concurrent goroutines based on Threads
	semaphore := make(chan struct{}, c.Threads)
	for _, link := range links {
		c.mutex.Lock()
		if _, ok := c.pagesContent[link]; ok {
			c.mutex.Unlock()
			continue
		}
		fullURL := toAbsoluteURL(pageURL, link)

		c.mutex.Unlock()
		validDomain := len(c.Selectors.Domains) == 0
		for _, domain := range c.Selectors.Domains {
			if getDomain(fullURL) == domain {
				validDomain = true
			}
		}
		if validDomain && strings.HasPrefix(fullURL, "https://") {
			if c.Threads > 1 {
				semaphore <- struct{}{}
				go func(url string) {
					defer func() {
						<-semaphore // Release the slot
						c.wg.Done() // Decrement counter after goroutine finishes
					}()
					c.wg.Add(1) // Increment counter before launching recursive call
					err := c.recursiveCrawl(url, currentDepth+1)
					if err != nil {
						fmt.Printf("Error crawling %s: %v\n", url, err)
					}
				}(fullURL)
			} else {
				err = c.recursiveCrawl(fullURL, currentDepth+1)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func containsSelectors(ids []string, classes []string, n *html.Node) bool {
	if len(ids) == 0 && len(classes) == 0 {
		return true
	}
	for _, targetId := range ids {
		if targetId == "" {
			continue
		}
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == targetId {
				return true
			}
		}
	}
	for _, targetClass := range classes {
		if targetClass == "" {
			continue
		}
		for _, attr := range n.Attr {
			if attr.Key == "class" && containsClass(attr.Val, targetClass) {
				return true
			}
		}
	}
	return false
}

// extractLinks extracts links within the specified element by id or class from the HTML content.
func extractLinks(htmlContent string, targetClasses, targetIDs []string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	var links []string
	var f func(*html.Node)
	inTargetElement := false

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if containsSelectors(targetIDs, targetClasses, n) {
				inTargetElement = true
				defer func() { inTargetElement = false }() // reset to false after leaving the element
			}
			if inTargetElement && n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return links, nil
}

// toAbsoluteURL converts a relative URL to an absolute URL based on the base URL.
func toAbsoluteURL(base, link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	if u.IsAbs() {
		return link
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return link
	}

	return baseURL.ResolveReference(u).String()
}

// getDomain returns the domain of a URL.
func getDomain(pageURL string) string {
	u, err := url.Parse(pageURL)
	if err != nil {
		return ""
	}
	return u.Host
}

// containsClass checks if a class attribute contains a specific class.
func containsClass(classAttr, targetClass string) bool {
	classes := strings.Fields(classAttr)
	for _, class := range classes {
		if class == targetClass {
			return true
		}
	}
	return false
}
