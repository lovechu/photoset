-- 创建管理员账号（安全提示）
-- 重要：不要使用弱口令！不要把本文件提交到版本库时携带任何真实密码哈希。
-- 使用前请先生成一个强密码的 bcrypt 哈希（cost>=10），例如用 Go：
--   go run scripts/genhash.go "你的强密码"
-- 然后将下面占位符 <BCRYPT_HASH_HERE> 替换为生成的哈希，并按需修改邮箱。

INSERT INTO users (nickname, email, password_hash, role, status, created_at, updated_at)
VALUES ('admin', 'admin@example.com', '<BCRYPT_HASH_HERE>', 'admin', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE role='admin', password_hash='<BCRYPT_HASH_HERE>';

SELECT id, nickname, email, role FROM users WHERE email='admin@example.com';
