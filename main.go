package main

import (
	"context"

	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/handler"
	"github.com/suzmii/ACMBot/render"
	"github.com/suzmii/ACMBot/scheduler"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("main")

func main() {
	logx.Init(subconfig.DefaultLogger)

	ctx := context.Background()
	logger.Info("开始初始化...")
	logger.Info("正在读取配置信息")
	cfg := config.LoadConfig()
	logger.Info("配置信息加载完成")

	logger.Info("正在初始化logger")
	logx.Init(cfg.Logger)
	logger.Info("logger初始化完成")

	logger.Info("正在初始化数据库")
	store, err := database.New(ctx, cfg.Database.DSN)
	if err != nil {
		logger.Fatalf("无法初始化数据库: %v", err)
	}

	logger.Info("正在初始化API")
	api, err := api.New(cfg.API)
	if err != nil {
		logger.Fatalf("无法初始化API: %v", err)
	}

	logger.Info("正在初始化Render")
	render, err := render.NewRender(cfg.Render)
	if err != nil {
		logger.Fatalf("无法初始化Render: %v", err)
	}

	logger.Info("正在初始化调度器")

	tasks := []scheduler.Task{
		scheduler.NewRaceUpdateTask(store, api),
	}
	scheduler, err := scheduler.NewScheduler(cfg.Scheduler, store, tasks)
	if err != nil {
		logger.Fatalf("无法初始化调度器: %v", err)
	}

	logger.Info("启动调度器")
	scheduler.Start()

	logger.Info("正在初始化响应函数")
	handler := handler.NewHandler(cfg.Handler, store, api, render)

	logger.Info("正在初始化ZeroBot")
	StartZeroBot(handler)

	select {}
}
