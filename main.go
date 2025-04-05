package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	_ "github.com/suzmii/ACMBot/app/bot/platforms/qq"
	"github.com/suzmii/ACMBot/app/model/db"
)

func main() {
	if err := db.MigrateAll(); err != nil {
		logrus.Fatal(err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
