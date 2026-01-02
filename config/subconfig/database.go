package subconfig

type Database struct {
	DSN string `mapstructure:"dsn" toml:"dsn"`
}

var DefaultDB = Database{
	DSN: "postgres://postgres:postgres@localhost:5432/postgres",
}
