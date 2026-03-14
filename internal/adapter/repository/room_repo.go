package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type roomRepository struct {
	db *pgxpool.Pool
}

// NewRoomRepository creates a new room repository
func NewRoomRepository(db *pgxpool.Pool) port.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	query := `INSERT INTO rooms (house_id, name, area, price, is_available, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, room.HouseID, room.Name, room.Area, room.Price, room.IsAvailable).
		Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	query := `SELECT id, house_id, name, area, price, is_available, created_at, updated_at FROM rooms WHERE id = $1`

	room := &domain.Room{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&room.ID, &room.HouseID, &room.Name, &room.Area, &room.Price, &room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}

	return room, nil
}

func (r *roomRepository) ListByHouseID(ctx context.Context, houseID int64) ([]*domain.Room, error) {
	query := `SELECT id, house_id, name, area, price, is_available, created_at, updated_at FROM rooms WHERE house_id = $1 ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := []*domain.Room{}
	for rows.Next() {
		r := &domain.Room{}
		if err := rows.Scan(&r.ID, &r.HouseID, &r.Name, &r.Area, &r.Price, &r.IsAvailable, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *roomRepository) Update(ctx context.Context, room *domain.Room) error {
	query := `UPDATE rooms SET name = $1, area = $2, price = $3, is_available = $4, updated_at = NOW() WHERE id = $5`

	commandTag, err := r.db.Exec(ctx, query, room.Name, room.Area, room.Price, room.IsAvailable, room.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("room not found")
	}

	return nil
}

func (r *roomRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM rooms WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("room not found")
	}

	return nil
}
