package config

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/suzmii/ACMBot/config/subconfig"
	"sync"
)

type Config struct {
	API      subconfig.API      `mapstructure:"api"`
	Database subconfig.Database `mapstructure:"database"`
	Logger   subconfig.Logger   `mapstructure:"logger"`
	ZeroBot  subconfig.ZeroBot  `mapstructure:"zerobot"`
}

var DefaultConfig = Config{
	API:      subconfig.API{},
	Database: subconfig.DefaultDB,
	Logger:   subconfig.DefaultLogger,
	ZeroBot:  subconfig.DefaultZeroBot,
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("toml")
		v.AddConfigPath(".")

		// 先尝试读取配置文件
		if err := v.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				logrus.Info("正在生成默认配置文件")

				v.Set("api", DefaultConfig.API)
				v.Set("database", DefaultConfig.Database)
				v.Set("logger", DefaultConfig.Logger)
				v.Set("zerobot", DefaultConfig.ZeroBot)

				if err := v.SafeWriteConfigAs("config.toml"); err != nil {
					logrus.Fatal("Failed to write default config: %v", err)
				}
				logrus.Fatal("已生成默认配置文件，请填写后再次运行")
			}
		}

		var cfg Config
		if err := v.Unmarshal(&cfg); err != nil {
			logrus.Fatalf("Error unmarshaling config: %v", err)
		}

		logrus.Infof("Loaded config from TOML: %+v", cfg)
		config = &cfg
	})

	return config
}
