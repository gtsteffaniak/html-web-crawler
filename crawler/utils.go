package crawler

import (
	"net/url"
	"strings"
)

// toAbsoluteURL converts a relative URL to an absolute URL based on the base URL.
// URL parsing errors are logged but the function continues with best-effort conversion.
func toAbsoluteURL(base, link string) string {
	// Handle protocol-relative URLs (starting with //)
	if strings.HasPrefix(link, "//") {
		baseURL, err := url.Parse(base)
		if err != nil {
			// Invalid base URL - return link as-is (best effort)
			return link
		}
		return baseURL.Scheme + ":" + link
	}
	u, err := url.Parse(link)
	if err != nil {
		// Invalid link URL - return as-is (best effort)
		return link
	}
	if u.IsAbs() {
		return link
	}
	if strings.HasPrefix(link, "/") {
		baseURL, err := url.Parse(base)
		if err != nil {
			// Invalid base URL - try to construct from domain
			return "https://" + getDomain(base) + link
		}
		return baseURL.Scheme + "://" + baseURL.Host + link
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		// Invalid base URL - return base as fallback
		return base
	}
	resolved := baseURL.ResolveReference(u)
	return resolved.String()
}

// getDomain returns the domain of a URL.
// Returns empty string if URL parsing fails (invalid URL).
func getDomain(pageURL string) string {
	u, err := url.Parse(pageURL)
	if err != nil {
		// Invalid URL - return empty string (caller should handle)
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
	// Handle protocol-relative URLs by checking if it starts with // or has a scheme
	if !strings.HasPrefix(fullURL, "https://") && !strings.HasPrefix(fullURL, "http://") && !strings.HasPrefix(fullURL, "//") {
		return false
	}
	// Convert protocol-relative URLs to absolute for domain checking
	if strings.HasPrefix(fullURL, "//") {
		fullURL = "https:" + fullURL
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
