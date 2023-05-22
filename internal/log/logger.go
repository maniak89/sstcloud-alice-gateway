package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func New(config Config) (zerolog.Logger, error) {

	var out io.Writer = os.Stdout

	if config.Console {
		out = zerolog.ConsoleWriter{Out: out}
	}

	logger := zerolog.New(out)

	if !config.DisableTimestamp {
		logger = logger.With().Timestamp().Logger()
	}

	if config.Caller {
		logger = logger.With().Caller().Logger()
	}

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		logger.Error().Err(err).Msg("Failed parse logging level")
		return logger, err
	}

	logger = logger.Level(lvl)

	return logger, nil
}
