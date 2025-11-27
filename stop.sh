#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}==> 停止后端服务...${NC}"
if [ -f logs/backend.pid ]; then
    kill $(cat logs/backend.pid) 2>/dev/null || true
    rm logs/backend.pid
fi
pkill -f "python3 work_v3/app.py" || true

echo -e "${YELLOW}==> 停止前端服务...${NC}"
if [ -f logs/frontend.pid ]; then
    kill $(cat logs/frontend.pid) 2>/dev/null || true
    rm logs/frontend.pid
fi
pkill -f "vite" || true

echo -e "${GREEN}✅ 所有服务已停止${NC}"
