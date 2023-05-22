package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/mappers"
	"mikrotik-alice-gateway/internal/models/alice"
	"mikrotik-alice-gateway/pkg/middleware/user"
)

func (s *service) Devices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)

	devices := s.deviceProvider.Routers(user.User(ctx))

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
