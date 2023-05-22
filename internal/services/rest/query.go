package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/mappers"
	"mikrotik-alice-gateway/internal/models/alice"
	"mikrotik-alice-gateway/pkg/middleware/user"
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
	devices := s.deviceProvider.Routers(user.User(ctx))

	aliceDevices := alice.Devices{
		UserID:  user.User(ctx),
		Devices: make([]alice.Device, 0, len(devices)),
	}
	for _, reqDev := range req.Devices {
		for _, dev := range devices {
			if dev.ID != mappers.ExtractDeviceID(reqDev.ID) {
				continue
			}
			for _, host := range dev.Hosts {
				if reqDev.ID != mappers.CreateHostDeviceID(dev.ID, host.ID) {
					continue
				}
				aliceDevices.Devices = append(aliceDevices.Devices, mappers.DeviceHostToAlice(dev, host))
				break
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
