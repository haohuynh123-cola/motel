package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"tro-go/internal/adapter/db/postgres"
	"tro-go/pkg/config"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	// 1. Load cấu hình và kết nối DB
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Không thể load config: %v", err)
	}

	ctx := context.Background()
	db, err := postgres.ConnectPool(ctx, cfg.DatabaseURL, cfg.MaxConns)
	if err != nil {
		log.Fatalf("Kết nối DB thất bại: %v", err)
	}
	defer db.Close()

	fmt.Println("🚀 Bắt đầu quá trình Seeding dữ liệu...")

	// 2. Tạo Nhà trọ (Houses)
	houseIDs := []int64{}
	for i := 1; i <= 5; i++ {
		var houseID int64
		query := `INSERT INTO houses (name, province, district, ward, address, created_at, updated_at)
                  VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`
		err := db.QueryRow(ctx, query,
			gofakeit.Name()+" Hostel",
			"TP. Hồ Chí Minh",
			gofakeit.City(),
			gofakeit.StreetName(),
			gofakeit.Address().Address,
		).Scan(&houseID)

		if err != nil {
			log.Printf("Lỗi tạo nhà trọ: %v", err)
			continue
		}
		houseIDs = append(houseIDs, houseID)
		fmt.Printf("✅ Đã tạo nhà trọ: %d\n", houseID)

		// 3. Mỗi nhà trọ tạo 10 phòng (Rooms)
		for j := 1; j <= 10; j++ {
			var roomID int64
			query := `INSERT INTO rooms (house_id, name, description, area, price, is_available, created_at, updated_at)
                      VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`
			err := db.QueryRow(ctx, query,
				houseID,
				fmt.Sprintf("Phòng %d%02d", i, j),
				gofakeit.Sentence(10),
				gofakeit.Float64Range(15, 40),           // Diện tích 15-40m2
				gofakeit.Float64Range(2000000, 6000000), // Giá 2tr-6tr
				true,
			).Scan(&roomID)
			if err != nil {
				log.Printf("Lỗi tạo phòng: %v", err)
			}
		}

		// 4. Tạo cấu hình điện nước cho nhà trọ
		queryUtility := `INSERT INTO utility_configs (house_id, electricity_price, water_price, trash_price, internet_price)
                         VALUES ($1, $2, $3, $4, $5)`
		db.Exec(ctx, queryUtility, houseID, 3500, 20000, 50000, 100000)
	}

	// 5. Tạo Khách thuê (Customers)
	customerIDs := []int64{}
	for i := 1; i <= 20; i++ {
		var customerID int64
		query := `INSERT INTO customers (full_name, identity_number, phone, email, address, gender, created_at, updated_at)
                  VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`
		err := db.QueryRow(ctx, query,
			gofakeit.Name(),
			gofakeit.DigitN(12), // Giả lập số CCCD
			gofakeit.Phone(),
			gofakeit.Email(),
			gofakeit.Address().Address,
			gofakeit.Gender(),
		).Scan(&customerID)
		if err == nil {
			customerIDs = append(customerIDs, customerID)
		}
	}
	fmt.Printf("✅ Đã tạo %d khách thuê\n", len(customerIDs))

	// 6. Tạo một vài Hợp đồng (Contracts) ngẫu nhiên
	for i := 0; i < 10; i++ {
		// Lấy ngẫu nhiên một khách và một phòng (Giả định ID phòng từ 1-50)
		randomCustomerID := customerIDs[rand.Intn(len(customerIDs))]
		randomRoomID := int64(rand.Intn(50) + 1)

		query := `INSERT INTO contracts (customer_id, room_id, start_date, deposit, monthly_rent, payment_day, status)
                  VALUES ($1, $2, $3, $4, $5, $6, 'active')`
		db.Exec(ctx, query,
			randomCustomerID,
			randomRoomID,
			time.Now().AddDate(0, 0, -rand.Intn(30)), // Ngày bắt đầu trong 30 ngày qua
			5000000,
			3000000,
			5,
		)

		// Cập nhật trạng thái phòng thành hết trống
		db.Exec(ctx, "UPDATE rooms SET is_available = false WHERE id = $1", randomRoomID)
	}

	fmt.Println("✨ Quá trình Seeding hoàn tất! Bạn đã có dữ liệu để học SQL.")
}
