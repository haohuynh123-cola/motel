package domain

import "time"

type ChatMessage struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID *int64    `json:"receiver_id,omitempty"` // NULL if it's for general support
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChatPayload struct {
	Type       string      `json:"type"`        // "message", "notification", "history"
	SenderID   int64       `json:"sender_id"`
	ReceiverID *int64      `json:"receiver_id"`
	Content    string      `json:"content"`
	Data       interface{} `json:"data"`
}
