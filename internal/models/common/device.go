package common

import (
	"fmt"
)

type Device struct {
	ID               string
	Name             string
	Model            string
	Enabled          bool
	Connected        bool
	Tempometer       *Tempometer
	AdditionalFields map[string]string
}

func (d *Device) String() string {
	return fmt.Sprintf("%s (%s %s)", d.Name, d.Model, d.ID)
}

type Tempometer struct {
	SetDegreesFloor int
	DegreesFloor    int
	DegreesAir      int
}
