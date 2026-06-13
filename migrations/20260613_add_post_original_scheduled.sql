-- 迁移：为帖子表添加 is_original 和 scheduled_at 字段
-- is_original: 标记帖子是否为原创内容
-- scheduled_at: 定时发布时间，为空表示立即发布

ALTER TABLE posts ADD COLUMN is_original BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE posts ADD COLUMN scheduled_at TIMESTAMP NULL DEFAULT NULL;
