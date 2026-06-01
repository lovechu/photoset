#!/bin/bash
cd /mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend

# 启动 frontend
tmux kill-session -t frontend 2>/dev/null
tmux new-session -d -s frontend "cd /mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend/frontend && npm run dev 2>/tmp/frontend_err.log"
echo "frontend started, exit: $?"

# 启动 admin
tmux kill-session -t admin 2>/dev/null
tmux new-session -d -s admin "cd /mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend/frontend-admin && npm run dev 2>/tmp/admin_err.log"
echo "admin started, exit: $?"

# 启动 backend
tmux kill-session -t backend 2>/dev/null
tmux new-session -d -s backend "cd /mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend && go run cmd/main.go >/tmp/backend_out.log 2>/tmp/backend_err.log"
echo "backend started, exit: $?"
