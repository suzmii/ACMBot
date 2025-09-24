package subconfig

type Logger struct {
	AlterToken string `mapstructure:"alter_token" toml:"alter_token"`
	Level      string `mapstructure:"level"       toml:"level"`
}

var DefaultLogger = Logger{
	AlterToken: "",
	Level:      "info",
}
