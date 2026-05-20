package internal

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("render-internal")

// asInt 兼容 playwright-go 把 JS number 解码为 float64 / int / int64 的几种情况。
func asInt(v interface{}) (int, error) {
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	case int64:
		return int(n), nil
	default:
		return 0, fmt.Errorf("expected number, got %T", v)
	}
}

type Render struct {
	playwright *playwright.Playwright
	browser    playwright.Browser
	ctx        playwright.BrowserContext
	pool       *PagePool

	closeOnce sync.Once
	closeErr  error
}

func (r *Render) RenderWithAutoSize(ctx context.Context, content bytes.Buffer) ([]byte, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	page, err := r.pool.Acquire()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire page: %w", err)
	}
	defer r.pool.Release(page)

	// 把 ctx 的 deadline 投影成 playwright 调用的超时上限
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, context.DeadlineExceeded
		}
		page.SetDefaultTimeout(float64(remaining.Milliseconds()))
	}

	if err = page.SetContent(
		content.String(),
		playwright.PageSetContentOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		},
	); err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// 等字体、图片和布局真正稳定
	_, err = page.Evaluate(string(ResourceRenderWaitAssets))
	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
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

	size, ok := result.([]interface{})
	if !ok || len(size) != 2 {
		return nil, fmt.Errorf("unexpected size result from evaluate: %#v", result)
	}
	width, err := asInt(size[0])
	if err != nil {
		return nil, fmt.Errorf("unexpected width: %w", err)
	}
	height, err := asInt(size[1])
	if err != nil {
		return nil, fmt.Errorf("unexpected height: %w", err)
	}

	if err := page.SetViewportSize(width, height); err != nil {
		return nil, err
	}

	// 等 viewport resize 生效
	_, _ = page.Evaluate(`() => new Promise(r => requestAnimationFrame(r))`)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
		err = fmt.Errorf("failed to install playwright: %w", err)
		logger.Error(err)
		return nil, err
	}
	logger.Info("安装完咯/或者你已经装了，自动跳过了")

	// init browser
	logger.Info("初始化playwright")
	var err error
	r.playwright, err = playwright.Run()
	if err != nil {
		err = fmt.Errorf("failed to start playwright: %w", err)
		logger.Error(err)
		return nil, err
	}
	logger.Info("初始化浏览器")
	r.browser, err = r.playwright.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(cfg.Headless),
	})
	if err != nil {
		err = fmt.Errorf("failed to start browser: %w", err)
		logger.Error(err)
		return nil, err
	}

	// init pool
	logger.Info("初始化PlaywrightContextPool")
	playctx, err := r.browser.NewContext(playwright.BrowserNewContextOptions{
		DeviceScaleFactor: playwright.Float(2.0)})
	if err != nil {
		return nil, fmt.Errorf("failed to create playwright context: %w", err)
	}
	r.ctx = playctx
	r.pool = NewPagePool(playctx, cfg.PoolSize)

	// InitTemplates
	logger.Info("Initializing templates")
	for name_, content := range templateContents {
		name := fmt.Sprintf("%v", name_)
		tmpl, err := template.New(name).Parse(*content)
		if err != nil {
			err = fmt.Errorf("failed to load template %s: %w", name, err)
			logger.Error(err)
			return nil, err
		}
		templates[name_] = tmpl
	}

	return r, nil
}

// Close 按相反顺序释放 context / browser / playwright。
// 可重入：多次调用安全，后续调用返回首次的错误。
func (r *Render) Close() error {
	r.closeOnce.Do(func() {
		if r.ctx != nil {
			if err := r.ctx.Close(); err != nil && r.closeErr == nil {
				r.closeErr = err
			}
		}
		if r.browser != nil {
			if err := r.browser.Close(); err != nil && r.closeErr == nil {
				r.closeErr = err
			}
		}
		if r.playwright != nil {
			if err := r.playwright.Stop(); err != nil && r.closeErr == nil {
				r.closeErr = err
			}
		}
	})
	return r.closeErr
}
