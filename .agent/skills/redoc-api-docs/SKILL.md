---
name: redoc-api-docs
description: Hướng dẫn viết và quản lý tài liệu API bằng ReDoc (OpenAPI 3.0.3). Sử dụng khi cần tạo mới, cập nhật hoặc chuẩn hóa tài liệu API cho Frontend, đảm bảo tính nhất quán về schema, error codes và security (JWT).
---

# Redoc API Documentation Skill

Skill này giúp bạn tạo ra tài liệu API chuyên nghiệp, dễ đọc và chính xác cho Frontend team sử dụng ReDoc. Nó tập trung vào việc viết file OpenAPI (YAML) chuẩn chỉnh, bao gồm các mô tả chi tiết, ví dụ thực tế và các schema dùng chung.

## Quy trình làm việc (Workflow)

### 1. Khởi tạo OpenAPI Spec
Nếu dự án chưa có file `openapi.yaml`, hãy sử dụng template mẫu từ `references/openapi-template.yaml` để bắt đầu. File này đã bao gồm:
- Cấu hình server (Local, Staging, Production).
- Định nghĩa Security (Bearer Auth/JWT).
- Các Schema lỗi phổ biến (400, 401, 404, 500).

### 2. Thêm hoặc Cập nhật Endpoint
Khi thêm một API mới, hãy đảm bảo:
- **Tags**: Phân loại theo tài nguyên (ví dụ: `Users`, `Houses`, `Rooms`).
- **Summary & Description**: Viết ngắn gọn ở Summary và chi tiết ở Description (giải thích logic nghiệp vụ nếu cần).
- **Parameters**: Định nghĩa rõ `in: path`, `in: query`, hoặc `in: header`.
- **Responses**: Luôn bao gồm ít nhất một success response (200/201) và các error responses phù hợp.

### 3. Định nghĩa Data Models (Components/Schemas)
Tránh viết schema trực tiếp trong endpoint. Hãy đưa chúng vào mục `components/schemas` để tái sử dụng:
- Sử dụng `camelCase` cho các thuộc tính (phù hợp với FE).
- Luôn có thuộc tính `example` cho mỗi field để ReDoc hiển thị mẫu dữ liệu sinh động.
- Sử dụng `required` để đánh dấu các trường bắt buộc.

### 4. Kiểm tra chất lượng (Quality Check)
Trước khi hoàn tất, hãy kiểm tra:
- Tên endpoint có tuân thủ RESTful (danh từ, số nhiều)?
- Các lỗi 4xx/5xx có đúng format của dự án không?
- JWT security đã được áp dụng cho các endpoint cần bảo mật chưa?

## Các mẫu tham khảo

- **Mẫu OpenAPI chuẩn**: Xem tại `references/openapi-template.yaml`.
- **Quy tắc đặt tên**:
    - Endpoint: `/v1/houses`, `/v1/houses/{id}/rooms`.
    - Schema: `HouseResponse`, `CreateHouseRequest`, `ErrorResponse`.

## Resources

### references/
- `openapi-template.yaml`: File mẫu OpenAPI 3.0.3 đầy đủ cấu trúc để bắt đầu dự án mới.
- `error-codes.md`: Danh sách các mã lỗi và thông điệp chuẩn của hệ thống.

### scripts/
- `validate-openapi.sh`: (Nếu có) Script để kiểm tra tính hợp lệ của file YAML bằng `spectral` hoặc `openapi-cli`.
