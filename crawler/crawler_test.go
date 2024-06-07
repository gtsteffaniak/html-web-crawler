package crawler

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCrawler(t *testing.T) {
	tests := []struct {
		name string
		want *Crawler
	}{
		{
			name: "Test New Crawler",
			want: &Crawler{
				pagesContent: make(map[string]string),
				Threads:      1,
				Timeout:      10,
				MaxDepth:     1,
				IgnoredUrls:  []string{},
				Selectors: Selectors{
					LinkTextPatterns: []string{},
					UrlPatterns:      []string{},
					ContentPatterns:  []string{},
					Classes:          []string{},
					Ids:              []string{},
					Domains:          []string{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCrawler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlSelectors(t *testing.T) {

	commonTests := map[string]string{
		"https://example.com": "example link",
		"https://taxhelp.com": "click here",
		"https://another.com": "tax help",
	}

	crawler := NewCrawler()

	tests := []struct {
		name  string
		s     *Selectors
		links map[string]string
		want  map[string]bool
	}{
		{
			name:  "Test without selectors",
			links: commonTests,
			s:     &crawler.Selectors,
			want: map[string]bool{
				"https://example.com": true,
				"https://taxhelp.com": true,
				"https://another.com": true,
			},
		},
		{
			name: "Single Test Link Text Patterns",
			s: &Selectors{
				LinkTextPatterns: []string{"tax"},
			},
			links: commonTests,
			want: map[string]bool{
				"https://example.com": false,
				"https://taxhelp.com": false,
				"https://another.com": true,
			},
		},
		{
			name: "Single Link URL Patterns",
			s: &Selectors{
				UrlPatterns: []string{"tax"},
			},
			links: commonTests,
			want: map[string]bool{
				"https://example.com": false,
				"https://taxhelp.com": true,
				"https://another.com": false,
			},
		},
		{
			name: "Single inclusive URL/Text Patterns",
			s: &Selectors{
				UrlPatterns:      []string{"tax"},
				LinkTextPatterns: []string{"tax"},
			},
			links: commonTests,
			want: map[string]bool{
				"https://example.com": false,
				"https://taxhelp.com": true,
				"https://another.com": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawler := NewCrawler()
			crawler.Selectors = *tt.s
			for link, text := range tt.links {
				assert.Contains(t, tt.want, link)
				got := crawler.linkTextCheck(link, text)
				if got != tt.want[link] {
					t.Errorf("mismatch for %v ", link)
				}
			}
		})
	}
}

func TestDomainSelectors(t *testing.T) {

	tests := []struct {
		name  string
		s     *Selectors
		links []string
		want  map[string]bool
	}{
		{
			name: "Check without domain selectors",
			s:    &Selectors{},
			links: []string{
				"",
				"/test",
				"https://",
				"https://example.com",
				"http://wifi.com",
			},
			want: map[string]bool{
				"":                    false,
				"/test":               false,
				"https://":            false,
				"https://example.com": true,
				"http://wifi.com":     true,
			},
		},
		{
			name: "Check with one domain selector",
			s: &Selectors{
				Domains: []string{"example.com"},
			},
			links: []string{
				"",
				"/test",
				"https://",
				"https://example.com",
				"http://wifi.com",
			},
			want: map[string]bool{
				"":                    false,
				"/test":               false,
				"https://":            false,
				"https://example.com": true,
				"http://wifi.com":     false,
			},
		},
		{
			name: "Check with multiple domain selectors",
			s: &Selectors{
				Domains: []string{"example.com", "wifi.com"},
			},
			links: []string{
				"",
				"/test",
				"https://",
				"https://example.com",
				"http://wifi.com",
			},
			want: map[string]bool{
				"":                    false,
				"/test":               false,
				"https://":            false,
				"https://example.com": true,
				"http://wifi.com":     true,
			},
		},
		{
			name: "Check with invalid domain selector",
			s: &Selectors{
				Domains: []string{"/test", ""},
			},
			links: []string{
				"",
				"/test",
				"https://",
				"https://example.com",
				"http://wifi.com",
			},
			want: map[string]bool{
				"":                    false,
				"/test":               false,
				"https://":            false,
				"https://example.com": false,
				"http://wifi.com":     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawler := NewCrawler()
			crawler.Selectors = *tt.s
			for _, link := range tt.links {
				assert.Contains(t, tt.want, link)
				got := crawler.validDomainCheck(link)
				if got != tt.want[link] {
					t.Errorf("mismatch for %v ", link)
				}
			}
		})
	}
}

func TestClassAndIdsSelectors(t *testing.T) {
	htmlTests := map[string]string{
		"firstHTML": `
		<body>
			<div>
				<a href="https://example.com" class="tax">example link</a>
				<a href="/relative" class="tax">example relative path</a>
				<a href="#" >example hash path</a>
				<a href="" >example link</a>
			</div>
		</body>
		`,
		"secondHTML": `
		<body>
			<div class="tax">
				<a href="https://testing.com">example link</a>
			</div>
		</body>
		`,
		"thirdHTML": `
		<body>
			<div id="tax">
				<a href="https://testing.com">example link</a>
				<a href="https://good.com" class="good">example good link</a>
			</div>
		</body>
		`,
		"fourthHTML": `
		<body>
			<div id="good">
				<a href="https://testing.com" id="tax">example link</a>
				<a href="https://goodwifi.com" class="good">click for wifi</a>
			</div>
		</body>
		`,
	}
	tests := []struct {
		name string
		s    *Selectors
		html map[string]string
		want map[string]map[string]string
	}{
		{
			name: "Test class selector",
			s: &Selectors{
				Classes: []string{"tax"},
			},
			html: htmlTests,
			want: map[string]map[string]string{
				"firstHTML": {
					"https://example.com": "example link",
					"/relative":           "example relative path",
				},
				"secondHTML": {
					"https://testing.com": "example link",
				},
				"thirdHTML":  {},
				"fourthHTML": {},
			},
		},
		{
			name: "Test ids selector",
			s: &Selectors{
				Ids: []string{"good"},
			},
			html: htmlTests,
			want: map[string]map[string]string{
				"firstHTML":  {},
				"secondHTML": {},
				"thirdHTML":  {},
				"fourthHTML": {
					"https://testing.com":  "example link",
					"https://goodwifi.com": "click for wifi",
				},
			},
		},
		{
			name: "Test class and id selectors",
			s: &Selectors{
				Classes: []string{"tax"},
				Ids:     []string{"good"},
			},
			html: htmlTests,
			want: map[string]map[string]string{
				"firstHTML": {
					"https://example.com": "example link",
					"/relative":           "example relative path",
				},
				"secondHTML": {
					"https://testing.com": "example link",
				},
				"thirdHTML": {},
				"fourthHTML": {
					"https://testing.com":  "example link",
					"https://goodwifi.com": "click for wifi",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCrawler()
			c.Selectors = *tt.s
			for key, html := range tt.html {
				assert.Contains(t, tt.want, key)
				got, _ := c.extractLinks(html)
				if !reflect.DeepEqual(got, tt.want[key]) {
					t.Errorf("\nmismatch for %v: \n > got %v,\n > want %v", key, got, tt.want[key])
				}
			}
		})
	}
}

func TestSingleSourceRun(t *testing.T) {
	c := NewCrawler()
	c.Threads = 10
	results, err := c.Crawl("https://www.gportal.link/blog/")
	assert.Equal(t, nil, err)
	assert.Greater(t, 3, len(results))
}

func TestMultipleSourceRun(t *testing.T) {
	c := NewCrawler()
	results, err := c.Crawl("https://www.apnews.com/")
	if err != nil {
		t.Errorf("Error running crawler: %v", err)
	}
	assert.Equal(t, nil, err)
	assert.Greater(t, 3, len(results))
}
