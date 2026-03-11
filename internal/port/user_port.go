package port

import (
	"context"
	"errors"

	"tro-go/internal/domain"
)

var ErrNotFound = errors.New("not found")
var ErrUsernameAlreadyExists = errors.New("username already exists")

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type UserUseCase interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, username, password string) (string, error) // Trả về Token JWT
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}
