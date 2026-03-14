package usecase

import (
	"context"
	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type roomUseCase struct {
	roomRepo             port.RoomRepository
	appRepo              port.AppointmentRepository
	notificationProvider port.NotificationProvider
}

func NewRoomUseCase(roomRepo port.RoomRepository, appRepo port.AppointmentRepository, notificationProvider port.NotificationProvider) port.RoomUseCase {
	return &roomUseCase{
		roomRepo:             roomRepo,
		appRepo:              appRepo,
		notificationProvider: notificationProvider,
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
	return u.notificationProvider.PublishEmail(ctx, domain.EmailNotification{
		To:      toEmail,
		Subject: "Nhắc nhở thanh toán tiền phòng",
		Body:    "Chào bạn, vui lòng kiểm tra và thanh toán tiền phòng đúng hạn.",
	})
}

func (u *roomUseCase) BookAppointment(ctx context.Context, appointment *domain.Appointment) error {
	return u.appRepo.Create(ctx, appointment)
}
