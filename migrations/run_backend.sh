#!/bin/bash
docker stop photoset-backend 2>/dev/null || true
docker rm photoset-backend 2>/dev/null || true
docker run -d \
  --name photoset-backend \
  --restart unless-stopped \
  --network photoset_photoset-network \
  -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_PORT=3306 \
  -e DB_USER=photoset \
  -e DB_PASSWORD=PhotoSet@Mysql2026 \
  -e DB_NAME=photoset \
  -e REDIS_HOST=redis \
  -e REDIS_PORT=6379 \
  -e REDIS_PASSWORD=PhotoSet@Redis2026 \
  -e JWT_SECRET=photoset-jwt-secret-2024 \
  -e STORAGE_TYPE=local \
  -e CORS_ALLOW_ORIGINS=tt.cy.mk,https://tt.cy.mk \
  -e GIN_MODE=release \
  -v photoset_uploads_volume:/app/uploads \
  -v /opt/photoset/logs:/app/logs \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --group-add 1001 \
  backend-backend:latest
