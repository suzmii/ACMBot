package main

import (
	"context"
	"time"

	"github.com/suzmii/ACMBot/handler"
	"github.com/suzmii/ACMBot/middleware"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func StartZeroBot(handler *handler.Handler) {
	// TODO: todo
	zero.OnMessage(zero.RegexRule(`\s*cf\s(.*)$`)).Handle(func(ctx *zero.Ctx) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO: 配置文件定义等待时间
		defer cancel()
		matched := ctx.State["regex_matched"].([]string)
		image, err := handler.GetCodeforcesUserProfileImage(c, matched[0])
		if err != nil {
			userMsg := middleware.HandleError(err)
			ctx.Send(message.Text(userMsg))
			return
		}
		ctx.Send(message.ImageBytes(image))
	})

}
