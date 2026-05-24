-- Community Module Tables Migration
-- Version: v1.0
-- Date: 2026-05-24
-- Description: Create all tables for community feature (posts, replies, likes, points, sensitive words, reports)

-- --------------------------------------------------------
-- Table structure for table `posts`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `posts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NOT NULL,
    `photoset_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `category` VARCHAR(20) NOT NULL DEFAULT 'discussion',
    `visibility` VARCHAR(20) NOT NULL DEFAULT 'public',
    `is_pinned` TINYINT(1) NOT NULL DEFAULT 0,
    `view_count` INT NOT NULL DEFAULT 0,
    `reply_count` INT NOT NULL DEFAULT 0,
    `like_count` INT NOT NULL DEFAULT 0,
    `status` VARCHAR(20) NOT NULL DEFAULT 'approved',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_category_visibility` (`category`, `visibility`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_is_pinned_created_at` (`is_pinned`, `created_at`),
    INDEX `idx_status` (`status`),
    CONSTRAINT `fk_posts_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `post_replies`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `post_replies` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `post_id` BIGINT UNSIGNED NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `content` TEXT NOT NULL,
    `parent_reply_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `like_count` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    INDEX `idx_post_id` (`post_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_parent_reply_id` (`parent_reply_id`),
    CONSTRAINT `fk_post_replies_post_id` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_replies_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_replies_parent_id` FOREIGN KEY (`parent_reply_id`) REFERENCES `post_replies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `post_likes`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `post_likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `post_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_post` (`user_id`, `post_id`),
    INDEX `idx_post_id` (`post_id`),
    CONSTRAINT `fk_post_likes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_likes_post_id` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `post_reply_likes`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `post_reply_likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `reply_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_reply` (`user_id`, `reply_id`),
    INDEX `idx_reply_id` (`reply_id`),
    CONSTRAINT `fk_post_reply_likes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_reply_likes_reply_id` FOREIGN KEY (`reply_id`) REFERENCES `post_replies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `user_points`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_points` (
    `user_id` BIGINT UNSIGNED NOT NULL,
    `points` INT NOT NULL DEFAULT 0,
    `level` INT NOT NULL DEFAULT 1,
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`user_id`),
    CONSTRAINT `fk_user_points_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `sensitive_words`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `sensitive_words` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `word` VARCHAR(100) NOT NULL,
    `replacement` VARCHAR(100) NOT NULL DEFAULT '***',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_word` (`word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `post_reports`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `post_reports` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `post_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `reply_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `reporter_id` BIGINT UNSIGNED NOT NULL,
    `reason` VARCHAR(500) NOT NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending',
    `handler_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `handled_at` DATETIME(3) NULL DEFAULT NULL,
    `handle_note` VARCHAR(500) NULL DEFAULT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_post_id` (`post_id`),
    INDEX `idx_reply_id` (`reply_id`),
    CONSTRAINT `fk_post_reports_post_id` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_reports_reply_id` FOREIGN KEY (`reply_id`) REFERENCES `post_replies` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_reports_reporter_id` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_post_reports_handler_id` FOREIGN KEY (`handler_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Table structure for table `user_point_logs`
-- Description: Track point change history for daily limit checking
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_point_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `points` INT NOT NULL,
    `action` VARCHAR(50) NOT NULL,
    `related_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    INDEX `idx_user_id_created_at` (`user_id`, `created_at`),
    INDEX `idx_action` (`action`),
    CONSTRAINT `fk_user_point_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Initial sensitive words (common examples)
-- --------------------------------------------------------
INSERT INTO `sensitive_words` (`word`, `replacement`, `is_active`) VALUES
('badword', '***', 1),
('spam', '***', 1),
('advertisement', '***', 1),
('fake', '***', 1);
