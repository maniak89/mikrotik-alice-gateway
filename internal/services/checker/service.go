package checker

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/device_provider"
	"mikrotik-alice-gateway/internal/models/common"
	"mikrotik-alice-gateway/internal/notifier"
	"mikrotik-alice-gateway/internal/storage"
)

type RouterFactory func(routerID, address, username, password string) device_provider.DeviceProvider

type service struct {
	config        Config
	storage       storage.Storage
	routerFactory RouterFactory
	notifier      notifier.Notifier
	cancelFunc    context.CancelFunc
	wg            sync.WaitGroup
	workers       map[string]*worker
	workersM      sync.Mutex
}

func New(config Config, storage storage.Storage, routerFactory RouterFactory, notifier notifier.Notifier) *service {
	return &service{
		config:        config,
		storage:       storage,
		routerFactory: routerFactory,
		notifier:      notifier,
		workers:       map[string]*worker{},
	}
}

func (s *service) Run(ctx context.Context, ready func()) error {
	logger := log.Ctx(ctx).With().Str("role", "checker").Logger()
	ctx = logger.WithContext(ctx)
	ctx, s.cancelFunc = context.WithCancel(ctx)
	defer s.cancelFunc()
	if err := s.processUpdates(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed process updates")
		return err
	}
	ready()
	<-ctx.Done()
	s.wg.Wait()
	return nil
}

func (s *service) Routers(userID string) []*common.Router {
	s.workersM.Lock()
	defer s.workersM.Unlock()
	var result []*common.Router
	for _, worker := range s.workers {
		if worker.storageRouter.UserID != userID {
			continue
		}
		result = append(result, &worker.stateRouter)
	}
	return result
}

func (s *service) processUpdates(ctx context.Context) error {
	logger := log.Ctx(ctx)
	routers, err := s.storage.Routers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed fetch routers")
		return err
	}
	s.workersM.Lock()
	defer s.workersM.Unlock()
	workers := make(map[string]*worker, len(s.workers))
	for _, router := range routers {
		worker, exist := s.workers[router.ID]
		if !exist || !worker.storageRouter.Equal(router) {
			if exist {
				worker.stop(ctx)
			}
			worker = newWorker(s.config, s.routerFactory(router.ID, router.Address, router.Username, router.Password), router, s.storage, s.notifier)
			s.wg.Add(1)
			go func() {
				defer func() {
					s.workersM.Lock()
					delete(s.workers, worker.storageRouter.ID)
					s.workersM.Unlock()
					s.wg.Done()
				}()
				worker.run(ctx)
			}()
		}
		workers[router.ID] = worker
	}
	for workerID, worker := range s.workers {
		if _, exists := workers[workerID]; exists {
			continue
		}
		worker.stop(ctx)
	}
	s.workers = workers
	return nil
}

func (s *service) Shutdown(ctx context.Context) error {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	return nil
}
