package authcontext

import (
	"context"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
)

// SetUserID : for testings
func SetUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, config.ContextUserIDKey, id)
}

func GetUserID(ctx context.Context) int64 {
	id, ok := ctx.Value(config.ContextUserIDKey).(int64)
	if !ok {
		return 0
	}
	return id
}