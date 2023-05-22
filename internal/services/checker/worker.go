package checker

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/device_provider"
	"mikrotik-alice-gateway/internal/models/common"
	storageModels "mikrotik-alice-gateway/internal/models/storage"
	"mikrotik-alice-gateway/internal/notifier"
	"mikrotik-alice-gateway/internal/storage"
)

type worker struct {
	config             Config
	provider           device_provider.DeviceProvider
	notifier           notifier.Notifier
	hostMap            map[string]*common.Host
	storageHostMap     map[string]*storageModels.Host
	storage            storage.Storage
	storageRouter      *storageModels.Router
	wg                 sync.WaitGroup
	cancelFunc         context.CancelFunc
	stateRouter        common.Router
	notifyCancelFuncs  map[string]context.CancelFunc
	notifyCancelFuncsM sync.Mutex
}

func newWorker(config Config, provider device_provider.DeviceProvider, router *storageModels.Router, storage storage.Storage, notifier notifier.Notifier) *worker {
	result := worker{
		config:            config,
		provider:          provider,
		notifier:          notifier,
		storage:           storage,
		hostMap:           make(map[string]*common.Host, len(router.Hosts)),
		storageHostMap:    make(map[string]*storageModels.Host, len(router.Hosts)),
		notifyCancelFuncs: make(map[string]context.CancelFunc, len(router.Hosts)),
		storageRouter:     router,
		stateRouter: common.Router{
			ID:     router.ID,
			UserID: router.UserID,
			Hosts:  make([]*common.Host, 0, len(router.Hosts)),
			Name:   router.Name,
		},
	}
	for _, storageHost := range router.Hosts {
		host := common.Host{
			ID:   storageHost.ID,
			Name: storageHost.Name,
		}
		result.hostMap[host.ID] = &host
		result.storageHostMap[host.ID] = storageHost
		result.stateRouter.Hosts = append(result.stateRouter.Hosts, &host)
	}
	return &result
}

func (w *worker) run(ctx context.Context) {
	logger := log.Ctx(ctx).With().Str("router_id", w.stateRouter.ID).Logger()
	ctx = logger.WithContext(ctx)
	ctx, w.cancelFunc = context.WithCancel(ctx)
	defer w.cancelFunc()
	w.wg.Add(2)
	go func() {
		defer w.wg.Done()
		r, err := w.provider.Resource(ctx)
		w.updateSystemInfo(ctx, r, err)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.config.SystemInfoCheckPeriod):
				r, err := w.provider.Resource(ctx)
				w.updateSystemInfo(ctx, r, err)
			}
		}
	}()
	go func() {
		defer w.wg.Done()
		r, err := w.provider.Leases(ctx)
		w.updateHosts(ctx, r, err)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.storageRouter.LeasePeriodCheck):
				r, err := w.provider.Leases(ctx)
				w.updateHosts(ctx, r, err)
			}
		}
	}()
	<-ctx.Done()
	w.wg.Wait()
}

func (w *worker) stop(ctx context.Context) {
	if w.cancelFunc != nil {
		w.cancelFunc()
	}
	w.notifyCancelFuncsM.Lock()
	defer w.notifyCancelFuncsM.Unlock()
	for _, c := range w.notifyCancelFuncs {
		c()
	}
}

func (w *worker) getRouter() *common.Router {
	return &w.stateRouter
}

func (w *worker) markAllOffline(ctx context.Context, err error) {
	logger := log.Ctx(ctx)
	logger.Error().Err(err).Msg("err state")
	w.stateRouter.Connected = false
	for _, h := range w.stateRouter.Hosts {
		storageHost := w.storageHostMap[h.ID]
		storageHost.IsOnline = false
		if err := w.storage.UpdateHost(ctx, storageHost); err != nil {
			logger.Error().Err(err).Msg("Failed update host")
		}
		if h.Online {
			h.Online = false
			w.notify(ctx, storageHost)
		}
	}
	return
}

func (w *worker) updateSystemInfo(ctx context.Context, resource device_provider.Resource, err error) {
	if err != nil {
		w.markAllOffline(ctx, err)
		return
	}
	w.stateRouter.SWVersion = resource.Version
	w.stateRouter.Manufacturer = resource.Platform
	w.stateRouter.Model = resource.BoardName
	w.stateRouter.Connected = true
}

func (w *worker) updateHosts(ctx context.Context, leases []device_provider.Lease, err error) {
	logger := log.Ctx(ctx)
	if err != nil {
		w.markAllOffline(ctx, err)
		return
	}
	for _, storageHost := range w.storageRouter.Hosts {
		logger := logger.With().Str("host_id", storageHost.ID).Logger()
		ctx := logger.WithContext(ctx)
		var found bool
		host := w.hostMap[storageHost.ID]
		for _, lease := range leases {
			if storageHost.Address.String == lease.Address ||
				storageHost.MacAddress.String == lease.MacAddress ||
				(storageHost.HostName.String == lease.HostName && lease.HostName != "") {
				found = true
				connected := lease.Status == device_provider.LeaseStatusBound
				storageHost.IsOnline = connected
				if connected {
					storageHost.LastOnline = time.Now()
				}
				if host.Online != connected {
					host.Online = connected
					w.notify(ctx, storageHost)
				}
				break
			}
		}
		if !found && host.Online {
			host.Online = false
			storageHost.IsOnline = false
			w.notify(ctx, storageHost)
		}
		if err := w.storage.UpdateHost(ctx, storageHost); err != nil {
			logger.Error().Err(err).Msg("Failed update host")
		}
	}
}

func (w *worker) notify(ctx context.Context, host *storageModels.Host) {
	logger := log.Ctx(ctx)
	w.notifyCancelFuncsM.Lock()
	defer w.notifyCancelFuncsM.Unlock()
	if host.IsOnline {
		if f, exists := w.notifyCancelFuncs[host.ID]; exists {
			f()
			delete(w.notifyCancelFuncs, host.ID)
		}
		if err := w.notifier.NotifyHostChanged(ctx, &w.stateRouter, w.hostMap[host.ID]); err != nil {
			logger.Error().Err(err).Msg("Failed notify")
		}
		return
	}
	w.wg.Add(1)
	wCtx, cancel := context.WithCancel(ctx)
	w.notifyCancelFuncs[host.ID] = cancel
	go func() {
		defer func() {
			w.notifyCancelFuncsM.Lock()
			defer w.notifyCancelFuncsM.Unlock()
			delete(w.notifyCancelFuncs, host.ID)
			w.wg.Done()
		}()
		select {
		case <-wCtx.Done():
		case <-time.After(host.OnlineTimeout):
			if err := w.notifier.NotifyHostChanged(ctx, &w.stateRouter, w.hostMap[host.ID]); err != nil {
				logger.Error().Err(err).Msg("Failed notify")
			}
		}
	}()
}
