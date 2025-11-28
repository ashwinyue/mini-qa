.PHONY: help build run test clean deps fmt lint

# 默认目标
help:
	@echo "Eino QA System - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build             - 编译项目"
	@echo "  make run               - 运行服务"
	@echo "  make test              - 运行测试"
	@echo "  make test-integration  - 运行集成测试"
	@echo "  make clean             - 清理构建产物"
	@echo "  make deps              - 下载依赖"
	@echo "  make fmt               - 格式化代码"
	@echo "  make lint              - 代码检查"
	@echo ""
	@echo "Milvus Commands:"
	@echo "  make milvus-up         - 启动 Milvus 服务"
	@echo "  make milvus-down       - 停止 Milvus 服务"
	@echo "  make milvus-logs       - 查看 Milvus 日志"
	@echo "  make milvus-clean      - 清理 Milvus 数据"

# 编译项目
build:
	@echo "Building..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build complete: bin/server"

# 运行服务
run:
	@echo "Starting server..."
	@go run cmd/server/main.go

# 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 运行测试（带覆盖率）
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# 代码检查
lint:
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Lint complete"

# 创建必要的目录
init-dirs:
	@mkdir -p data/sqlite
	@mkdir -p logs
	@echo "Directories created"

# Milvus 相关命令
milvus-up:
	@echo "Starting Milvus..."
	@docker-compose -f docker-compose.milvus.yml up -d
	@echo "Milvus started at localhost:19530"

milvus-down:
	@echo "Stopping Milvus..."
	@docker-compose -f docker-compose.milvus.yml down
	@echo "Milvus stopped"

milvus-logs:
	@docker-compose -f docker-compose.milvus.yml logs -f standalone

milvus-clean:
	@echo "Cleaning Milvus data..."
	@docker-compose -f docker-compose.milvus.yml down -v
	@rm -rf volumes/
	@echo "Milvus data cleaned"

# 运行集成测试（需要 Milvus 运行）
test-integration:
	@echo "Running integration tests..."
	@go test -v ./internal/infrastructure/repository/milvus/...
	@echo "Integration tests complete"
