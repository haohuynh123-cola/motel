package usecase

import (
	"context"
	"fmt"
	"math"
	"tro-go/internal/domain"
	"tro-go/internal/port"
	"tro-go/pkg/contextutil"
)

type houseUseCase struct {
	houseRepo port.HouseRepository
}

func NewHouseUseCase(houseRepo port.HouseRepository) port.HouseUseCase {
	return &houseUseCase{
		houseRepo: houseRepo,
	}
}

func (u *houseUseCase) CreateHouse(ctx context.Context, house *domain.House) error {
	userID, ok := contextutil.GetUserID(ctx)
	if !ok {
		return fmt.Errorf("không tìm thấy thông tin người dùng trong hệ thống")
	}
	fmt.Printf("User %d đang tạo nhà trọ: %s\n", userID, house.Name)
	return u.houseRepo.Create(ctx, house)
}

func (u *houseUseCase) GetHouse(ctx context.Context, id int64) (*domain.House, error) {
	return u.houseRepo.GetByID(ctx, id)
}

func (u *houseUseCase) ListHouses(ctx context.Context, page, limit int) (*port.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	houses, err := u.houseRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := u.houseRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return &port.Pagination{
		Data:        houses,
		Total:       total,
		CurrentPage: page,
		LastPage:    lastPage,
		Limit:       limit,
	}, nil
}

func (u *houseUseCase) UpdateHouse(ctx context.Context, house *domain.House) error {
	return u.houseRepo.Update(ctx, house)
}

func (u *houseUseCase) DeleteHouse(ctx context.Context, id int64) error {
	return u.houseRepo.Delete(ctx, id)
}
