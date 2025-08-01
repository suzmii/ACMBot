package qq

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/internal/errs"
	"github.com/suzmii/ACMBot/internal/model"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	zMsg "github.com/wdvxdr1123/ZeroBot/message"
	"strings"
	"time"
)

func NewZeroBotAdapter() *ZeroBotAdapter {
	cfg := config.LoadConfig().ZeroBot
	wsHost := cfg.Host
	wsPort := cfg.Port
	wsToken := cfg.Token

	if wsHost == "" || wsPort == 0 {
		logrus.Fatal("未设置websocket端口与主机")
	}

	var zeroCfg = zero.Config{
		NickName:      []string{"ACMBot"},
		CommandPrefix: cfg.CommandPrefix,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(
				fmt.Sprintf("ws://%s:%d", wsHost, wsPort),
				wsToken,
			),
		},
		RingLen: 0,
	}

	zeroCfg.Driver = append(zeroCfg.Driver)

	return &ZeroBotAdapter{
		zeroCfg: zeroCfg,
		menuText: `以下是功能列表：所有命令都要加上前缀` + "`" + cfg.CommandPrefix + "`" + `哦🥰
0. help(或菜单)，输出本消息

1. cf/at [username]，用于查询codeforces/atcoder用户的基本信息

2. rating(或rt) [username]，用于查询codeforces用户的rating变化曲线

3. 近期[比赛,atc,nk,lg,cf]，用于查询近期的比赛数据，数据来源于clist.by`,
	}
}

type ZeroBotAdapter struct {
	zeroCfg  zero.Config
	menuText string
}

type ZeroBotMessager struct {
	z *zero.Ctx
}

func newZeroBotMessager(z *zero.Ctx) *ZeroBotMessager {
	return &ZeroBotMessager{z: z}
}

func (z *ZeroBotMessager) Send(message model.Message) {
	switch m := message.(type) {
	case model.TextMessage:
		z.z.Send(m.Text)
	case model.ImageMessage:
		z.z.Send(zMsg.ImageBytes(m.Image))
	default:
		z.z.Send("呃呃呃，这个消息居然发不出来？(Not implemented)")
	}
}

func (z ZeroBotAdapter) Bind(events []model.Event) {
	for _, event := range events {
		commands := event.Commands
		handler := event.Handler

		zeroHandler := func(zCtx *zero.Ctx) {
			ctx_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			ctx := &model.Context{
				Ctx:      ctx_,
				Adapter:  z,
				Messager: newZeroBotMessager(zCtx),
				Params:   strings.Fields(zCtx.State["args"].(string)),
				UserID:   zCtx.Event.UserID,
				Username: zCtx.Event.Sender.NickName,
			}
			start := time.Now()
			err := handler(ctx)
			end := time.Now()
			logrus.Debug("handler cost:", end.Sub(start))
			if err != nil {
				var friendlyError errs.FriendlyError
				if errors.As(err, &friendlyError) {
					zCtx.Send(friendlyError.Error())
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					zCtx.Send("请求超时，再试一次？")
					return
				}
				logrus.Error(err)
				zCtx.Send("寄！出问题了！")
				return
			}
		}

		for _, command := range commands {
			zero.OnCommand(command).Handle(zeroHandler)
		}
	}
}

func (z ZeroBotAdapter) Start() {
	zero.Run(&z.zeroCfg)
}

func (z ZeroBotAdapter) Platform() string {
	return "QQ"
}
