package subconfig

type Database struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	AutoCreateDB bool   `mapstructure:"auto_create_db"`
	AutoMigrate  bool   `mapstructure:"auto_migrate"`
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
