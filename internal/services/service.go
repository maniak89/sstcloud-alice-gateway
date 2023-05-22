package services

import (
	"context"
)

type Service interface {
	Run(ctx context.Context, ready func()) error
	Shutdown(ctx context.Context) error
}
