package port

import (
	"context"
	"tro-go/internal/domain"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *domain.ChatMessage) error
	GetMessageHistory(ctx context.Context, userID, otherID int64, limit, offset int) ([]domain.ChatMessage, error)
	GetConversations(ctx context.Context, userID int64) ([]domain.ChatMessage, error) // Get latest msg for each contact
}

type ChatUseCase interface {
	SendMessage(ctx context.Context, senderID, receiverID int64, content string) (*domain.ChatMessage, error)
	GetHistory(ctx context.Context, userID, otherID int64, limit, offset int) ([]domain.ChatMessage, error)
}
