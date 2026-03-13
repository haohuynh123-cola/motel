package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type chatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) port.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) SaveMessage(ctx context.Context, msg *domain.ChatMessage) error {
	query := `INSERT INTO chat_messages (sender_id, receiver_id, content, created_at)
	          VALUES ($1, $2, $3, NOW()) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, msg.SenderID, msg.ReceiverID, msg.Content).
		Scan(&msg.ID, &msg.CreatedAt)
}

func (r *chatRepository) GetMessageHistory(ctx context.Context, userID, otherID int64, limit, offset int) ([]domain.ChatMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, is_read, created_at
		FROM chat_messages
		WHERE (sender_id = $1 AND receiver_id = $2)
		   OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(ctx, query, userID, otherID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var msg domain.ChatMessage
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.IsRead, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *chatRepository) GetConversations(ctx context.Context, userID int64) ([]domain.ChatMessage, error) {
	// Query to get the last message for each unique conversation partner
	query := `
		WITH LastMessages AS (
			SELECT DISTINCT ON (
				CASE WHEN sender_id < receiver_id THEN sender_id ELSE receiver_id END,
				CASE WHEN sender_id < receiver_id THEN receiver_id ELSE sender_id END
			) *
			FROM chat_messages
			WHERE sender_id = $1 OR receiver_id = $1
			ORDER BY 
				CASE WHEN sender_id < receiver_id THEN sender_id ELSE receiver_id END,
				CASE WHEN sender_id < receiver_id THEN receiver_id ELSE sender_id END,
				created_at DESC
		)
		SELECT id, sender_id, receiver_id, content, is_read, created_at
		FROM LastMessages
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var msg domain.ChatMessage
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.IsRead, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
