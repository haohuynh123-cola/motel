package domain

import "time"

type Customer struct {
	ID             int64     `json:"id"`
	FullName       string    `json:"full_name"`
	IdentityNumber string    `json:"identity_number"` // CMND/CCCD
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`         // Quê quán/Địa chỉ thường trú
	Birthday       time.Time `json:"birthday"`
	Gender         string    `json:"gender"`          // Nam/Nữ
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
