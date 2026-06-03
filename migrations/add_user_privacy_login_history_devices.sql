-- User Privacy Settings, Login History, and Devices Migration
-- Version: v1.0
-- Date: 2026-06-03
-- Description: Create tables for user privacy settings, login history, and device management

-- --------------------------------------------------------
-- Table structure for table `user_privacy_settings`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_privacy_settings` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `show_profile` TINYINT(1) NOT NULL DEFAULT 1,
    `show_posts` TINYINT(1) NOT NULL DEFAULT 1,
    `show_favorites` TINYINT(1) NOT NULL DEFAULT 1,
    `allow_search` TINYINT(1) NOT NULL DEFAULT 1,
    `allow_message` TINYINT(1) NOT NULL DEFAULT 1,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_id` (`user_id`),
    CONSTRAINT `fk_user_privacy_settings_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `login_history`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `login_history` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `ip` VARCHAR(50) NOT NULL,
    `ip_location` VARCHAR(100) NULL DEFAULT NULL,
    `user_agent` TEXT NULL DEFAULT NULL,
    `device` VARCHAR(100) NULL DEFAULT NULL,
    `browser` VARCHAR(100) NULL DEFAULT NULL,
    `os` VARCHAR(100) NULL DEFAULT NULL,
    `login_type` VARCHAR(20) NOT NULL DEFAULT 'password',
    `success` TINYINT(1) NOT NULL DEFAULT 1,
    `fail_reason` VARCHAR(200) NULL DEFAULT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_user_id_created_at` (`user_id`, `created_at`),
    CONSTRAINT `fk_login_history_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `user_devices`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_devices` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `device_id` VARCHAR(100) NOT NULL,
    `device_name` VARCHAR(100) NULL DEFAULT NULL,
    `device_type` VARCHAR(20) NOT NULL,
    `os` VARCHAR(50) NULL DEFAULT NULL,
    `browser` VARCHAR(50) NULL DEFAULT NULL,
    `ip` VARCHAR(50) NULL DEFAULT NULL,
    `ip_location` VARCHAR(100) NULL DEFAULT NULL,
    `user_agent` TEXT NULL DEFAULT NULL,
    `last_active_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `is_active` TINYINT(1) NOT NULL DEFAULT 1,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_device` (`user_id`, `device_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_is_active` (`is_active`),
    INDEX `idx_last_active_at` (`last_active_at`),
    CONSTRAINT `fk_user_devices_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Create event to clean up old login history (older than 90 days)
-- Note: This requires MySQL Event Scheduler to be enabled
-- --------------------------------------------------------
-- SET GLOBAL event_scheduler = ON;
-- CREATE EVENT IF NOT EXISTS cleanup_login_history
-- ON SCHEDULE EVERY 1 DAY
-- DO
--   DELETE FROM login_history WHERE created_at < NOW() - INTERVAL 90 DAY;

-- --------------------------------------------------------
-- Create event to clean up inactive devices (older than 30 days)
-- --------------------------------------------------------
-- CREATE EVENT IF NOT EXISTS cleanup_inactive_devices
-- ON SCHEDULE EVERY 1 DAY
-- DO
--   DELETE FROM user_devices WHERE is_active = false AND last_active_at < NOW() - INTERVAL 30 DAY;