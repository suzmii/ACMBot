package main

import (
	"github.com/suzmii/ACMBot/internal/adapter/qq"
	"github.com/suzmii/ACMBot/internal/database"
	"github.com/suzmii/ACMBot/internal/handler"
	"github.com/suzmii/ACMBot/internal/render"
	"github.com/suzmii/ACMBot/internal/util/logger"
	"os"
	"os/signal"
	"syscall"
)

import _ "net/http/pprof"
import "net/http"

func main() {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	logger.Init()
	database.Init()
	render.Init()

	qqAdapter := qq.NewZeroBotAdapter()
	qqAdapter.Bind(handler.Events)
	qqAdapter.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
