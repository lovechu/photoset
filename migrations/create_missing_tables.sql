-- 创建 site_settings 表
CREATE TABLE IF NOT EXISTS site_settings (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `key` VARCHAR(100) NOT NULL,
    value TEXT,
    `group` VARCHAR(50) NOT NULL DEFAULT 'general',
    UNIQUE KEY uk_key (`key`),
    KEY idx_group (`group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认配置
INSERT INTO site_settings (`key`, value, `group`) VALUES
('site_name', 'PhotoSet', 'general'),
('site_description', '精美套图分享平台', 'general'),
('contact_email', 'admin@photoset.com', 'general'),
('footer_text', '© 2024 PhotoSet. All rights reserved.', 'general')
ON DUPLICATE KEY UPDATE value = VALUES(value);

-- 创建 categories 表
CREATE TABLE IF NOT EXISTS categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    description VARCHAR(200) DEFAULT '',
    sort_order INT DEFAULT 0,
    UNIQUE KEY uk_name (name),
    UNIQUE KEY uk_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
