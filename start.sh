#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}==> 构建 FAISS 索引...${NC}"
cd work_v3 && ../.venv_new/bin/python rag-train.py
echo -e "${YELLOW}==> 为默认租户创建索引副本...${NC}"
mkdir -p tenants/default/faiss_index
cp faiss_index/* tenants/default/faiss_index/
cd ..
echo -e "${GREEN}✅ 索引构建完成${NC}"

# 创建日志目录
mkdir -p logs

echo -e "\n${YELLOW}==> 启动后端服务 (http://localhost:8000)...${NC}"
nohup ./.venv_new/bin/python work_v3/app.py > logs/backend.log 2>&1 &
echo $! > logs/backend.pid
sleep 2
echo -e "${GREEN}✅ 后端服务已启动 (PID: $(cat logs/backend.pid))${NC}"

echo -e "\n${YELLOW}==> 启动前端服务 (http://localhost:5174)...${NC}"
cd react-go-admin
nohup npm run dev > ../logs/frontend.log 2>&1 &
echo $! > ../logs/frontend.pid
cd ..
sleep 2
echo -e "${GREEN}✅ 前端服务已启动 (PID: $(cat logs/frontend.pid))${NC}"

echo -e "\n=========================================="
echo -e "${GREEN}✅ 服务启动完成！${NC}"
echo -e "=========================================="
echo -e "前端地址: ${YELLOW}http://localhost:5174${NC}"
echo -e "后端地址: ${YELLOW}http://localhost:8000${NC}"
echo -e "健康检查: ${YELLOW}http://localhost:8000/health${NC}"
echo -e ""
echo -e "查看日志:"
echo -e "  后端: tail -f logs/backend.log"
echo -e "  前端: tail -f logs/frontend.log"
echo -e ""
echo -e "停止服务: ./stop.sh 或 make stop"
echo -e "=========================================="
