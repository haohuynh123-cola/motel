package usecase

import (
	"context"
	"encoding/json"
	"log"
	"tro-go/internal/domain"
	"tro-go/pkg/email"
	"tro-go/pkg/kafka"
)

type NotificationWorker struct {
	consumer    *kafka.Consumer
	emailSender email.EmailSender
}

func NewNotificationWorker(consumer *kafka.Consumer, emailSender email.EmailSender) *NotificationWorker {
	return &NotificationWorker{
		consumer:    consumer,
		emailSender: emailSender,
	}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	log.Println("Notification Worker started, listening to Kafka...")
	w.consumer.ReadMessage(ctx, func(ctx context.Context, msg []byte) error {
		var notification domain.EmailNotification
		if err := json.Unmarshal(msg, &notification); err != nil {
			return err
		}

		log.Printf("Worker: Sending email to %s...", notification.To)
		return w.emailSender.Send(notification.To, notification.Subject, notification.Body)
	})
}
