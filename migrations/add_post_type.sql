-- Add post_type column to posts table
-- Version: v1.1
-- Date: 2026-05-25
-- Description: Add post_type column to support multiple post types (dynamic, article, question, suggest, quick)

-- Add post_type column with default value 'dynamic' for backward compatibility
ALTER TABLE `posts` ADD COLUMN `post_type` VARCHAR(20) NOT NULL DEFAULT 'dynamic' AFTER `category`;

-- Update existing posts to have appropriate post_type based on their characteristics
-- Posts with title are likely articles, others are dynamic
UPDATE `posts` SET `post_type` = 'article' WHERE `title` != '' AND `title` IS NOT NULL;
UPDATE `posts` SET `post_type` = 'dynamic' WHERE `title` = '' OR `title` IS NULL;

-- Add index for post_type for better query performance
ALTER TABLE `posts` ADD INDEX `idx_post_type` (`post_type`);
