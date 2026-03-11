package domain

import "time"

// Room represents a room within a boarding house.
type Room struct {
	ID          int64     `json:"id"`
	HouseID     int64     `json:"house_id"`
	Name        string    `json:"name"`
	Area        float64   `json:"area"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
