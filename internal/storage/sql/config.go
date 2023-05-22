package sql

type Config struct {
	ConnectionString string `env:"DB_CONNECTION_STRING,required"`
	LogOnlyErrors    bool   `env:"LOG_ONLY_ERROR,default=true"`
}
