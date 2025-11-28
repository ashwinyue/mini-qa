package http

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestNewServer 测试创建服务器
func TestNewServer(t *testing.T) {
	router := gin.New()
	config := DefaultServerConfig()

	server := NewServer(router, config)

	assert.NotNil(t, server)
	assert.NotNil(t, server.router)
	assert.NotNil(t, server.httpServer)
	assert.NotNil(t, server.config)
}

// TestDefaultServerConfig 测试默认服务器配置
func TestDefaultServerConfig(t *testing.T) {
	config := DefaultServerConfig()

	assert.Equal(t, "0.0.0.0", config.Host)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.Equal(t, 1<<20, config.MaxHeaderBytes)
}

// TestServerConfigDefaults 测试服务器配置默认值
func TestServerConfigDefaults(t *testing.T) {
	router := gin.New()
	config := &ServerConfig{
		Host: "localhost",
		Port: 9090,
		// 其他字段为零值
	}

	server := NewServer(router, config)

	// 验证默认值被设置
	assert.Equal(t, 30*time.Second, server.config.ReadTimeout)
	assert.Equal(t, 30*time.Second, server.config.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.config.IdleTimeout)
	assert.Equal(t, 10*time.Second, server.config.ShutdownTimeout)
	assert.Equal(t, 1<<20, server.config.MaxHeaderBytes)
}

// TestServerAddress 测试服务器地址
func TestServerAddress(t *testing.T) {
	router := gin.New()

	tests := []struct {
		name         string
		host         string
		port         int
		expectedAddr string
	}{
		{
			name:         "default address",
			host:         "0.0.0.0",
			port:         8080,
			expectedAddr: "0.0.0.0:8080",
		},
		{
			name:         "localhost",
			host:         "localhost",
			port:         9090,
			expectedAddr: "localhost:9090",
		},
		{
			name:         "custom port",
			host:         "127.0.0.1",
			port:         3000,
			expectedAddr: "127.0.0.1:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ServerConfig{
				Host: tt.host,
				Port: tt.port,
			}

			server := NewServer(router, config)
			assert.Equal(t, tt.expectedAddr, server.httpServer.Addr)
		})
	}
}

// TestServerShutdown 测试服务器关闭
func TestServerShutdown(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	config := &ServerConfig{
		Host:            "localhost",
		Port:            0, // 使用随机端口
		ShutdownTimeout: 5 * time.Second,
	}

	server := NewServer(router, config)

	// 在 goroutine 中启动服务器
	go func() {
		_ = server.Start()
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

// TestServerGetRouter 测试获取路由器
func TestServerGetRouter(t *testing.T) {
	router := gin.New()
	config := DefaultServerConfig()

	server := NewServer(router, config)

	assert.Equal(t, router, server.GetRouter())
}

// TestServerGetHTTPServer 测试获取 HTTP 服务器
func TestServerGetHTTPServer(t *testing.T) {
	router := gin.New()
	config := DefaultServerConfig()

	server := NewServer(router, config)

	httpServer := server.GetHTTPServer()
	assert.NotNil(t, httpServer)
	assert.Equal(t, router, httpServer.Handler)
}

// TestServerTimeouts 测试服务器超时配置
func TestServerTimeouts(t *testing.T) {
	router := gin.New()
	config := &ServerConfig{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	server := NewServer(router, config)

	assert.Equal(t, 10*time.Second, server.httpServer.ReadTimeout)
	assert.Equal(t, 15*time.Second, server.httpServer.WriteTimeout)
	assert.Equal(t, 20*time.Second, server.httpServer.IdleTimeout)
}

// TestServerMaxHeaderBytes 测试最大请求头大小
func TestServerMaxHeaderBytes(t *testing.T) {
	router := gin.New()
	config := &ServerConfig{
		Host:           "localhost",
		Port:           8080,
		MaxHeaderBytes: 2 << 20, // 2 MB
	}

	server := NewServer(router, config)

	assert.Equal(t, 2<<20, server.httpServer.MaxHeaderBytes)
}

// TestServerStartError 测试服务器启动错误
func TestServerStartError(t *testing.T) {
	router := gin.New()

	// 创建第一个服务器
	config1 := &ServerConfig{
		Host: "localhost",
		Port: 18888, // 使用固定端口
	}
	server1 := NewServer(router, config1)

	// 在 goroutine 中启动第一个服务器
	go func() {
		_ = server1.Start()
	}()

	// 等待第一个服务器启动
	time.Sleep(100 * time.Millisecond)

	// 尝试在同一端口启动第二个服务器
	config2 := &ServerConfig{
		Host: "localhost",
		Port: 18888, // 相同端口
	}
	server2 := NewServer(router, config2)

	// 启动应该失败
	errChan := make(chan error, 1)
	go func() {
		errChan <- server2.Start()
	}()

	// 等待错误
	select {
	case err := <-errChan:
		assert.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("expected error but got timeout")
	}

	// 清理：关闭第一个服务器
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server1.Shutdown(ctx)
}

// TestServerShutdownTimeout 测试关闭超时
func TestServerShutdownTimeout(t *testing.T) {
	router := gin.New()

	// 添加一个慢速处理器
	router.GET("/slow", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	config := &ServerConfig{
		Host:            "localhost",
		Port:            0,               // 随机端口
		ShutdownTimeout: 1 * time.Second, // 短超时
	}

	server := NewServer(router, config)

	// 在 goroutine 中启动服务器
	go func() {
		_ = server.Start()
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 发起慢速请求
	go func() {
		// 这个请求会被超时中断
		_, _ = http.Get(fmt.Sprintf("http://localhost:%d/slow", config.Port))
	}()

	// 等待请求开始
	time.Sleep(100 * time.Millisecond)

	// 尝试关闭服务器（应该在超时后返回）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	err := server.Shutdown(ctx)
	duration := time.Since(start)

	// 关闭应该在超时时间内完成
	assert.NoError(t, err)
	assert.Less(t, duration, 3*time.Second)
}

// TestServerMultipleShutdown 测试多次关闭
func TestServerMultipleShutdown(t *testing.T) {
	router := gin.New()
	config := &ServerConfig{
		Host: "localhost",
		Port: 0,
	}

	server := NewServer(router, config)

	// 在 goroutine 中启动服务器
	go func() {
		_ = server.Start()
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 第一次关闭
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()
	err1 := server.Shutdown(ctx1)
	assert.NoError(t, err1)

	// 第二次关闭（应该立即返回）
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	err2 := server.Shutdown(ctx2)
	// 第二次关闭可能返回错误，但不应该 panic
	_ = err2
}

// BenchmarkServerStartShutdown 基准测试服务器启动和关闭
func BenchmarkServerStartShutdown(b *testing.B) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &ServerConfig{
			Host: "localhost",
			Port: 0, // 随机端口
		}

		server := NewServer(router, config)

		// 启动服务器
		go func() {
			_ = server.Start()
		}()

		// 等待启动
		time.Sleep(10 * time.Millisecond)

		// 关闭服务器
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = server.Shutdown(ctx)
		cancel()
	}
}
