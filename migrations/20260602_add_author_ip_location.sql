-- 迁移：为帖子表和回帖表添加 author_ip_location 字段
-- 用于存储发帖/回帖时的实时IP归属地（省份级别），而非使用用户表的缓存值

ALTER TABLE posts ADD COLUMN author_ip_location VARCHAR(100) NOT NULL DEFAULT '';
ALTER TABLE post_replies ADD COLUMN author_ip_location VARCHAR(100) NOT NULL DEFAULT '';
