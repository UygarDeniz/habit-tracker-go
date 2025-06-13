package middleware

import (
	"context"
	"errors"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)

// GetUserIDFromContext safely extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

func IsAuthenticated(ctx context.Context) bool {
	_, err := GetUserIDFromContext(ctx)
	return err == nil
}
