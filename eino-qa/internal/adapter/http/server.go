package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server HTTP 服务器
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *ServerConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host            string        // 监听地址
	Port            int           // 监听端口
	ReadTimeout     time.Duration // 读取超时
	WriteTimeout    time.Duration // 写入超时
	IdleTimeout     time.Duration // 空闲超时
	ShutdownTimeout time.Duration // 优雅关闭超时
	MaxHeaderBytes  int           // 最大请求头大小
	EnablePprof     bool          // 是否启用 pprof
	EnableMetrics   bool          // 是否启用指标
}

// NewServer 创建 HTTP 服务器
// 需求: 6.1 - Gin HTTP 服务器在配置的端口上监听请求
func NewServer(router *gin.Engine, config *ServerConfig) *Server {
	// 设置默认值
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 30 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 30 * time.Second
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 60 * time.Second
	}
	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 10 * time.Second
	}
	if config.MaxHeaderBytes == 0 {
		config.MaxHeaderBytes = 1 << 20 // 1 MB
	}

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		router:     router,
		httpServer: httpServer,
		config:     config,
	}
}

// Start 启动服务器
// 需求: 6.1 - 服务器启动
func (s *Server) Start() error {
	fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)

	// 启动服务器
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭服务器
// 需求: 6.1 - 服务器优雅关闭
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")

	// 优雅关闭服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	fmt.Println("HTTP server stopped")
	return nil
}

// Run 运行服务器并处理优雅关闭
// 需求: 6.1 - 服务器启动和优雅关闭
func (s *Server) Run() error {
	// 创建错误通道
	errChan := make(chan error, 1)

	// 在 goroutine 中启动服务器
	go func() {
		if err := s.Start(); err != nil {
			errChan <- err
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号或错误
	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		fmt.Printf("Received signal: %v\n", sig)
	}

	// 创建关闭超时 context
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	// 优雅关闭服务器
	return s.Shutdown(ctx)
}

// GetRouter 获取路由器
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetHTTPServer 获取底层 HTTP 服务器
func (s *Server) GetHTTPServer() *http.Server {
	return s.httpServer
}

// DefaultServerConfig 创建默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:            "0.0.0.0",
		Port:            8080,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1 MB
	}
}
