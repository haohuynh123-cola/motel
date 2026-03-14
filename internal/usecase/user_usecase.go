package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type userUseCase struct {
	userRepo  port.UserRepository
	jwtSecret string
}

func NewUserUseCase(userRepo port.UserRepository, jwtSecret string) port.UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUseCase) Register(ctx context.Context, user *domain.User) error {
	userExist, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, port.ErrNotFound) {
		return err
	}

	if userExist != nil {
		return port.ErrUsernameAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return u.userRepo.Create(ctx, user)
}

func (u *userUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("Login error: user %s not found: %v\n", username, err)
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Login error: password mismatch for user %s\n", username)
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"id":          user.ID,
		"username":    user.Username,
		"permissions": user.Permissions,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
func (u *userUseCase) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (u *userUseCase) ListUsers(ctx context.Context) (*port.ApiResponse, error) {
	users, err := u.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	// Xóa password trước khi trả về
	for _, user := range users {
		user.Password = ""
	}
	return &port.ApiResponse{
		Status: true,
		Data:   users,
	}, nil
}
