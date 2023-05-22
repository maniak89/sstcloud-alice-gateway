package sst

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/pkg/sst"
)

type Config struct {
	sst.Config
	Password string
	EMail    string
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
		Username: c.config.EMail,
		Password: c.config.Password,
		EMail:    c.config.EMail,
		Language: sst.LangRu,
	})
	return err
}

func (c *Client) Houses(ctx context.Context) ([]*device_provider.House, error) {
	houses, err := c.cl.Houses(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*device_provider.House, 0, len(houses))
	for _, h := range houses {
		result = append(result, &device_provider.House{
			ID:             h.ID,
			Name:           h.Name,
			DeviceProvider: c,
		})
	}
	return result, nil
}

func (c *Client) Devices(ctx context.Context, house *device_provider.House) ([]*device_provider.Device, error) {
	devices, err := c.cl.Devices(ctx, house.ID)
	if err != nil {
		return nil, err
	}
	result := make([]*device_provider.Device, 0, len(devices))
	now := time.Now()
	for _, device := range devices {
		device.Name = house.Name + " " + device.Name
		if device.Type != sst.MCS350 && device.Type != sst.MCS300 {
			log.Ctx(ctx).Warn().Str("type", device.Type.String()).Str("name", device.Name).Msg("Not supported type")
			continue
		}

		result = append(result, &device_provider.Device{
			ID:    device.ID,
			House: house,
			IDStr: fmt.Sprintf("%d_%d", house.ID, device.ID),
			Name:  house.Name + " " + device.Name,
			Tempometer: device_provider.Tempometer{
				DegreesFloor:             device.TermParsedConfiguration.CurrentTemperature.TemperatureFloor,
				DegreesAir:               device.TermParsedConfiguration.CurrentTemperature.TemperatureAir,
				SetDegreesFloor:          device.TermParsedConfiguration.Settings.TemperatureManual,
				ChangedAtDegreesFloor:    now,
				ChangedAtDegreesAir:      now,
				ChangedAtSetDegreesFloor: now,
			},
			Model:     device.Type.String(),
			Enabled:   device.TermParsedConfiguration.Settings.Status == sst.DeviceStatusOn,
			Connected: device.IsConnected,
			UpdatedAt: now,
		})
	}
	return result, nil
}

func (c *Client) SetTemperature(ctx context.Context, device *device_provider.Device, temp int) error {
	if err := c.cl.PowerStatus(ctx, device.House.ID, device.ID, true); err != nil {
		return err
	}
	if err := c.cl.Temperature(ctx, device.House.ID, device.ID, temp); err != nil {
		return err
	}
	return nil
}

func (c *Client) PowerStatus(ctx context.Context, device *device_provider.Device, power bool) error {
	if err := c.cl.PowerStatus(ctx, device.House.ID, device.ID, power); err != nil {
		return err
	}
	return nil
}

func (c *Client) EMail() string {
	return c.config.EMail
}

func (c *Client) Password() string {
	return c.config.Password
}
