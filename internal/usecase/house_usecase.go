package usecase

import (
	"context"
	"fmt"
	"time"
	"tro-go/internal/domain"
	"tro-go/internal/port"
	"tro-go/pkg/contextutil"
	"tro-go/pkg/email"
)

type houseUseCase struct {
	houseRepo port.HouseRepository
}

// NewHouseUseCase creates a new instance of HouseUseCase
func NewHouseUseCase(houseRepo port.HouseRepository) port.HouseUseCase {
	return &houseUseCase{
		houseRepo: houseRepo,
	}
}

func (u *houseUseCase) CreateHouse(ctx context.Context, house *domain.House) error {
	// Lấy userID từ context một cách xịn xò
	userID, ok := contextutil.GetUserID(ctx)
	if !ok {
		return fmt.Errorf("không tìm thấy thông tin người dùng trong hệ thống")
	}

	// In ra để debug (Sau này bạn có thể gán house.OwnerID = userID)
	fmt.Printf("User %d đang tạo nhà trọ: %s\n", userID, house.Name)

	return u.houseRepo.Create(ctx, house)
}

func (u *houseUseCase) GetHouse(ctx context.Context, id int64) (*domain.House, error) {
	return u.houseRepo.GetByID(ctx, id)
}

func (u *houseUseCase) ListHouses(ctx context.Context) ([]*domain.House, error) {
	return u.houseRepo.List(ctx)
}

func (u *houseUseCase) UpdateHouse(ctx context.Context, house *domain.House) error {
	return u.houseRepo.Update(ctx, house)
}

func (u *houseUseCase) DeleteHouse(ctx context.Context, id int64) error {
	return u.houseRepo.Delete(ctx, id)
}

type roomUseCase struct {
	roomRepo    port.RoomRepository
	emailSender email.EmailSender
}

// NewRoomUseCase creates a new instance of RoomUseCase
func NewRoomUseCase(roomRepo port.RoomRepository, emailSender email.EmailSender) port.RoomUseCase {
	return &roomUseCase{
		roomRepo:    roomRepo,
		emailSender: emailSender,
	}
}

func (u *roomUseCase) CreateRoom(ctx context.Context, room *domain.Room) error {
	return u.roomRepo.Create(ctx, room)
}

func (u *roomUseCase) GetRoom(ctx context.Context, id int64) (*domain.Room, error) {
	return u.roomRepo.GetByID(ctx, id)
}

func (u *roomUseCase) ListRoomsByHouse(ctx context.Context, houseID int64) ([]*domain.Room, error) {
	return u.roomRepo.ListByHouseID(ctx, houseID)
}

func (u *roomUseCase) UpdateRoom(ctx context.Context, room *domain.Room) error {
	return u.roomRepo.Update(ctx, room)
}

func (u *roomUseCase) DeleteRoom(ctx context.Context, id int64) error {
	return u.roomRepo.Delete(ctx, id)
}

func (u *roomUseCase) SendPaymentReminder(ctx context.Context, id int64, toEmail string) error {
	// 1. Lấy thông tin phòng để biết giá tiền và tên
	room, err := u.roomRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("không tìm thấy phòng: %w", err)
	}

	// 2. Chuẩn bị dữ liệu gửi Email
	amountStr := fmt.Sprintf("%.0f", room.Price)
	// Hạn chót là ngày mùng 5 của tháng hiện tại
	now := time.Now()
	dueDate := fmt.Sprintf("05/%02d/%d", now.Month(), now.Year())

	// 3. Gọi module email để gửi (truyền cứng tên người nhận tạm thời)
	err = u.emailSender.SendReminderEmail(toEmail, "Khách thuê phòng", room.Name, amountStr, dueDate)
	if err != nil {
		return fmt.Errorf("lỗi khi gửi email: %w", err)
	}

	return nil
}
