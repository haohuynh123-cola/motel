-- Phục hồi lại cột role
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'member';

-- Cập nhật lại role từ bảng user_roles
UPDATE users u
SET role = r.slug
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE u.id = ur.user_id;

-- Xoá các bảng RBAC
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
