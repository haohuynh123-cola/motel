package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context, handler func(ctx context.Context, msg []byte) error) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		
		if err := handler(ctx, m.Value); err != nil {
			log.Printf("error handling message: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
