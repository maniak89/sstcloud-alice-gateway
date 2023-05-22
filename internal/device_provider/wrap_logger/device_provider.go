package wrap_logger

import (
	"context"
	"strconv"
	"sync"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/models/storage"
)

type wrapper struct {
	child   device_provider.DeviceProvider
	logger  Logger
	linkID  string
	userID  string
	isInit  bool
	isInitM sync.Mutex
}

type Logger interface {
	Log(ctx context.Context, linkID string, level storage.LogLevel, msg string)
}

func New(child device_provider.DeviceProvider, userID, linkID string, logger Logger) device_provider.DeviceProvider {
	return &wrapper{
		child:  child,
		logger: logger,
		linkID: linkID,
		userID: userID,
	}
}

func (w *wrapper) insure(ctx context.Context) error {
	w.isInitM.Lock()
	defer w.isInitM.Unlock()
	if w.isInit {
		return nil
	}
	if err := w.child.Init(ctx); err != nil {
		w.logger.Log(ctx, w.linkID, storage.Error, err.Error())
		return err
	}
	w.logger.Log(ctx, w.linkID, storage.Info, "Success connected")
	w.isInit = true
	return nil
}

func (w *wrapper) Init(ctx context.Context) error {
	return w.insure(ctx)
}

func (w *wrapper) Houses(ctx context.Context) ([]*device_provider.House, error) {
	if err := w.insure(ctx); err != nil {
		return nil, err
	}
	result, err := w.child.Houses(ctx)
	if err != nil {
		w.logger.Log(ctx, w.linkID, storage.Error, "Failed get houses: "+err.Error())
		return nil, err
	}
	for _, r := range result {
		r.DeviceProvider = w
		r.UserID = w.userID
	}
	return result, nil
}
func (w *wrapper) Devices(ctx context.Context, house *device_provider.House) ([]*device_provider.Device, error) {
	if err := w.insure(ctx); err != nil {
		return nil, err
	}
	result, err := w.child.Devices(ctx, house)
	if err != nil {
		w.logger.Log(ctx, w.linkID, storage.Error, "Failed get devices: "+err.Error())
		return nil, err
	}
	w.logger.Log(ctx, w.linkID, storage.Info, "Success get devices. Total "+strconv.Itoa(len(result)))
	return result, nil
}

func (w *wrapper) SetTemperature(ctx context.Context, device *device_provider.Device, temp int) error {
	if err := w.insure(ctx); err != nil {
		return err
	}
	if err := w.child.SetTemperature(ctx, device, temp); err != nil {
		w.logger.Log(ctx, w.linkID, storage.Error, "Failed set temp: "+err.Error())
		return err
	}
	w.logger.Log(ctx, w.linkID, storage.Info, "Success set temp on device "+device.String()+" to "+strconv.Itoa(temp))
	return nil
}

func (w *wrapper) PowerStatus(ctx context.Context, device *device_provider.Device, power bool) error {
	if err := w.insure(ctx); err != nil {
		return err
	}
	if err := w.child.PowerStatus(ctx, device, power); err != nil {
		w.logger.Log(ctx, w.linkID, storage.Error, "Failed set power status: "+err.Error())
		return err
	}
	w.logger.Log(ctx, w.linkID, storage.Info, "Success set power status on device "+device.String()+" to "+strconv.FormatBool(power))
	return nil
}
