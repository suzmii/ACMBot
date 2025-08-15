package subconfig

type Sync struct {
	RaceRefreshHours int `mapstructure:"RaceRefreshHours"`
	RaceKeepHours    int `mapstructure:"RaceKeepHours"`
	// Codeforces related TTLs
	CodeforcesUserRefreshHours       int `mapstructure:"CodeforcesUserRefreshHours"`
	CodeforcesRatingRefreshHours     int `mapstructure:"CodeforcesRatingRefreshHours"`
	CodeforcesSubmissionRefreshHours int `mapstructure:"CodeforcesSubmissionRefreshHours"`
}

var DefaultSync = Sync{
	RaceRefreshHours:                 24,
	RaceKeepHours:                    24,
	CodeforcesUserRefreshHours:       4,
	CodeforcesRatingRefreshHours:     4,
	CodeforcesSubmissionRefreshHours: 24,
}
