package port

import (
	"context"
	"tro-go/internal/domain"
)

type ContractRepository interface {
	Create(ctx context.Context, contract *domain.Contract) error
	GetByID(ctx context.Context, id int64) (*domain.Contract, error)
	List(ctx context.Context) ([]*domain.Contract, error)
	ListByHouseID(ctx context.Context, houseID int64) ([]*domain.Contract, error)
	UpdateStatus(ctx context.Context, id int64, status domain.ContractStatus) error
}

type ContractUseCase interface {
	CreateContract(ctx context.Context, contract *domain.Contract) error
	GetContract(ctx context.Context, id int64) (*ApiResponse, error)
	ListAllContracts(ctx context.Context) (*ApiResponse, error)
	ListContractsByHouse(ctx context.Context, houseID int64) (*ApiResponse, error)
}
