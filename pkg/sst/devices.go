package sst

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type DeviceType int

const (
	MCS300 DeviceType = iota
	MCS350
	NeptunProWWiFi
	OKElectro
	EquationProWWiFi
	ThermoregulatorEcoSmart25
)

func (d DeviceType) String() string {
	switch d {
	case MCS300:
		return "MCS 300"
	case MCS350:
		return "MCS 350"
	case NeptunProWWiFi:
		return "Neptun ProW+WiFi"
	case OKElectro:
		return "OK Electro"
	case EquationProWWiFi:
		return "Equation ProW+WiFi"
	case ThermoregulatorEcoSmart25:
		return "Thermoregulator EcoSmart 25"
	}
	return "Unknown device"
}

type Device struct {
	ActiveNetwork              int       `json:"active_network"`
	ChartTemperatureComfort    int       `json:"chart_temperature_comfort"`
	ChartTemperatureEconomical int       `json:"chart_temperature_economical"`
	Configuration              string    `json:"configuration"`
	CreatedAt                  time.Time `json:"created_at"`
	House                      int       `json:"house"`
	ID                         int       `json:"id"`
	IsActive                   bool      `json:"is_active"`
	IsConnected                bool      `json:"is_connected"`
	LineNames                  []string  `json:"line_names"`
	LinesEnabled               []bool    `json:"lines_enabled"`
	MacAddress                 string    `json:"mac_address"`
	Name                       string    `json:"name"`
	ParsedConfiguration        string    `json:"parsed_configuration"`
	Power                      int       `json:"power"`
	PowerRelayTime             string    `json:"power_relay_time"`
	PreviousMode               string    `json:"previous_mode"`
	SpecificSettings           struct{}  `json:"specific_settings"`
	TimeSetting                struct {
		Device            int        `json:"device"`
		ID                int        `json:"id"`
		VacationTimeRange [][]string `json:"vacation_time_range"`
		WorkdayTimeRange  [][]string `json:"workday_time_range"`
	} `json:"time_setting"`
	Timeout                 int        `json:"timeout"`
	Type                    DeviceType `json:"type"`
	UpdatedAt               time.Time  `json:"updated_at"`
	WirelessSensorsNames    []string
	TermParsedConfiguration *DeviceTermParsedConfiguration `json:"-"`
}

type DeviceMode string

const (
	DeviceModeManual DeviceMode = "manual"
)

type DeviceStatus string

const (
	DeviceStatusOn  DeviceStatus = "on"
	DeviceStatusOff DeviceStatus = "off"
)

type DeviceStatusSelect string

const (
	DeviceStatusSelected   DeviceStatusSelect = "selected"
	DeviceStatusUnselected DeviceStatusSelect = "unselected"
)

type DeviceStatusAccess string

const (
	DeviceStatusAvailable DeviceStatusAccess = "available"
)

type DeviceTermParsedConfiguration struct {
	Detector int `json:"detector"`
	Settings struct {
		Mode         DeviceMode   `json:"mode"`
		Status       DeviceStatus `json:"status"`
		SelfTraining struct {
			Air        DeviceStatusSelect `json:"air"`
			Floor      DeviceStatusSelect `json:"floor"`
			Status     DeviceStatus       `json:"status"`
			OpenWindow DeviceStatusSelect `json:"open_window"`
		} `json:"self_training"`
		TemperatureAir           int `json:"temperature_air"`
		TemperatureManual        int `json:"temperature_manual"`
		TemperatureVacation      int `json:"temperature_vacation"`
		TemperatureCorrectionAir int `json:"temperature_correction_air"`
	} `json:"settings"`
	DeviceID           string             `json:"device_id"`
	MacAddress         string             `json:"mac_address"`
	RelayStatus        DeviceStatusSelect `json:"relay_status"`
	SignalLevel        int                `json:"signal_level"`
	AccessStatus       DeviceStatusAccess `json:"access_status"`
	CurrentTemperature struct {
		Event            int `json:"event"`
		DayOfWeek        int `json:"day_of_week"`
		TemperatureAir   int `json:"temperature_air"`
		TemperatureFloor int `json:"temperature_floor"`
	} `json:"current_temperature"`
	OpenWindowMinutes int `json:"open_window_minutes"`
}

func (c *Client) Devices(ctx context.Context, house int) ([]Device, error) {
	var result []Device
	if err := c.sendRequest(ctx, http.MethodGet, fmt.Sprintf("/houses/%d/devices/", house), nil, &result); err != nil {
		return nil, err
	}
	for i := range result {
		if result[i].ParsedConfiguration == "" {
			continue
		}
		switch result[i].Type {
		case MCS350, MCS300:
			var parsed DeviceTermParsedConfiguration
			if err := json.Unmarshal([]byte(result[i].ParsedConfiguration), &parsed); err != nil {
				log.Ctx(ctx).Error().Err(err).Str("configuration", result[i].ParsedConfiguration).Msg("Failed parse additional configuration")
				return nil, err
			}
			result[i].TermParsedConfiguration = &parsed
		}
	}
	return result, nil
}

func (c *Client) Temperature(ctx context.Context, house, device, temperature int) error {
	return c.sendRequest(ctx, http.MethodPost, fmt.Sprintf("/houses/%d/devices/%d/temperature/", house, device), struct {
		TemperatureManual int `json:"temperature_manual"`
	}{
		TemperatureManual: temperature,
	}, nil)
}

func (c *Client) PowerStatus(ctx context.Context, house, device int, power bool) error {
	status := DeviceStatusOff
	if power {
		status = DeviceStatusOn
	}
	return c.sendRequest(ctx, http.MethodPost, fmt.Sprintf("/houses/%d/devices/%d/status/", house, device), struct {
		Status DeviceStatus `json:"status"`
	}{
		Status: status,
	}, nil)
}
