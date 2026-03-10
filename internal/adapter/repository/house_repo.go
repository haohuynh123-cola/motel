package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type houseRepository struct {
	db *pgxpool.Pool
}

// NewHouseRepository creates a new house repository
func NewHouseRepository(db *pgxpool.Pool) port.HouseRepository {
	return &houseRepository{db: db}
}

func (r *houseRepository) Create(ctx context.Context, house *domain.House) error {
	query := `INSERT INTO houses (name, address, created_at, updated_at)
	          VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, house.Name, house.Address).Scan(&house.ID, &house.CreatedAt, &house.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *houseRepository) GetByID(ctx context.Context, id int64) (*domain.House, error) {
	query := `SELECT id, name, address, created_at, updated_at FROM houses WHERE id = $1`

	house := &domain.House{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&house.ID, &house.Name, &house.Address, &house.CreatedAt, &house.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("house not found")
		}
		return nil, err
	}

	return house, nil
}

func (r *houseRepository) List(ctx context.Context) ([]*domain.House, error) {
	query := `SELECT id, name, address, created_at, updated_at FROM houses ORDER BY id DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var houses []*domain.House
	for rows.Next() {
		h := &domain.House{}
		if err := rows.Scan(&h.ID, &h.Name, &h.Address, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		houses = append(houses, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return houses, nil
}

func (r *houseRepository) Update(ctx context.Context, house *domain.House) error {
	query := `UPDATE houses SET name = $1, address = $2, updated_at = NOW() WHERE id = $3`

	commandTag, err := r.db.Exec(ctx, query, house.Name, house.Address, house.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("house not found")
	}

	return nil
}

func (r *houseRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM houses WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("house not found")
	}

	return nil
}
