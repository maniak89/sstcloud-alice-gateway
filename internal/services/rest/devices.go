package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/maniak89/sstcloud-alice-gateway/internal/mappers"
	"github.com/maniak89/sstcloud-alice-gateway/internal/models/alice"
	"github.com/rs/zerolog/log"
)

func (s *service) Devices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)
	devices, err := s.deviceProvider.Devices(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed get devices")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, _, _ := jwtauth.FromContext(ctx)

	aliceDevices := alice.Devices{
		UserID:  token.Subject(),
		Devices: make([]alice.Device, len(devices)),
	}
	for i, d := range devices {
		aliceDevices.Devices[i] = mappers.DeviceToAlice(d)
	}

	if err := json.NewEncoder(w).Encode(alice.Response{
		RequestID: r.Header.Get(xRequestID),
		Payload:   aliceDevices,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
