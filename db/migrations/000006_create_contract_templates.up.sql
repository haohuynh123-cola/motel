-- Bảng mẫu hợp đồng (Contract Templates)
CREATE TABLE IF NOT EXISTS contract_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL, -- Tên mẫu: Hợp đồng chung cư mini, Hợp đồng phòng giá rẻ...
    content TEXT NOT NULL,      -- Các điều khoản như: Tiền điện, Giờ giấc, Nội quy...
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Thêm mẫu mặc định vào hệ thống luôn để bạn dùng thử
INSERT INTO contract_templates (name, content) VALUES 
('Mẫu Hợp Đồng Tiêu Chuẩn', 'Điều 1: Tiền nhà đóng đúng hạn. Điều 2: Không tụ tập gây ồn ào sau 23h. Điều 3: Giữ gìn vệ sinh chung...');

-- Cập nhật bảng contracts để tham chiếu tới mẫu nào đã dùng (tùy chọn nhưng nên có)
ALTER TABLE contracts ADD COLUMN template_id INTEGER REFERENCES contract_templates(id);
