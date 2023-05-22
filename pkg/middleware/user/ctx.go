package user

import (
	"context"
)

type cctx struct{}

var vCtx cctx

func withContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, vCtx, userID)
}

func User(ctx context.Context) string {
	if user, ok := ctx.Value(vCtx).(string); ok {
		return user
	}
	return ""
}
