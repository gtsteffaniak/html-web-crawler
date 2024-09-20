package crawler

import (
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

// simpleSearch searches for match and returns partial text results list
func simpleSearch(s, text string, bufferLen int) []string {
	found := []string{}
	if strings.Contains(text, s) {
		index := strings.Index(text, s)
		if index != -1 {
			start := index - bufferLen
			if start < 0 {
				start = 0
			}
			end := index + len(s) + bufferLen
			if end > len(text) {
				end = len(text)
			}
			cleaned := strings.ReplaceAll(text[start:end], "\n", "")
			found = append(found, cleaned)
		}
	}
	return found
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
