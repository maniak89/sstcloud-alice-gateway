package rest

import (
	"github.com/maniak89/sstcloud-alice-gateway/internal/services/rest/handlers/oauth2"
)

type Config struct {
	Address string `env:"HTTP_ADDRESS,default=:80"`
	OAUTH2  oauth2.Config
}
