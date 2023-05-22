package sst

import (
	"context"
	"net/http"
	"time"
)

type Behaviour string

const (
	BehaviourVAC = "VAC"
)

type House struct {
	ID          int       `json:"id"`
	InHome      bool      `json:"in_home"`
	Name        string    `json:"name"`
	ReportDate  int       `json:"report_date"`
	Timezone    string    `json:"timezone"`
	UID         string    `json:"uid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Behaviour   Behaviour `json:"behaviour"`
	CloseValves int       `json:"close_valves"`
	Users       []int     `json:"users"`
	Workdays    struct {
		CurrentDay  string `json:"current_day"`
		CurrentWeek int    `json:"current_week"`
		House       int    `json:"house"`
		ID          int    `json:"id"`
		IsCustom    bool   `json:"is_custom"`
		// etc....
	} `json:"workdays"`
}

func (c *Client) Houses(ctx context.Context) ([]House, error) {
	var result []House
	if err := c.sendRequest(ctx, http.MethodGet, "/houses/", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
