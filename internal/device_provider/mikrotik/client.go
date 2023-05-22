package mikrotik

import (
	"context"
	"errors"
	"sync"

	"github.com/go-routeros/routeros"
	"github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/device_provider"
)

type client struct {
	config     Config
	client     *routeros.Client
	wg         sync.WaitGroup
	wCtx       context.Context
	cancelFunc context.CancelFunc
	lock       sync.Mutex
}

func New(config Config) *client {
	return &client{
		config: config,
	}
}

func (c *client) Init(ctx context.Context) error {
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Logger()
	cl, err := routeros.Dial(c.config.Address, c.config.Username, c.config.Password)
	if err != nil {
		logger.Error().Err(err).Msg("Failed connect to router")
		return err
	}
	c.client = cl
	c.wCtx, c.cancelFunc = context.WithCancel(ctx)
	return nil
}

func (c *client) Stop(ctx context.Context) {
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	c.wg.Wait()
}

func (c *client) Resource(ctx context.Context) (device_provider.Resource, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Logger()
	resp, err := c.client.Run("/system/resource/print")
	if err != nil {
		logger.Error().Err(err).Msg("Failed exec interface print")
		return device_provider.Resource{}, err
	}
	for _, r := range resp.Re {
		return device_provider.Resource{
			Version:   r.Map["version"],
			BoardName: r.Map["board-name"],
			Platform:  r.Map["platform"],
		}, nil
	}
	return device_provider.Resource{}, errors.New("ome wrong")
}

type IFace struct {
	ID       string
	Name     string
	Type     string
	Running  bool
	Disabled bool
}

func (c *client) Interfaces(ctx context.Context) ([]IFace, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Logger()
	resp, err := c.client.Run("/interface/print")
	if err != nil {
		logger.Error().Err(err).Msg("Failed exec interface print")
		return nil, err
	}
	result := make([]IFace, 0, len(resp.Re))
	for _, r := range resp.Re {
		result = append(result, IFace{
			ID:       r.Map[".id"],
			Type:     r.Map["type"],
			Name:     r.Map["name"],
			Running:  r.Map["running"] == "true",
			Disabled: r.Map["disabled"] == "true",
		})
	}
	logger.Debug().Interface("resp", resp).Msg("given")
	return result, nil
}

func (c *client) Leases(ctx context.Context) ([]device_provider.Lease, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Logger()
	resp, err := c.client.Run("/ip/dhcp-server/lease/print")
	if err != nil {
		logger.Error().Err(err).Msg("Failed exec interface print")
		return nil, err
	}
	result := make([]device_provider.Lease, 0, len(resp.Re))
	for _, r := range resp.Re {
		result = append(result, device_provider.Lease{
			Address:    r.Map["address"],
			Status:     device_provider.LeaseStatus(r.Map["status"]),
			HostName:   r.Map["host-name"],
			MacAddress: r.Map["mac-address"],
		})
	}
	logger.Debug().Interface("resp", resp).Msg("given")
	return result, nil
}

func (c *client) InterfaceListenRunning(ctx context.Context, name string) (<-chan bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Str("name", name).Logger()
	l, err := c.client.Listen("/interface/" + name + "/running/listen")
	if err != nil {
		logger.Error().Err(err).Msg("Failed execute interface running listen")
		return nil, err
	}
	result := make(chan bool)
	c.wg.Add(1)
	go func() {
		defer func() {
			c.wg.Done()
			close(result)
		}()
		for {
			select {
			case <-ctx.Done():
				logger.Debug().Err(ctx.Err()).Msg("Closed by user context")
				return
			case <-c.wCtx.Done():
				logger.Debug().Err(ctx.Err()).Msg("Closed by device context")
				return
			case reply := <-l.Chan():
				logger.Debug().Interface("reply", reply).Msg("get messages")
			}
		}
	}()
	return result, nil
}
