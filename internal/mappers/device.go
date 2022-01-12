package mappers

import (
	"github.com/maniak89/sstcloud-alice-gateway/internal/models/alice"
	"github.com/maniak89/sstcloud-alice-gateway/internal/models/common"
)

const (
	minTemp = 12
	maxTemp = 45
)

func DeviceToAlice(device common.Device) []alice.Device {
	result := alice.Device{
		ID:   device.ID,
		Name: device.Name,
		DeviceInfo: &alice.DeviceInfo{
			Model: device.Model,
		},
		CustomData: device.AdditionalFields,
	}
	var additionalDevices []alice.Device
	if device.Tempometer != nil {
		result.Type = alice.DeviceTypeThermostat
		result.Capabilities = []interface{}{
			alice.CapabilityOnOff{
				Type:        alice.CapabilityTypeOnOff,
				Retrievable: true,
				Parameters: alice.CapabilityOnOffParameters{
					Split: false,
				},
				State: alice.CapabilityOnOffState{
					Instance: alice.CapabilityOnOffInstanceOn,
					Value:    device.Enabled,
				},
			},
			alice.CapabilityRange{
				Type:        alice.CapabilityTypeRange,
				Retrievable: true,
				Parameters: alice.CapabilityRangeParametersTemperature{
					Instance:     alice.CapabilityRangeInstanceTemperature,
					Unit:         alice.UnitCelsius,
					RandomAccess: true,
					Range: alice.CapabilityRangeParametersRange{
						Max:       maxTemp,
						Min:       minTemp,
						Precision: 1,
					},
				},
				State: alice.CapabilityRangeStateTemperature{
					Instance: alice.CapabilityRangeInstanceTemperature,
					Value:    float32(device.Tempometer.SetDegreesFloor),
				},
			},
		}

		additionalDevices = append(additionalDevices, alice.Device{
			ID:   device.ID + "_air",
			Name: device.Name + " температура воздуха",
			DeviceInfo: &alice.DeviceInfo{
				Model: device.Model,
			},
			CustomData: device.AdditionalFields,
			Type:       alice.DeviceTypeSensor,
			Properties: []interface{}{
				alice.PropertiesFloat{
					Type:        alice.PropertiesTypeFloat,
					Retrievable: true,
					Parameters: alice.PropertiesFloatParametersTemperature{
						Instance: alice.PropertiesFloatParametersInstanceTemperature,
						Unit:     alice.UnitCelsius,
					},
					State: alice.PropertiesFloatState{
						Instance: alice.PropertiesFloatParametersInstanceTemperature,
						Value:    float32(device.Tempometer.DegreesAir),
					},
				},
			},
		}, alice.Device{
			ID:   device.ID + "_floor",
			Name: device.Name + " температура пола",
			DeviceInfo: &alice.DeviceInfo{
				Model: device.Model,
			},
			CustomData: device.AdditionalFields,
			Type:       alice.DeviceTypeSensor,
			Properties: []interface{}{
				alice.PropertiesFloat{
					Type:        alice.PropertiesTypeFloat,
					Retrievable: true,
					Parameters: alice.PropertiesFloatParametersTemperature{
						Instance: alice.PropertiesFloatParametersInstanceTemperature,
						Unit:     alice.UnitCelsius,
					},
					State: alice.PropertiesFloatState{
						Instance: alice.PropertiesFloatParametersInstanceTemperature,
						Value:    float32(device.Tempometer.DegreesFloor),
					},
				},
			},
		})
	}
	return append([]alice.Device{result}, additionalDevices...)
}
