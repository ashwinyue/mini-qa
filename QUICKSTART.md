# 快速启动指南

## 前置要求

- Python 3.9+
- Node.js 16+
- 已配置 `.env` 文件中的 `DASHSCOPE_API_KEY`

## 一键启动

### 方式一：使用脚本（推荐）

```bash
# 首次使用：安装依赖
make install

# 启动服务（自动构建索引）
./start.sh

# 停止服务
./stop.sh
```

### 方式二：使用 Makefile

```bash
# 首次使用：安装依赖
make install

# 启动服务（自动构建索引）
make start

# 停止服务
make stop
```

服务启动后：
- 前端：http://localhost:5173
- 后端：http://localhost:8000
- 健康检查：http://localhost:8000/health

## 常用命令

```bash
# 查看所有可用命令
make help

# 查看服务状态
make status

# 查看实时日志
make logs

# 测试后端 API
make test

# 停止所有服务
make stop

# 清理缓存和日志
make clean
```

## 单独操作

```bash
# 仅构建索引
make build-index

# 仅启动后端
make start-backend

# 仅启动前端
make start-frontend
```

## 日志文件

- 后端日志：`logs/backend.log`
- 前端日志：`logs/frontend.log`

## 故障排查

### 端口被占用

```bash
# 查看端口占用
lsof -i :8000
lsof -i :5173

# 强制停止
make stop
```

### 索引构建失败

```bash
# 检查 .env 文件是否配置了 DASHSCOPE_API_KEY
cat .env

# 手动构建索引
cd work_v3
../venv/bin/python3 rag-train.py
```

### 服务无法启动

```bash
# 查看详细日志
tail -f logs/backend.log
tail -f logs/frontend.log

# 检查依赖是否安装完整
make install
```

## 开发模式

如果需要在开发模式下运行（可以看到实时输出）：

```bash
# 后端
cd work_v3
source ../venv/bin/activate
python3 app.py

# 前端（新终端）
cd frontend
npm run dev
```
