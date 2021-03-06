package oauth2

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	jwt2 "github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"

	"github.com/maniak89/sstcloud-alice-gateway/internal/services/rest/handlers/oauth2/mongo"
)

type service struct {
	config Config
	server *server.Server
}

func New(config Config) *service {
	return &service{
		config: config,
	}
}

func (s *service) Init(ctx context.Context) error {
	logger := log.Ctx(ctx)
	if err := s.config.Validate(); err != nil {
		return err
	}
	if !s.config.Enabled {
		return nil
	}
	manager := manage.NewDefaultManager()
	tokenStore, err := mongo.New(ctx, s.config.Storage)
	if err != nil {
		logger.Error().Err(err).Msg("Failed init token storage")
		return err
	}
	manager.MapTokenStorage(tokenStore)

	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(s.config.Key), jwt.GetSigningMethod(s.config.AuthAlgo)))

	clientStore := store.NewClientStore()
	if err := clientStore.Set(s.config.UserID, &models.Client{
		ID:     s.config.UserID,
		Secret: s.config.UserSecret,
		Domain: s.config.UserDomain,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed init client storage")
		return err
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	s.server = srv
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Error().Err(err).Msg("internal error")
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Error().Err(re.Error).Msg("response error")
	})

	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		logger := log.Ctx(r.Context())
		switch r.Method {
		case http.MethodGet:
			userID = r.URL.Query().Get("client_id")
			if userID == "" {
				return "", errors.New("empty client id")
			}
			return userID, nil
		}
		blob, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Failed read body")
			return "", err
		}
		logger.Debug().Bytes("body", blob).Interface("headers", r.Header).Msg("Can't do auth")
		return "", errors.New("unknown how fix it")
	})
	return nil
}

func (s *service) Authorize(w http.ResponseWriter, r *http.Request) {
	if err := s.server.HandleAuthorizeRequest(w, r); err != nil {
		log.Ctx(r.Context()).Warn().Err(err).Msg("Failed authorize")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *service) Token(w http.ResponseWriter, r *http.Request) {
	if err := s.server.HandleTokenRequest(w, r); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("Failed authorize")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *service) Verify(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.server.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(jwtauth.NewContext(r.Context(), tokenWrapper{token}, nil))
		handler.ServeHTTP(w, r)
	})
}

type tokenWrapper struct {
	oauth2.TokenInfo
}

func (t tokenWrapper) Audience() []string {
	return nil
}

func (t tokenWrapper) Expiration() time.Time {
	return time.Now().Add(t.TokenInfo.GetAccessExpiresIn())
}

func (t tokenWrapper) IssuedAt() time.Time {
	return t.TokenInfo.GetAccessCreateAt()
}

func (t tokenWrapper) Issuer() string {
	return "self"
}

func (t tokenWrapper) JwtID() string {
	return t.TokenInfo.GetAccess()
}

func (t tokenWrapper) NotBefore() time.Time {
	return time.Time{}
}

func (t tokenWrapper) Subject() string {
	return t.TokenInfo.GetUserID()
}

func (t tokenWrapper) PrivateClaims() map[string]interface{} {
	return nil
}

func (t tokenWrapper) Get(s string) (interface{}, bool) {
	return nil, false
}

func (t tokenWrapper) Set(s string, i interface{}) error {
	return nil
}

func (t tokenWrapper) Remove(s string) error {
	return nil
}

func (t tokenWrapper) Clone() (jwt2.Token, error) {
	return t, nil
}

func (t tokenWrapper) Iterate(ctx context.Context) jwt2.Iterator {
	return nil
}

func (t tokenWrapper) Walk(ctx context.Context, visitor jwt2.Visitor) error {
	return nil
}

func (t tokenWrapper) AsMap(ctx context.Context) (map[string]interface{}, error) {
	return nil, nil
}
