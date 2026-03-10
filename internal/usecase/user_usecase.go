package usecase

import (
	"context"
	"errors"
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
	// Kiểm tra username đã tồn tại chưa
	_, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err == nil {
		return errors.New("username already exists")
	}

	// Mã hoá (Hash) mật khẩu trước khi lưu vào database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return u.userRepo.Create(ctx, user)
}

func (u *userUseCase) Login(ctx context.Context, username, password string) (string, error) {
	// 1. Tìm user theo username
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. So sánh mật khẩu user nhập vào với mật khẩu đã mã hoá trong DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Nếu đúng, tạo JWT Token
	claims := jwt.MapClaims{
		"id":          user.ID,
		"username":    user.Username,
		"permissions": user.Permissions,
		"exp":         time.Now().Add(time.Hour * 72).Unix(), // Hết hạn sau 3 ngày
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
	// Không trả về mật khẩu
	user.Password = ""
	return user, nil
}
