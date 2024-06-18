package playwright

import (
	"embed"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

//go:embed stealth.min.js
var staticFiles embed.FS

// BrowserOptions Playwright browser options
type BrowserOptions struct {
	ProxyServer     *string
	EnableStealthJs *bool
	Language        *string
	viewportWidth   *int
	viewportHeight  *int
}

type CloseBrowserHandler func() error

func CreateBrowser() (playwright.BrowserContext, CloseBrowserHandler, error) {
	// -----------------------------------------------------------
	// open browser
	// -----------------------------------------------------------
	return NewBrowserContextWithOptions(BrowserOptions{
		EnableStealthJs: playwright.Bool(true),
		viewportWidth:   playwright.Int(1920),
		viewportHeight:  playwright.Int(1080 - 35),
	})
}

func GetHtmlContent(pageURL string) (string, error) {
	browserCtx, browserClose, err := CreateBrowser()
	if err != nil {
		return "", fmt.Errorf("could not create virtual browser context: %w", err)
	}
	defer (func() {
		_ = browserClose()
	})()
	// create new page
	page, err := browserCtx.NewPage()
	if err != nil {
		return "", fmt.Errorf("could not create page: %w", err)
	}
	// open homepage, input keyword and search
	if _, err = page.Goto(pageURL, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateNetworkidle}); err != nil {
		return "", fmt.Errorf("could not goto: %w", err)
	}
	htmlString, err := page.Content()
	return htmlString, err
}

// NewBrowserContextWithOptions creates a new browser context with options
func NewBrowserContextWithOptions(opt BrowserOptions) (playwright.BrowserContext, CloseBrowserHandler, error) {

	// start playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("could not start playwright: %w", err)
	}

	// launch browser
	launchOptions := playwright.BrowserTypeLaunchOptions{}
	launchOptions.Headless = playwright.Bool(true)

	browser, err := pw.Chromium.Launch(launchOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("could not launch browser: %w", err)
	}
	// set locale
	pLocate := "en-US"

	// set viewport
	pViewport := &playwright.Size{
		Width:  1920,
		Height: 1080,
	}
	if opt.viewportWidth != nil && opt.viewportHeight != nil {
		pViewport = &playwright.Size{
			Width:  *opt.viewportWidth,
			Height: *opt.viewportHeight,
		}
	}
	// create context
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Locale:   &pLocate,
		Viewport: pViewport,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not create context: %w", err)
	}

	// add init stealth.min.js
	if opt.EnableStealthJs != nil && *opt.EnableStealthJs {
		stealthJs, err := staticFiles.ReadFile("stealth.min.js")
		if err != nil {
			return nil, nil, fmt.Errorf("could not read stealth.min.js: %w", err)
		}
		stealthJsScript := playwright.Script{
			Content: playwright.String(string(stealthJs)),
		}
		if err = context.AddInitScript(stealthJsScript); err != nil {
			return nil, nil, fmt.Errorf("could not add stealth.min.js: %w", err)
		}
	}

	// create close handler
	closeHandler := func() error {
		if err := browser.Close(); err != nil {
			return fmt.Errorf("could not close browser: %w", err)
		}
		if err := pw.Stop(); err != nil {
			return fmt.Errorf("could not stop Playwright: %w", err)
		}
		return nil
	}

	return context, closeHandler, nil

}
