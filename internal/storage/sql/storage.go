package sql

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects"

	storageModels "sstcloud-alice-gateway/internal/models/storage"
	storagePkg "sstcloud-alice-gateway/internal/storage"
)

type storage struct {
	config     Config
	connection *sql.DB
	db         *reform.Querier
	reformDB   reform.DBTXContext
	driver     string
}

func New(config Config) *storage {
	return &storage{
		config: config,
	}
}

func (s *storage) Connect(ctx context.Context) error {
	logger := log.Ctx(ctx)
	parsedConnectionString, err := url.Parse(s.config.ConnectionString)
	if err != nil {
		logger.Error().Err(err).Str("connection_string", s.config.ConnectionString).Msg("Failed parse connection string")
		return err
	}
	s.driver = parsedConnectionString.Scheme
	var found bool
	for _, driver := range sql.Drivers() {
		if driver == s.driver {
			found = true
			break
		}
	}
	if !found {
		err := errors.New("not supported db driver: " + s.driver)
		logger.Error().Err(err).Msg("")
		return err
	}
	connectionString := s.config.ConnectionString
	// sqlite3 не работает, если в начале connectionString - схема.
	// а без нее не работает postgres(не парсит аргументы)
	if s.driver == "sqlite3" {
		connectionString = s.config.ConnectionString[len(s.driver)+3:]
	}
	sqlDB, err := sql.Open(s.driver, connectionString)
	if err != nil {
		logger.Error().Err(err).Msg("Failed open sql connection")
		return err
	}
	s.connection = sqlDB

	t := reform.NewDB(sqlDB, dialects.ForDriver(s.driver), reform.NewPrintfLogger(logger.Printf))
	s.db = t.Querier
	s.reformDB = t
	return nil
}

func (s *storage) Disconnect(ctx context.Context) error {
	logger := log.Ctx(ctx)
	if s.connection == nil {
		err := storagePkg.ErrInvalidState
		logger.Error().Err(err).Msg("Invalid state of connection")
		return err
	}
	if err := s.connection.Close(); err != nil {
		logger.Error().Err(err).Msg("Failed close of connection")
		return err
	}
	s.connection = nil
	return nil
}

func (s *storage) Links(ctx context.Context, userID string) ([]*storageModels.Link, error) {
	logger := log.Ctx(ctx)
	rows, err := s.db.SelectAllFrom(storageModels.LinkTable, "where user_id = "+s.db.Placeholder(1), userID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed find links")
		return nil, err
	}
	result := make([]*storageModels.Link, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.(*storageModels.Link))
	}
	return result, nil
}

func (s *storage) Log(ctx context.Context, linkID string, level storageModels.LogLevel, msg string) {
	logger := log.Ctx(ctx).With().Str("link_id", linkID).Str("level", string(level)).Str("msg", msg).Logger()
	if err := s.db.WithContext(ctx).Insert(&storageModels.Log{
		LinkID:  linkID,
		Level:   level,
		Message: msg,
		Time:    time.Now(),
	}); err != nil {
		logger.Error().Err(err).Msg("Failed add log")
		return
	}
	return
}
