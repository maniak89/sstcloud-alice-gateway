package sql

type Config struct {
	ConnectionString string `env:"DB_CONNECTION_STRING,required"`
}
