package internal

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/playwright-community/playwright-go"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("render-internal")

type Render struct {
	playwright *playwright.Playwright
	browser    playwright.Browser
	pool       *PagePool
}

func (r *Render) RenderWithAutoSize(ctx context.Context, content bytes.Buffer) ([]byte, error) {

	page := r.pool.Acquire()
	var err error

	if err = page.SetContent(
		content.String(),
		playwright.PageSetContentOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		},
	); err != nil {
		return nil, err
	}

	// 等字体、布局真正稳定
	_, err = page.Evaluate(`() => {
		return document.fonts.ready.then(() => {
			return new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))
		})
	}`)
	if err != nil {
		return nil, err
	}

	// 获取尺寸（返回数组，避免 map）
	result, err := page.Evaluate(`() => {
		const el = document.getElementById('background') || document.getElementById('main') || document.body;
		const rect = el.getBoundingClientRect();
		return [Math.ceil(rect.width), Math.ceil(rect.height)];
	}`)
	if err != nil {
		return nil, err
	}

	size := result.([]interface{})
	width := int(size[0].(int))
	height := int(size[1].(int))

	if err := page.SetViewportSize(width, height); err != nil {
		return nil, err
	}

	// 等 viewport resize 生效
	_, _ = page.Evaluate(`() => new Promise(r => requestAnimationFrame(r))`)

	return page.Screenshot(playwright.PageScreenshotOptions{
		Type: playwright.ScreenshotTypePng,
	})
}

func New(cfg subconfig.Render) (*Render, error) {
	r := &Render{}

	// init driver
	logger.Info("正在安装playwright, 请耐心等待......")
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	}); err != nil {
		err = fmt.Errorf("Failed to install playwright: %v", err)
		logger.Error(err)
		return nil, err
	}
	logger.Info("安装完咯/或者你已经装了，自动跳过了")

	// init browser
	logger.Info("初始化playwright")
	var err error
	r.playwright, err = playwright.Run()
	if err != nil {
		err = fmt.Errorf("Failed to start playwright: %v", err)
		logger.Error(err)
		return nil, err
	}
	logger.Info("初始化浏览器")
	r.browser, err = r.playwright.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(cfg.Headless),
	})
	if err != nil {
		err = fmt.Errorf("Failed to start browser: %v", err)
		logger.Error(err)
		return nil, err
	}

	// init pool
	logger.Info("初始化PlaywrightContextPool")
	playctx, err := r.browser.NewContext(playwright.BrowserNewContextOptions{
		DeviceScaleFactor: playwright.Float(2.0)})
	if err != nil {
		return nil, fmt.Errorf("failed to create playwright context")
	}
	r.pool = NewPagePool(playctx, cfg.PoolSize)

	// InitTemplates
	logger.Info("Initializing templates")
	for name_, content := range templateContents {
		name := fmt.Sprintf("%v", name_)
		tmpl, err := template.New(name).Parse(*content)
		if err != nil {
			err = fmt.Errorf("Failed to load template %s: %v", name, err)
			logger.Error(err)
			return nil, err
		}
		templates[name_] = tmpl
	}

	return r, nil
}
