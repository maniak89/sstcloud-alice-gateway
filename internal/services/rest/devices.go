package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/mappers"
	"sstcloud-alice-gateway/internal/models/alice"
	"sstcloud-alice-gateway/internal/models/common"
	"sstcloud-alice-gateway/pkg/middleware/user"
)

func (s *service) devices(ctx context.Context) ([]common.Device, map[string]device_provider.DeviceProvider, error) {
	logger := log.Ctx(ctx)
	providers, err := s.fetchDeviceProviders(ctx, user.User(ctx))
	if err != nil {
		logger.Error().Err(err).Msg("Failed get device providers")
		return nil, nil, err
	}
	var devices []common.Device
	deviceMap := map[string]device_provider.DeviceProvider{}
	var resultErr error
	for _, provider := range providers {
		devs, err := provider.Devices(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Failed get devices")
			resultErr = multierror.Append(resultErr, err)
			continue
		}
		for _, dev := range devs {
			deviceMap[dev.ID] = provider
		}
		devices = append(devices, devs...)
	}
	if mErr, ok := resultErr.(*multierror.Error); ok {
		if mErr.Len() == len(providers) {
			logger.Error().Err(mErr).Msg("Failed get device from all providers")
			return nil, nil, mErr
		}
	}
	return devices, deviceMap, nil
}

func (s *service) Devices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)

	devices, _, err := s.devices(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aliceDevices := alice.Devices{
		UserID:  user.User(ctx),
		Devices: make([]alice.Device, 0, len(devices)),
	}
	for _, d := range devices {
		aliceDevices.Devices = append(aliceDevices.Devices, mappers.DeviceToAlice(d)...)
	}

	if err := json.NewEncoder(w).Encode(alice.Response{
		RequestID: r.Header.Get(xRequestID),
		Payload:   aliceDevices,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
