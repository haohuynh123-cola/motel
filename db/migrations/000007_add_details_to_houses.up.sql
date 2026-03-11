-- Thêm các cột địa chỉ chi tiết và hình ảnh cho bảng houses
ALTER TABLE houses 
ADD COLUMN IF NOT EXISTS province VARCHAR(100),
ADD COLUMN IF NOT EXISTS district VARCHAR(100),
ADD COLUMN IF NOT EXISTS ward VARCHAR(100),
ADD COLUMN IF NOT EXISTS images TEXT[]; -- Mảng chuỗi trong Postgres

-- Thêm mô tả và hình ảnh cho bảng rooms
ALTER TABLE rooms 
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS images TEXT[];
