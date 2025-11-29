# Eino QA System 启动指南

## 前置要求

### 1. 环境依赖

- Go 1.21 或更高版本
- Milvus 2.4+ (向量数据库)
- DashScope API Key (阿里云通义千问)

### 2. 安装 Milvus

使用 Docker Compose 快速启动 Milvus：

```bash
# 下载 Milvus docker-compose 配置
wget https://github.com/milvus-io/milvus/releases/download/v2.4.0/milvus-standalone-docker-compose.yml -O docker-compose.milvus.yml

# 启动 Milvus
docker-compose -f docker-compose.milvus.yml up -d

# 验证 Milvus 运行状态
docker-compose -f docker-compose.milvus.yml ps
```

Milvus 默认监听端口：
- gRPC: 19530
- HTTP: 9091

## 快速启动

### 1. 配置环境变量

复制 `.env.example` 为 `.env` 并填写配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```bash
# DashScope API Key (必填)
DASHSCOPE_API_KEY=your_dashscope_api_key_here

# API Keys for vector management endpoints (可选)
API_KEY_1=your_api_key_1_here
API_KEY_2=your_api_key_2_here
```

### 2. 配置系统参数

编辑 `config/config.yaml` 根据需要调整配置：

```yaml
server:
  port: 8080
  mode: debug  # debug, release

dashscope:
  api_key: ${DASHSCOPE_API_KEY}
  chat_model: qwen-turbo
  embed_model: text-embedding-v2
  embedding_dimension: 1536

milvus:
  host: localhost
  port: 19530

database:
  base_path: ./data/db

# ... 其他配置
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 构建应用

```bash
# 构建二进制文件
make build

# 或者直接使用 go build
go build -o bin/server ./cmd/server
```

### 5. 启动服务

```bash
# 使用 Makefile
make run

# 或者直接运行二进制文件
./bin/server
```

服务启动后会输出：

```
Loading configuration from: config/config.yaml
Configuration loaded successfully
Initializing dependency injection container...
Container initialized successfully
System components:
  - Logger: info level, json format
  - DashScope: model=qwen-turbo, embed=text-embedding-v2
  - Milvus: localhost:19530
  - Database: ./data/db
  - Server: port=8080, mode=debug
Starting HTTP server on port 8080...
```

## 验证服务

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：

```json
{
  "status": "healthy",
  "timestamp": "2024-11-29T10:00:00Z",
  "components": {
    "milvus": {
      "status": "healthy"
    },
    "database": {
      "status": "healthy"
    }
  }
}
```

### 2. 测试对话接口

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "你好",
    "tenant_id": "default"
  }'
```

预期响应：

```json
{
  "answer": "你好！我是智能客服助手，有什么可以帮助您的吗？",
  "route": "direct",
  "session_id": "session_xxx",
  "metadata": {
    "intent": "direct",
    "confidence": 0.95,
    "duration_ms": 234
  }
}
```

### 3. 测试向量管理接口

添加向量（需要 API Key）：

```bash
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key_1_here" \
  -d '{
    "texts": [
      "Python 课程包含基础语法、数据结构、面向对象编程等内容",
      "Go 语言课程涵盖并发编程、网络编程、微服务开发等主题"
    ],
    "tenant_id": "default"
  }'
```

## 目录结构

```
eino-qa/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── domain/                  # 领域层
│   ├── usecase/                 # 用例层
│   ├── adapter/                 # 适配层
│   └── infrastructure/          # 基础设施层
│       ├── config/              # 配置管理
│       ├── logger/              # 日志系统
│       ├── metrics/             # 指标收集
│       ├── ai/                  # AI 组件
│       ├── repository/          # 仓储实现
│       ├── tenant/              # 多租户管理
│       └── container/           # 依赖注入容器
├── config/
│   └── config.yaml              # 配置文件
├── data/
│   └── db/                      # SQLite 数据库文件
├── bin/                         # 编译后的二进制文件
├── .env                         # 环境变量（不提交到 git）
├── .env.example                 # 环境变量示例
├── Makefile                     # 构建脚本
└── go.mod                       # Go 模块定义
```

## 常用命令

### 开发模式

```bash
# 运行服务（自动重启）
make dev

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint
```

### 生产模式

```bash
# 构建生产版本
make build-prod

# 运行生产版本
./bin/server
```

## 多租户使用

系统支持多租户隔离，每个租户拥有独立的：
- Milvus Collection（向量数据）
- SQLite 数据库文件（订单、会话数据）

### 创建租户

租户会在首次使用时自动创建，只需在请求中指定 `tenant_id`：

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant1" \
  -d '{
    "query": "Python 课程有哪些内容？"
  }'
```

或者在请求体中指定：

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Python 课程有哪些内容？",
    "tenant_id": "tenant1"
  }'
```

### 租户数据隔离

- Milvus Collection: `kb_{tenant_id}`
- SQLite 数据库: `./data/db/{tenant_id}.db`

## 优雅关闭

服务支持优雅关闭，按 `Ctrl+C` 或发送 `SIGTERM` 信号：

```bash
# 发送 SIGTERM 信号
kill -TERM <pid>
```

关闭过程：
1. 停止接收新请求
2. 等待现有请求完成（最多 10 秒）
3. 关闭数据库连接
4. 关闭 Milvus 连接
5. 释放其他资源

## 故障排查

### 1. Milvus 连接失败

错误信息：
```
failed to connect to milvus: connection refused
```

解决方法：
- 检查 Milvus 是否运行：`docker-compose -f docker-compose.milvus.yml ps`
- 检查端口是否正确：默认 19530
- 检查防火墙设置

### 2. DashScope API 调用失败

错误信息：
```
failed to initialize chat model: invalid api key
```

解决方法：
- 检查 `.env` 文件中的 `DASHSCOPE_API_KEY` 是否正确
- 验证 API Key 是否有效：访问 [DashScope 控制台](https://dashscope.console.aliyun.com/)

### 3. 数据库文件权限错误

错误信息：
```
failed to open database: permission denied
```

解决方法：
- 检查 `data/db` 目录权限：`chmod 755 data/db`
- 确保应用有写入权限

### 4. 端口被占用

错误信息：
```
failed to start server: address already in use
```

解决方法：
- 修改 `config/config.yaml` 中的 `server.port`
- 或者停止占用端口的进程：`lsof -ti:8080 | xargs kill`

## 日志查看

### 标准输出日志

默认情况下，日志输出到标准输出（stdout）：

```bash
./bin/server | jq .
```

### 文件日志

修改 `config/config.yaml` 启用文件日志：

```yaml
logging:
  output: file
  file_path: ./logs/app.log
```

查看日志：

```bash
tail -f logs/app.log | jq .
```

## 性能监控

### 指标接口

访问指标接口查看系统运行状态：

```bash
curl http://localhost:8080/health/metrics
```

返回指标包括：
- 请求总数
- 平均响应时间
- 各意图类型的请求分布
- 错误率

### 健康检查

- 存活检查：`GET /health/live`
- 就绪检查：`GET /health/ready`
- 完整健康检查：`GET /health`

## 下一步

- 阅读 [API 文档](docs/api_specification.md)
- 查看 [架构设计](docs/architecture_design.md)
- 了解 [部署指南](docs/deployment_guide.md)
- 参考 [示例代码](examples/)

## 获取帮助

如有问题，请：
1. 查看日志文件
2. 检查配置文件
3. 参考故障排查章节
4. 提交 Issue
