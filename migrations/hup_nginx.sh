#!/bin/bash
echo "=== 获取宝塔 nginx 主进程 PID ==="
NGINX_PID=$(ps aux | grep "nginx: master" | grep -v grep | awk '{print $2}' | head -1)
echo "nginx master PID: $NGINX_PID"
echo ""
echo "=== 发送 HUP 信号 reload ==="
kill -HUP $NGINX_PID
sleep 1
echo ""
echo "--- 测试 ---"
curl -s http://localhost/admin/ -I 2>&1 | head -3
curl -s http://localhost/admin/assets/index-7ut3ajdM.js -I 2>&1 | head -3
curl -s http://localhost/admin -I 2>&1 | head -3
