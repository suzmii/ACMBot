package usererr

import (
	"fmt"
	"strings"

	"github.com/suzmii/ACMBot/errorx"
)

func ErrUserNotFound(handles ...string) error {
	if len(handles) == 0 {
		return errorx.NewUserError("未找到用户")
	}
	return errorx.NewUserError(fmt.Sprintf("未找到用户%s", strings.Join(handles, ",")))
}

var ErrTooManyRequests = errorx.NewUserError("请求太多了，服务器要炸了")
var ErrCodeforcesAPIFailure = errorx.NewUserError("Codeforces炸了")
var ErrNoRatingRecords = errorx.NewUserError("该用户暂无rating记录")
var ErrNoRandomProblemMatched = errorx.NewUserError("没有找到符合条件的题目，换个筛选试试")
var ErrInvalidProblemID = errorx.NewUserError("题号格式不正确，示例：/problem 2209F")
