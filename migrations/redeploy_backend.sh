#!/bin/bash
set -e

echo "=== 1. 强制停止容器 ==="
docker stop -t 5 photoset-backend 2>/dev/null || true
docker wait photoset-backend 2>/dev/null || true
sleep 2

echo "=== 2. 确认容器状态 ==="
docker ps -a --filter name=photoset-backend --format "{{.Names}} {{.Status}}"

echo "=== 3. 复制新二进制 ==="
docker cp /tmp/photoset_build/main photoset-backend:/app/main
echo "cp done"

echo "=== 4. 检查容器内文件 ==="
docker run --rm --volumes-from photoset-backend alpine ls -lh /app/main

echo "=== 5. 启动容器 ==="
docker start photoset-backend
sleep 4

echo "=== 6. 验证 ==="
docker ps --filter name=photoset-backend --format "{{.Names}} {{.Status}}"
docker logs photoset-backend --tail 5 2>&1

echo "=== 7. API 测试 ==="
sleep 2
curl -s http://127.0.0.1:8080/api/health
