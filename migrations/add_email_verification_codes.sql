-- 邮箱验证码表
CREATE TABLE IF NOT EXISTS `email_verification_codes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `email` varchar(100) NOT NULL COMMENT '邮箱地址',
  `code` varchar(6) NOT NULL COMMENT '6位验证码',
  `used` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已使用',
  `expire` datetime NOT NULL COMMENT '过期时间',
  `purpose` varchar(20) NOT NULL DEFAULT 'bind' COMMENT '用途: verify-注册验证, bind-绑定邮箱',
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`),
  KEY `idx_expire` (`expire`),
  KEY `idx_email_purpose` (`email`, `purpose`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮箱验证码表';
