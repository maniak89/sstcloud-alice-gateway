package rest

type Config struct {
	Address string `env:"HTTP_ADDRESS,default=:80"`
}
