package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/handler"
	"github.com/suzmii/ACMBot/middleware"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func StartZeroBot(handler *handler.Handler) {
	// 读取配置
	cfg := config.LoadConfig()

	// Codeforces 用户资料查询
	zero.OnMessage(zero.RegexRule(`\s*cf\s+(\S+)\s*$`)).Handle(func(ctx *zero.Ctx) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO: 配置文件定义等待时间
		defer cancel()
		matched := ctx.State["regex_matched"].([]string)
		username := matched[1] // 获取第一个非空白字符序列作为用户名

		// 去除特殊字符，只保留字母、数字和下划线
		username = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
				return r
			}
			return -1
		}, username)

		if username == "" {
			ctx.Send(message.Text("请输入有效的 Codeforces 用户名"))
			return
		}

		image, err := handler.GetCodeforcesUserProfileImage(c, username)
		if err != nil {
			userMsg := middleware.HandleError(err)
			ctx.Send(message.Text(userMsg))
			return
		}
		ctx.Send(message.ImageBytes(image))
	})

	// 近期比赛查询（支持所有平台和特定平台）
	zero.OnMessage(zero.RegexRule(`\s*(race|比赛|近期比赛|近期\s*(.+))\s*$`)).Handle(func(ctx *zero.Ctx) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		matched := ctx.State["regex_matched"].([]string)

		// 将中文平台名称映射到 resource 常量
		resourceMap := map[string]string{
			"codeforces": "codeforces.com",
			"cf":         "codeforces.com",
			"atcoder":    "atcoder.jp",
			"at":         "atcoder.jp",
			"leetcode":   "leetcode.com",
			"lc":         "leetcode.com",
			"力扣":         "leetcode.com",
			"lg":         "luogu.com.cn",
			"洛谷":         "luogu.com.cn",
			"nk":         "ac.nowcoder.com",
			"牛客":         "ac.nowcoder.com",
		}

		// 检查是否指定了平台
		platformName := matched[2] // 匹配 "近期\s*(.+)" 的捕获组
		if platformName != "" {
			// 特定平台查询
			resource, ok := resourceMap[platformName]
			if !ok {
				ctx.Send(message.Text("不支持的平台，支持的平台：codeforces/cf, atcoder, leetcode/力扣, 洛谷, 牛客"))
				return
			}

			races, err := handler.GetUpcomingRace(c, resource)
			if err != nil {
				userMsg := middleware.HandleError(err)
				ctx.Send(message.Text(userMsg))
				return
			}
			if len(races) == 0 {
				ctx.Send(message.Text(fmt.Sprintf("暂无 %s 的即将开始的比赛", platformName)))
				return
			}

			// 格式化输出比赛信息
			var msg strings.Builder
			msg.WriteString(fmt.Sprintf("📅 %s 近期比赛\n\n", platformName))
			for _, race := range races {
				msg.WriteString(fmt.Sprintf("🏆 %s\n", race.Title))
				msg.WriteString(fmt.Sprintf("⏰ 开始: %s\n", race.Start.Format("2006-01-02 15:04")))
				msg.WriteString(fmt.Sprintf("⏰ 结束: %s\n", race.End.Format("2006-01-02 15:04")))
				msg.WriteString(fmt.Sprintf("🔗 链接: %s\n\n", race.Link))
			}
			ctx.Send(message.Text(msg.String()))
		} else {
			// 所有平台查询
			races, err := handler.GetUpcomingRaces(c)
			if err != nil {
				userMsg := middleware.HandleError(err)
				ctx.Send(message.Text(userMsg))
				return
			}
			if len(races) == 0 {
				ctx.Send(message.Text("暂无即将开始的比赛"))
				return
			}
			// 格式化输出比赛信息
			var msg strings.Builder
			msg.WriteString("📅 近期比赛（未来7天）\n\n")
			for _, race := range races {
				msg.WriteString(fmt.Sprintf("🏆 %s\n", race.Title))
				msg.WriteString(fmt.Sprintf("📌 平台: %s\n", race.Resource))
				msg.WriteString(fmt.Sprintf("⏰ 开始: %s\n", race.Start.Format("2006-01-02 15:04")))
				msg.WriteString(fmt.Sprintf("🔗 链接: %s\n\n", race.Link))
			}
			ctx.Send(message.Text(msg.String()))
		}
	})

	// 启动 ZeroBot 服务器
	zero.Run(&zero.Config{
		NickName:      []string{"ACMBot"},
		CommandPrefix: cfg.ZeroBot.CommandPrefix,
		SuperUsers:    []int64{},
		Driver: []zero.Driver{
			driver.NewWebSocketClient(fmt.Sprintf("ws://%s:%d", cfg.ZeroBot.Host, cfg.ZeroBot.Port), cfg.ZeroBot.Token),
		},
	})
}
