package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/handler"
	"github.com/suzmii/ACMBot/middleware"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// formatRaceInfo 格式化比赛信息
func formatRaceInfo(race database.Race) string {
	now := time.Now()
	var status, statusEmoji string
	var timeInfo string

	if race.Start.After(now) {
		// 尚未开始
		status = "此比赛尚未开始"
		statusEmoji = "🕣"
		duration := race.Start.Sub(now)
		days := int(duration.Hours()) / 24
		hours := int(duration.Hours()) % 24
		minutes := int(duration.Minutes()) % 60
		timeInfo = fmt.Sprintf("距离开始: %02d天%02d小时%02d分钟", days, hours, minutes)
	} else if race.End.After(now) {
		// 已开始
		status = "此比赛已开始"
		statusEmoji = "🕘"
		duration := race.End.Sub(now)
		days := int(duration.Hours()) / 24
		hours := int(duration.Hours()) % 24
		minutes := int(duration.Minutes()) % 60
		timeInfo = fmt.Sprintf("距离结束: %02d天%02d小时%02d分钟", days, hours, minutes)
	} else {
		// 已结束
		status = "此比赛已结束"
		statusEmoji = "🏁"
		timeInfo = "比赛已结束"
	}

	// 计算持续时间
	duration := race.End.Sub(race.Start)
	durationDays := int(duration.Hours()) / 24
	durationHours := int(duration.Hours()) % 24
	durationMinutes := int(duration.Minutes()) % 60
	var durationStr string
	if durationDays > 0 {
		durationStr = fmt.Sprintf("%d天%d小时%d分钟", durationDays, durationHours, durationMinutes)
	} else if durationHours > 0 {
		durationStr = fmt.Sprintf("%d小时%d分钟", durationHours, durationMinutes)
	} else {
		durationStr = fmt.Sprintf("%d分钟", durationMinutes)
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("%s%s%s\n", statusEmoji, status, statusEmoji))
	msg.WriteString(fmt.Sprintf("比赛来源: %s\n", race.Resource))
	msg.WriteString(fmt.Sprintf("比赛名称: %s\n", race.Title))
	msg.WriteString(fmt.Sprintf("%s\n", timeInfo))
	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	msg.WriteString(fmt.Sprintf("开始时间: %s %s\n", race.Start.Format("2006-01-02 15:04:05"), weekdays[race.Start.Weekday()]))
	msg.WriteString(fmt.Sprintf("持续时间: %s\n", durationStr))

	return msg.String()
}

func StartZeroBot(handler *handler.Handler) {
	// 读取配置
	cfg := config.LoadConfig()

	// AtCoder 用户资料查询
	zero.OnMessage(zero.RegexRule(`\s*at\s+(\S+)\s*$`)).Handle(func(ctx *zero.Ctx) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		matched := ctx.State["regex_matched"].([]string)
		username := matched[1]

		username = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
				return r
			}
			return -1
		}, username)

		if username == "" {
			ctx.Send(message.Text("请输入有效的 AtCoder 用户名"))
			return
		}

		image, err := handler.GetAtcoderUserProfileImage(c, username)
		if err != nil {
			userMsg := middleware.HandleError(err)
			ctx.Send(message.Text(userMsg))
			return
		}

		ctx.Send(message.ImageBytes(image))
	})

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

	// Codeforces rating 曲线图查询
	zero.OnMessage(zero.RegexRule(`\s*(rt|rating)\s+(\S+)\s*$`)).Handle(func(ctx *zero.Ctx) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO: 配置文件定义等待时间
		defer cancel()
		matched := ctx.State["regex_matched"].([]string)
		username := matched[2] // 获取用户名

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

		image, err := handler.GetCodeforcesRatingImage(c, username)
		if err != nil {
			userMsg := middleware.HandleError(err)
			ctx.Send(message.Text(userMsg))
			return
		}
		ctx.Send(message.ImageBytes(image))
	})

	// 近期比赛查询（支持所有平台和特定平台）
	zero.OnMessage(zero.RegexRule(`\s*(race|比赛|近期比赛|近期\s*(.+))\s*(\d*)\s*$`)).Handle(func(ctx *zero.Ctx) {
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

		// 解析页码
		page := 1
		pageStr := matched[3] // 匹配页码
		if pageStr != "" {
			fmt.Sscanf(pageStr, "%d", &page)
			if page < 1 {
				page = 1
			}
		}

		const pageSize = 4 // 每页最多展示4条

		// 检查是否指定了平台
		platformName := matched[2] // 匹配 "近期\s*(.+)" 的捕获组

		var races []database.Race
		var err error

		// 获取比赛数据
		if platformName != "" {
			// 特定平台查询
			resource, ok := resourceMap[platformName]
			if !ok {
				ctx.Send(message.Text("不支持的平台，支持的平台：codeforces/cf, atcoder, leetcode/力扣, 洛谷, 牛客"))
				return
			}
			races, err = handler.GetUpcomingRace(c, resource)
			if err != nil {
				userMsg := middleware.HandleError(err)
				ctx.Send(message.Text(userMsg))
				return
			}
			if len(races) == 0 {
				ctx.Send(message.Text(fmt.Sprintf("暂无 %s 的即将开始的比赛", platformName)))
				return
			}
		} else {
			// 所有平台查询
			races, err = handler.GetUpcomingRaces(c)
			if err != nil {
				userMsg := middleware.HandleError(err)
				ctx.Send(message.Text(userMsg))
				return
			}
			if len(races) == 0 {
				ctx.Send(message.Text("暂无即将开始的比赛"))
				return
			}
		}

		// 按开始时间离当前时间的绝对值排序
		now := time.Now()
		sort.Slice(races, func(i, j int) bool {
			distI := races[i].Start.Sub(now).Abs()
			distJ := races[j].Start.Sub(now).Abs()
			return distI < distJ
		})

		// 计算分页
		totalPages := (len(races) + pageSize - 1) / pageSize
		if page > totalPages {
			page = totalPages
		}
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > len(races) {
			end = len(races)
		}
		pageRaces := races[start:end]

		// 格式化输出比赛信息
		raceStrs := make([]string, len(pageRaces))
		for i, race := range pageRaces {
			raceStrs[i] = formatRaceInfo(race)
		}

		var msg strings.Builder
		if platformName != "" {
			msg.WriteString(fmt.Sprintf("📅 %s 近期比赛（第%d页/共%d页）\n\n", platformName, page, totalPages))
		} else {
			msg.WriteString("📅 近期比赛\n\n")
		}

		msg.WriteString(strings.Join(raceStrs, "\n"))

		if platformName != "" {
			msg.WriteString(fmt.Sprintf("━━━━━━━━━━━━\n第%d页/共%d页\n输入「近期%s <页码>」查看其他页", page, totalPages, platformName))
		} else {
			msg.WriteString(fmt.Sprintf("━━━━━━━━━━━━\n第%d页/共%d页\n输入「近期比赛 <页码>」查看其他页", page, totalPages))
		}

		ctx.Send(message.Text(msg.String()))
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
