package port

import (
	"context"
	"tro-go/internal/domain"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (*domain.DashboardStats, error)
}

type DashboardUseCase interface {
	GetDashboardStats(ctx context.Context) (*ApiResponse, error)
}
