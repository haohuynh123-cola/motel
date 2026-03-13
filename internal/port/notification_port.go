package port

import (
	"context"
	"tro-go/internal/domain"
)

type NotificationProvider interface {
	PublishEmail(ctx context.Context, notification domain.EmailNotification) error
}
