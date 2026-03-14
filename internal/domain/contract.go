package domain

import "time"

type ContractStatus string

const (
	ContractActive   ContractStatus = "active"
	ContractEnded    ContractStatus = "ended"
	ContractCanceled ContractStatus = "canceled"
)

type Contract struct {
	ID             int64          `json:"id"`
	CustomerID     int64          `json:"customer_id"`
	RoomID         int64          `json:"room_id"`
	StartDate      time.Time      `json:"start_date"`
	EndDate        time.Time      `json:"end_date"`
	Deposit        float64        `json:"deposit"`      // Tiền cọc
	MonthlyRent    float64        `json:"monthly_rent"` // Giá thuê hàng tháng
	PaymentDay     int            `json:"payment_day"`  // Ngày đóng tiền hàng tháng (1-31)
	Status         ContractStatus `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
}
