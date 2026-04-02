package database

import (
	"errors"
	"fmt"

	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/internal/database/gen"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	mysqldriver "github.com/go-sql-driver/mysql"
)

func Init() *gorm.DB {
	cfg := config.LoadConfig().Database
	var db *gorm.DB
	if cfg.Name == "" {
		logrus.Fatalf("database name is empty")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		if !cfg.AutoCreateDB {
			logrus.Fatalf("failed to connect to DB: %v", err)
		}

		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1049 {
			logrus.Warn(fmt.Sprintf("DataBase %s NOT exist, Creating", cfg.Name))
			// DataBase [DBName] Not Found

			/*
				1. Connect to server without select DB
				2. Create DB
				3. Use it
			*/

			// Connect to server without DB
			dsnNoDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.Username, cfg.Password, cfg.Host, cfg.Port)
			db, err = gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
			if err != nil {
				logrus.Fatalf("Failed to Open DataBase	 while create DB: %v", err)
			}

			// Create DB
			err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Name)).Error
			if err != nil {
				logrus.Fatalf("Failed to create DataBase: %v", err)
			}
			logrus.Infof("Create DB %s Successfully", cfg.Name)

			// Use it
			SQLUseDB := fmt.Sprintf(`
					USE %s
				`, cfg.Name)
			err = db.Exec(SQLUseDB).Error
			if err != nil {
				logrus.Fatalf("Failed to use database %v: %v", cfg.Name, err)
			}
		}
	}

	if cfg.AutoCreateDB {
		err = db.AutoMigrate(AllModels...)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	gen.SetDefault(db)

	return db
}
