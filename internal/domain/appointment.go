package domain

import "time"

type Appointment struct {
	ID              int64     `json:"id"`
	RoomID          int64     `json:"room_id"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerPhone   string    `json:"customer_phone"`
	AppointmentDate time.Time `json:"appointment_date"`
	Note            string    `json:"note"`
	Status          string    `json:"status"` // pending, confirmed, cancelled
	CreatedAt       time.Time `json:"created_at"`
}
