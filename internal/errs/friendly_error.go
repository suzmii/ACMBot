package errs

import (
	"fmt"
)

type FriendlyError struct {
	text string
}

func (e FriendlyError) Error() string {
	return e.text
}

func NewFriendlyError(message string) FriendlyError {
	return FriendlyError{
		text: message,
	}
}

func ErrUserNotFound(name string) error {
	return FriendlyError{
		text: fmt.Sprintf("%s? 没听过诶", name),
	}
}

var (
	ErrNoRatingRecords         = NewFriendlyError("这位貌似没得rating记录")
	ErrNoUsername              = NewFriendlyError("你要查谁？")
	ErrGroupOnly               = NewFriendlyError("这指令得进群才好使")
	ErrBadPlatform             = NewFriendlyError("还没适配")
	ErrOrganizationUnmatched   = NewFriendlyError("绑定前请前往https://datasync.com/settings/social,将Organization设置为`ACMBot`")
	ErrUsernameHasBindByOthers = NewFriendlyError("这账号被别人绑了，如需换绑请先解绑")
	ErrIllegalUsername         = NewFriendlyError("有非法字符捏")
	_                          = NewFriendlyError("本软件为开源软件，遵循GPLv2协议，如果你获取本软件的途径中支付了费用，那你可能是受骗了")
	_                          = NewFriendlyError("如果你是开发者，欢迎review我们的代码，并提出宝贵意见，如果你有什么建议和意见，也欢迎提Issue或PR告诉我们")
)
