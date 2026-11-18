-- 重置管理员密码（安全提示）
-- 重要：不要使用弱口令！请先用 go run scripts/genhash.go "你的强密码" 生成 bcrypt 哈希，
-- 再将下面的 <BCRYPT_HASH_HERE> 替换为生成的哈希后执行。
UPDATE users SET password_hash='<BCRYPT_HASH_HERE>' WHERE email='admin@example.com';
SELECT id, nickname, email, role FROM users WHERE email='admin@example.com';
