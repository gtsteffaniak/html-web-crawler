package crawler

import (
	"reflect"
	"testing"
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
