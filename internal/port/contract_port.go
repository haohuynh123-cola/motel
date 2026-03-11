package port

import (
	"context"
	"tro-go/internal/domain"
)

type ContractRepository interface {
	Create(ctx context.Context, contract *domain.Contract) error
	GetByID(ctx context.Context, id int64) (*domain.Contract, error)
	UpdateStatus(ctx context.Context, id int64, status domain.ContractStatus) error
}

type ContractUseCase interface {
	CreateContract(ctx context.Context, contract *domain.Contract) error
	GetContract(ctx context.Context, id int64) (*domain.Contract, error)
}
