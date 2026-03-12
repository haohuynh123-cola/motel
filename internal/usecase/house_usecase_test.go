package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"tro-go/internal/domain"
	"tro-go/pkg/contextutil"
)

// ---------------------------------------------------------
// 1. TẠO MOCK REPOSITORY (KHO CHỨA GIẢ)
// ---------------------------------------------------------
// MockHouseRepo là một struct giả lập, nó sẽ thoả mãn interface port.HouseRepository
type MockHouseRepo struct {
	// Khai báo các hàm ảo (func) để chúng ta có thể tuỳ chỉnh kết quả trả về trong từng Test Case
	MockCreate  func(ctx context.Context, house *domain.House) error
	MockGetByID func(ctx context.Context, id int64) (*domain.House, error)
	// (Để ngắn gọn, tạm thời chỉ mock 2 hàm cần dùng. Go cho phép struct không implement đủ interface NẾU chúng ta trick một chút,
	// nhưng cách tốt nhất là implement đủ các hàm của interface).
}

// Implement các hàm bắt buộc của port.HouseRepository
func (m *MockHouseRepo) Create(ctx context.Context, house *domain.House) error {
	return m.MockCreate(ctx, house)
}
func (m *MockHouseRepo) GetByID(ctx context.Context, id int64) (*domain.House, error) {
	return m.MockGetByID(ctx, id)
}
func (m *MockHouseRepo) List(ctx context.Context, cursor, limit int) ([]*domain.House, error) {
	return nil, nil // Tạm thời bỏ qua
}
func (m *MockHouseRepo) Update(ctx context.Context, house *domain.House) error {
	return nil // Tạm thời bỏ qua
}
func (m *MockHouseRepo) Delete(ctx context.Context, id int64) error {
	return nil // Tạm thời bỏ qua
}

// ---------------------------------------------------------
// 2. VIẾT UNIT TEST BẰNG "TABLE-DRIVEN TESTS"
// ---------------------------------------------------------

func TestHouseUseCase_CreateHouse(t *testing.T) {
	// Table-Driven Tests: Định nghĩa các kịch bản (cases)
	testCases := []struct {
		name          string                        // Tên kịch bản test
		inputHouse    *domain.House                 // Dữ liệu đầu vào
		mockSetup     func(mockRepo *MockHouseRepo) // Cài đặt hành vi giả của DB
		expectedError error                         // Kết quả mong đợi (Có lỗi hay không?)
	}{
		{
			name: "Thành công - Tạo nhà hợp lệ",
			inputHouse: &domain.House{
				Name:    "Nhà Trọ Số 1",
				Address: "Hà Nội",
			},
			mockSetup: func(mockRepo *MockHouseRepo) {
				// Giả lập DB lưu thành công (trả về error = nil)
				mockRepo.MockCreate = func(ctx context.Context, house *domain.House) error {
					house.ID = 1 // Giả lập DB tự sinh ID
					house.CreatedAt = time.Now()
					return nil
				}
			},
			expectedError: nil,
		},
		{
			name: "Thất bại - Lỗi từ Database",
			inputHouse: &domain.House{
				Name:    "Nhà Trọ Lỗi",
				Address: "Hồ Chí Minh",
			},
			mockSetup: func(mockRepo *MockHouseRepo) {
				// Giả lập DB bị sập hoặc lỗi (trả về error)
				mockRepo.MockCreate = func(ctx context.Context, house *domain.House) error {
					return errors.New("database connection lost")
				}
			},
			expectedError: errors.New("database connection lost"),
		},
	}

	// ---------------------------------------------------------
	// 3. CHẠY CÁC TEST CASES
	// ---------------------------------------------------------
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// A. Khởi tạo Mock Repo
			mockRepo := &MockHouseRepo{}

			// B. Cài đặt kịch bản cho Mock Repo dựa vào Table Test
			tc.mockSetup(mockRepo)

			// C. Khởi tạo UseCase (Bơm Repo giả vào)
			useCase := NewHouseUseCase(mockRepo)

			// D. Thực thi hàm cần test
			// Giả lập user có ID = 1 đang gọi hàm
			ctx := contextutil.WithUserID(context.Background(), 1)
			err := useCase.CreateHouse(ctx, tc.inputHouse)

			// E. Kiểm tra kết quả (Assert)
			// So sánh lỗi thực tế (err) với lỗi mong đợi (tc.expectedError)
			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("Mong đợi lỗi '%v' nhưng không có lỗi nào xảy ra", tc.expectedError)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Mong đợi lỗi '%v' nhưng nhận được '%v'", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Không mong đợi lỗi nhưng nhận được: %v", err)
				}
				// Nếu thành công, kiểm tra xem ID đã được DB (giả) gắn vào chưa
				if tc.inputHouse.ID == 0 {
					t.Errorf("Mong đợi ID được cập nhật sau khi tạo, nhưng ID vẫn bằng 0")
				}
			}
		})
	}
}
