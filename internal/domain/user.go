package domain

import "time"

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password,omitempty"`
	FullName    string    `json:"full_name"`
	Permissions []string  `json:"permissions"` // Danh sách các slug quyền hạn
	CreatedAt   time.Time `json:"created_at"`
}
