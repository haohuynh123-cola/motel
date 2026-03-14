package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) port.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	// Dùng Transaction vì chúng ta phải insert vào 2 bảng (users và user_roles)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO users (username, password, full_name, created_at)
	          VALUES ($1, $2, $3, NOW()) RETURNING id, created_at`

	err = tx.QueryRow(ctx, query, user.Username, user.Password, user.FullName).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	// Tự động gán quyền 'member' khi tạo tài khoản mới
	roleQuery := `INSERT INTO user_roles (user_id, role_id)
				  SELECT $1, id FROM roles WHERE slug = 'member'`
	_, err = tx.Exec(ctx, roleQuery, user.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	// Dùng JOIN để gom tất cả các quyền của User thành 1 mảng (array_agg)
	query := `
		SELECT u.id, u.username, u.password, u.full_name, u.created_at,
		       COALESCE(array_agg(p.slug) FILTER (WHERE p.slug IS NOT NULL), '{}') as permissions
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r_table ON ur.role_id = r_table.id
		LEFT JOIN role_permissions rp ON r_table.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE u.username = $1
		GROUP BY u.id
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.FullName, &user.CreatedAt, &user.Permissions,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, port.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.password, u.full_name, u.created_at,
		       COALESCE(array_agg(p.slug) FILTER (WHERE p.slug IS NOT NULL), '{}') as permissions
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r_table ON ur.role_id = r_table.id
		LEFT JOIN role_permissions rp ON r_table.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE u.id = $1
		GROUP BY u.id
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.FullName, &user.CreatedAt, &user.Permissions,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) List(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.created_at,
		       COALESCE(array_agg(p.slug) FILTER (WHERE p.slug IS NOT NULL), '{}') as permissions
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r_table ON ur.role_id = r_table.id
		LEFT JOIN role_permissions rp ON r_table.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		GROUP BY u.id
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		u := &domain.User{}
		err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.CreatedAt, &u.Permissions)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
