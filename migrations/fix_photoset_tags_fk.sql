-- 修复 photoset_tags 表外键约束
-- 问题：photoset_id 错误引用 tags(id)，tag_id 错误引用 photosets(id)
-- 修复：删除错误外键，添加正确外键

SET FOREIGN_KEY_CHECKS = 0;

-- 删除 photoset_tags 表上所有外键（不管名字叫什么）
-- 先查找所有外键名
SELECT CONCAT('ALTER TABLE `photoset_tags` DROP FOREIGN KEY `', CONSTRAINT_NAME, '`;') AS drop_sql
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE TABLE_NAME = 'photoset_tags'
  AND TABLE_SCHEMA = DATABASE()
  AND REFERENCED_TABLE_NAME IS NOT NULL;

-- 手动执行上面查出的 DROP 语句，或直接用下面的通用方法：

-- 删除外键约束（尝试常见命名）
-- 如果上面查出的具体名字，替换下面的名字即可
-- ALTER TABLE `photoset_tags` DROP FOREIGN KEY `fk_photoset_tags_tag`;
-- ALTER TABLE `photoset_tags` DROP FOREIGN KEY `fk_photoset_tags_photoset`;

-- 删除 photoset_id 和 tag_id 上的索引（如果存在）
-- 注意：删除外键后可能需要重新添加索引

-- 添加正确的外键约束
ALTER TABLE `photoset_tags`
  ADD CONSTRAINT `fk_photoset_tags_photoset` FOREIGN KEY (`photoset_id`) REFERENCES `photosets` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_photoset_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE;

SET FOREIGN_KEY_CHECKS = 1;
