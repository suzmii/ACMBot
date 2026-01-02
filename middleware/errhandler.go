package middleware

import (
	"errors"

	"github.com/suzmii/ACMBot/errorx"
)

// 将错误转换为用户友好的信息
func HandleError(err error) string {
	var usererr errorx.UserError
	if ok := errors.As(err, &usererr); ok {
		return usererr.Hint
	}
	return "出了点问题...暂时没法处理你的请求"
}
