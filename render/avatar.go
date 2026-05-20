package render

import (
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// avatarHTTPClient 用于把头像 URL 拉成 data URI 之后再嵌入模板，
// 避免浏览器在渲染时还要走一次网络。
var avatarHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

const avatarMaxBytes = 2 << 20 // 2 MiB，正常头像不会更大

// resolveAvatar 拉取 url 并返回可嵌入 URL 上下文的 data URI。
// 空字符串、已是 data URI、拉取失败都会原样返回，让模板退化到原有行为
// （浏览器自己请求或不显示），不阻断渲染。
//
// 返回 template.URL 而非 string 是必要的：html/template 会把
// data: 开头的字符串当作不安全 scheme 替换成 "#ZgotmplZ"，
// 用 template.URL 显式声明可信。
func resolveAvatar(ctx context.Context, url string) template.URL {
	if url == "" || strings.HasPrefix(url, "data:") {
		return template.URL(url)
	}
	dataURI, err := fetchAsDataURI(ctx, url)
	if err != nil {
		logger.Warnf("inline avatar %s 失败，回退到 URL: %v", url, err)
		return template.URL(url)
	}
	return template.URL(dataURI)
}

// resolveAvatars 并行拉取多个头像，失败的位置回退到原 URL。
func resolveAvatars(ctx context.Context, urls []string) []template.URL {
	out := make([]template.URL, len(urls))
	var wg sync.WaitGroup
	for i, u := range urls {
		wg.Add(1)
		go func(i int, u string) {
			defer wg.Done()
			out[i] = resolveAvatar(ctx, u)
		}(i, u)
	}
	wg.Wait()
	return out
}

func fetchAsDataURI(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := avatarHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, avatarMaxBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > avatarMaxBytes {
		return "", fmt.Errorf("body exceeds %d bytes", avatarMaxBytes)
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" || strings.HasPrefix(ct, "application/octet-stream") {
		ct = http.DetectContentType(data)
	}
	return "data:" + ct + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

