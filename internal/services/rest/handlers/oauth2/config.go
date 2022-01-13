package oauth2

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Enabled               bool   `env:"OAUTH2_ENABLED,default=false"`
	Key                   string `env:"OAUTH2_KEY,required"`
	AuthAlgo              string `env:"OAUTH2_JWT_ALGO,default=HS256"`
	AuthVerifyKeyInBase64 bool   `env:"OAUTH2_KEY_IN_BASE64"`
	UserID                string `env:"OAUTH2_USER_ID"`
	UserSecret            string `env:"OAUTH2_USER_SECRET"`
	UserDomain            string `env:"OAUTH2_USER_DOMAIN,default=http://localhost"`
	TokenFile             string `env:"OAUTH2_TOKEN_STORE_FILE,default=tokens"`
}

func (c Config) Validate() error {
	if c.Key == "" {
		return errors.New("empty OAUTH2_KEY")
	}
	if !c.Enabled {
		return nil
	}
	if c.UserID == "" {
		return errors.New("empty OAUTH2_USER_ID")
	}
	if c.UserSecret == "" {
		return errors.New("empty OAUTH2_USER_SECRET")
	}
	return nil
}

func (c Config) GetAuthKey(ctx context.Context) ([]byte, error) {
	key := []byte(c.Key)
	if c.AuthVerifyKeyInBase64 {
		k, err := base64.StdEncoding.DecodeString(c.Key)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("key", c.Key).Msg("Failed base64 decode")
			return nil, err
		}
		key = k
	}
	return key, nil
}
