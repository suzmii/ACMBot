package subconfig

type Sync struct {
	RaceRefreshHours int `mapstructure:"RaceRefreshHours"`
	RaceKeepHours    int `mapstructure:"RaceKeepHours"`
}

var DefaultSync = Sync{
	RaceRefreshHours: 24,
	RaceKeepHours:    24,
}
