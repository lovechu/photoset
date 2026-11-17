-- 补齐后端代码已实现但缺少迁移脚本的表
-- 修复：这些表缺失导致对应接口返回 500（如发帖后收藏检查 /community/posts/:id/favorite/check）

-- 1. 帖子收藏表
CREATE TABLE IF NOT EXISTS `post_favorites` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `post_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_post` (`user_id`,`post_id`),
  KEY `idx_post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='帖子收藏表';

-- 2. 用户积分日志表
CREATE TABLE IF NOT EXISTS `user_point_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `points` int NOT NULL,
  `action` varchar(50) NOT NULL,
  `related_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户积分变动日志表';

-- 3. 推送 Token 表
CREATE TABLE IF NOT EXISTS `push_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(500) NOT NULL,
  `platform` varchar(20) NOT NULL,
  `device_id` varchar(100) NOT NULL DEFAULT '',
  `device_name` varchar(100) NOT NULL DEFAULT '',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `last_used_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `idx_user_platform` (`user_id`,`platform`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='推送Token表';
