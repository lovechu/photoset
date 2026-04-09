INSERT INTO users (nickname, email, password_hash, role, status, created_at, updated_at)
VALUES ('admin', 'admin@photoset.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE role='admin', password_hash='$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy';

SELECT id, nickname, email, role FROM users WHERE email='admin@photoset.dev';
