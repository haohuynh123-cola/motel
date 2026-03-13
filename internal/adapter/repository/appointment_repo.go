package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type appointmentRepository struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) port.AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) Create(ctx context.Context, app *domain.Appointment) error {
	query := `INSERT INTO appointments (room_id, customer_name, customer_email, customer_phone, appointment_date, note, status, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id, created_at`
	
	if app.Status == "" {
		app.Status = "pending"
	}

	return r.db.QueryRow(ctx, query, app.RoomID, app.CustomerName, app.CustomerEmail, app.CustomerPhone, app.AppointmentDate, app.Note, app.Status).
		Scan(&app.ID, &app.CreatedAt)
}

func (r *appointmentRepository) GetByID(ctx context.Context, id int64) (*domain.Appointment, error) {
	query := `SELECT id, room_id, customer_name, customer_email, customer_phone, appointment_date, note, status, created_at 
	          FROM appointments WHERE id = $1`
	
	app := &domain.Appointment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&app.ID, &app.RoomID, &app.CustomerName, &app.CustomerEmail, &app.CustomerPhone, &app.AppointmentDate, &app.Note, &app.Status, &app.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, port.ErrNotFound
		}
		return nil, err
	}
	return app, nil
}

func (r *appointmentRepository) ListByRoomID(ctx context.Context, roomID int64) ([]*domain.Appointment, error) {
	query := `SELECT id, room_id, customer_name, customer_email, customer_phone, appointment_date, note, status, created_at 
	          FROM appointments WHERE room_id = $1 ORDER BY appointment_date DESC`
	
	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		app := &domain.Appointment{}
		err := rows.Scan(&app.ID, &app.RoomID, &app.CustomerName, &app.CustomerEmail, &app.CustomerPhone, &app.AppointmentDate, &app.Note, &app.Status, &app.CreatedAt)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, app)
	}
	return appointments, nil
}
