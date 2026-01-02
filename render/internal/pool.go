package internal

import (
	"context"
	"errors"
	"sync"

	"github.com/playwright-community/playwright-go"
)

var ErrAcquireTimeout = errors.New("context pool acquire timeout")

type ContextPool struct {
	browser playwright.Browser

	max    int
	idle   chan playwright.BrowserContext
	mu     sync.Mutex
	closed bool

	ctxOptions playwright.BrowserNewContextOptions
}

func (p *ContextPool) Acquire(ctx context.Context) (playwright.BrowserContext, error) {
	// 先尝试直接取空闲的
	select {
	case c := <-p.idle:
		return c, nil
	default:
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("context pool is closed")
	}

	// 当前创建数 = max - idle容量剩余
	current := p.max - cap(p.idle) + len(p.idle)
	if current < p.max {
		p.mu.Unlock()
		return p.browser.NewContext(p.ctxOptions)
	}
	p.mu.Unlock()

	// 否则阻塞等待
	select {
	case c := <-p.idle:
		return c, nil
	case <-ctx.Done():
		return nil, ErrAcquireTimeout
	}
}

func (p *ContextPool) Release(ctx playwright.BrowserContext) {
	if ctx == nil {
		return
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		ctx.Close()
		return
	}
	p.mu.Unlock()

	select {
	case p.idle <- ctx:
		// 正常归还
	default:
		// 理论上不应该发生，保险起见
		ctx.Close()
	}
}

func (p *ContextPool) Discard(ctx playwright.BrowserContext) {
	if ctx != nil {
		ctx.Close()
	}
}

func (p *ContextPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	close(p.idle)
	p.mu.Unlock()

	for ctx := range p.idle {
		ctx.Close()
	}
}

func NewContextPool(browser playwright.Browser, max int, options playwright.BrowserNewContextOptions) *ContextPool {
	if max <= 0 || max > 64 {
		logger.Warnf("渲染池大小不合理(%d)，调整为8", max)
		max = 8
	}
	return &ContextPool{
		browser:    browser,
		max:        max,
		idle:       make(chan playwright.BrowserContext, max),
		ctxOptions: options,
	}
}

type PagePool struct {
	ch chan playwright.Page
}

func NewPagePool(ctx playwright.BrowserContext, size int) *PagePool {
	ch := make(chan playwright.Page, size)
	for i := 0; i < size; i++ {
		page, err := ctx.NewPage()
		if err != nil {
			panic(err)
		}
		ch <- page
	}
	return &PagePool{ch: ch}
}
func (p *PagePool) Acquire() playwright.Page {
	return <-p.ch
}

func (p *PagePool) Release(page playwright.Page) {
	p.ch <- page
}
