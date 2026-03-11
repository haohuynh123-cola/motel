package usecase

import (
	"context"
	"fmt"
	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type contractUseCase struct {
	contractRepo port.ContractRepository
	roomRepo     port.RoomRepository
	customerRepo port.CustomerRepository
}

func NewContractUseCase(contractRepo port.ContractRepository, roomRepo port.RoomRepository, customerRepo port.CustomerRepository) port.ContractUseCase {
	return &contractUseCase{
		contractRepo: contractRepo,
		roomRepo:     roomRepo,
		customerRepo: customerRepo,
	}
}

func (u *contractUseCase) CreateContract(ctx context.Context, contract *domain.Contract) error {
	// 1. Kiểm tra khách thuê có tồn tại không
	_, err := u.customerRepo.GetByID(ctx, contract.CustomerID)
	if err != nil {
		return fmt.Errorf("không tìm thấy khách thuê: %w", err)
	}

	// 2. Kiểm tra phòng còn trống không
	room, err := u.roomRepo.GetByID(ctx, contract.RoomID)
	if err != nil {
		return fmt.Errorf("không tìm thấy phòng: %w", err)
	}

	if !room.IsAvailable {
		return fmt.Errorf("phòng hiện đang có người ở, không thể ký hợp đồng mới")
	}

	// 3. Tạo hợp đồng
	contract.Status = domain.ContractActive
	err = u.contractRepo.Create(ctx, contract)
	if err != nil {
		return err
	}

	// 4. Cập nhật trạng thái phòng sang "Hết phòng" (IsAvailable = false)
	room.IsAvailable = false
	err = u.roomRepo.Update(ctx, room)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật trạng thái phòng: %w", err)
	}

	return nil
}

func (u *contractUseCase) GetContract(ctx context.Context, id int64) (*domain.Contract, error) {
	return u.contractRepo.GetByID(ctx, id)
}
