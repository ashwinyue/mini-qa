package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/container"

	"github.com/joho/godotenv"
)

const (
	defaultConfigPath = "config/config.yaml"
	shutdownTimeout   = 10 * time.Second
)

func main() {
	// 1. 加载环境变量
	// 需求: 1.1 - 环境变量管理
	if err := loadEnv(); err != nil {
		log.Printf("Warning: %v", err)
	}

	// 2. 加载配置
	// 需求: 1.1 - 配置文件结构
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 3. 初始化依赖注入容器
	// 需求: 1.1, 6.1 - 组件初始化
	c, err := initContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer closeContainer(c)

	// 4. 启动 HTTP 服务器
	// 需求: 6.1 - HTTP 服务器启动
	if err := runServer(c); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// loadEnv 加载环境变量
func loadEnv() error {
	// 尝试从当前目录加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// 尝试从父目录加载
		if err := godotenv.Load("../.env"); err != nil {
			return fmt.Errorf(".env file not found, using system environment variables")
		}
	}
	log.Println("Environment variables loaded")
	return nil
}

// loadConfig 加载配置文件
func loadConfig() (*config.Config, error) {
	// 从环境变量或默认路径获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	// 如果配置文件不存在，尝试从父目录查找
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		parentPath := filepath.Join("..", configPath)
		if _, err := os.Stat(parentPath); err == nil {
			configPath = parentPath
		}
	}

	log.Printf("Loading configuration from: %s", configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}

// initContainer 初始化依赖注入容器
func initContainer(cfg *config.Config) (*container.Container, error) {
	log.Println("Initializing dependency injection container...")

	c, err := container.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	log.Println("Container initialized successfully")
	log.Printf("System components:")
	log.Printf("  - Logger: %s level, %s format", cfg.Logging.Level, cfg.Logging.Format)
	log.Printf("  - DashScope: model=%s, embed=%s", cfg.DashScope.ChatModel, cfg.DashScope.EmbedModel)
	log.Printf("  - Milvus: %s:%d", cfg.Milvus.Host, cfg.Milvus.Port)
	log.Printf("  - Database: %s", cfg.Database.BasePath)
	log.Printf("  - Server: port=%d, mode=%s", cfg.Server.Port, cfg.Server.Mode)

	return c, nil
}

// runServer 运行 HTTP 服务器
func runServer(c *container.Container) error {
	// 创建错误通道
	errChan := make(chan error, 1)

	// 在 goroutine 中启动服务器
	// 需求: 6.1 - 服务器启动
	go func() {
		log.Printf("Starting HTTP server on port %d...", c.Config.Server.Port)
		if err := c.Server.Start(); err != nil {
			errChan <- err
		}
	}()

	// 等待中断信号
	// 需求: 6.1 - 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号或错误
	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
		log.Println("Initiating graceful shutdown...")
	}

	// 创建关闭超时 context
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// 优雅关闭服务器
	if err := c.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server shutdown completed")
	return nil
}

// closeContainer 关闭容器，释放资源
func closeContainer(c *container.Container) {
	if c == nil {
		return
	}

	log.Println("Closing container and releasing resources...")

	if err := c.Close(); err != nil {
		log.Printf("Error closing container: %v", err)
	} else {
		log.Println("Container closed successfully")
	}
}
