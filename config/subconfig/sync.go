package subconfig

type Sync struct {
	RaceRefreshHours int `mapstructure:"RaceRefreshHours"`
}

var DefaultSync = Sync{
	RaceRefreshHours: 24,
}
