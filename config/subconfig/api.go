package subconfig

type API struct {
	CodeforcesKey      string `mapstructure:"codeforces_key"`
	CodeforcesSecret   string `mapstructure:"codeforces_secret"`
	ClistAuthenticated string `mapstructure:"clist_authenticated"`
}
