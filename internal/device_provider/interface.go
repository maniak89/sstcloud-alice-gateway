package device_provider

import (
	"context"

	"sstcloud-alice-gateway/internal/models/common"
)

type DeviceProvider interface {
	Init(ctx context.Context) error
	Devices(ctx context.Context) ([]common.Device, error)
	SetTemperature(ctx context.Context, device common.Device, temp int) error
	PowerStatus(ctx context.Context, device common.Device, power bool) error
	EMail() string
	Password() string
}
