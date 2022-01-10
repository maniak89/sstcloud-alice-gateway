package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/maniak89/sstcloud-alice-gateway/internal/log"
	"github.com/oklog/run"
	zerolog "github.com/rs/zerolog/log"

	"github.com/maniak89/sstcloud-alice-gateway/internal/sst"
)

type config struct {
	Logger log.Config
	SST    sst.Config
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

	sstClient := sst.New(cfg.SST)
	if err := sstClient.Init(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed init sst client")
	}

	sstClient.Devices(ctx)

	logger.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		logger.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	logger.Info().Msg("The service is stopped")
}
