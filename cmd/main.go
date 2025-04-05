package main

import (
	"github.com/suzmii/ACMBot/pkg/model/db"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	_ "github.com/suzmii/ACMBot/internal/bot/platforms/qq"
)

func main() {
	if err := db.MigrateAll(); err != nil {
		logrus.Fatal(err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
