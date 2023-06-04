package util

import "context"

type ctxKey string

var (
	keyUserID ctxKey = "ctx:user_id"
)

func SetUserIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

func GetUserIDFromCtx(ctx context.Context) string {
	val := ctx.Value(keyUserID)
	if val == nil {
		return ""
	}

	if accID, ok := val.(string); ok {
		return accID
	}

	return ""
}
