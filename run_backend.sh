#!/bin/bash
cd /mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend
go run cmd/main.go > /tmp/backend_out.log 2> /tmp/backend_err.log
echo "EXIT:$?" >> /tmp/backend_out.log
