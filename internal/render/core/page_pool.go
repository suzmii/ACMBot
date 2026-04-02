package core

import "github.com/playwright-community/playwright-go"

var pagePool []playwright.Page // 用于存储已创建的页面实例

// GetPage 从 pagePool 中获取页面，复用现有页面，如果没有可用页面，则创建新页面
func GetPage() (playwright.Page, error) {
	if len(pagePool) > 0 {
		// 如果池中有可用页面，取出一个并返回
		page := pagePool[len(pagePool)-1]
		pagePool = pagePool[:len(pagePool)-1] // 从池中移除该页面

		err := page.SetContent("")
		if err != nil {
			return nil, err
		}
		return page, nil
	}

	// 如果池中没有可用页面，则创建一个新页面
	return _browser.NewPage(playwright.BrowserNewPageOptions{
		DeviceScaleFactor: &[]float64{2.0}[0],
	})
}

// ReleasePage 将页面归还到池中
func ReleasePage(page playwright.Page) {
	pagePool = append(pagePool, page)
}
