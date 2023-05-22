package services

import (
	"context"
	"sync"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
)

type OrderRunner struct {
	prev chan struct{}
}

func (o *OrderRunner) SetupService(ctx context.Context, srv Service, role string, g *run.Group) error {
	logger := log.Ctx(ctx).With().Str("role", role).Logger()
	ctx = logger.WithContext(ctx)
	chPrev := o.prev
	chNext := make(chan struct{})
	var once sync.Once
	closer := func() {
		once.Do(func() {
			close(chNext)
		})
	}
	g.Add(func() error {
		if chPrev != nil {
			<-chPrev
		}
		logger.Info().Msg("Running the service...")
		defer logger.Info().Msg("˜the service is stopped")
		if err := ctx.Err(); err != nil {
			return nil
		}
		return srv.Run(ctx, closer)
	}, func(error) {
		closer()
		logger.Info().Msg("Shutdowning the service...")
		defer logger.Info().Msg("˜the service is shutdown")
		cCtx := logger.WithContext(context.Background())
		if err := srv.Shutdown(cCtx); err != nil {
			logger.Error().Err(err).Msg("Cannot shutdown the service properly")
		}
	})
	o.prev = chNext
	return nil
}
