package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/mappers"
	"sstcloud-alice-gateway/internal/models/alice"
	"sstcloud-alice-gateway/pkg/middleware/user"
)

func (s *service) Devices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)

	devices := s.deviceProvider.Devices(user.User(ctx))

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
