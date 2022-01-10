package sst

import (
	"context"
	"encoding/json"

	"github.com/maniak89/sstcloud-alice-gateway/pkg/sst"
	"github.com/rs/zerolog/log"
)

type Config struct {
	sst.Config
	Username string       `env:"SST_USERNAME"`
	Password string       `env:"SST_PASSWORD,required"`
	EMail    string       `env:"SST_EMAIL,required"`
	Lang     sst.Language `env:"SST_LANGUAGE,default=ru"`
}

type Client struct {
	cl     *sst.Client
	config Config
}

func New(config Config) *Client {
	return &Client{
		cl:     sst.New(config.Config),
		config: config,
	}
}

func (c *Client) Init(ctx context.Context) error {
	_, err := c.cl.Login(ctx, sst.LoginRequest{
		Username: c.config.Username,
		Password: c.config.Password,
		EMail:    c.config.EMail,
		Language: c.config.Lang,
	})
	return err
}

func (c *Client) Devices(ctx context.Context) error {
	logger := log.Ctx(ctx)
	houses, err := c.cl.Houses(ctx)
	if err != nil {
		return err
	}
	globalDevices := map[int]sst.Device{}
	for _, house := range houses {
		devices, err := c.cl.Devices(ctx, house.ID)
		if err != nil {
			return err
		}
		for _, device := range devices {
			if len(houses) > 0 {
				device.Name = house.Name + " " + device.Name
			}
			globalDevices[device.ID] = device
		}
	}
	blob, _ := json.Marshal(globalDevices)
	logger.Trace().RawJSON("devices", blob).Msg("Houses given")
	return nil
}
