package user

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-Id")
		ctx := r.Context()
		ctx = withContext(ctx, user)
		ctx = log.Ctx(ctx).With().Str("user_id", user).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
