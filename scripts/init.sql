-- NATMap MySQL 数据库初始化脚本
-- 从 Cloudflare D1 迁移到 MySQL

-- 创建数据库
CREATE DATABASE IF NOT EXISTS natmap CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE natmap;

-- 租户表
CREATE TABLE IF NOT EXISTS tenants (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_name VARCHAR(255) NOT NULL COMMENT '租户名称',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_name (tenant_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 应用表
CREATE TABLE IF NOT EXISTS apps (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL COMMENT '租户ID',
    app_name VARCHAR(255) NOT NULL COMMENT '应用名称',
    description TEXT COMMENT '应用描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_app (tenant_id, app_name),
    KEY idx_tenant_id (tenant_id),
    CONSTRAINT fk_apps_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用表';

-- 映射表
CREATE TABLE IF NOT EXISTS mappings (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL COMMENT '租户ID',
    app_id INT UNSIGNED NOT NULL COMMENT '应用ID',
    public_ip VARCHAR(45) NOT NULL COMMENT '公网IP地址',
    public_port INT UNSIGNED NOT NULL COMMENT '公网端口',
    local_ip VARCHAR(45) COMMENT '本地IP地址',
    local_port INT UNSIGNED COMMENT '本地端口',
    protocol VARCHAR(10) DEFAULT 'tcp' COMMENT '协议类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_app (tenant_id, app_id),
    KEY idx_tenant_id (tenant_id),
    KEY idx_app_id (app_id),
    KEY idx_updated_at (updated_at),
    CONSTRAINT fk_mappings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_mappings_app FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='NAT映射表';

-- 用户表（用于管理后台认证）
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role ENUM('admin', 'user') DEFAULT 'user' COMMENT '角色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 插入默认管理员用户（密码: admin123）
-- 使用 bcrypt 哈希: $2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEsKYGsFqkNQJ3fKzG
INSERT INTO users (username, password_hash, role) VALUES 
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEsKYGsFqkNQJ3fKzG', 'admin')
ON DUPLICATE KEY UPDATE id=id;

-- 插入测试数据（可选）
-- 租户: 喀什信海通电子科技有限公司
INSERT INTO tenants (tenant_name) VALUES ('喀什信海通电子科技有限公司') ON DUPLICATE KEY UPDATE id=id;
SET @tenant_id = LAST_INSERT_ID();

-- 应用: http2tunnel1
INSERT INTO apps (tenant_id, app_name, description) 
SELECT @tenant_id, 'http2tunnel1', 'HTTP2隧道应用'
WHERE @tenant_id > 0
ON DUPLICATE KEY UPDATE id=id;
