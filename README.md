# HTML Web Crawler

A Golang library to crawl the web for links and information.

## About

This web crawler was initially conceived in Python -- a language I deemed suitable for these type of tasks. However, upon realizing I would need multithreaded processing to make it fast enough, I quickly realized it would more aptly benefit from go's native concurrency.

In stark contrast to the Python implementation, the Go counterpart -- even without leveraging concurrency -- astoundingly outperformed its predecessor. Processing the identical task in under 4 seconds, Go showcased an 8-fold acceleration over Python's 32-second execution time, while consuming considerably fewer resources.

The decision to opt for Python over Go is an interesting topic which I intend to delve into extensively on my blog. In the meantime, I have this as a library, ready to integrate it seamlessly into another of my projects.

# How to use

## CLI

First, install or download the program
```
// Make sure your go bin is in your path if installing via go
go install github.com/gtsteffaniak/html-web-crawler@latest
```

```
html-web-crawler crawl --urls https://apnews.com/
```

Use `--help` to see more options:

```
Usage: html-web-crawler <command>

A Golang library and CLI to crawl the web for links and information.

Commands:
  crawl      Gather URLs that match search criteria, crawling recursively.
             Returns full HTML content.
  collect    Intensive collection of specific items (images, search terms, etc).
             Does not return full HTML.
  install    Install Chrome browser for JavaScript-enabled scraping.

Run "html-web-crawler <command> --help" for more information on a command.
```

### Available Flags

Each command has organized flags grouped by purpose:

**Global Flags** (all commands):
- `--urls` / `-u`: URLs to crawl (required)
- `--silent` / `-s`: Hide all output
- `--threads`: Number of concurrent URLs (default: 1)
- `--timeout`: Request timeout in seconds (default: 10)

**Crawl Settings**:
- `--max-depth`: Maximum crawl depth (default: 2)
- `--max-links`: Limit number of pages (0 = unlimited)
- `--js-depth`: Depth for JavaScript rendering (default: 0)

**Selectors** (filter which links to follow):
- `--class-selectors`: HTML classes to target
- `--id-selectors`: HTML ids to target
- `--domains`: Allowed domains
- `--exclude-domains`: Blocked domains
- `--link-text`: Link text patterns
- `--url-patterns`: URL patterns
- `--content`: Required content terms
- `--exclude-urls`: URLs to ignore

**Search Options**:
- `--search-any`: OR search patterns
- `--search-all`: AND search patterns

**Collection Options** (collect command only):
- `--filetypes`: File types to collect (images, pdf, video, etc.)


## Example CMD commands and purpose

Get all links on a given page without crawling further:
```bash
html-web-crawler crawl --urls https://apnews.com/ --max-depth 1
```

Crawl with multiple threads for faster processing:
```bash
html-web-crawler crawl --urls https://apnews.com/ --threads 10 --max-depth 3
```

Query DuckDuckGo search with JavaScript enabled:

*Note: JavaScript requires Chrome to be installed and `CHROME_EXECUTABLE` environment variable set.*
```bash
html-web-crawler collect \
  --urls "https://duckduckgo.com/?t=h_&q=puppies&iax=images&ia=images" \
  --js-depth 1 \
  --filetypes images
```

Collect pages that include specific text:
```bash
html-web-crawler collect \
  --urls "https://gportal.link/blog" \
  --search-any "my search string"
```

Collect images from crawled pages:
```bash
html-web-crawler collect \
  --urls "https://gportal.link/blog" \
  --filetypes images
```

Crawl only specific domains:
```bash
html-web-crawler crawl \
  --urls https://apnews.com/hub/earthquakes \
  --domains apnews.com \
  --max-depth 3
```

Filter by HTML class selectors:
```bash
html-web-crawler crawl \
  --urls https://apnews.com/ \
  --class-selectors "PageList-items-item,article-list" \
  --threads 20
```

## Include as a module in your go program

Note: you can also see [ai-earthquake-tracker](https://github.com/gtsteffaniak/ai-earthquake-tracker) as an example.

```
package main

import (
	"fmt"

	"github.com/gtsteffaniak/html-web-crawler/crawler"
)

func main() {
	Crawler := crawler.NewCrawler()
	// add crawling html selector classes
	Crawler.Selectors.Classes = []string{"PageList-items-item"}
	// Allow 50 consecutive pages to crawl at a time
	Crawler.Threads = 50
	// Crawl starting with a given url
	crawledData, _ := Crawler.Crawl("https://apnews.com/hub/earthquakes")
	// Print/Process the crawled data
	fmt.Println("Total: ", len(crawledData))
}
```
