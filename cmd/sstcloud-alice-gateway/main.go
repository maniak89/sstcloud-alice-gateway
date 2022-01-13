package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/maniak89/sstcloud-alice-gateway/internal/log"
	"github.com/maniak89/sstcloud-alice-gateway/internal/services"
	"github.com/maniak89/sstcloud-alice-gateway/internal/services/rest"
	"github.com/maniak89/sstcloud-alice-gateway/internal/services/sst"
	"github.com/oklog/run"
	zerolog "github.com/rs/zerolog/log"
)

type config struct {
	Logger log.Config
	SST    sst.Config
	Rest   rest.Config
}

const signalChLen = 10

func main() {
	var cfg config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot decode config envs")
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot init logger")
	}

	ctx, cancel := context.WithCancel(logger.WithContext(context.Background()))

	g := &run.Group{}
	{
		stop := make(chan os.Signal, signalChLen)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}

	orderRunner := services.OrderRunner{}

	sstClient := sst.New(cfg.SST)
	if err := sstClient.Init(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed init sst client")
	}

	restService, err := rest.New(ctx, cfg.Rest, logger.With().Str("role", "rest").Logger(), sstClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create rest service")
	}
	if err := orderRunner.SetupService(ctx, restService, "rest", g); err != nil {
		logger.Fatal().Err(err).Msg("Failed setup rest service")
	}

	logger.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		logger.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	logger.Info().Msg("The service is stopped")
}
