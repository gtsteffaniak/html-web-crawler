package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"

	"github.com/gtsteffaniak/html-web-crawler/browser"
)

// FetchHTML retrieves the HTML content of the given URL.
func (c *Crawler) FetchHTML(pageURL string, javascriptEnabled bool) (string, error) {
	switch c.mode {
	case "crawl":
		if !c.Silent {
			fmt.Println("fetching", pageURL)
		}
	case "collect":
		// nothing yet
	}
	if javascriptEnabled {
		html, err := browser.GetHtmlContent(pageURL)
		if err != nil {
			fmt.Println(err)
		}
		return html, err
	} else {
		return c.requestPage(pageURL)
	}
}

func (c *Crawler) requestPage(pageURL string) (string, error) {
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
	htmlString := string(bodyBytes)
	if len(c.SearchAny) > 0 {
		for _, s := range c.SearchAny {
			results := simpleSearch(s, htmlString, 30)
			if len(results) == 0 {
				continue
			}
			if !c.Silent {
				fmt.Println("\n=== found matches : ", pageURL)
				for _, r := range results {
					fmt.Println(r)
				}
			}
			if slices.Contains(c.Selectors.Collections, "html") {
				c.mutex.Lock()
				c.collectedItems = append(c.collectedItems, pageURL)
				c.mutex.Unlock()
			}
		}
	}
	return htmlString, nil
}

func (c *Crawler) containsSelectors(n *html.Node) bool {
	if len(c.Selectors.Ids) == 0 && len(c.Selectors.Classes) == 0 {
		return true
	}
	for _, targetId := range c.Selectors.Ids {
		if targetId == "" {
			continue
		}
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == targetId {
				return true
			}
		}
	}
	for _, targetClass := range c.Selectors.Classes {
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
func (c *Crawler) extractLinks(htmlContent string) (map[string]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	links := make(map[string]string)
	var f func(*html.Node)
	inTargetElement := false
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if c.containsSelectors(n) {
				inTargetElement = true
				defer func() { inTargetElement = false }() // reset to false after leaving the element
			}
			if inTargetElement && n.Data == "a" {
				var linkText string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						linkURL := attr.Val
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							if c.Type == html.TextNode {
								linkText += c.Data
							}
						}
						links[linkURL] = strings.TrimSpace(linkText)
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

// extractLinks extracts links within the specified element by id or class from the HTML content.
func (c *Crawler) extractItems(htmlContent, pageUrl string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	items := []string{}
	var f func(*html.Node)
	inTargetElement := false
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if c.containsSelectors(n) {
				inTargetElement = true
				defer func() { inTargetElement = false }() // reset to false after leaving the element
			}
			if inTargetElement {
				items = append(items, c.performSearch(n, pageUrl)...)
			}
		}
		if !inTargetElement {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)
	return items, nil
}

// nodeToString converts an html.Node to a string.
func nodeToString(n *html.Node) (string, error) {
	var buf bytes.Buffer
	err := html.Render(&buf, n)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// performSearch searches for items in the given html node.
func (c *Crawler) performSearch(n *html.Node, pageUrl string) []string {
	items := []string{}
	htmlString, err := nodeToString(n)
	if err != nil {
		fmt.Println("error converting node to string", err)
	}
	for _, re := range c.regexPatterns {
		foundItems := re.FindAllString(htmlString, -1)
		for _, url := range foundItems {
			if strings.HasPrefix(url, "http") {
				split := strings.Split(url, "https://")
				if len(split) > 1 {
					// Access the last element using index len(split)-1
					url = "https://" + split[len(split)-1]
				}
			} else {
				url = toAbsoluteURL(pageUrl, url)
			}
			if c.validDomainCheck(url) {
				items = append(items, url)
			}
		}
	}
	return items
}
