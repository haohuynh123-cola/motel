package usecase

import (
	"context"
	"errors"
	"testing"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

// MockUserRepository implements port.UserRepository for testing
type MockUserRepository struct {
	GetByUsernameFunc func(ctx context.Context, username string) (*domain.User, error)
	CreateFunc        func(ctx context.Context, user *domain.User) error
	GetByIDFunc       func(ctx context.Context, id int64) (*domain.User, error)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return m.GetByUsernameFunc(ctx, username)
}
func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.CreateFunc(ctx, user)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.GetByIDFunc(ctx, id)
}

func TestUserUseCase_Register(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "secret"

	t.Run("Success when user does not exist", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			GetByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
				return nil, port.ErrNotFound
			},
			CreateFunc: func(ctx context.Context, user *domain.User) error {
				return nil
			},
		}
		uc := NewUserUseCase(mockRepo, jwtSecret)

		user := &domain.User{Username: "newuser", Password: "password123"}
		err := uc.Register(ctx, user)

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
	})

	t.Run("Fail when username already exists", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			GetByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
				return &domain.User{ID: 1, Username: "existinguser"}, nil
			},
		}
		uc := NewUserUseCase(mockRepo, jwtSecret)

		user := &domain.User{Username: "existinguser", Password: "password123"}
		err := uc.Register(ctx, user)

		if !errors.Is(err, port.ErrUsernameAlreadyExists) {
			t.Errorf("Expected ErrUsernameAlreadyExists, got %v", err)
		}
	})

	t.Run("Fail when database error occurs", func(t *testing.T) {
		dbErr := errors.New("db connection failed")
		mockRepo := &MockUserRepository{
			GetByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
				return nil, dbErr
			},
		}
		uc := NewUserUseCase(mockRepo, jwtSecret)

		user := &domain.User{Username: "anyuser", Password: "password123"}
		err := uc.Register(ctx, user)

		if !errors.Is(err, dbErr) {
			t.Errorf("Expected %v, got %v", dbErr, err)
		}
	})
}
