package domain

type RoomStats struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
	Occupied  int64 `json:"occupied"`
}

type DashboardStats struct {
	TotalHouses    int64     `json:"total_houses"`
	TotalCustomers int64     `json:"total_customers"`
	TotalContracts int64     `json:"total_contracts"`
	RoomStats      RoomStats `json:"room_stats"`
	TotalRevenue   float64   `json:"total_revenue"` // Tổng doanh thu dự kiến từ hợp đồng active
}
