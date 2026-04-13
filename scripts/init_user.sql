-- 初始化默认管理员用户
-- 密码: admin123 (使用 bcrypt 加密)

INSERT INTO users (username, password_hash, role, created_at, updated_at)
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhmM6JGKpS4G3R1G2JH8YpfB0Bqy', 'admin', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
