package repository

import (
	"context"
	"fmt"
	"tro-go/internal/domain"
	"tro-go/internal/port"

	"github.com/jackc/pgx/v5/pgxpool"
)

type customerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) port.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `INSERT INTO customers (full_name, identity_number, phone, email, address, birthday, gender)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		customer.FullName,
		customer.IdentityNumber,
		customer.Phone,
		customer.Email,
		customer.Address,
		customer.Birthday,
		customer.Gender,
	).Scan(&customer.ID, &customer.CreatedAt)

	if err != nil {
		return fmt.Errorf("không thể tạo khách thuê: %w", err)
	}
	return nil
}

func (r *customerRepository) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	customer := &domain.Customer{}
	query := `SELECT id, full_name, identity_number, phone, email, address, birthday, gender, created_at FROM customers WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&customer.ID, &customer.FullName, &customer.IdentityNumber,
		&customer.Phone, &customer.Email, &customer.Address,
		&customer.Birthday, &customer.Gender, &customer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) List(ctx context.Context) ([]*domain.Customer, error) {
	query := `SELECT id, full_name, identity_number, phone, email, address, birthday, gender, created_at FROM customers`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		c := &domain.Customer{}
		err := rows.Scan(&c.ID, &c.FullName, &c.IdentityNumber, &c.Phone, &c.Email, &c.Address, &c.Birthday, &c.Gender, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `UPDATE customers SET full_name=$1, phone=$2, email=$3, address=$4, birthday=$5, gender=$6, updated_at=NOW() WHERE id=$7`
	_, err := r.db.Exec(ctx, query, customer.FullName, customer.Phone, customer.Email, customer.Address, customer.Birthday, customer.Gender, customer.ID)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
