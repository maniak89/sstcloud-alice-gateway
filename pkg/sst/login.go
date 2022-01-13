package sst

import (
	"context"
	"net/http"
)

type Language string

const (
	LangRu Language = "ru"
	LangEn Language = "en"
)

type LoginRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	EMail    string   `json:"email"`
	Language Language `json:"language"`
}

type LoginResponse struct {
	Key string `json:"key"`
}

func (c *Client) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	if request.Username == "" {
		request.Username = request.EMail
	}
	if request.Language == "" {
		request.Language = LangEn
	}
	var response LoginResponse
	if err := c.sendRequest(ctx, http.MethodPost, "/auth/login/", request, &response); err != nil {
		return nil, err
	}
	c.token = &response.Key
	return &response, nil
}
