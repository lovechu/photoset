-- PhotoSet 数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS photoset CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE photoset;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    nickname VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1-active,0-inactive',
    last_login_at DATETIME(3) NULL,
    UNIQUE KEY idx_users_email (email),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试管理员账号 (密码: admin123)
INSERT INTO users (created_at, updated_at, nickname, email, password_hash, role, status) VALUES
(NOW(), NOW(), 'admin', 'admin@photoset.com', '$2a$10$rJZK5x5Y5Y5Y5Y5Y5Y5Y5e1x5x5x5x5x5x5x5x5x5x5x5x5x', 'admin', 1);

