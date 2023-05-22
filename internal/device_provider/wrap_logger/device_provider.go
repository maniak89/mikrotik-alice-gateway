package wrap_logger

import (
	"context"
	"strconv"
	"sync"

	"mikrotik-alice-gateway/internal/device_provider"
	"mikrotik-alice-gateway/internal/models/storage"
)

type wrapper struct {
	child    device_provider.DeviceProvider
	logger   Logger
	routerID string
	isInit   bool
	isInitM  sync.Mutex
}

type Logger interface {
	Log(ctx context.Context, routerID string, level storage.LogLevel, msg string)
}

func New(child device_provider.DeviceProvider, routerID string, logger Logger) device_provider.DeviceProvider {
	return &wrapper{
		child:    child,
		logger:   logger,
		routerID: routerID,
	}
}

func (w *wrapper) insure(ctx context.Context) error {
	w.isInitM.Lock()
	defer w.isInitM.Unlock()
	if w.isInit {
		return nil
	}
	if err := w.child.Init(ctx); err != nil {
		w.logger.Log(ctx, w.routerID, storage.Error, err.Error())
		return err
	}
	w.logger.Log(ctx, w.routerID, storage.Info, "Success connected")
	w.isInit = true
	return nil
}

func (w *wrapper) Init(ctx context.Context) error {
	return w.insure(ctx)
}

func (w *wrapper) Leases(ctx context.Context) ([]device_provider.Lease, error) {
	if err := w.insure(ctx); err != nil {
		return nil, err
	}
	result, err := w.child.Leases(ctx)
	if err != nil {
		w.logger.Log(ctx, w.routerID, storage.Error, "Failed get leases: "+err.Error())
		return nil, err
	}
	w.logger.Log(ctx, w.routerID, storage.Info, "Success get leases. Total "+strconv.Itoa(len(result)))
	return result, nil
}

func (w *wrapper) Resource(ctx context.Context) (device_provider.Resource, error) {
	if err := w.insure(ctx); err != nil {
		return device_provider.Resource{}, err
	}
	result, err := w.child.Resource(ctx)
	if err != nil {
		w.logger.Log(ctx, w.routerID, storage.Error, "Failed get resource: "+err.Error())
		return device_provider.Resource{}, err
	}
	w.logger.Log(ctx, w.routerID, storage.Info, "Success get resource")
	return result, nil
}
