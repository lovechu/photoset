INSERT INTO pages (slug, title, content_md, user_id, status, created_at, updated_at)
VALUES 
  ('about', '关于我们', '# 关于我们\n\n我们是摄影爱好者社区，致力于分享美丽瞬间。', 1, 'published', UNIX_TIMESTAMP()*1000, UNIX_TIMESTAMP()*1000),
  ('terms', '使用协议', '# 使用协议\n\n请遵守平台规则。', 1, 'published', UNIX_TIMESTAMP()*1000, UNIX_TIMESTAMP()*1000),
  ('privacy', '隐私政策', '# 隐私政策\n\n我们非常重视您的隐私。', 1, 'published', UNIX_TIMESTAMP()*1000, UNIX_TIMESTAMP()*1000),
  ('faq', '常见问题', '# 常见问题\n\n### 如何上传作品？\n请注册并进入创作者后台。', 1, 'published', UNIX_TIMESTAMP()*1000, UNIX_TIMESTAMP()*1000),
  ('contact', '联系我们', '# 联系我们\n\n邮箱：support@photoset.io', 1, 'published', UNIX_TIMESTAMP()*1000, UNIX_TIMESTAMP()*1000)
ON DUPLICATE KEY UPDATE slug = slug;