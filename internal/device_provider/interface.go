package device_provider

import (
	"context"
)

type DeviceProvider interface {
	Init(ctx context.Context) error
	Houses(ctx context.Context) ([]*House, error)
	Devices(ctx context.Context, house *House) ([]*Device, error)
	SetTemperature(ctx context.Context, device *Device, temp int) error
	PowerStatus(ctx context.Context, device *Device, power bool) error
}
