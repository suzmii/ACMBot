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
	ctx        playwright.BrowserContext
	pool       *PagePool
}

func (r *Render) RenderWithAutoSize(ctx context.Context, content bytes.Buffer) ([]byte, error) {

	page, err := r.pool.Acquire()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire page: %w", err)
	}
	defer r.pool.Release(page)

	if err = page.SetContent(
		content.String(),
		playwright.PageSetContentOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		},
	); err != nil {
		return nil, err
	}

	// 等字体、图片和布局真正稳定
	_, err = page.Evaluate(string(ResourceRenderWaitAssets))
	if err != nil {
		return nil, err
	}

	// 获取尺寸（返回数组，避免 map）
	result, err := page.Evaluate(`() => {
		const el = document.getElementById('background') || document.getElementById('main') || document.body;
		const rect = el.getBoundingClientRect();
		const width = Math.max(
			Math.ceil(rect.width),
			el.scrollWidth,
			document.documentElement.scrollWidth,
			document.body ? document.body.scrollWidth : 0,
		);
		const height = Math.max(
			Math.ceil(rect.height),
			el.scrollHeight,
			document.documentElement.scrollHeight,
			document.body ? document.body.scrollHeight : 0,
		);
		return [width, height];
	}`)
	if err != nil {
		return nil, err
	}

	size := result.([]interface{})
	width := int(size[0].(float64))
	height := int(size[1].(float64))

	if err := page.SetViewportSize(width, height); err != nil {
		return nil, err
	}

	// 等 viewport resize 生效
	_, _ = page.Evaluate(`() => new Promise(r => requestAnimationFrame(r))`)

	target := page.Locator("#background")
	count, err := target.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		target = page.Locator("#main")
		count, err = target.Count()
		if err != nil {
			return nil, err
		}
	}
	if count == 0 {
		target = page.Locator("body")
	}

	return target.Screenshot(playwright.LocatorScreenshotOptions{
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
	r.ctx = playctx
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

// Close 按相反顺序释放 context / browser / playwright。
// 调用后 Render 不再可用。
func (r *Render) Close() error {
	var firstErr error
	if r.ctx != nil {
		if err := r.ctx.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if r.browser != nil {
		if err := r.browser.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if r.playwright != nil {
		if err := r.playwright.Stop(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
