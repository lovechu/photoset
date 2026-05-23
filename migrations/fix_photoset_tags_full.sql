-- =============================================
-- 修复 photoset_tags 外键 + 重建关联数据
-- =============================================

SET FOREIGN_KEY_CHECKS = 0;

-- 1. 查看当前外键约束名（运行后记住输出的名字）
SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE TABLE_NAME = 'photoset_tags' AND TABLE_SCHEMA = DATABASE() AND REFERENCED_TABLE_NAME IS NOT NULL;

-- 2. 删除所有错误的外键（用上面查出的实际名字替换）
-- ALTER TABLE `photoset_tags` DROP FOREIGN KEY `查出的外键名1`;
-- ALTER TABLE `photoset_tags` DROP FOREIGN KEY `查出的外键名2`;

-- 3. 清空并重建 photoset_tags（因为数据全部是错的）
TRUNCATE TABLE `photoset_tags`;

-- 4. 添加正确的外键
ALTER TABLE `photoset_tags`
  ADD CONSTRAINT `fk_photoset_tags_photoset` FOREIGN KEY (`photoset_id`) REFERENCES `photosets` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_photoset_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE;

SET FOREIGN_KEY_CHECKS = 1;

-- 5. 重建关联：根据 photosets.title 和 tags.name 建立关联
--    假设套图的 title 包含模特名，标签 name 也是模特名
INSERT IGNORE INTO `photoset_tags` (`photoset_id`, `tag_id`, `created_at`, `updated_at`)
SELECT p.id, t.id, NOW(), NOW()
FROM `photosets` p
INNER JOIN `tags` t ON p.title LIKE CONCAT('%', t.name, '%')
WHERE NOT EXISTS (
  SELECT 1 FROM `photoset_tags` pt WHERE pt.photoset_id = p.id AND pt.tag_id = t.id
);

-- 6. 如果上面的模糊匹配不够精确，也可以按分类字段关联
-- INSERT IGNORE INTO `photoset_tags` (`photoset_id`, `tag_id`, `created_at`, `updated_at`)
-- SELECT p.id, t.id, NOW(), NOW()
-- FROM `photosets` p
-- INNER JOIN `tags` t ON p.category = t.name
-- WHERE NOT EXISTS (
--   SELECT 1 FROM `photoset_tags` pt WHERE pt.photoset_id = p.id AND pt.tag_id = t.id
-- );

-- 7. 验证结果
SELECT pt.photoset_id, p.title, pt.tag_id, t.name AS tag_name
FROM `photoset_tags` pt
JOIN `photosets` p ON pt.photoset_id = p.id
JOIN `tags` t ON pt.tag_id = t.id
ORDER BY pt.photoset_id;
