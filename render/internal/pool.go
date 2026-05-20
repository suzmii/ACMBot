package internal

import (
	"github.com/playwright-community/playwright-go"
)

// PagePool 用作并发渲染信号量 + Page 工厂。
// 每次 Acquire 都从 BrowserContext 创建一个全新 Page，Release 时立即关闭，
// 避免 v8 heap、echarts 实例、detached DOM 等状态跨渲染累积。
type PagePool struct {
	ctx playwright.BrowserContext
	sem chan struct{}
}

func NewPagePool(ctx playwright.BrowserContext, size int) *PagePool {
	if size <= 0 || size > 64 {
		logger.Warnf("渲染池大小不合理(%d)，调整为8", size)
		size = 8
	}
	return &PagePool{
		ctx: ctx,
		sem: make(chan struct{}, size),
	}
}

func (p *PagePool) Acquire() (playwright.Page, error) {
	p.sem <- struct{}{}
	page, err := p.ctx.NewPage()
	if err != nil {
		<-p.sem
		return nil, err
	}
	return page, nil
}

// Release 必须只在 Acquire 成功后调用。
func (p *PagePool) Release(page playwright.Page) {
	if page != nil {
		_ = page.Close()
	}
	<-p.sem
}
