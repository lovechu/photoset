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

-- 套图表
CREATE TABLE IF NOT EXISTS photosets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    title VARCHAR(200) NOT NULL,
    cover VARCHAR(500) NOT NULL,
    description TEXT,
    is_free TINYINT NOT NULL DEFAULT 1 COMMENT '1-free,0-paid',
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    user_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT 'draft,published,pending',
    INDEX idx_photosets_user_id (user_id),
    INDEX idx_photosets_status (status),
    INDEX idx_photosets_deleted_at (deleted_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 图片表
CREATE TABLE IF NOT EXISTS photos (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    photoset_id BIGINT UNSIGNED NOT NULL,
    url VARCHAR(500) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    INDEX idx_photos_photoset_id (photoset_id),
    INDEX idx_photos_deleted_at (deleted_at),
    FOREIGN KEY (photoset_id) REFERENCES photosets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    name VARCHAR(50) NOT NULL,
    UNIQUE KEY idx_tags_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 套图标签关联表
CREATE TABLE IF NOT EXISTS photoset_tags (
    photoset_id BIGINT UNSIGNED NOT NULL,
    tag_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (photoset_id, tag_id),
    INDEX idx_photoset_tags_tag_id (tag_id),
    FOREIGN KEY (photoset_id) REFERENCES photosets(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 收藏表
CREATE TABLE IF NOT EXISTS favorites (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    photoset_id BIGINT UNSIGNED NOT NULL,
    UNIQUE KEY uk_user_photoset (user_id, photoset_id),
    INDEX idx_favorites_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (photoset_id) REFERENCES photosets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试管理员账号 (密码: admin123)
INSERT INTO users (created_at, updated_at, nickname, email, password_hash, role, status) VALUES
(NOW(), NOW(), 'admin', 'admin@photoset.com', '$2a$10$rJZK5x5Y5Y5Y5Y5Y5Y5Y5e1x5x5x5x5x5x5x5x5x5x5x5x5x', 'admin', 1);


