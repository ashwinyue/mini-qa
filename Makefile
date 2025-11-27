.PHONY: help install build-index start stop clean status logs test

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make install     - 安装所有依赖（Python + Node.js）"
	@echo "  make build-index - 构建 FAISS 索引（包括默认租户）"
	@echo "  make start       - 一键启动前后端服务"
	@echo "  make stop        - 停止所有服务"
	@echo "  make status      - 查看服务运行状态"
	@echo "  make logs        - 实时查看服务日志"
	@echo "  make test        - 测试后端 API"
	@echo "  make clean       - 清理临时文件和缓存"

# 安装依赖
install:
	@echo "==> 创建 Python 虚拟环境..."
	python3 -m venv venv
	@echo "==> 安装 Python 依赖..."
	./venv/bin/pip install -q fastapi uvicorn langchain langchain-community langchain-openai langgraph faiss-cpu python-dotenv pydantic dashscope redis mcp
	@echo "==> 安装前端依赖..."
	cd frontend && npm install
	@echo "✅ 依赖安装完成"

# 构建 FAISS 索引
build-index:
	@echo "==> 构建 FAISS 索引..."
	cd work_v3 && ../venv/bin/python3 rag-train.py
	@echo "==> 为默认租户创建索引副本..."
	mkdir -p work_v3/tenants/default/faiss_index
	cp work_v3/faiss_index/* work_v3/tenants/default/faiss_index/
	@echo "✅ 索引构建完成"

# 一键启动所有服务
start:
	@./start.sh

# 停止所有服务
stop:
	@./stop.sh

# 查看服务状态
status:
	@echo "==> 服务运行状态:"
	@echo ""
	@if [ -f logs/backend.pid ] && ps -p $$(cat logs/backend.pid) > /dev/null 2>&1; then \
		echo "✅ 后端服务运行中 (PID: $$(cat logs/backend.pid))"; \
	else \
		echo "❌ 后端服务未运行"; \
	fi
	@if [ -f logs/frontend.pid ] && ps -p $$(cat logs/frontend.pid) > /dev/null 2>&1; then \
		echo "✅ 前端服务运行中 (PID: $$(cat logs/frontend.pid))"; \
	else \
		echo "❌ 前端服务未运行"; \
	fi
	@echo ""
	@echo "==> 端口占用情况:"
	@lsof -i :8000 -i :5173 2>/dev/null || echo "无服务占用端口"

# 查看日志
logs:
	@echo "==> 实时查看日志 (Ctrl+C 退出)..."
	@if [ -f logs/backend.log ] && [ -f logs/frontend.log ]; then \
		tail -f logs/backend.log logs/frontend.log; \
	elif [ -f logs/backend.log ]; then \
		tail -f logs/backend.log; \
	elif [ -f logs/frontend.log ]; then \
		tail -f logs/frontend.log; \
	else \
		echo "❌ 没有找到日志文件"; \
	fi

# 测试后端 API
test:
	@echo "==> 测试后端健康检查..."
	@curl -s http://localhost:8000/health | python3 -m json.tool || echo "❌ 后端服务未响应"
	@echo ""
	@echo "==> 测试聊天接口..."
	@curl -s -X POST http://localhost:8000/chat \
		-H "Content-Type: application/json" \
		-d '{"query":"这门课程适合零基础吗？","user_id":"test"}' | python3 -m json.tool || echo "❌ 聊天接口测试失败"

# 清理临时文件
clean:
	@echo "==> 清理 Python 缓存..."
	find work_v3 -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find work_v3 -type f -name "*.pyc" -delete 2>/dev/null || true
	@echo "==> 清理日志文件..."
	rm -f logs/*.log logs/*.pid 2>/dev/null || true
	rm -f work_v3/logs/*.log 2>/dev/null || true
	@echo "✅ 清理完成"
