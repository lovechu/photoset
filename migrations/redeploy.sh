#!/bin/bash
sshpass -p 'gajaEWUU9245' ssh -o StrictHostKeyChecking=no root@103.149.93.111 << 'ENDSSH'
set -e
echo "=== Check binary timestamp in container ==="
docker exec photoset-backend ls -la /app/main
echo ""
echo "=== Check /tmp/photoset_build/main timestamp ==="
ls -la /tmp/photoset_build/main
echo ""
echo "=== Copy fresh binary to container ==="
docker cp /tmp/photoset_build/main photoset-backend:/app/main
echo ""
echo "=== Restart container ==="
docker restart photoset-backend
sleep 3
echo ""
echo "=== Test detail API ==="
curl -s http://127.0.0.1:8080/api/photosets/4
ENDSSH
