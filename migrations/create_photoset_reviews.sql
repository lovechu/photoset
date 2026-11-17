-- 套图评价表
-- 修复：评价功能代码已实现但缺少建表迁移，导致 GET /api/photosets/:id/reviews 返回 500
CREATE TABLE IF NOT EXISTS `photoset_reviews` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '评价用户ID',
  `photoset_id` bigint unsigned NOT NULL COMMENT '套图ID',
  `rating` int NOT NULL COMMENT '评分 1-5',
  `tags` varchar(500) NOT NULL DEFAULT '' COMMENT '评价标签，逗号分隔',
  `content` text COMMENT '评价内容',
  `is_anonymous` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否匿名',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_photoset` (`user_id`,`photoset_id`),
  KEY `idx_photoset_id` (`photoset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='套图评价表';
