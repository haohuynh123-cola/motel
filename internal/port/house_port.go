package port

import (
	"context"

	"tro-go/internal/domain"
)

// HouseRepository defines the interface for house data access
type HouseRepository interface {
	Create(ctx context.Context, house *domain.House) error
	GetByID(ctx context.Context, id int64) (*domain.House, error)
	List(ctx context.Context) ([]*domain.House, error)
	Update(ctx context.Context, house *domain.House) error
	Delete(ctx context.Context, id int64) error
}

// RoomRepository defines the interface for room data access
type RoomRepository interface {
	Create(ctx context.Context, room *domain.Room) error
	GetByID(ctx context.Context, id int64) (*domain.Room, error)
	ListByHouseID(ctx context.Context, houseID int64) ([]*domain.Room, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id int64) error
}

// HouseUseCase defines the business logic interface for houses
type HouseUseCase interface {
	CreateHouse(ctx context.Context, house *domain.House) error
	GetHouse(ctx context.Context, id int64) (*domain.House, error)
	ListHouses(ctx context.Context) ([]*domain.House, error)
	UpdateHouse(ctx context.Context, house *domain.House) error
	DeleteHouse(ctx context.Context, id int64) error
}

// RoomUseCase defines the business logic interface for rooms
type RoomUseCase interface {
	CreateRoom(ctx context.Context, room *domain.Room) error
	GetRoom(ctx context.Context, id int64) (*domain.Room, error)
	ListRoomsByHouse(ctx context.Context, houseID int64) ([]*domain.Room, error)
	UpdateRoom(ctx context.Context, room *domain.Room) error
	DeleteRoom(ctx context.Context, id int64) error
	SendPaymentReminder(ctx context.Context, id int64, toEmail string) error
}
