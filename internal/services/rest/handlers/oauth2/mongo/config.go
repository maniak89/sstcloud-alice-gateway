package mongo

type Config struct {
	URI  string `env:"MONGO_DB_URI"`
	Name string `env:"MONGO_DB_NAME"`
}
