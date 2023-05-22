package alice

import (
	"time"
)

type Config struct {
	SkillID        string        `env:"YANDEX_ALICE_SKILL_ID,required"`
	Address        string        `env:"YANDEX_ALICE_ADDRESS,default=https://dialogs.yandex.net"`
	RequestTimeout time.Duration `env:"YANDEX_ALICE_TIMEOUT,default=5s"`
	OAuth2Token    string        `env:"YANDEX_OAUTH2_TOKEN,required"`
}
