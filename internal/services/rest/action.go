package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/mappers"
	"sstcloud-alice-gateway/internal/models/alice"
	"sstcloud-alice-gateway/pkg/middleware/user"
)

func (s *service) Action(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)
	var req alice.ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed unmarshal data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	devices, deviceMap, err := s.devices(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aliceDevices := alice.Devices{
		UserID:  user.User(ctx),
		Devices: make([]alice.Device, 0, len(devices)),
	}
	for _, reqDev := range req.Payload.Devices {
		for _, dev := range devices {
			if dev.ID != reqDev.ID {
				continue
			}
			aliceDevice := alice.Device{
				ID: dev.ID,
			}
			logger := log.With().Str("device_id", dev.ID).Logger()
			for _, capability := range reqDev.Capabilities {
				logger := logger.With().Str("capability_type", string(capability.Type)).Logger()

				actionResult := alice.ActionResult{
					Status: alice.ActionResultStatusDone,
				}
				switch capability.Type {
				case alice.CapabilityTypeOnOff:
					if err := deviceMap[reqDev.ID].PowerStatus(ctx, mappers.DeviceFromAlice(reqDev), capability.State.Value.(bool)); err != nil {
						logger.Error().Err(err).Msg("Failed set status")
						actionResult = alice.ActionResult{
							Status:           alice.ActionResultStatusError,
							ErrorCode:        alice.ErrorCodeDeviceUnreachable,
							ErrorDescription: err.Error(),
						}
					}
				case alice.CapabilityTypeRange:
					if capability.State.Instance != string(alice.PropertiesFloatParametersInstanceTemperature) {
						actionResult = alice.ActionResult{
							Status:           alice.ActionResultStatusError,
							ErrorCode:        alice.ErrorCodeInvalidAction,
							ErrorDescription: fmt.Sprintf("unknown action %s", capability.State.Instance),
						}
					} else {
						value := int(capability.State.Value.(float64))
						if capability.State.Relative {
							value = dev.Tempometer.SetDegreesFloor + int(capability.State.Value.(float64))
						}
						if value > mappers.MaxTemp || value < mappers.MinTemp {
							actionResult = alice.ActionResult{
								Status:           alice.ActionResultStatusError,
								ErrorCode:        alice.ErrorCodeInvalidAction,
								ErrorDescription: fmt.Sprintf("value %d not in range %d-%d", value, mappers.MinTemp, mappers.MaxTemp),
							}
						} else if err := deviceMap[reqDev.ID].SetTemperature(ctx, mappers.DeviceFromAlice(reqDev), value); err != nil {
							logger.Error().Err(err).Msg("Failed set status")
							actionResult = alice.ActionResult{
								Status:           alice.ActionResultStatusError,
								ErrorCode:        alice.ErrorCodeDeviceUnreachable,
								ErrorDescription: err.Error(),
							}
						}
					}
				default:
					actionResult = alice.ActionResult{
						Status:           alice.ActionResultStatusError,
						ErrorCode:        alice.ErrorCodeInvalidAction,
						ErrorDescription: fmt.Sprintf("unknown action %s", capability.Type),
					}
				}
				aliceDevice.Capabilities = append(aliceDevice.Capabilities, alice.CapabilityResponse{
					Type: capability.Type,
					State: alice.CapabilityResponseState{
						Instance:     capability.State.Instance,
						ActionResult: actionResult,
					},
				})
			}
			aliceDevices.Devices = append(aliceDevices.Devices, aliceDevice)
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
