package mappers

import (
	"strings"

	"sstcloud-alice-gateway/internal/models/alice"
	"sstcloud-alice-gateway/internal/models/common"
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
		}

		additionalDevices = append(additionalDevices, alice.Device{
			ID:   createDeviceID(device.ID, AdditionalSensorAir),
			Name: device.Name + " температура воздуха",
			DeviceInfo: &alice.DeviceInfo{
				Model: device.Model,
			},
			CustomData: mapMux(device.AdditionalFields, map[string]string{
				AdditionalSensor: AdditionalSensorAir,
			}),
			Type: alice.DeviceTypeSensor,
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
			ID:   createDeviceID(device.ID, AdditionalSensorFloor),
			Name: device.Name + " температура пола",
			DeviceInfo: &alice.DeviceInfo{
				Model: device.Model,
			},
			CustomData: mapMux(device.AdditionalFields, map[string]string{
				AdditionalSensor: AdditionalSensorFloor,
			}),
			Type: alice.DeviceTypeSensor,
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

func ExtractDeviceID(str string) string {
	parts := strings.Split(str, "_")
	if len(parts) < 2 {
		return str
	}
	lastPart := parts[len(parts)-1]
	if lastPart == AdditionalSensorFloor || lastPart == AdditionalSensorAir {
		return strings.Join(parts[0:len(parts)-1], "_")
	}
	return str
}

func createDeviceID(str1, str2 string) string {
	return strings.Join([]string{str1, str2}, "_")
}

func DeviceFromAlice(device alice.DeviceRequest) common.Device {
	return common.Device{
		ID:               device.ID,
		AdditionalFields: device.CustomData,
	}
}
