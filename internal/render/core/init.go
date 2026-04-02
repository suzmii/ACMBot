package core

import (
	"fmt"
	"html/template"

	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

func Init() {
	initTemplates()
	initBrowser()
	InitDriver()
	initPool()
}

func initPool() {
	page, err := _browser.NewPage(playwright.BrowserNewPageOptions{
		DeviceScaleFactor: &[]float64{2.0}[0],
	})
	if err != nil {
		panic(err)
	}
	pagePool = append(pagePool, page)
}

func InitDriver() {
	log.Info("正在安装playwright, 请耐心等待......")
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	}); err != nil {
		log.Fatalf("Failed to install playwright: %v", err)
	}
	log.Info("安装完咯/或者你已经装了，自动跳过了")
}

func initBrowser() {
	var err error
	_playwright, err = playwright.Run()
	if err != nil {
		log.Fatalf("Failed to start playwright: %v", err)
	}

	_browser, err = _playwright.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		//Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("Failed to launch Chromium: %v", err)
	}
}

func initTemplates() {
	log.Info("Initializing templates")
	for name_, content := range templateContents {
		name := fmt.Sprintf("%v", name_)
		tmpl, err := template.New(name).Parse(*content)
		if err != nil {
			log.Fatalf("Failed to load template %s: %v", name, err)
		}
		templates[name_] = tmpl
	}
}
