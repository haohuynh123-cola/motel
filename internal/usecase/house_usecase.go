package usecase

import (
	"context"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type houseUseCase struct {
	houseRepo port.HouseRepository
}

// NewHouseUseCase creates a new instance of HouseUseCase
func NewHouseUseCase(houseRepo port.HouseRepository) port.HouseUseCase {
	return &houseUseCase{
		houseRepo: houseRepo,
	}
}

func (u *houseUseCase) CreateHouse(ctx context.Context, house *domain.House) error {
	return u.houseRepo.Create(ctx, house)
}

func (u *houseUseCase) GetHouse(ctx context.Context, id int64) (*domain.House, error) {
	return u.houseRepo.GetByID(ctx, id)
}

func (u *houseUseCase) ListHouses(ctx context.Context) ([]*domain.House, error) {
	return u.houseRepo.List(ctx)
}

func (u *houseUseCase) UpdateHouse(ctx context.Context, house *domain.House) error {
	return u.houseRepo.Update(ctx, house)
}

func (u *houseUseCase) DeleteHouse(ctx context.Context, id int64) error {
	return u.houseRepo.Delete(ctx, id)
}

type roomUseCase struct {
	roomRepo port.RoomRepository
}

// NewRoomUseCase creates a new instance of RoomUseCase
func NewRoomUseCase(roomRepo port.RoomRepository) port.RoomUseCase {
	return &roomUseCase{
		roomRepo: roomRepo,
	}
}

func (u *roomUseCase) CreateRoom(ctx context.Context, room *domain.Room) error {
	return u.roomRepo.Create(ctx, room)
}

func (u *roomUseCase) GetRoom(ctx context.Context, id int64) (*domain.Room, error) {
	return u.roomRepo.GetByID(ctx, id)
}

func (u *roomUseCase) ListRoomsByHouse(ctx context.Context, houseID int64) ([]*domain.Room, error) {
	return u.roomRepo.ListByHouseID(ctx, houseID)
}

func (u *roomUseCase) UpdateRoom(ctx context.Context, room *domain.Room) error {
	return u.roomRepo.Update(ctx, room)
}

func (u *roomUseCase) DeleteRoom(ctx context.Context, id int64) error {
	return u.roomRepo.Delete(ctx, id)
}
