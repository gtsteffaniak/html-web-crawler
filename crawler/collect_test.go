package crawler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractItems(t *testing.T) {
	htmlTests := map[string]string{
		"firstHTML": `
		<body id="good">
			<div >
				<img src="/url/image.png" alt="example link">
			</div>
		</body>
		`,
		"secondHTML": `
		<body>
			<div class="tax">
				<p>check out the this cool image.jpg:</p>
				<img src="https://testing.com/image.jpg" alt="example link">
			</div>
		</body>
		`,
		"thirdHTML": `
		<body>
			<div>
				<p>check out the this cool image.jpg:</p>
				<img src="https://testing.com/image.jpg" alt="example link">
			</div>
			<p> if you go to https://image/cat.svg you will find another image</p>
		</body>
		`,
	}
	tests := []struct {
		name string
		s    *Selectors
		html map[string]string
		want map[string][]string
	}{
		{
			name: "Test class selector",
			s: &Selectors{
				Collections: []string{"images"},
				Classes:     []string{"tax"},
				Ids:         []string{},
			},
			html: htmlTests,
			want: map[string][]string{
				"firstHTML": {},
				"secondHTML": {
					"https://testing.com/image.jpg",
				},
				"thirdHTML": {},
			},
		},
		{
			name: "Test ids selector",
			s: &Selectors{
				Collections: []string{"images"},
				Ids:         []string{"good"},
				Classes:     []string{},
			},
			html: htmlTests,
			want: map[string][]string{
				"firstHTML": {
					"https://www.domain.com/url/image.png",
				},
				"secondHTML": {},
				"thirdHTML":  {},
			},
		},
		{
			name: "Test regex outside of element",
			s: &Selectors{
				Collections: []string{"images"},
				Ids:         []string{},
				Classes:     []string{},
			},
			html: htmlTests,
			want: map[string][]string{
				"firstHTML": {
					"https://www.domain.com/url/image.png",
				},
				"secondHTML": {
					"https://testing.com/image.jpg",
				},
				"thirdHTML": {
					"https://testing.com/image.jpg",
					"https://image/cat.svg",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCrawler()
			c.Selectors = *tt.s
			c.mode = "collect"
			c.compileCollections()
			for key, html := range tt.html {
				assert.Contains(t, tt.want, key)
				got, _ := c.extractItems(html, "https://www.domain.com")
				if !reflect.DeepEqual(got, tt.want[key]) {
					t.Errorf("\nmismatch for %v: \n > got %v,\n > want %v", key, got, tt.want[key])
				}
			}
		})
	}
}

func Benchmark_collectionSearch(b *testing.B) {
	// Pick a representative test case from tests
	testHtml := `
	<body id="good">
		<div >
			<img src="/url/image.png" alt="example link">
		</div>
	</body>`
	c := NewCrawler()
	c.Selectors = Selectors{
		Collections: []string{"images"},
		Classes:     []string{"tax"},
		Ids:         []string{},
	}
	c.mode = "collect"
	c.compileCollections()
	for i := 0; i < b.N; i++ {
		_, _ = c.extractItems(testHtml, "https://www.domain.com")
	}
}

func TestSingleSourceRunCollectHtml(t *testing.T) {
	// First test: default Collections=["html"] should collect page URLs
	c := NewCrawler()
	c.Threads = 1
	c.MaxDepth = 1
	c.MaxLinks = 3
	results, err := c.Collect("https://www.cnn.com/")
	fmt.Println(err)
	for _, result := range results {
		fmt.Println(result)
	}
	assert.Equal(t, nil, err)
	// With MaxLinks=3 and MaxDepth=1, we should get at least the starting URL plus some links
	assert.GreaterOrEqual(t, len(results), 1, "Should collect at least the starting page URL")

}
func TestSingleSourceRunCollectImages(t *testing.T) {
	// First test: default Collections=["html"] should collect page URLs
	c := NewCrawler()
	c.MaxDepth = 1
	c.Threads = 10
	c.MaxLinks = 3
	c.Selectors.Collections = []string{"images"}
	results, err := c.Collect("https://www.cnn.com/")
	assert.Equal(t, nil, err)
	assert.GreaterOrEqual(t, len(results), 5)
}
