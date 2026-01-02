package config

import (
	"errors"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/suzmii/ACMBot/config/subconfig"
)

type Config struct {
	API       subconfig.API       `mapstructure:"api"`
	Database  subconfig.Database  `mapstructure:"database"`
	Logger    subconfig.Logger    `mapstructure:"logger"`
	ZeroBot   subconfig.ZeroBot   `mapstructure:"zerobot"`
	Scheduler subconfig.Scheduler `mapstructure:"scheduler"`
	Handler   subconfig.Handler   `mapstructure:"handler"`
	Render    subconfig.Render    `mapstructure:"render"`
}

var DefaultConfig = Config{
	API:       subconfig.API{},
	Database:  subconfig.DefaultDB,
	Logger:    subconfig.DefaultLogger,
	ZeroBot:   subconfig.DefaultZeroBot,
	Scheduler: subconfig.DefaultScheduler,
	Handler:   subconfig.DefaultHandler,
	Render:    subconfig.DefaultRenderConfig,
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
				v.Set("scheduler", DefaultConfig.Scheduler)
				v.Set("handler", DefaultConfig.Handler)
				v.Set("render", DefaultConfig.Render)

				if err := v.SafeWriteConfigAs("config.toml"); err != nil {
					logrus.Fatal("Failed to write default config:", err)
				}
				logrus.Info("已生成默认配置文件，请填写后再次运行")
				os.Exit(0)
			}
		}

		var cfg Config
		if err := v.Unmarshal(&cfg); err != nil {
			logrus.Fatalf("Error unmarshaling config: %v", err)
		}

		logrus.Debugf("从Toml文件中读取到了如下配置: %+#v", cfg)
		config = &cfg
	})

	return config
}
