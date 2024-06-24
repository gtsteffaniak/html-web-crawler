package crawler

import (
	"fmt"
	"net/url"
	"strings"
)

// toAbsoluteURL converts a relative URL to an absolute URL based on the base URL.
func toAbsoluteURL(base, link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	if u.IsAbs() {
		return link
	}
	if strings.HasPrefix(link, "/") {
		base = "https://" + getDomain(base) + link
	}
	return base
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

func simpleSearch(s, text, url string) {
	if strings.Contains(text, s) {
		parts := strings.Split(text, s)
		firstPart := parts[0]
		secondPart := parts[1]

		// Get last 10 characters of the first element
		last10First := firstPart
		if len(firstPart) > 20 {
			last10First = firstPart[len(firstPart)-20:]
		}

		// Get first 10 characters of the second element
		first10Second := secondPart
		if len(secondPart) > 20 {
			first10Second = secondPart[:20]
		}
		// Print the results
		fmt.Printf("\n%s%s%s : %s\n", last10First, s, first10Second, url)
	}
}

func (c *Crawler) linkTextCheck(link, linkText string) bool {
	if len(c.Selectors.UrlPatterns) == 0 && len(c.Selectors.LinkTextPatterns) == 0 {
		return true
	}
	for _, pattern := range c.Selectors.UrlPatterns {
		if strings.Contains(link, pattern) {
			return true
		}
	}
	for _, pattern := range c.Selectors.LinkTextPatterns {
		if strings.Contains(linkText, pattern) {
			return true
		}
	}
	return false
}

func (c *Crawler) validDomainCheck(fullURL string) bool {
	if !(strings.HasPrefix(fullURL, "https://") || strings.HasPrefix(fullURL, "http://")) {
		return false
	}
	domain := getDomain(fullURL)
	if domain == "" {
		return false
	}
	for _, d := range c.Selectors.ExcludeDomains {
		if d == "" {
			continue
		}
		if strings.HasSuffix(domain, d) {
			return false
		}
	}
	if len(c.Selectors.Domains) == 0 {
		return true
	}
	for _, d := range c.Selectors.Domains {
		if d == "" {
			continue
		}
		if strings.HasSuffix(domain, d) {
			return true
		}
	}
	return false
}
