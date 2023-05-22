package device_provider

import (
	"context"
	"fmt"
	"time"
)

type Device struct {
	House            *House
	ID               int
	IDStr            string
	Name             string
	Model            string
	Enabled          bool
	Connected        bool
	Tempometer       Tempometer
	AdditionalFields map[string]string
	UpdatedAt        time.Time
}

func (d *Device) String() string {
	return fmt.Sprintf("%s (%s %s)", d.Name, d.Model, d.ID)
}

type Tempometer struct {
	SetDegreesFloor          int
	ChangedAtSetDegreesFloor time.Time
	DegreesFloor             int
	ChangedAtDegreesFloor    time.Time
	DegreesAir               int
	ChangedAtDegreesAir      time.Time
}

func (d *Device) SetTemperature(ctx context.Context, temp int) error {
	return d.House.DeviceProvider.SetTemperature(ctx, d, temp)
}

func (d *Device) PowerStatus(ctx context.Context, power bool) error {
	return d.House.DeviceProvider.PowerStatus(ctx, d, power)
}
