# 快速开始指南

本指南将帮助你快速启动 Eino QA System 并进行基本测试。

## 前置要求

- Go 1.23+
- Docker 和 Docker Compose（用于运行 Milvus）
- DashScope API Key（阿里云通义千问）

## 步骤 1: 克隆项目

```bash
git clone <repository-url>
cd eino-qa
```

## 步骤 2: 安装依赖

```bash
make deps
```

## 步骤 3: 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入你的 DashScope API Key：

```bash
DASHSCOPE_API_KEY=your_actual_api_key_here
API_KEY_1=test_api_key_1
API_KEY_2=test_api_key_2
```

## 步骤 4: 启动 Milvus

使用 Docker Compose 启动 Milvus 向量数据库：

```bash
make milvus-up
```

等待约 30 秒，确保 Milvus 完全启动。你可以查看日志：

```bash
make milvus-logs
```

看到 "Milvus Proxy successfully started" 表示启动成功。

## 步骤 5: 初始化目录

```bash
make init-dirs
```

这将创建必要的数据和日志目录。

## 步骤 6: 运行测试

### 运行单元测试

```bash
make test
```

### 运行 Milvus 集成测试

```bash
make test-integration
```

如果测试通过，说明 Milvus 集成工作正常。

## 步骤 7: 编译项目

```bash
make build
```

编译后的可执行文件位于 `bin/server`。

## 步骤 8: 运行服务

```bash
make run
```

服务将在 `http://localhost:8080` 启动。

## 测试 API

### 健康检查

```bash
curl http://localhost:8080/health
```

### 添加向量（知识库文档）

```bash
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test_api_key_1" \
  -d '{
    "texts": [
      "Python 是一种高级编程语言，适合初学者学习",
      "Go 语言以其高性能和并发特性而闻名"
    ],
    "tenant_id": "default"
  }'
```

### 对话查询

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Python 适合初学者吗？",
    "tenant_id": "default",
    "session_id": "test-session-001"
  }'
```

## 常见问题

### Q: Milvus 启动失败

**A:** 检查 Docker 是否正常运行，端口 19530 是否被占用：

```bash
docker ps
lsof -i :19530
```

### Q: 测试失败 "Milvus not available"

**A:** 确保 Milvus 服务正在运行：

```bash
make milvus-logs
```

### Q: API Key 验证失败

**A:** 确保请求头中包含正确的 API Key：

```bash
-H "X-API-Key: test_api_key_1"
```

## 停止服务

### 停止应用服务

按 `Ctrl+C` 停止运行的服务。

### 停止 Milvus

```bash
make milvus-down
```

### 清理 Milvus 数据（可选）

```bash
make milvus-clean
```

## 下一步

- 查看 [README.md](README.md) 了解更多功能
- 查看 [设计文档](.kiro/specs/eino-qa-system/design.md) 了解架构
- 查看 [Milvus 集成文档](internal/infrastructure/repository/milvus/README.md) 了解向量数据库使用

## 开发建议

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make lint
```

### 测试覆盖率

```bash
make test-coverage
```

这将生成 `coverage.html` 文件，可以在浏览器中查看。

## 多租户测试

系统支持多租户隔离，每个租户有独立的向量 Collection：

```bash
# 租户 A 添加文档
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test_api_key_1" \
  -d '{
    "texts": ["租户 A 的文档"],
    "tenant_id": "tenant_a"
  }'

# 租户 B 添加文档
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test_api_key_1" \
  -d '{
    "texts": ["租户 B 的文档"],
    "tenant_id": "tenant_b"
  }'

# 租户 A 查询（只能查到自己的文档）
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "查询文档",
    "tenant_id": "tenant_a"
  }'
```

## 故障排查

### 查看日志

应用日志位于 `logs/app.log`（如果配置为文件输出）。

### 查看 Milvus 状态

```bash
docker exec -it milvus-standalone bash
# 在容器内
curl http://localhost:9091/healthz
```

### 重置环境

如果遇到问题，可以完全重置环境：

```bash
make clean
make milvus-clean
make init-dirs
make milvus-up
# 等待 30 秒
make build
make run
```

## 获取帮助

如有问题，请查看：
- [项目文档](README.md)
- [设计文档](.kiro/specs/eino-qa-system/design.md)
- [Milvus 官方文档](https://milvus.io/docs)
