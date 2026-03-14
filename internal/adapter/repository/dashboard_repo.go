package repository

import (
	"context"
	"tro-go/internal/domain"
	"tro-go/internal/port"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) port.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	// 1. Đếm tổng số nhà
	err := r.db.QueryRow(ctx, "SELECT COUNT(id) FROM houses").Scan(&stats.TotalHouses)
	if err != nil {
		return nil, err
	}

	// 2. Đếm tổng số khách thuê
	err = r.db.QueryRow(ctx, "SELECT COUNT(id) FROM customers").Scan(&stats.TotalCustomers)
	if err != nil {
		return nil, err
	}

	// 3. Đếm tổng số hợp đồng đang active
	err = r.db.QueryRow(ctx, "SELECT COUNT(id) FROM contracts WHERE status = 'active'").Scan(&stats.TotalContracts)
	if err != nil {
		return nil, err
	}

	// 4. Thống kê phòng
	queryRooms := `
		SELECT
			COUNT(id) as total,
			COUNT(id) FILTER (WHERE is_available = true) as available,
			COUNT(id) FILTER (WHERE is_available = false) as occupied
		FROM rooms
	`
	err = r.db.QueryRow(ctx, queryRooms).Scan(&stats.RoomStats.Total, &stats.RoomStats.Available, &stats.RoomStats.Occupied)
	if err != nil {
		return nil, err
	}

	// 5. Tính tổng doanh thu tháng dự kiến (tổng monthly_rent của hợp đồng active)
	err = r.db.QueryRow(ctx, "SELECT COALESCE(SUM(monthly_rent), 0) FROM contracts WHERE status = 'active'").Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
