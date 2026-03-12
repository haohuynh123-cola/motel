# Tro-Go API 🏠

Tro-Go là một hệ thống Backend API hiệu năng cao được xây dựng bằng ngôn ngữ **Go (Golang)**, chuyên dùng để quản lý hệ thống Nhà trọ và Phòng trọ.

Dự án này được thiết kế theo tiêu chuẩn công nghiệp, tuân thủ nghiêm ngặt **Kiến trúc Sạch (Clean Architecture)** và tích hợp quy trình **CI/CD** hiện đại.

---

## 🏗 Kiến trúc Hệ thống (Clean Architecture)

Luồng dữ liệu của hệ thống được chia thành 4 tầng rành mạch:

1. **Adapter / Handler:** Xử lý HTTP Request/Response (Echo Framework).
2. **UseCase:** Chứa toàn bộ logic nghiệp vụ (Auth, Email, RBAC).
3. **Port:** Định nghĩa các Interface để kết nối các tầng.
4. **Adapter / Repository:** Làm việc trực tiếp với Database (PostgreSQL).

---

## 🛠 Công nghệ sử dụng (Tech Stack)

*   **Backend:** Go 1.24+ (Web framework Echo v4).
*   **Database:** PostgreSQL 15 + `jackc/pgx/v5`.
*   **Migration:** `golang-migrate` (Tự động tạo bảng khi khởi động).
*   **Bảo mật:** JWT (JSON Web Tokens) + Bcrypt (Mã hóa mật khẩu).
*   **Phân quyền:** RBAC (Role-Based Access Control) mô hình 5 bảng linh hoạt.
*   **Email:** Giao thức SMTP (Tích hợp sẵn template HTML).
*   **DevOps:** Docker, GitHub Actions (CI/CD), Nginx.

---

## 🌟 Tính năng nổi bật

*   ✅ **Xác thực đa lớp:** Đăng ký, Đăng nhập và quản lý phiên làm việc qua JWT.
*   ✅ **Phân quyền chi tiết:** Phân quyền theo từng hành động (VD: Chỉ admin mới được xoá nhà, member chỉ được xem).
*   ✅ **Gửi Email tự động:** API nhắc nợ tiền phòng hàng tháng, tự động lấy thông tin giá tiền và gửi mail HTML chuyên nghiệp.
*   ✅ **Kiểm thử tự động:** Hệ thống Unit Test tích hợp sẵn kỹ thuật Mocking để test logic không cần DB.
*   ✅ **CI/CD chuẩn:** Tự động chạy Test và Deploy lên VPS mỗi khi Push code lên Github.

---

## ⚙️ Cấu hình Hệ thống (.env)

Tạo file `.env` tại thư mục gốc với các thông số sau:

```env
# Cổng chạy API
APP_PORT=8080

# Kết nối Database
DATABASE_URL=postgres://user:pass@db:5432/dbname?sslmode=disable
DB_MAX_CONNS=20

# Bảo mật JWT
JWT_SECRET=your_super_secret_key

# Cấu hình Gửi Email (Gmail SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password_16_chars
```

---

## 🚀 Quy trình Phát triển & Triển khai

### 1. Chạy dưới máy cá nhân (Local)
Dự án tích hợp **Air** để tự động Reload code khi bạn nhấn Save:
```bash
docker compose up
```

### 2. Chạy Unit Test
```bash
go test -v ./...
```

### 3. Tự động triển khai (CI/CD)
Dự án đã cấu hình sẵn GitHub Actions tại `.github/workflows/deploy.yml`.
*   Mỗi khi bạn **Push** lên nhánh `main`, Github sẽ tự động chạy **Unit Test**.
*   Nếu Test vượt qua (PASS), hệ thống sẽ SSH vào VPS và tự động **Deploy** phiên bản mới nhất.

---

## 📚 Tài liệu API (Endpoints)

| Method | Endpoint | Quyền hạn | Mô tả |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Public | Đăng ký tài khoản |
| `POST` | `/api/v1/auth/login` | Public | Đăng nhập lấy Token |
| `GET` | `/api/v1/auth/me` | User | Lấy thông tin cá nhân |
| `POST` | `/api/v1/houses` | User | Tạo nhà trọ mới |
| `DELETE` | `/api/v1/houses/:id` | **Admin** | Xoá nhà trọ (Chỉ Admin) |
| `POST` | `/api/v1/rooms/:id/remind` | User | Gửi mail nhắc nợ tiền phòng |

---

## 🤝 Hỗ trợ
Nếu có bất kỳ câu hỏi nào về kiến trúc hoặc cách cấu hình, vui lòng liên hệ Ban quản trị dự án.
