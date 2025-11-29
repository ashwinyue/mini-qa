.PHONY: help build build-prod run test clean deps fmt lint docker-build docker-up docker-down

# 默认目标
help:
	@echo "Eino QA System - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build             - 编译项目（开发版本）"
	@echo "  make build-prod        - 编译项目（生产版本）"
	@echo "  make run               - 运行服务"
	@echo "  make test              - 运行测试"
	@echo "  make test-coverage     - 运行测试（带覆盖率）"
	@echo "  make test-integration  - 运行集成测试"
	@echo "  make clean             - 清理构建产物"
	@echo "  make deps              - 下载依赖"
	@echo "  make fmt               - 格式化代码"
	@echo "  make lint              - 代码检查"
	@echo "  make init-dirs         - 创建必要的目录"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build      - 构建 Docker 镜像"
	@echo "  make docker-up         - 启动所有服务（Docker Compose）"
	@echo "  make docker-down       - 停止所有服务"
	@echo "  make docker-logs       - 查看应用日志"
	@echo ""
	@echo "Milvus:"
	@echo "  make milvus-up         - 启动 Milvus 服务"
	@echo "  make milvus-down       - 停止 Milvus 服务"
	@echo "  make milvus-logs       - 查看 Milvus 日志"
	@echo "  make milvus-clean      - 清理 Milvus 数据"

# 编译项目（开发版本）
build:
	@echo "Building (development)..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build complete: bin/server"

# 编译项目（生产版本）
build-prod:
	@echo "Building (production)..."
	@CGO_ENABLED=1 go build -a -installsuffix cgo -ldflags="-w -s" -o bin/server cmd/server/main.go
	@echo "Production build complete: bin/server"

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

# Docker 相关命令
docker-build:
	@echo "Building Docker image..."
	@docker build -t eino-qa:latest .
	@echo "Docker image built: eino-qa:latest"

docker-up:
	@echo "Starting all services with Docker Compose..."
	@docker-compose up -d
	@echo "All services started"
	@echo "Application: http://localhost:8080"
	@echo "Milvus: localhost:19530"
	@echo "MinIO Console: http://localhost:9001"

docker-down:
	@echo "Stopping all services..."
	@docker-compose down
	@echo "All services stopped"

docker-logs:
	@docker-compose logs -f eino-qa

docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down -v
	@docker rmi eino-qa:latest
	@echo "Docker resources cleaned"
