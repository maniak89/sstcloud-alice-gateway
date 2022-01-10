package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/maniak89/sstcloud-alice-gateway/internal/models"
)

type service struct {
	config         Config
	srv            *http.Server
	deviceProvider DeviceProvider
}

type DeviceProvider interface {
	Devices(ctx context.Context) ([]models.Device, error)
}

func New(config Config, log zerolog.Logger, provider DeviceProvider) *service {
	r := chi.NewRouter()
	r.Use(
		hlog.NewHandler(log),
		middleware.Recoverer,
	)
	service := service{
		config:         config,
		deviceProvider: provider,
		srv:            &http.Server{Addr: config.Address, Handler: r},
	}
	r.Route("/v1.0", func(r chi.Router) {
		r.Head("/", service.Health)
	})
	return &service
}

func (s *service) Run(ctx context.Context, ready func()) error {
	logger := log.Ctx(ctx)
	logger.Info().Str("address", s.srv.Addr).Msg("Start listening")
	defer func() {
		logger.Info().Msg("Stop listening")
	}()
	ready()
	if err := s.srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		logger.Error().Err(err).Msg("Failed start listening")
		return err
	}

	return nil
}

func (s *service) Shutdown(ctx context.Context) error {
	logger := log.Ctx(ctx)

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed shutdown")
		return err
	}

	return nil
}
