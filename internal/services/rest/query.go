package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/mappers"
	"sstcloud-alice-gateway/internal/models/alice"
	"sstcloud-alice-gateway/pkg/middleware/user"
)

func (s *service) Query(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)
	var req alice.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed unmarshal data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	devices, _, err := s.devices(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed get devices")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aliceDevices := alice.Devices{
		UserID:  user.User(ctx),
		Devices: make([]alice.Device, 0, len(devices)),
	}
	for _, reqDev := range req.Devices {
		for _, dev := range devices {
			if dev.ID != mappers.ExtractDeviceID(reqDev.ID) {
				continue
			}
			newDevices := mappers.DeviceToAlice(dev)
			for _, newDev := range newDevices {
				if newDev.ID != reqDev.ID {
					continue
				}
				aliceDevices.Devices = append(aliceDevices.Devices, newDev)
			}
			break
		}
	}
	if err := json.NewEncoder(w).Encode(alice.Response{
		RequestID: r.Header.Get(xRequestID),
		Payload:   aliceDevices,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
