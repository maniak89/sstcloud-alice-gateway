package notifier

import (
	"context"

	"sstcloud-alice-gateway/internal/device_provider"
)

type Notifier interface {
	NotifyDevicesChanged(ctx context.Context, house *device_provider.House, device []*device_provider.Device) error
}
