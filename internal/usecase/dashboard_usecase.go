package usecase

import (
	"context"
	"tro-go/internal/port"
)

type dashboardUseCase struct {
	dashboardRepo port.DashboardRepository
}

func NewDashboardUseCase(dashboardRepo port.DashboardRepository) port.DashboardUseCase {
	return &dashboardUseCase{dashboardRepo: dashboardRepo}
}

func (u *dashboardUseCase) GetDashboardStats(ctx context.Context) (*port.ApiResponse, error) {
	stats, err := u.dashboardRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &port.ApiResponse{
		Status: true,
		Data:   stats,
	}, nil
}
