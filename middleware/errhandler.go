package middleware

import (
	"context"
	"errors"

	"github.com/suzmii/ACMBot/errorx"
)

// 将错误转换为用户友好的信息
func HandleError(err error) string {
	var usererr errorx.UserError
	if ok := errors.As(err, &usererr); ok {
		return usererr.Hint
	}

	if errors.Is(err, context.Canceled) {
		return "请求超时，再试一次？"
	}

	return "出了点问题...暂时没法处理你的请求"
}
