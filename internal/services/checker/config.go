package checker

import (
	"time"
)

type Config struct {
	RequestPeriod time.Duration `env:"REQUEST_PERIOD,default=5m"`
}
