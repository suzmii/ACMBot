package render

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"html/template"
)

var (
	_playwright *playwright.Playwright
	_browser    playwright.Browser
	pagePool    []playwright.Page // 用于存储已创建的页面实例

	templates = make(map[Template]*template.Template)
)

type Template int

const (
	TemplateCodeforcesProfileV1 Template = iota + 1
	TemplateCodeforcesProfileV2
	TemplateCodeforcesRatingRecords
	TemplateAtcoderProfile
	TemplateQQGroupRank
)

var (
	//go:embed templates/codeforces_profile_v1.gohtml
	TemplateContentCodeforcesProfileV1 string
	//go:embed templates/codeforces_profile_v2.gohtml
	TemplateContentCodeforcesProfileV2 string
	//go:embed templates/codeforces_rating_change.gohtml
	TemplateContentCodeforcesRatingRecords string
	//go:embed templates/atcoder_profile.gohtml
	TemplateContentAtcoderProfile string
	//go:embed templates/qq_group_rank.gohtml
	TemplateContentQQGroupRank string
)

var templateContents = map[Template]*string{
	TemplateCodeforcesProfileV1:     &TemplateContentCodeforcesProfileV1,
	TemplateCodeforcesProfileV2:     &TemplateContentCodeforcesProfileV2,
	TemplateCodeforcesRatingRecords: &TemplateContentCodeforcesRatingRecords,
	TemplateAtcoderProfile:          &TemplateContentAtcoderProfile,
	TemplateQQGroupRank:             &TemplateContentQQGroupRank,
}

type HtmlOptions struct {
	Path string
	HTML string
}

func Init() {
	initTemplates()
	initBrowser()
	InitDriver()
	initPool()
}

func initPool() {
	page, err := _browser.NewPage()
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

// GetTemplate returns a compiled template by name.
func GetTemplate(name Template) *template.Template {
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("Template %v not found", name)
	}

	return tmpl
}

// GetPage 从 pagePool 中获取页面，复用现有页面，如果没有可用页面，则创建新页面
func GetPage(opt playwright.BrowserNewPageOptions) (playwright.Page, error) {
	if len(pagePool) > 0 {
		// 如果池中有可用页面，取出一个并返回
		page := pagePool[len(pagePool)-1]
		pagePool = pagePool[:len(pagePool)-1] // 从池中移除该页面
		err := page.SetViewportSize(opt.Viewport.Width, opt.Viewport.Height)
		if err != nil {
			return nil, err
		}
		err = page.SetContent("")
		if err != nil {
			return nil, err
		}
		return page, nil
	}

	// 如果池中没有可用页面，则创建一个新页面
	return _browser.NewPage(opt)
}

// ReleasePage 将页面归还到池中
func ReleasePage(page playwright.Page) {
	pagePool = append(pagePool, page)
}

func Html(PageOpt *playwright.BrowserNewPageOptions, content bytes.Buffer) ([]byte, error) {
	// 获取页面，复用池中的页面或创建新页面
	page, err := GetPage(*PageOpt)
	if err != nil {
		return nil, err
	}
	defer ReleasePage(page) // 使用完页面后，归还到池中

	if err = page.SetContent(content.String(), playwright.PageSetContentOptions{WaitUntil: playwright.WaitUntilStateNetworkidle}); err != nil {
		return nil, err
	}

	data, err := page.Screenshot(playwright.PageScreenshotOptions{
		//FullPage: &[]bool{true}[0],
		Type: playwright.ScreenshotTypePng,
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}
