package main

import (
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/database"
)

func main() {
	err := database.Init().AutoMigrate(database.AllModels...)
	if err != nil {
		logrus.Fatal(err)
	}
}
