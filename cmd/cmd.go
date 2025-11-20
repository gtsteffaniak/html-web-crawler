package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gtsteffaniak/html-web-crawler/crawler"
	"github.com/gtsteffaniak/html-web-crawler/version"
)

// CLI defines the overall command structure
type CLI struct {
	Version bool       `short:"V" name:"version" help:"Show version information."`
	Crawl   CrawlCmd   `cmd:"" help:"Gather URLs that match search criteria, crawling recursively. Returns full HTML content."`
	Collect CollectCmd `cmd:"" help:"Intensive collection of specific items (images, search terms, etc). Does not return full HTML."`
	Install InstallCmd `cmd:"" help:"Install Chrome browser for JavaScript-enabled scraping."`
}

// GlobalFlags are flags shared across all commands
type GlobalFlags struct {
	URLs    []string `required:"" name:"urls" help:"URLs to crawl (comma-separated or multiple --urls flags)." short:"u"`
	Silent  bool     `name:"silent" help:"Hide all output; only use exit codes." short:"s"`
	Threads int      `name:"threads" help:"Number of concurrent URLs to check when crawling." default:"1"`
	Timeout int      `name:"timeout" help:"Timeout in seconds for each HTTP request." default:"10"`
}

// CrawlSettings control the crawling behavior
type CrawlSettings struct {
	MaxDepth int `name:"max-depth" help:"Maximum depth for pages to crawl. 1 = only links from given URLs." default:"2"`
	MaxLinks int `name:"max-links" help:"Limit crawling to a number of pages (0 = unlimited)." default:"0"`
	JsDepth  int `name:"js-depth" help:"Depth to use JavaScript rendering (requires Chrome)." default:"0"`
}

// Selectors control which links to follow and content to collect
type Selectors struct {
	ClassSelectors []string `name:"class-selectors" help:"HTML classes that links must be inside to crawl (OR condition with ids)." placeholder:"class1,class2"`
	IdSelectors    []string `name:"id-selectors" help:"HTML ids that links must be inside to crawl (OR condition with classes)." placeholder:"id1,id2"`
	Domains        []string `name:"domains" help:"Exact match domains to crawl." placeholder:"example.com,example2.com"`
	ExcludeDomains []string `name:"exclude-domains" help:"Exact match domains NOT to crawl." placeholder:"spam.com,ads.com"`
	LinkText       []string `name:"link-text" help:"Link text patterns to crawl (OR condition with URL patterns)." placeholder:"pattern1,pattern2"`
	URLPatterns    []string `name:"url-patterns" help:"URL patterns to crawl (OR condition with link text)." placeholder:"pattern1,pattern2"`
	Content        []string `name:"content" help:"Terms that must exist in page contents." placeholder:"term1,term2"`
	ExcludeURLs    []string `name:"exclude-urls" help:"URLs to ignore." placeholder:"url1,url2"`
}

// SearchOptions control content searching
type SearchOptions struct {
	SearchAny []string `name:"search-any" help:"Search for any pattern (OR operator)." placeholder:"term1,term2"`
	SearchAll []string `name:"search-all" help:"Search must contain all patterns (AND operator)." placeholder:"term1,term2"`
}

// CollectionOptions control what to collect
type CollectionOptions struct {
	FileTypes []string `name:"filetypes" help:"File types to collect (pdf, docx, doc, images, video, audio, etc)." placeholder:"images,pdf"`
}

// CrawlCmd crawls URLs and returns full HTML content
type CrawlCmd struct {
	GlobalFlags
	CrawlSettings
	Selectors
	SearchOptions
}

// CollectCmd collects specific items from URLs
type CollectCmd struct {
	GlobalFlags
	CrawlSettings
	Selectors
	SearchOptions
	CollectionOptions
}

// InstallCmd installs Chrome for JavaScript rendering
type InstallCmd struct{}

// Run executes the crawl command
func (c *CrawlCmd) Run(ctx *kong.Context) error {
	if !c.Silent {
		log.Printf("Starting crawl with %d thread(s)...", c.Threads)
	}

	crawler := c.buildCrawler()

	urls := c.expandURLs()
	result, err := crawler.Crawl(urls...)
	if err != nil {
		return fmt.Errorf("crawl failed: %w", err)
	}

	if !c.Silent {
		log.Printf("Crawled %d pages", len(result))
	}

	ctx.Bind(result)
	return nil
}

