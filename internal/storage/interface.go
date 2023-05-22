package storage

import (
	"context"
	"errors"

	"sstcloud-alice-gateway/internal/models/storage"
)

var ErrInvalidState = errors.New("invalid state")

type Storage interface {
	Links(ctx context.Context, userID string) ([]*storage.Link, error)
	Log(ctx context.Context, linkID string, level storage.LogLevel, msg string)
}
