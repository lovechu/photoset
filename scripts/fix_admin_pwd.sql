UPDATE users SET password_hash='$2a$10$6HTkZ8Cz..1LFbO8nxBmquouXGrgeRwjalJFWOI8uzw5RHlaiJx2C' WHERE email='admin@photoset.dev';
SELECT id, nickname, email, role FROM users WHERE email='admin@photoset.dev';
