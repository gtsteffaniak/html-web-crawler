package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// FetchHTML retrieves the HTML content of the given URL.
func (c *Crawler) FetchHTML(pageURL string) (string, error) {
	switch {
	case len(c.SearchAny) > 0:
		fmt.Print(".")
	case c.mode == "collect":
		// do nothing
	default:
		fmt.Println("fetching", pageURL)
	}

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
			simpleSearch(s, htmlString, pageURL)
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
func (c *Crawler) extractItems(htmlContent string) ([]string, error) {
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
				htmlString, err := nodeToString(n)
				if err != nil {
					fmt.Println("error converting node to string", err)
				}
				for _, i := range c.Selectors.Collections {
					regex, exists := collectionTypes[i]
					if !exists {
						regex = fmt.Sprintf(`([https?:]|\/)[^\s'"]+\.(?:%v)`, i)
					}
					items = append(items, regexSearch(regex, htmlString)...)
				}
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
