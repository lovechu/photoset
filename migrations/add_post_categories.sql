CREATE TABLE IF NOT EXISTS `post_categories` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `key` VARCHAR(64) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `description` VARCHAR(256) DEFAULT '',
    `color` VARCHAR(7) DEFAULT '#409EFF',
    `icon` VARCHAR(64) DEFAULT '',
    `sort_order` INT DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert initial 4 categories (safe to re-run, ignores duplicates)
INSERT IGNORE INTO `post_categories` (`key`, `name`, `sort_order`, `color`) VALUES
('discussion', '讨论', 1, '#409EFF'),
('qa', '问答', 2, '#67C23A'),
('showcase', '作品展示', 3, '#E6A23C'),
('suggestion', '建议', 4, '#F56C6C');
