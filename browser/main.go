package browser

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

var (
	chromeExec = os.Getenv("CHROME_EXECUTABLE")
)

func GetHtmlContent(pageURL string) (string, error) {
	b := launcher.NewBrowser()
	if b.Validate() != nil && chromeExec == "" {
		log.Fatal(`Attempted to use javascript engine, but no chromium browser was found.
		You can fix this two ways:
		1. installing chromium and set CHROME_EXECUTABLE path to chromium executable.
		2. running the "./html-web-cralwer install" command to automatically install.
`)
	}
	u := launcher.New().Bin(chromeExec).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(pageURL).Timeout(time.Second)
	content, err := page.HTML()
	if err != nil {
		fmt.Println(err)
	}

	return string(content), nil
}

func Install() error {
	p, err := launcher.NewBrowser().Get()
	utils.E(err)

	fmt.Println(p)
	return nil
}
