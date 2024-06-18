package crawler

import "sync"

type Crawler struct {
	Threads     int
	Timeout     int
	MaxDepth    int
	MaxLinks    int
	SearchAny   []string
	SearchAll   []string
	IgnoredUrls []string
	Selectors   Selectors
	JsDepth     int
	// private fields
	pagesContent   map[string]string
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
}

func NewCrawler() *Crawler {
	return &Crawler{
		pagesContent: make(map[string]string),
		Threads:      1,
		Timeout:      10,
		MaxDepth:     2,
		MaxLinks:     0,
		JsDepth:      0,
		SearchAny:    []string{},
		IgnoredUrls:  []string{},
		Selectors: Selectors{
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
