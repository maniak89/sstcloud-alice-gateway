package mappers

import (
	"strings"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/models/alice"
)

const (
	MinTemp = 12
	MaxTemp = 45
)

const (
	AdditionalSensor      = "sensor"
	AdditionalSensorAir   = "air"
	AdditionalSensorFloor = "floor"
)

func DeviceToAlice(device *device_provider.Device) []alice.Device {
	result := []alice.Device{{
		ID:   device.IDStr,
		Name: device.Name,
		DeviceInfo: &alice.DeviceInfo{
			Model: device.Model,
		},
		CustomData: device.AdditionalFields,
		Type:       alice.DeviceTypeThermostat,
		Capabilities: []interface{}{
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
					Unit:         alice.PropertyParameterUnitCelsius,
					RandomAccess: true,
					Range: alice.CapabilityRangeParametersRange{
						Max:       MaxTemp,
						Min:       MinTemp,
						Precision: 1,
					},
				},
				State: alice.CapabilityRangeStateTemperature{
					Instance: alice.CapabilityRangeInstanceTemperature,
					Value:    float32(device.Tempometer.SetDegreesFloor),
				},
			},
		},
	},
		{
			ID:   createDeviceID(device.IDStr, AdditionalSensorAir),
			Name: device.Name + " температура воздуха",
			DeviceInfo: &alice.DeviceInfo{
				Model: device.Model,
			},
			CustomData: mapMux(device.AdditionalFields, map[string]string{
				AdditionalSensor: AdditionalSensorAir,
			}),
			Type: alice.DeviceTypeSensor,
			Properties: []alice.Property{
				{
					Type:        alice.PropertyTypeFloat,
					Retrievable: true,
					Reportable:  true,
					Parameters: alice.PropertyParameter{
						Instance: alice.PropertyParameterInstanceTemperature,
						Unit:     alice.PropertyParameterUnitCelsius,
					},
					State: alice.PayloadStateDevicePropertiesState{
						Instance: alice.PropertyParameterInstanceTemperature,
						Value:    device.Tempometer.DegreesAir,
					},
					LastUpdated:    device.UpdatedAt,
					StateChangedAt: device.Tempometer.ChangedAtDegreesAir,
				},
			},
		},
		{
			ID:   createDeviceID(device.IDStr, AdditionalSensorFloor),
			Name: device.Name + " температура пола",
			DeviceInfo: &alice.DeviceInfo{
				Model: device.Model,
			},
			CustomData: mapMux(device.AdditionalFields, map[string]string{
				AdditionalSensor: AdditionalSensorFloor,
			}),
			Type: alice.DeviceTypeSensor,
			Properties: []alice.Property{
				{
					Type:        alice.PropertyTypeFloat,
					Retrievable: true,
					Reportable:  true,
					Parameters: alice.PropertyParameter{
						Instance: alice.PropertyParameterInstanceTemperature,
						Unit:     alice.PropertyParameterUnitCelsius,
					},
					State: alice.PayloadStateDevicePropertiesState{
						Instance: alice.PropertyParameterInstanceTemperature,
						Value:    device.Tempometer.DegreesFloor,
					},
					LastUpdated:    device.UpdatedAt,
					StateChangedAt: device.Tempometer.ChangedAtDegreesFloor,
				},
			},
		},
	}
	return result
}

func mapMux(m1, m2 map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result
}

func createDeviceID(str1, str2 string) string {
	return strings.Join([]string{str1, str2}, "_")
}
