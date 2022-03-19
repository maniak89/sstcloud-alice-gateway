package mongo

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionToken = "token"
)

type storage struct {
	cli  *mongo.Client
	db   *mongo.Database
	coll *mongo.Collection
}

func New(ctx context.Context, config Config) (*storage, error) {
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed init mongo connection")
		return nil, err
	}
	db := cli.Database(config.Name)
	return &storage{
		cli:  cli,
		db:   db,
		coll: db.Collection(collectionToken),
	}, nil
}

func (s *storage) Close(ctx context.Context) error {
	if s.db != nil {
		if err := s.cli.Disconnect(ctx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed close db connection")
			return err
		}
	}
	return nil
}

func (s *storage) Create(ctx context.Context, info oauth2.TokenInfo) error {
	_, err := s.coll.InsertOne(ctx, info)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed paste token")
		return err
	}
	return nil
}

func (s *storage) RemoveByCode(ctx context.Context, code string) error {
	logger := log.Ctx(ctx).With().Str("code", code).Logger()
	_, err := s.coll.DeleteMany(ctx, bson.D{{"Code", code}})
	if err != nil {
		logger.Error().Err(err).Msg("Failed delete by code")
		return err
	}
	return nil
}

func (s *storage) RemoveByAccess(ctx context.Context, access string) error {
	logger := log.Ctx(ctx).With().Str("access", access).Logger()
	_, err := s.coll.DeleteMany(ctx, bson.D{{"Access", access}})
	if err != nil {
		logger.Error().Err(err).Msg("Failed delete by access")
		return err
	}
	return nil
}

func (s *storage) RemoveByRefresh(ctx context.Context, refresh string) error {
	logger := log.Ctx(ctx).With().Str("refresh", refresh).Logger()
	_, err := s.coll.DeleteMany(ctx, bson.D{{"Refresh", refresh}})
	if err != nil {
		logger.Error().Err(err).Msg("Failed delete by refresh")
		return err
	}
	return nil
}

func (s *storage) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	logger := log.Ctx(ctx).With().Str("code", code).Logger()
	var token models.Token
	if err := s.coll.FindOne(ctx, bson.D{{"Code", code}}).Decode(&token); err != nil {
		logger.Error().Err(err).Msg("Failed get by code")
		return nil, err
	}
	return &token, nil
}

func (s *storage) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	logger := log.Ctx(ctx).With().Str("access", access).Logger()
	var token models.Token
	if err := s.coll.FindOne(ctx, bson.D{{"Access", access}}).Decode(&token); err != nil {
		logger.Error().Err(err).Msg("Failed get by access")
		return nil, err
	}
	return &token, nil
}

func (s *storage) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	logger := log.Ctx(ctx).With().Str("refresh", refresh).Logger()
	var token models.Token
	if err := s.coll.FindOne(ctx, bson.D{{"Refresh", refresh}}).Decode(&token); err != nil {
		logger.Error().Err(err).Msg("Failed get by refresh")
		return nil, err
	}
	return &token, nil
}
