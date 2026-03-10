# Tro-Go API 🏠

Tro-Go là một hệ thống Backend API hiệu năng cao được xây dựng bằng ngôn ngữ **Go (Golang)**, chuyên dùng để quản lý hệ thống Nhà trọ và Phòng trọ. 

Dự án này được thiết kế theo tiêu chuẩn của hệ thống lớn, tuân thủ nghiêm ngặt **Kiến trúc Sạch (Clean Architecture / Hexagonal Architecture)**, giúp mã nguồn dễ dàng bảo trì, mở rộng và viết Unit Test độc lập với các framework bên ngoài.

---

## 🏗 Kiến trúc Hệ thống (Clean Architecture)

Luồng dữ liệu của hệ thống được chia thành 4 tầng rành mạch, đi từ ngoài vào trong:

1. **Adapter / Handler (Tầng Giao tiếp):** Sử dụng `Echo framework` để xử lý HTTP Request/Response. Chịu trách nhiệm nhận JSON, kiểm tra tính hợp lệ cơ bản và trả về kết quả cho Client.
2. **UseCase (Tầng Nghiệp vụ):** Trái tim của ứng dụng. Chứa toàn bộ các quy tắc kinh doanh (Business Logic) như kiểm tra quyền hạn, mã hóa mật khẩu, tính toán. Tầng này hoàn toàn không biết sự tồn tại của Echo hay PostgreSQL.
3. **Port (Tầng Giao diện/Hợp đồng):** Chứa các `Interface`. Đóng vai trò là cầu nối (Dependency Injection) giữa UseCase và Repository, giúp dễ dàng Mock data khi viết Test.
4. **Adapter / Repository (Tầng Dữ liệu):** Xử lý giao tiếp trực tiếp với cơ sở dữ liệu (PostgreSQL) thông qua các câu lệnh SQL thuần được tối ưu hóa bằng thư viện `pgx`.

---

## 🛠 Công nghệ sử dụng (Tech Stack)

*   **Ngôn ngữ:** Go 1.24+
*   **Web Framework:** [Echo v4](https://echo.labstack.com/) (Siêu nhẹ, tốc độ cao).
*   **Database:** PostgreSQL 15.
*   **Database Driver:** `jackc/pgx/v5` (Driver mạnh mẽ nhất cho Postgres trong Go, hỗ trợ Connection Pooling).
*   **Quản lý Database Schema:** `golang-migrate` (Tự động chạy Migration khi khởi động Server).
*   **Bảo mật & Phân quyền:** JWT (JSON Web Tokens) kết hợp với mô hình **RBAC** (Role-Based Access Control) 5 bảng.
*   **Môi trường (DevOps):** Docker & Docker Compose (Hỗ trợ Hot-Reloading với `Air`).

---

## 🌟 Tính năng cốt lõi

*   **Xác thực (Authentication):** Đăng ký, Đăng nhập, trả về Token JWT an toàn.
*   **Phân quyền (Authorization - RBAC):** Hệ thống phân quyền động dựa trên Role và Permission (VD: `admin` có quyền `house:delete`, `member` thì không). Bảo vệ API bằng Middleware tự chế.
*   **Quản lý Nhà (Houses):** CRUD danh sách các tòa nhà trọ.
*   **Quản lý Phòng (Rooms):** CRUD danh sách phòng trọ, liên kết trực tiếp với từng Tòa nhà (Foreign Key).

---

## ⚙️ Cấu hình Hệ thống (.env)

Hệ thống đọc cấu hình tự động thông qua file `.env`. Nếu chạy bằng Docker, bạn chỉ cần tạo một file `.env` ở thư mục gốc của dự án.

```env
# Cổng chạy API Server
APP_PORT=8080

# Chuỗi kết nối đến PostgreSQL (Sử dụng tên service 'db' nếu chạy trong Docker)
DATABASE_URL=postgres://postgres:postgrespassword@db:5432/tro_go?sslmode=disable

# Số lượng kết nối tối đa giữ trong Pool (Tối ưu chịu tải)
DB_MAX_CONNS=20

# Khóa bí mật dùng để ký và giải mã JWT Token (Đổi mã này trên Production!)
JWT_SECRET=my-super-secret-key-change-it-in-production
```

---

## 🚀 Hướng dẫn chạy dự án (Local Development)

Dự án đã được cấu hình sẵn môi trường phát triển cực kỳ tiện lợi với **Docker** và công cụ **Air** (Hot-reloading).

### Yêu cầu tiên quyết
*   Máy tính đã cài đặt [Docker](https://www.docker.com/) và [Docker Compose](https://docs.docker.com/compose/).

### Khởi động Server
Bạn chỉ cần mở Terminal tại thư mục gốc của dự án và chạy duy nhất 1 lệnh:

```bash
docker compose up
```

Lệnh này sẽ thực hiện các việc sau một cách tự động:
1.  Khởi tạo một container PostgreSQL database.
2.  Tự động chạy các file trong `db/migrations/` để tạo các bảng RBAC, Houses, Rooms.
3.  Tự động tạo tài khoản `admin` mẫu và các nhóm quyền (Role).
4.  Khởi động API Server ở địa chỉ: `http://localhost:8080`.
5.  **Theo dõi (Watch):** Bất cứ khi nào bạn sửa và lưu một file `.go`, server sẽ tự động biên dịch và khởi động lại trong tích tắc (nhờ công cụ Air).

### Dừng Server
```bash
docker compose down
```
*(Lưu ý: Dữ liệu Database được bảo vệ an toàn trong Docker Volume `postgres_data`. Dù bạn có tắt hay bật lại, dữ liệu nhà trọ, user vẫn được giữ nguyên).*

---

## 📚 Tổ chức Thư mục (Directory Structure)

```text
tro-go/
├── cmd/
│   └── api/
│       └── main.go           # Điểm khởi đầu của ứng dụng, thiết lập DI, Router, kết nối DB.
├── db/
│   └── migrations/           # Chứa các file SQL phục vụ việc nâng cấp/tạo bảng tự động.
├── internal/                 # Các package nội bộ, không thể import từ dự án khác.
│   ├── adapter/
│   │   ├── handler/          # HTTP Controllers (Nhận Request -> Đẩy xuống UseCase -> Trả JSON).
│   │   └── repository/       # Làm việc trực tiếp với Postgres (Chứa lệnh SQL thuần).
│   ├── domain/               # Định nghĩa các Struct thực thể (Entities) dùng chung.
│   ├── port/                 # Định nghĩa các Interface kết nối giữa các tầng.
│   └── usecase/              # Trái tim dự án: Chứa logic nghiệp vụ (Auth, Tính toán, RBAC).
├── pkg/
│   └── config/               # Hàm đọc biến môi trường từ file .env.
├── .air.toml                 # Cấu hình tính năng Hot-Reloading cho môi trường Dev.
├── docker-compose.yml        # Thiết lập môi trường Dev bằng Docker.
├── docker-compose.prod.yml   # Thiết lập môi trường Production (Sử dụng Image siêu nhẹ).
└── Dockerfile                # Script build ứng dụng thành file nhị phân (Binary) cho Production.
```