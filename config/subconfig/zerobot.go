package subconfig

type ZeroBot struct {
	CommandPrefix string `mapstructure:"command_prefix"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Token         string `mapstructure:"token"`
}

var DefaultZeroBot = ZeroBot{
	CommandPrefix: "/",
	Host:          "localhost",
	Port:          15630,
	Token:         "",
}