// Run executes the collect command
func (col *CollectCmd) Run(ctx *kong.Context) error {
	if !col.Silent {
		log.Printf("Starting collection with %d thread(s)...", col.Threads)
	}

	crawler := col.buildCrawler()

	urls := col.expandURLs()
	result, err := crawler.Collect(urls...)
	if err != nil {
		return fmt.Errorf("collection failed: %w", err)
	}

	if !col.Silent {
		log.Printf("Collected %d items", len(result))
	}

	ctx.Bind(result)
	return nil
}

// Run executes the install command
func (i *InstallCmd) Run(ctx *kong.Context) error {
	fmt.Println("Chrome installation:")
	fmt.Println("  Note: Consider installing via your native package manager,")
	fmt.Println("        then set 'CHROME_EXECUTABLE' in the environment")
	fmt.Println()
	fmt.Println("Example installations:")
	fmt.Println("  macOS:   brew install --cask google-chrome")
	fmt.Println("  Ubuntu:  sudo apt install google-chrome-stable")
	fmt.Println("  Fedora:  sudo dnf install google-chrome-stable")
	return nil
}

// buildCrawler creates a crawler instance from command flags
func (c *CrawlCmd) buildCrawler() *crawler.Crawler {
	cr := crawler.NewCrawler()
	cr.Threads = c.Threads
	cr.Timeout = c.Timeout
	cr.MaxDepth = c.MaxDepth
	cr.MaxLinks = c.MaxLinks
	cr.JsDepth = c.JsDepth
	cr.Silent = c.Silent
	cr.SearchAny = c.SearchAny
	cr.SearchAll = c.SearchAll

	cr.Selectors.Ids = c.IdSelectors
	cr.Selectors.Classes = c.ClassSelectors
	cr.Selectors.Domains = c.Domains
	cr.Selectors.ExcludeDomains = c.ExcludeDomains
	cr.Selectors.LinkTextPatterns = c.LinkText
	cr.Selectors.UrlPatterns = c.URLPatterns
	cr.Selectors.ContentPatterns = c.Content
	cr.Selectors.ExcludedUrls = c.ExcludeURLs

	return cr
}

// buildCrawler creates a crawler instance from command flags
func (col *CollectCmd) buildCrawler() *crawler.Crawler {
	cr := crawler.NewCrawler()
	cr.Threads = col.Threads
	cr.Timeout = col.Timeout
	cr.MaxDepth = col.MaxDepth
	cr.MaxLinks = col.MaxLinks
	cr.JsDepth = col.JsDepth
	cr.Silent = col.Silent
	cr.SearchAny = col.SearchAny
	cr.SearchAll = col.SearchAll

	cr.Selectors.Ids = col.IdSelectors
	cr.Selectors.Classes = col.ClassSelectors
	cr.Selectors.Domains = col.Domains
	cr.Selectors.ExcludeDomains = col.ExcludeDomains
	cr.Selectors.LinkTextPatterns = col.LinkText
	cr.Selectors.UrlPatterns = col.URLPatterns
	cr.Selectors.ContentPatterns = col.Content
	cr.Selectors.ExcludedUrls = col.ExcludeURLs
	cr.Selectors.Collections = col.FileTypes

	return cr
}

// expandURLs handles comma-separated URLs in addition to multiple --urls flags
func (c *GlobalFlags) expandURLs() []string {
	var urls []string
	for _, u := range c.URLs {
		// Split by comma to support both formats
		split := strings.Split(u, ",")
		for _, s := range split {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				urls = append(urls, trimmed)
			}
		}
	}
	return urls
}

// getVersion returns the version string for display
func getVersion() string {
	v := version.Version
	if v == "" {
		v = "dev"
	}
	if version.CommitSHA != "" {
		v += fmt.Sprintf(" (commit: %s)", version.CommitSHA)
	}
	return v
}

// Execute parses arguments and runs the appropriate command
func Execute() (interface{}, error) {
	// Check for version flag before parsing (Kong requires a command otherwise)
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-V" {
			fmt.Println(getVersion())
			return nil, nil
		}
	}

	var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("html-web-crawler"),
		kong.Description("A Golang library and CLI to crawl the web for links and information."),
		kong.UsageOnError(),
		kong.Vars{
			"version": getVersion(),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
	)

	// Handle version flag (in case it's used with a command)
	if cli.Version {
		fmt.Println(getVersion())
		return nil, nil
	}

	err := ctx.Run(ctx)
	if err != nil {
		return nil, err
	}

	// Extract bound data from context
	var result interface{}
	ctx.Bind(&result)

	return result, nil
}
