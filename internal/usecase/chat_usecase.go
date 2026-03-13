package usecase

import (
	"context"
	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type chatUseCase struct {
	repo port.ChatRepository
}

func NewChatUseCase(repo port.ChatRepository) port.ChatUseCase {
	return &chatUseCase{repo: repo}
}

func (u *chatUseCase) SendMessage(ctx context.Context, senderID, receiverID int64, content string) (*domain.ChatMessage, error) {
	msg := &domain.ChatMessage{
		SenderID:   senderID,
		ReceiverID: &receiverID,
		Content:    content,
	}
	err := u.repo.SaveMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (u *chatUseCase) GetHistory(ctx context.Context, userID, otherID int64, limit, offset int) ([]domain.ChatMessage, error) {
	return u.repo.GetMessageHistory(ctx, userID, otherID, limit, offset)
}
