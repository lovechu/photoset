#!/bin/bash
set -e
cd /tmp/photoset_build

echo "=== 重新编译 ==="
docker run --rm \
  -v /tmp/photoset_build:/build \
  -w /build/photoset \
  golang:1.22-alpine \
  sh -c 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -o /build/main ./cmd/main.go 2>&1'

echo "BUILD_OK"
ls -lh /tmp/photoset_build/main

echo "=== 打包新镜像 ==="
cp /tmp/photoset_build/main /tmp/main
cat > /tmp/Dockerfile.backend2 << 'EOF'
FROM backend-backend:latest
COPY main /app/main
EOF

cd /tmp && docker build -f Dockerfile.backend2 -t backend-backend:latest . 2>&1 | tail -3

echo "=== docker-compose 重启 ==="
cd /opt/photoset && docker compose up -d backend 2>&1 | tail -5

sleep 4

echo "=== 验证 ==="
curl -s http://127.0.0.1:8080/api/health
echo ""
curl -s http://127.0.0.1:8080/uploads/photos/20260418013905/04 -I 2>&1 | head -3
