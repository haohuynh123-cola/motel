package usecase

import (
	"context"
	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type customerUseCase struct {
	customerRepo port.CustomerRepository
}

func NewCustomerUseCase(customerRepo port.CustomerRepository) port.CustomerUseCase {
	return &customerUseCase{customerRepo: customerRepo}
}

func (u *customerUseCase) RegisterCustomer(ctx context.Context, customer *domain.Customer) error {
	return u.customerRepo.Create(ctx, customer)
}

func (u *customerUseCase) GetCustomer(ctx context.Context, id int64) (*domain.Customer, error) {
	return u.customerRepo.GetByID(ctx, id)
}

func (u *customerUseCase) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	return u.customerRepo.List(ctx)
}
