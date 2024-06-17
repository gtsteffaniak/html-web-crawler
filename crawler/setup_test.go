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
				Threads:      1,
				Timeout:      10,
				MaxDepth:     1,
				MaxLinks:     0,
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
