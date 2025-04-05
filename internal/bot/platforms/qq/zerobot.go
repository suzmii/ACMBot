package qq

import (
	"errors"
	"fmt"
	bot2 "github.com/suzmii/ACMBot/internal/bot"
	myMsg "github.com/suzmii/ACMBot/internal/bot/message"
	"github.com/suzmii/ACMBot/internal/context"
	"github.com/suzmii/ACMBot/pkg/errs"
	"strings"
	"time"

	"github.com/suzmii/ACMBot/conf"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

type qqContext struct {
	context.Context
	zCtx *zero.Ctx
}

type ctxOption func(*qqContext)

func withZeroCtx(zCtx *zero.Ctx) ctxOption {
	return func(ctx *qqContext) {
		ctx.zCtx = zCtx
	}
}

func newQQContext(opts ...ctxOption) *qqContext {
	res := &qqContext{
		Context: context.Context{
			Platform:  context.PlatformQQ,
			StepValue: nil,
		},
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (c *qqContext) GetCallerInfo() context.CallerInfo {
	result := context.CallerInfo{
		ID:       c.zCtx.Event.UserID,
		NickName: c.zCtx.Event.Sender.NickName,
	}

	gid := c.zCtx.Event.GroupID

	if gid != 0 {
		gInfo := c.zCtx.GetGroupInfo(gid, false)
		result.Group = context.GroupInfo{
			ID:          gid,
			Name:        gInfo.Name,
			MemberCount: gInfo.MemberCount,
		}
	}

	return result
}

func (c *qqContext) GetContextType() context.Platform {
	return context.PlatformQQ
}

func (c *qqContext) Send(msg myMsg.Message) {
	c.zCtx.Send(msg.ToZeroMessage())
}

func (c *qqContext) SendError(err error) {
	for _, user := range zeroCfg.SuperUsers {
		c.zCtx.SendPrivateMessage(user, err.Error())
	}
}

func (c *qqContext) Params() []string {
	argStr := c.zCtx.State["args"].(string)
	return strings.Fields(argStr)
}

var (
	zeroCfg zero.Config
)

// TODO: 把配置转移到bot层级
func init() {
	cfg := conf.GetConfig().Bot
	zeroCfg = zero.Config{
		NickName:      cfg.NickName,
		CommandPrefix: bot2.CommandPrefix,
		SuperUsers:    cfg.SuperUsers,
		Driver:        []zero.Driver{},
	}

	for _, cfg := range cfg.WS {
		zeroCfg.Driver = append(zeroCfg.Driver, driver.NewWebSocketClient(
			fmt.Sprintf("ws://%s:%d", cfg.Host, cfg.Port),
			cfg.Token))
	}

	zero.OnRequest(func(ctx *zero.Ctx) bool {
		return ctx.Event.RequestType == "group"
	}).Handle(func(ctx *zero.Ctx) {
		ctx.SetGroupAddRequest(ctx.Event.Flag, ctx.Event.SubType, true, "")
		for _, user := range cfg.SuperUsers {
			ctx.SendPrivateMessage(user, fmt.Sprintf("已自动同意加群邀请: %d", ctx.Event.GroupID))
		}
		ctx.SendPrivateMessage(ctx.Event.UserID, fmt.Sprintf("已自动同意加群邀请: %d", ctx.Event.GroupID))
		ctx.SendPrivateMessage(ctx.Event.UserID, "要去新地方了呢~\n如果大家不知道如何使用，可以用`#help`, `#菜单`呼出功能列表哦")
	})

	zero.OnRequest(func(ctx *zero.Ctx) bool {
		return ctx.Event.RequestType == "friend"
	}).Handle(func(ctx *zero.Ctx) {
		ctx.SetFriendAddRequest(ctx.Event.Flag, true, "")
		for _, user := range cfg.SuperUsers {
			ctx.SendPrivateMessage(user, fmt.Sprintf("已自动同意好友请求: %d", ctx.Event.UserID))
		}
		go func() {
			time.Sleep(5 * time.Second)
			ctx.SendPrivateMessage(ctx.Event.UserID, "很高兴认识你~ \n"+bot2.MenuText)
		}()
	})

	for _, command := range bot2.Commands {
		commands := command.Commands
		handler := command.Handler
		zeroHandler := func(ctx *zero.Ctx) {
			qCtx := newQQContext(withZeroCtx(ctx))
			c := &context.Context{
				ApiCaller: qCtx,
				Platform:  qCtx.Platform,
			}
			err := handler(c)
			if err == nil {
				return
			}
			qCtx.Send(myMsg.Text(err.Error()))
			var internalError errs.InternalError
			if errors.As(err, &internalError) {
				qCtx.SendError(err)
			}
		}

		for _, command := range commands {
			zero.OnCommand(command).Handle(zeroHandler)
		}
	}
	zero.Run(&zeroCfg)
}
