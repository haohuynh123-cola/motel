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
	"github.com/jackc/pgx/v5"
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

	fmt.Println("🚀 Bắt đầu quá trình Seeding 10 TRIỆU dữ liệu siêu tốc...")
	startTime := time.Now()

	// Số lượng cần tạo
	totalCustomers := 10_000_000
	totalHouses := 10_000
	batchSize := 100_000 // Mỗi lần nạp 100k dòng để không bị sập RAM

	// ---------------------------------------------------------

	// BƯỚC 1: CÀO 10 TRIỆU KHÁCH HÀNG (CUSTOMERS)

	// ---------------------------------------------------------

	fmt.Printf("\n⏳ Đang tạo %d Khách hàng...\n", totalCustomers)

	// Dùng 1 biến đếm tịnh tiến để đảm bảo CCCD không bao giờ trùng nhau

	baseIdentityNumber := 100000000000 // Bắt đầu từ 100 tỷ

	for i := 0; i < totalCustomers; i += batchSize {

		// Tạo một mảng 2 chiều chứa dữ liệu của 100,000 khách hàng

		rows := make([][]interface{}, 0, batchSize)

		for j := 0; j < batchSize; j++ {

			// Số CCCD là độc nhất: Base + (Batch Index * Batch Size) + Loop Index

			uniqueCCCD := fmt.Sprintf("%d", baseIdentityNumber+i+j)

			rows = append(rows, []interface{}{

				gofakeit.Name(),

				uniqueCCCD, // CCCD tuyệt đối không trùng

				gofakeit.Phone(),

				gofakeit.Email(),

				gofakeit.Address().Address,

				gofakeit.Gender(),

				time.Now(),

				time.Now(),
			})

		}

		// Dùng CopyFrom: Đổ thẳng 100,000 dòng vào Postgres trong 1 nhịp
		_, err := db.CopyFrom(
			ctx,
			pgx.Identifier{"customers"},
			[]string{"full_name", "identity_number", "phone", "email", "address", "gender", "created_at", "updated_at"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			log.Fatalf("Lỗi khi CopyFrom Customers tại batch %d: %v", i, err)
		}

		fmt.Printf("   ✅ Đã nạp %d / %d khách hàng...\n", i+batchSize, totalCustomers)
	}

	// ---------------------------------------------------------
	// BƯỚC 2: CÀO 10,000 NHÀ TRỌ (HOUSES)
	// (Không nên tạo 10 triệu nhà vì phi thực tế, 10k nhà là rất lớn rồi)
	// ---------------------------------------------------------
	fmt.Printf("\n⏳ Đang tạo %d Nhà trọ...\n", totalHouses)
	houseRows := make([][]interface{}, 0, totalHouses)
	for i := 0; i < totalHouses; i++ {
		houseRows = append(houseRows, []interface{}{
			gofakeit.Company() + " Motel",
			gofakeit.State(),
			gofakeit.City(),
			gofakeit.StreetName(),
			gofakeit.Address().Address,
			time.Now(),
			time.Now(),
		})
	}
	_, err = db.CopyFrom(
		ctx,
		pgx.Identifier{"houses"},
		[]string{"name", "province", "district", "ward", "address", "created_at", "updated_at"},
		pgx.CopyFromRows(houseRows),
	)
	if err != nil {
		log.Fatalf("Lỗi khi CopyFrom Houses: %v", err)
	}
	fmt.Printf("   ✅ Đã nạp %d Nhà trọ.\n", totalHouses)

	// ---------------------------------------------------------
	// BƯỚC 3: TẠO PHÒNG TRỌ (100 Phòng cho mỗi Nhà = 1 Triệu Phòng)
	// ---------------------------------------------------------
	totalRooms := totalHouses * 100
	fmt.Printf("\n⏳ Đang tạo %d Phòng trọ...\n", totalRooms)

	for i := 0; i < totalRooms; i += batchSize {
		roomRows := make([][]interface{}, 0, batchSize)
		for j := 0; j < batchSize; j++ {
			// ID nhà trọ random từ 1 đến 10,000
			houseID := int64(rand.Intn(totalHouses) + 1)

			roomRows = append(roomRows, []interface{}{
				houseID,
				fmt.Sprintf("Phòng %d", gofakeit.Number(100, 999)),
				gofakeit.Sentence(5),
				gofakeit.Float64Range(15, 50),
				gofakeit.Float64Range(1500000, 8000000), // 1.5tr đến 8tr
				true,
				time.Now(),
				time.Now(),
			})
		}

		_, err := db.CopyFrom(
			ctx,
			pgx.Identifier{"rooms"},
			[]string{"house_id", "name", "description", "area", "price", "is_available", "created_at", "updated_at"},
			pgx.CopyFromRows(roomRows),
		)
		if err != nil {
			log.Fatalf("Lỗi khi CopyFrom Rooms tại batch %d: %v", i, err)
		}
		fmt.Printf("   ✅ Đã nạp %d / %d phòng trọ...\n", i+batchSize, totalRooms)
	}

	duration := time.Since(startTime)
	fmt.Printf("\n🎉 HOÀN TẤT SEEDING! Tổng thời gian chạy: %v\n", duration)
}
