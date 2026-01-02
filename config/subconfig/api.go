package subconfig

type API struct {
	CodeforcesKey      string `mapstructure:"codeforces_key"      toml:"codeforces_key"`
	CodeforcesSecret   string `mapstructure:"codeforces_secret"   toml:"codeforces_secret"`
	ClistAuthenticated string `mapstructure:"clist_authenticated" toml:"clist_authenticated"`
}
