package usecase

import (
	"context"
	"math"
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

func (u *customerUseCase) GetCustomer(ctx context.Context, id int64) (*port.ApiResponse, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &port.ApiResponse{Status: true, Data: customer}, nil
}

func (u *customerUseCase) ListCustomers(ctx context.Context, page, limit int) (*port.ApiResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	customers, err := u.customerRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := u.customerRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return &port.ApiResponse{
		Status: true,
		Data:   customers,
		Meta: &port.Meta{
			Total:       total,
			CurrentPage: page,
			LastPage:    lastPage,
			Limit:       limit,
		},
	}, nil
}

func (u *customerUseCase) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
	return u.customerRepo.Update(ctx, customer)
}

func (u *customerUseCase) DeleteCustomer(ctx context.Context, id int64) error {
	return u.customerRepo.Delete(ctx, id)
}
