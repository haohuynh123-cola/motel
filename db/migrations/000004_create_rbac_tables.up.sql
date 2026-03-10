-- Tạo bảng Vai trò (Roles)
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL
);

-- Tạo bảng Quyền hạn (Permissions)
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL
);

-- Bảng trung gian Role - Permission
CREATE TABLE role_permissions (
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Bảng trung gian User - Role
CREATE TABLE user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- 1. Insert các Vai trò cơ bản
INSERT INTO roles (name, slug) VALUES
('Quản trị viên', 'admin'),
('Thành viên', 'member');

-- 2. Insert các Quyền hạn chi tiết
INSERT INTO permissions (name, slug) VALUES
('Xem nhà', 'house:read'),
('Tạo nhà', 'house:create'),
('Sửa nhà', 'house:update'),
('Xoá nhà', 'house:delete'),
('Xem phòng', 'room:read'),
('Tạo phòng', 'room:create'),
('Sửa phòng', 'room:update'),
('Xoá phòng', 'room:delete');

-- 3. Cấp TẤT CẢ quyền cho Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE slug = 'admin'), id FROM permissions;

-- 4. Cấp quyền cơ bản (Đọc, Tạo) cho Member
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE slug = 'member'), id FROM permissions
WHERE slug IN ('house:read', 'house:create', 'room:read', 'room:create');

-- 5. Chuyển đổi dữ liệu cũ: Cấp Role cho các User đang có trong DB dựa vào cột 'role' cũ
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.slug = u.role
WHERE u.role IS NOT NULL;

-- 6. Xoá cột 'role' cũ đi vì không cần nữa
ALTER TABLE users DROP COLUMN IF EXISTS role;
