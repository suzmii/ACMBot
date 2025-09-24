package subconfig

type Database struct {
	Host         string `mapstructure:"host"           toml:"host"`
	Port         int    `mapstructure:"port"           toml:"port"`
	Username     string `mapstructure:"username"       toml:"username"`
	Password     string `mapstructure:"password"       toml:"password"`
	Name         string `mapstructure:"name"           toml:"name"`
	AutoCreateDB bool   `mapstructure:"auto_create_db" toml:"auto_create_db"`
	AutoMigrate  bool   `mapstructure:"auto_migrate"   toml:"auto_migrate"`
}

var DefaultDB = Database{
	Host:         "localhost",
	Port:         3306,
	Username:     "root",
	Password:     "password",
	Name:         "ACMBot",
	AutoCreateDB: false,
	AutoMigrate:  true,
}
