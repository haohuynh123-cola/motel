package repository

import (
	"context"
	"fmt"
	"tro-go/internal/domain"
	"tro-go/internal/port"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contractRepository struct {
	db *pgxpool.Pool
}

func NewContractRepository(db *pgxpool.Pool) port.ContractRepository {
	return &contractRepository{db: db}
}

func (r *contractRepository) Create(ctx context.Context, contract *domain.Contract) error {
	query := `INSERT INTO contracts (customer_id, room_id, start_date, end_date, deposit, monthly_rent, payment_day, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		contract.CustomerID,
		contract.RoomID,
		contract.StartDate,
		contract.EndDate,
		contract.Deposit,
		contract.MonthlyRent,
		contract.PaymentDay,
		contract.Status,
	).Scan(&contract.ID, &contract.CreatedAt)

	if err != nil {
		return fmt.Errorf("không thể tạo hợp đồng: %w", err)
	}
	return nil
}

func (r *contractRepository) GetByID(ctx context.Context, id int64) (*domain.Contract, error) {
	c := &domain.Contract{}
	query := `SELECT id, customer_id, room_id, start_date, end_date, deposit, monthly_rent, payment_day, status, created_at FROM contracts WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.CustomerID, &c.RoomID, &c.StartDate, &c.EndDate,
		&c.Deposit, &c.MonthlyRent, &c.PaymentDay, &c.Status, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *contractRepository) List(ctx context.Context) ([]*domain.Contract, error) {
	query := `SELECT id, customer_id, room_id, start_date, end_date, deposit, monthly_rent, payment_day, status, created_at FROM contracts`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contracts := []*domain.Contract{}
	for rows.Next() {
		c := &domain.Contract{}
		err := rows.Scan(
			&c.ID, &c.CustomerID, &c.RoomID, &c.StartDate, &c.EndDate,
			&c.Deposit, &c.MonthlyRent, &c.PaymentDay, &c.Status, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, c)
	}
	return contracts, nil
}

func (r *contractRepository) ListByHouseID(ctx context.Context, houseID int64) ([]*domain.Contract, error) {
	query := `SELECT c.id, c.customer_id, c.room_id, c.start_date, c.end_date, c.deposit, c.monthly_rent, c.payment_day, c.status, c.created_at 
              FROM contracts c
              JOIN rooms r ON c.room_id = r.id
              WHERE r.house_id = $1`

	rows, err := r.db.Query(ctx, query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contracts := []*domain.Contract{}
	for rows.Next() {
		c := &domain.Contract{}
		err := rows.Scan(
			&c.ID, &c.CustomerID, &c.RoomID, &c.StartDate, &c.EndDate,
			&c.Deposit, &c.MonthlyRent, &c.PaymentDay, &c.Status, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, c)
	}
	return contracts, nil
}

func (r *contractRepository) UpdateStatus(ctx context.Context, id int64, status domain.ContractStatus) error {
	query := `UPDATE contracts SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}
