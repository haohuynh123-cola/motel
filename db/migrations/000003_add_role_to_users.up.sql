ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'member';
-- Cập nhật user admin hiện tại (nếu có) lên quyền admin
UPDATE users SET role = 'admin' WHERE username = 'admin';
