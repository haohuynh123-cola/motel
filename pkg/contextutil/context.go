package contextutil

import (
	"context"
)

// contextKey là kiểu riêng tư để tránh trùng lặp key với các thư viện khác
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "user_role"
)

// WithUserID nhét userID vào context
func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID lấy userID từ context ra
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

// WithRole nhét role vào context
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

// GetRole lấy role từ context ra
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
