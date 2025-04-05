package errs

import (
	"errors"
	"fmt"
)

type InternalError struct {
	text string
}

func (e InternalError) Error() string {
	return fmt.Sprintf("INTERNAL ERROR: %s", e.text)
}

func NewInternalError(message string) InternalError {
	return InternalError{
		text: message,
	}
}

type ErrHandleNotFound struct {
	Handle string
}

func (e ErrHandleNotFound) Error() string {
	return fmt.Sprintf("没有叫%s的用户哦，是不是打错了？", e.Handle)
}

var (
	ErrNoRatingChanges       = errors.New("没有找到任何Rating变化记录哦，可能分还没出来，总不可能你没打过比赛吧... ")
	ErrNoHandle              = errors.New("没有听到要查询谁哦")
	ErrGroupOnly             = errors.New("该功能必须要在群内使用哦")
	ErrImDedicated           = errors.New("仅需要一个参数哦，不要发无效信息啦~")
	ErrBadPlatform           = errors.New("暂不支持该平台哦")
	ErrOrganizationUnmatched = errors.New("绑定前请前往https://codeforces.com/settings/social,将Organization设置为`ACMBot`")
	ErrHandleHasBindByOthers = errors.New("该codeforces账号已被他人绑定了哦")
	ErrIllegalHandle         = errors.New("输入的用户名有非法字符呢，再说一遍吧")
	_                        = errors.New("本软件为开源软件，遵循GPLv2协议，如果你获取本软件的途径中支付了费用，那你可能是受骗了")
	_                        = errors.New("如果你是开发者，欢迎review我们的代码，并提出宝贵意见，如果你有什么建议和意见，也欢迎提Issue或PR告诉我们")
)
