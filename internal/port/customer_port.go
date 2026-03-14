package port

import (
	"context"
	"tro-go/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id int64) (*domain.Customer, error)
	List(ctx context.Context, offset, limit int) ([]*domain.Customer, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id int64) error
}

type CustomerUseCase interface {
	RegisterCustomer(ctx context.Context, customer *domain.Customer) error
	GetCustomer(ctx context.Context, id int64) (*ApiResponse, error)
	ListCustomers(ctx context.Context, page, limit int) (*ApiResponse, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) error
	DeleteCustomer(ctx context.Context, id int64) error
}

