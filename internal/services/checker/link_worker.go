package checker

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	storageModels "sstcloud-alice-gateway/internal/models/storage"
	"sstcloud-alice-gateway/internal/notifier"
)

type linkWorker struct {
	config     Config
	provider   device_provider.DeviceProvider
	link       *storageModels.Link
	notifier   notifier.Notifier
	workerMap  map[int]*houseWorker
	workerMapM sync.Mutex
	wg         sync.WaitGroup
	cancelFunc context.CancelFunc
}

func newLinkWorker(config Config, provider device_provider.DeviceProvider, link *storageModels.Link, notifier notifier.Notifier) *linkWorker {
	result := linkWorker{
		config:    config,
		provider:  provider,
		notifier:  notifier,
		link:      link,
		workerMap: map[int]*houseWorker{},
	}
	return &result
}

func (w *linkWorker) run(ctx context.Context) {
	logger := log.Ctx(ctx).With().Str("link_id", w.link.ID).Logger()
	ctx = logger.WithContext(ctx)
	ctx, w.cancelFunc = context.WithCancel(ctx)
	defer func() {
		w.cancelFunc()
		w.wg.Wait()
	}()

	r, err := w.provider.Houses(ctx)
	w.updateHouses(ctx, r, err)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(w.config.RequestPeriod):
			r, err := w.provider.Houses(ctx)
			w.updateHouses(ctx, r, err)
		}
	}
}

func (w *linkWorker) stop(ctx context.Context) {
	if w.cancelFunc != nil {
		w.cancelFunc()
	}
	w.workerMapM.Lock()
	defer w.workerMapM.Unlock()
	for _, c := range w.workerMap {
		c.stop(ctx)
	}
}

func (w *linkWorker) getState() []*device_provider.Device {
	var result []*device_provider.Device
	w.workerMapM.Lock()
	defer w.workerMapM.Unlock()
	for _, c := range w.workerMap {
		result = append(result, c.getState()...)
	}
	return result
}

func (w *linkWorker) markAllOffline(ctx context.Context, err error) {
	logger := log.Ctx(ctx)
	w.workerMapM.Lock()
	defer w.workerMapM.Unlock()
	logger.Error().Err(err).Msg("err state")
	for _, worker := range w.workerMap {
		worker.markAllOffline(ctx)
	}
	return
}

func (w *linkWorker) updateHouses(ctx context.Context, houses []*device_provider.House, err error) {
	logger := log.Ctx(ctx)
	if err != nil {
		w.markAllOffline(ctx, err)
		return
	}
	workerMap := make(map[int]*houseWorker)
	for _, house := range houses {
		logger := logger.With().Int("house_id", house.ID).Logger()
		ctx := logger.WithContext(ctx)
		worker, exists := w.workerMap[house.ID]
		if !exists {
			worker = newHouseWorker(w.config, w.provider, house, w.notifier)
			w.wg.Add(1)
			go func() {
				defer func() {
					w.workerMapM.Lock()
					delete(w.workerMap, worker.house.ID)
					defer func() {
						w.workerMapM.Unlock()
						w.wg.Done()
					}()
					w.stop(ctx)
				}()
				w.run(ctx)
			}()
		}
		workerMap[house.ID] = worker
	}
	w.workerMapM.Lock()
	defer w.workerMapM.Unlock()
	for k, v := range w.workerMap {
		if _, exists := workerMap[k]; exists {
			continue
		}
		v.markAllOffline(ctx)
		v.stop(ctx)
	}
}
