package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/maniak89/sstcloud-alice-gateway/internal/models/common"
	"github.com/maniak89/sstcloud-alice-gateway/internal/services/rest/handlers/oauth2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

type service struct {
	config         Config
	srv            *http.Server
	deviceProvider DeviceProvider
}

type DeviceProvider interface {
	Devices(ctx context.Context) ([]common.Device, error)
}

const xRequestID = "X-Request-Id"

func New(ctx context.Context, config Config, log zerolog.Logger, provider DeviceProvider) (*service, error) {
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
	)
	service := service{
		config:         config,
		deviceProvider: provider,
		srv:            &http.Server{Addr: config.Address, Handler: r},
	}
	oauthSrv := oauth2.New(config.OAUTH2)
	if err := oauthSrv.Init(ctx); err != nil {
		return nil, err
	}
	if config.OAUTH2.Enabled {
		r.Route("/oauth2", func(r chi.Router) {
			r.Mount("/authorize", http.HandlerFunc(oauthSrv.Authorize))
			r.Mount("/token", http.HandlerFunc(oauthSrv.Token))
		})
	}
	key, err := config.OAUTH2.GetAuthKey(ctx)
	if err != nil {
		return nil, err
	}
	r.Route("/v1.0", func(r chi.Router) {
		if config.OAUTH2.Enabled {
			r.Use(oauthSrv.Verify)
		} else {
			r.Use(jwtauth.Verifier(jwtauth.New(config.OAUTH2.AuthAlgo, nil, key)))
		}
		r.Use(
			jwtauth.Authenticator,
		)
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
