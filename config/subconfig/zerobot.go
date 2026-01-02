package subconfig

type ZeroBot struct {
	CommandPrefix string `mapstructure:"command_prefix" toml:"command_prefix"`
	Host          string `mapstructure:"host"           toml:"host"`
	Port          int    `mapstructure:"port"           toml:"port"`
	Token         string `mapstructure:"token"          toml:"token"`
}

var DefaultZeroBot = ZeroBot{
	CommandPrefix: "/",
	Host:          "localhost",
	Port:          15630,
	Token:         "",
}
