package checker

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/notifier"
	"sstcloud-alice-gateway/internal/storage"
)

type DeviceFactory func(userID, linkID, username, password string) device_provider.DeviceProvider

type service struct {
	config        Config
	storage       storage.Storage
	deviceFactory DeviceFactory
	notifier      notifier.Notifier
	cancelFunc    context.CancelFunc
	wg            sync.WaitGroup
	workers       map[string]*linkWorker
	workersM      sync.Mutex
}

func New(config Config, storage storage.Storage, deviceFactory DeviceFactory, notifier notifier.Notifier) *service {
	return &service{
		config:        config,
		storage:       storage,
		deviceFactory: deviceFactory,
		notifier:      notifier,
		workers:       map[string]*linkWorker{},
	}
}

func (s *service) Run(ctx context.Context, ready func()) error {
	logger := log.Ctx(ctx).With().Str("role", "checker").Logger()
	ctx = logger.WithContext(ctx)
	ctx, s.cancelFunc = context.WithCancel(ctx)
	defer func() {
		s.cancelFunc()
		s.wg.Wait()
	}()
	if err := s.processUpdates(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed process updates")
		return err
	}
	ready()
	<-ctx.Done()
	return nil
}

func (s *service) Devices(userID string) []*device_provider.Device {
	s.workersM.Lock()
	defer s.workersM.Unlock()
	var result []*device_provider.Device
	for _, worker := range s.workers {
		if worker.link.UserID != userID {
			continue
		}
		result = append(result, worker.getState()...)
	}
	return result
}

func (s *service) processUpdates(ctx context.Context) error {
	logger := log.Ctx(ctx)
	links, err := s.storage.Links(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed fetch links")
		return err
	}
	s.workersM.Lock()
	defer s.workersM.Unlock()
	workers := make(map[string]*linkWorker, len(s.workers))
	for _, link := range links {
		worker, exist := s.workers[link.ID]
		if !exist || !worker.link.Equal(link) {
			if exist {
				worker.stop(ctx)
			}
			worker = newLinkWorker(s.config, s.deviceFactory(link.UserID, link.ID, link.SSTEmail, link.SSTPassword), link, s.notifier)
			s.wg.Add(1)
			go func() {
				defer func() {
					s.workersM.Lock()
					defer func() {
						s.workersM.Unlock()
						s.wg.Done()
					}()
					delete(s.workers, worker.link.ID)
				}()
				worker.run(ctx)
			}()
		}
		workers[link.ID] = worker
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
