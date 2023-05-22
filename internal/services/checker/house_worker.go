package checker

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/notifier"
)

type houseWorker struct {
	config           Config
	provider         device_provider.DeviceProvider
	notifier         notifier.Notifier
	cancelFunc       context.CancelFunc
	state            []*device_provider.Device
	stateM           sync.Mutex
	house            *device_provider.House
	notifyCancelFunc context.CancelFunc
}

func newHouseWorker(config Config, provider device_provider.DeviceProvider, house *device_provider.House, notifier notifier.Notifier) *houseWorker {
	return &houseWorker{
		config:   config,
		provider: provider,
		notifier: notifier,
		house:    house,
	}
}

func (w *houseWorker) run(ctx context.Context) {
	logger := log.Ctx(ctx).With().Int("house_id", w.house.ID).Logger()
	ctx = logger.WithContext(ctx)
	ctx, w.cancelFunc = context.WithCancel(ctx)
	defer w.cancelFunc()

	r, err := w.provider.Devices(ctx, w.house)
	w.updateDevices(ctx, r, err)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(w.config.RequestPeriod):
			r, err := w.provider.Devices(ctx, w.house)
			w.updateDevices(ctx, r, err)
		}
	}
}

func (w *houseWorker) stop(ctx context.Context) {
	if w.cancelFunc != nil {
		w.cancelFunc()
	}
}

func (w *houseWorker) getState() []*device_provider.Device {
	w.stateM.Lock()
	defer w.stateM.Unlock()
	result := make([]*device_provider.Device, 0, len(w.state))
	for _, s := range w.state {
		result = append(result, s)
	}
	return result
}

func (w *houseWorker) updateDevices(ctx context.Context, devices []*device_provider.Device, err error) {
	if err != nil {
		return
	}
	savedDeviceMap := map[int]*device_provider.Device{}
	for _, device := range w.getState() {
		savedDeviceMap[device.ID] = device
	}
	var notify []*device_provider.Device
	for _, device := range devices {
		savedDevice, exists := savedDeviceMap[device.ID]
		if !exists {
			continue
		}
		changed := savedDevice.Enabled != device.Enabled ||
			savedDevice.Connected != device.Connected ||
			savedDevice.Name != device.Name ||
			savedDevice.House.Name != device.House.Name

		if savedDevice.Tempometer.DegreesFloor != device.Tempometer.DegreesFloor {
			changed = true
		} else {
			device.Tempometer.ChangedAtDegreesFloor = savedDevice.Tempometer.ChangedAtDegreesFloor
		}
		if savedDevice.Tempometer.DegreesAir != device.Tempometer.DegreesAir {
			changed = true
		} else {
			device.Tempometer.ChangedAtDegreesAir = savedDevice.Tempometer.ChangedAtDegreesAir
		}
		if savedDevice.Tempometer.SetDegreesFloor != device.Tempometer.SetDegreesFloor {
			changed = true
		} else {
			device.Tempometer.ChangedAtSetDegreesFloor = savedDevice.Tempometer.ChangedAtSetDegreesFloor
		}
		if changed {
			notify = append(notify, device)
		}
	}
	w.stateM.Lock()
	w.state = devices
	w.stateM.Unlock()
	w.notify(ctx, notify)
}

func (w *houseWorker) markAllOffline(ctx context.Context) {
	states := w.getState()
	notify := make([]*device_provider.Device, 0, len(states))
	for _, state := range states {
		if !state.Connected {
			continue
		}
		state.Connected = false
		notify = append(notify, state)
	}
	w.notify(ctx, notify)
}

func (w *houseWorker) notify(ctx context.Context, devices []*device_provider.Device) {
	logger := log.Ctx(ctx)
	if err := w.notifier.NotifyDevicesChanged(ctx, w.house, devices); err != nil {
		logger.Error().Err(err).Msg("Failed notify")
	}
	return
}
