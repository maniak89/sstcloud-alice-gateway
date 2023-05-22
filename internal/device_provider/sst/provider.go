package sst

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/models/common"
	"sstcloud-alice-gateway/pkg/sst"
)

const (
	house_id  = "house_id"
	device_id = "device_id"
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

func (c *Client) Devices(ctx context.Context) ([]common.Device, error) {
	houses, err := c.cl.Houses(ctx)
	if err != nil {
		return nil, err
	}
	globalDevices := map[int]struct{}{}
	var result []common.Device
	for _, house := range houses {
		devices, err := c.cl.Devices(ctx, house.ID)
		if err != nil {
			return nil, err
		}
		for _, device := range devices {
			if _, exist := globalDevices[device.ID]; exist {
				continue
			}
			if len(houses) > 0 {
				device.Name = house.Name + " " + device.Name
			}
			if device.Type != sst.MCS350 && device.Type != sst.MCS300 {
				log.Ctx(ctx).Warn().Str("type", device.Type.String()).Str("name", device.Name).Msg("Not supported type")
				continue
			}

			globalDevices[device.ID] = struct{}{}
			result = append(result, common.Device{
				ID:   fmt.Sprintf("%d_%d", house.ID, device.ID),
				Name: device.Name,
				AdditionalFields: map[string]string{
					house_id:  strconv.Itoa(house.ID),
					device_id: strconv.Itoa(device.ID),
				},
				Tempometer: &common.Tempometer{
					DegreesFloor:    device.TermParsedConfiguration.CurrentTemperature.TemperatureFloor,
					DegreesAir:      device.TermParsedConfiguration.CurrentTemperature.TemperatureAir,
					SetDegreesFloor: device.TermParsedConfiguration.Settings.TemperatureManual,
				},
				Model:     device.Type.String(),
				Enabled:   device.TermParsedConfiguration.Settings.Status == sst.DeviceStatusOn,
				Connected: device.IsConnected,
			})
		}
	}
	return result, nil
}

func (c *Client) SetTemperature(ctx context.Context, device common.Device, temp int) error {
	houseID, deviceID, err := extractAdditionalFields(ctx, device)
	if err != nil {
		return err
	}
	if err := c.cl.PowerStatus(ctx, houseID, deviceID, true); err != nil {
		return err
	}
	if err := c.cl.Temperature(ctx, houseID, deviceID, temp); err != nil {
		return err
	}
	return nil
}

func (c *Client) PowerStatus(ctx context.Context, device common.Device, power bool) error {
	houseID, deviceID, err := extractAdditionalFields(ctx, device)
	if err != nil {
		return err
	}
	if err := c.cl.PowerStatus(ctx, houseID, deviceID, power); err != nil {
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

func extractAdditionalFields(ctx context.Context, device common.Device) (int, int, error) {
	logger := log.Ctx(ctx).With().Str("device_id", device.ID).Str("name", device.Name).Logger()
	houseId, err := strconv.Atoi(device.AdditionalFields[house_id])
	if err != nil {
		logger.Error().Err(err).Msg("Failed parse house id")
		return 0, 0, err
	}
	deviceId, err := strconv.Atoi(device.AdditionalFields[device_id])
	if err != nil {
		logger.Error().Err(err).Msg("Failed parse device id")
		return 0, 0, err
	}
	return houseId, deviceId, nil
}
