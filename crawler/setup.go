package crawler

import (
	"regexp"
	"sync"
)

type Crawler struct {
	Threads   int
	Timeout   int
	MaxDepth  int
	MaxLinks  int
	SearchAny []string
	SearchAll []string
	Selectors Selectors
	JsDepth   int
	// private fields
	pagesContent   map[string]string
	regexPatterns  []regexp.Regexp
	collectedItems []string
	mutex          sync.Mutex
	wg             sync.WaitGroup
	mode           string
}

type Selectors struct {
	Collections      []string
	Classes          []string
	Ids              []string
	Domains          []string
	UrlPatterns      []string
	LinkTextPatterns []string
	ContentPatterns  []string
	ExcludeDomains   []string
	ExcludedUrls     []string
}

func NewCrawler() *Crawler {
	return &Crawler{
		pagesContent: make(map[string]string),
		Threads:      1,  // single threaded by default
		Timeout:      10, // 10 seconds
		MaxDepth:     2,  // default is provided urls and follow any links on that page
		MaxLinks:     0,  // unlimited
		JsDepth:      0,  // javascript disabled by default
		SearchAny:    []string{},
		SearchAll:    []string{},
		Selectors: Selectors{
			ExcludedUrls:     []string{},
			Collections:      []string{"images"},
			LinkTextPatterns: []string{},
			UrlPatterns:      []string{},
			ContentPatterns:  []string{},
			Classes:          []string{},
			Ids:              []string{},
			Domains:          []string{},
			ExcludeDomains:   []string{},
		},
	}
}
