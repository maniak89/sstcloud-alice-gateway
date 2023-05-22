package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"sstcloud-alice-gateway/internal/device_provider"
	"sstcloud-alice-gateway/internal/storage"
	"sstcloud-alice-gateway/pkg/middleware/user"
)

type service struct {
	config         Config
	srv            *http.Server
	storage        storage.Storage
	deviceProvider DeviceProvider
}

type DeviceProvider interface {
	Devices(userID string) []*device_provider.Device
}

const xRequestID = "X-Request-Id"

func New(ctx context.Context, config Config, log zerolog.Logger, storage storage.Storage, deviceProvider DeviceProvider) (*service, error) {
	r := chi.NewRouter()
	r.Use(
		hlog.NewHandler(log),
		hlog.MethodHandler("method"),
		hlog.URLHandler("url"),
		hlog.RequestIDHandler("x_request_id", xRequestID),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			zerolog.Ctx(ctx).Trace().Str("method", r.Method).Str("url", r.URL.String()).Str("x_request_id", r.Header.Get(xRequestID)).Int("status", status).Int("size", size).Dur("duration", duration).Msg("request processed")
		}),
		middleware.Recoverer,
		user.Middleware,
	)
	service := service{
		config:         config,
		deviceProvider: deviceProvider,
		srv:            &http.Server{Addr: config.Address, Handler: r},
		storage:        storage,
	}

	r.Route("/v1.0", func(r chi.Router) {
		r.Head("/", service.Health)
		r.Route("/user/devices", func(r chi.Router) {
			r.Get("/", service.Devices)
			r.Post("/query", service.Query)
			r.Post("/action", service.Action)
		})
	})

	return &service, nil
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
