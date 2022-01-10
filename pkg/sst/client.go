package sst

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Client struct {
	cl     *http.Client
	config Config

	token *string
}

type Config struct {
	URL     string        `env:"SST_URL,default=https://api.sst-cloud.com"`
	Timeout time.Duration `env:"SST_TIMEOUT,default=5s"`
}

func New(config Config) *Client {
	return &Client{
		config: config,
		cl: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *Client) sendRequest(ctx context.Context, method, uri string, in, out interface{}) error {
	logger := zerolog.Ctx(ctx).With().Str("method", method).Str("uri", uri).Logger()
	var body io.Reader
	if in != nil {
		blob, err := json.Marshal(in)
		if err != nil {
			logger.Error().Err(err).Msg("Failed marshal request body")
			return err
		}
		body = bytes.NewReader(blob)
	}

	req, err := http.NewRequest(method, c.config.URL+uri, body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed create request object")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != nil {
		req.Header.Set("Authorization", "Token "+*c.token)
	}

	resp, err := c.cl.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed make request")
		return err
	}

	if resp.StatusCode >= 400 {
		blobResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Failed read response")
			return err
		}
		err = errors.New(string(blobResponse))
		logger.Error().Err(err).Msg("Error response")
		return err
	}

	if out == nil {
		return nil
	}

	blobResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed read response")
		return err
	}

	if err := json.Unmarshal(blobResponse, out); err != nil {
		logger.Error().Err(err).Bytes("response", blobResponse).Msg("Failed decode response")
		return err
	}

	return nil
}
