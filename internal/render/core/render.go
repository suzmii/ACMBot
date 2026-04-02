package core

import (
	"bytes"
	_ "embed"
	"errors"

	"github.com/playwright-community/playwright-go"
)

var (
	_playwright *playwright.Playwright
	_browser    playwright.Browser
)

// HtmlAutoSize 自动获取内容尺寸进行截图
func HtmlAutoSize(content bytes.Buffer) ([]byte, error) {

	page, err := GetPage()
	if err != nil {
		return nil, err
	}
	defer ReleasePage(page) // 使用完页面后，归还到池中

	if err = page.SetContent(content.String(), playwright.PageSetContentOptions{WaitUntil: playwright.WaitUntilStateNetworkidle}); err != nil {
		return nil, err
	}

	// 获取页面内容的实际尺寸
	result, err := page.Evaluate(`() => {
		// 获取主要内容容器
		const backgroundElement = document.getElementById('background');
		const mainElement = document.getElementById('main');
		
		let targetElement = backgroundElement || mainElement;
		
		if (targetElement) {
			// 使用getBoundingClientRect获取实际渲染尺寸
			const rect = targetElement.getBoundingClientRect();
			return { 
				width: Math.ceil(rect.width), 
				height: Math.ceil(rect.height) 
			};
		}
		
		// fallback: 使用传统方法
		const body = document.body;
		const html = document.documentElement;
		const height = Math.max(
			body.scrollHeight, body.offsetHeight,
			html.clientHeight, html.scrollHeight, html.offsetHeight
		);
		const width = Math.max(
			body.scrollWidth, body.offsetWidth,
			html.clientWidth, html.scrollWidth, html.offsetWidth
		);
		return { width, height };
	}`)
	if err != nil {
		return nil, err
	}

	// 类型断言获取尺寸
	if sizeMap, ok := result.(map[string]interface{}); ok {
		var width, height int

		switch w := sizeMap["width"].(type) {
		case float64:
			width = int(w)
		case int:
			width = w
		default:
			return nil, errors.New("渲染失败: 无法计算页面宽度") // 默认宽度
		}

		switch h := sizeMap["height"].(type) {
		case float64:
			height = int(h)
		case int:
			height = h
		default:
			return nil, errors.New("渲染失败: 无法计算页面高度") // 默认高度
		}

		// 设置页面视口为实际内容尺寸
		err = page.SetViewportSize(width, height)
		if err != nil {
			return nil, err
		}
	}

	data, err := page.Screenshot(playwright.PageScreenshotOptions{
		Type: playwright.ScreenshotTypePng,
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}
