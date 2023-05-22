package alice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/mappers"
	"mikrotik-alice-gateway/internal/models/alice"
	"mikrotik-alice-gateway/internal/models/common"
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

func (c *client) NotifyHostChanged(ctx context.Context, router *common.Router, host *common.Host) error {
	logger := log.Ctx(ctx)
	obj := mappers.DeviceHostToAliceEvent(router, host)
	blob, err := json.Marshal(alice.State{
		TS: time.Now().Unix(),
		Payload: alice.PayloadState{
			UserID:  router.UserID,
			Devices: []alice.PayloadStateDevice{obj},
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed marshal body")
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.callbackAddress, bytes.NewReader(blob))
	if err != nil {
		logger.Error().Err(err).Msg("Failed create request object")
		return err
	}
	req = req.WithContext(ctx)
	req.URL.Query().Set("Content-Type", "application/json")
	req.URL.Query().Set("Authorization", "Bearer "+c.config.OAuth2Token)
	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed make request")
		return err
	}
	logger.Debug().Str("status", resp.Status).Msg("status")
	return nil
}
