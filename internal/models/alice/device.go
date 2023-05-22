package alice

import (
	"time"
)

type DeviceType string

const (
	DeviceTypeThermostat DeviceType = "devices.types.thermostat"
	DeviceTypeSensor     DeviceType = "devices.types.sensor"
)

type Device struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Room         string            `json:"room,omitempty"`
	Type         DeviceType        `json:"type"`
	CustomData   map[string]string `json:"custom_data,omitempty"`
	Capabilities []interface{}     `json:"capabilities,omitempty"`
	Properties   []Property        `json:"properties,omitempty"`
	DeviceInfo   *DeviceInfo       `json:"device_info,omitempty"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HWVersion    string `json:"hw_version,omitempty"`
	SWVersion    string `json:"sw_version,omitempty"`
}

type CapabilityType string

const (
	CapabilityTypeOnOff CapabilityType = "devices.capabilities.on_off"
	CapabilityTypeRange CapabilityType = "devices.capabilities.range"
)

type ActionResultStatus string

const (
	ActionResultStatusDone  ActionResultStatus = "DONE"
	ActionResultStatusError ActionResultStatus = "ERROR"
)

type ErrorCode string

const (
	ErrorCodeDeviceUnreachable ErrorCode = "DEVICE_UNREACHABLE"
	ErrorCodeInvalidAction     ErrorCode = "INVALID_ACTION"
	ErrorCodeInvalidValue      ErrorCode = "INVALID_VALUE"
)

type ActionResult struct {
	Status           ActionResultStatus `json:"status"`
	ErrorCode        ErrorCode          `json:"error_code,omitempty"`
	ErrorDescription string             `json:"error_description,omitempty"`
}

type CapabilityResponse struct {
	Type  CapabilityType          `json:"type"`
	State CapabilityResponseState `json:"state"`
}

type CapabilityResponseState struct {
	Instance     string       `json:"instance"`
	ActionResult ActionResult `json:"action_result"`
}

type CapabilityOnOff struct {
	Type        CapabilityType            `json:"type"`
	Retrievable bool                      `json:"retrievable"`
	Parameters  CapabilityOnOffParameters `json:"parameters"`
	State       CapabilityOnOffState      `json:"state"`
}

type CapabilityOnOffParameters struct {
	Split bool `json:"split"`
}

type CapabilityOnOffState struct {
	Instance CapabilityOnOffInstance `json:"instance"`
	Value    bool                    `json:"value"`
}

type CapabilityOnOffInstance string

const (
	CapabilityOnOffInstanceOn  CapabilityOnOffInstance = "on"
	CapabilityOnOffInstanceOff CapabilityOnOffInstance = "off"
)

type CapabilityRange struct {
	Type        CapabilityType `json:"type"`
	Retrievable bool           `json:"retrievable"`
	Parameters  interface{}    `json:"parameters"`
	State       interface{}    `json:"state"`
}

type CapabilityRangeStateTemperature struct {
	Instance CapabilityRangeInstance `json:"instance"`
	Value    float32                 `json:"value"`
}

type CapabilityRangeParametersTemperature struct {
	Instance     CapabilityRangeInstance        `json:"instance"`
	Unit         PropertyParameterUnit          `json:"unit"`
	RandomAccess bool                           `json:"random_access"`
	Range        CapabilityRangeParametersRange `json:"range"`
}

type CapabilityRangeInstance string

const (
	CapabilityRangeInstanceTemperature CapabilityRangeInstance = "temperature"
)

type CapabilityRangeParametersRange struct {
	Min       float32 `json:"min,omitempty"`
	Max       float32 `json:"max,omitempty"`
	Precision float32 `json:"precision,omitempty"`
}

type PropertyType string

const (
	PropertyTypeFloat PropertyType = "devices.properties.float"
)

type PropertyParameterValueValue string

const (
	PropertyParameterInstanceGasDetected    PropertyParameterValueValue = "detected"
	PropertyParameterInstanceGasNotDetected PropertyParameterValueValue = "not_detected"
)

type PropertyParameterValue struct {
	Value PropertyParameterValueValue `json:"value"`
	Name  string                      `json:"name"`
}

type PropertyParameterInstance string

const (
	PropertyParameterInstanceGas         = "gas"
	PropertyParameterInstanceTemperature = "temperature"
)

type PropertyParameterUnit string

const (
	PropertyParameterUnitUnknown PropertyParameterUnit = ""
	PropertyParameterUnitCelsius PropertyParameterUnit = "unit.celsius"
)

type PropertyParameter struct {
	Instance PropertyParameterInstance `json:"instance"`
	Events   []PropertyParameterValue  `json:"events,omitempty"`
	Unit     PropertyParameterUnit     `json:"unit,omitempty"`
}

type Property struct {
	Type           PropertyType                      `json:"type"`
	Retrievable    bool                              `json:"retrievable"`
	Reportable     bool                              `json:"reportable"`
	Parameters     PropertyParameter                 `json:"parameters"`
	State          PayloadStateDevicePropertiesState `json:"state"`
	StateChangedAt time.Time                         `json:"state_changed_at"`
	LastUpdated    time.Time                         `json:"last_updated"`
}

type DeviceRequest struct {
	ID           string            `json:"id"`
	CustomData   map[string]string `json:"custom_data"`
	Capabilities []struct {
		Type  CapabilityType `json:"type"`
		State struct {
			Instance string      `json:"instance"`
			Value    interface{} `json:"value"`
			Relative bool        `json:"relative"`
		}
	} `json:"capabilities"`
}
