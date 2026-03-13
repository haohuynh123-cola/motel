package kafka

import (
	"context"
	"tro-go/internal/domain"
	"tro-go/internal/port"
	pkgKafka "tro-go/pkg/kafka"
)

type notificationAdapter struct {
	producer *pkgKafka.Producer
}

func NewNotificationAdapter(producer *pkgKafka.Producer) port.NotificationProvider {
	return &notificationAdapter{producer: producer}
}

func (a *notificationAdapter) PublishEmail(ctx context.Context, notification domain.EmailNotification) error {
	return a.producer.SendMessage(ctx, "email-notification", notification)
}
