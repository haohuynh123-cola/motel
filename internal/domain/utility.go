package domain

type UtilityConfig struct {
	HouseID       int64   `json:"house_id"`
	ElectricityPrice float64 `json:"electricity_price"` // Giá điện mỗi kWh
	WaterPrice       float64 `json:"water_price"`       // Giá nước mỗi khối
	TrashPrice       float64 `json:"trash_price"`       // Tiền rác hàng tháng
	InternetPrice    float64 `json:"internet_price"`    // Tiền mạng hàng tháng
}

type UtilityUsage struct {
	ID              int64   `json:"id"`
	RoomID          int64   `json:"room_id"`
	Month           int     `json:"month"`
	Year            int     `json:"year"`
	ElectricityBegin float64 `json:"electricity_begin"` // Số điện cũ
	ElectricityEnd   float64 `json:"electricity_end"`   // Số điện mới
	WaterBegin       float64 `json:"water_begin"`       // Số nước cũ
	WaterEnd         float64 `json:"water_end"`         // Số nước mới
}
