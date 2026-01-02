package subconfig

type Render struct {
	Headless bool `mapstructure:"headless"`
	PoolSize int  `mapstructure:"poolSize"`
}

var DefaultRenderConfig = Render{
	Headless: true,
	PoolSize: 8,
}
