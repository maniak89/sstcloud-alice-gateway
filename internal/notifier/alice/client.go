package alice

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/mappers"
	"sstcloud-alice-gateway/internal/models/alice"
)

type client struct {
	config          Config
	callbackAddress string
	client          *http.Client
}

func New(config Config) *client {
	return &client{
		config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
		callbackAddress: config.Address + "/api/v1/skills/" + config.SkillID + "/callback/state",
	}
}

func (c *client) NotifyDevicesChanged(ctx context.Context, house *device_provider.House, device []*device_provider.Device) error {
	logger := log.Ctx(ctx)
	devices := make([]alice.PayloadStateDevice, 0, len(device))
	for _, dev := range device {
		objs := mappers.DeviceToAlice(dev)
		for _, obj := range objs {
			if len(obj.Properties) == 0 && len(obj.Capabilities) == 0 {
				continue
			}
			props := make([]alice.PayloadStateDeviceProperties, 0, len(obj.Properties))
			for _, prop := range obj.Properties {
				props = append(props, alice.PayloadStateDeviceProperties{
					Type:  prop.Type,
					State: prop.State,
				})
			}
			devices = append(devices, alice.PayloadStateDevice{
				ID:           obj.ID,
				Properties:   props,
				Capabilities: obj.Capabilities,
			})
		}
	}
	blob, err := json.Marshal(alice.State{
		TS: time.Now().Unix(),
		Payload: alice.PayloadState{
			UserID:  house.UserID,
			Devices: devices,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed marshal body")
		return err
	}
	logger.Trace().Str("url", c.callbackAddress).Bytes("blob", blob).Msg("Prepare request body")
	req, err := http.NewRequest(http.MethodPost, c.callbackAddress, bytes.NewReader(blob))
	if err != nil {
		logger.Error().Err(err).Msg("Failed create request object")
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "OAuth "+c.config.OAuth2Token)
	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed make request")
		return err
	}
	if resp.StatusCode != http.StatusOK {
		blob, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Failed read body")
		}
		logger.Error().Str("status", resp.Status).Bytes("response", blob).Msg("status")
		return nil
	}
	logger.Debug().Str("status", resp.Status).Msg("status")

	return nil
}
