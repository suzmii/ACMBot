package subconfig

type Logger struct {
	AlterToken string `mapstructure:"alter_token"`
	Level      string `mapstructure:"level"`
}

var DefaultLogger = Logger{
	AlterToken: "",
	Level:      "info",
}
