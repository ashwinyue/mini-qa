package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	// 可以注入依赖来检查各个组件的健康状态
	checkMilvus    func(ctx context.Context) error
	checkDB        func(ctx context.Context) error
	checkDashScope func(ctx context.Context) error
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// WithMilvusCheck 设置 Milvus 健康检查函数
func (h *HealthHandler) WithMilvusCheck(check func(ctx context.Context) error) *HealthHandler {
	h.checkMilvus = check
	return h
}

// WithDBCheck 设置数据库健康检查函数
func (h *HealthHandler) WithDBCheck(check func(ctx context.Context) error) *HealthHandler {
	h.checkDB = check
	return h
}

// WithDashScopeCheck 设置 DashScope 健康检查函数
func (h *HealthHandler) WithDashScopeCheck(check func(ctx context.Context) error) *HealthHandler {
	h.checkDashScope = check
	return h
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status     string                     `json:"status"`
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
}

// ComponentHealth 组件健康状态
type ComponentHealth struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HandleHealth 处理健康检查请求
// GET /health
// 需求: 7.5
func (h *HealthHandler) HandleHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response := &HealthResponse{
		Status:     "healthy",
		Timestamp:  time.Now().Format(time.RFC3339),
		Components: make(map[string]ComponentHealth),
	}

	// 检查各个组件的健康状态
	allHealthy := true

	// 检查 Milvus
	if h.checkMilvus != nil {
		if err := h.checkMilvus(ctx); err != nil {
			response.Components["milvus"] = ComponentHealth{
				Status:  "unhealthy",
				Message: err.Error(),
			}
			allHealthy = false
		} else {
			response.Components["milvus"] = ComponentHealth{
				Status: "healthy",
			}
		}
	}

	// 检查数据库
	if h.checkDB != nil {
		if err := h.checkDB(ctx); err != nil {
			response.Components["database"] = ComponentHealth{
				Status:  "unhealthy",
				Message: err.Error(),
			}
			allHealthy = false
		} else {
			response.Components["database"] = ComponentHealth{
				Status: "healthy",
			}
		}
	}

	// 检查 DashScope
	if h.checkDashScope != nil {
		if err := h.checkDashScope(ctx); err != nil {
			response.Components["dashscope"] = ComponentHealth{
				Status:  "unhealthy",
				Message: err.Error(),
			}
			allHealthy = false
		} else {
			response.Components["dashscope"] = ComponentHealth{
				Status: "healthy",
			}
		}
	}

	// 设置整体状态
	if !allHealthy {
		response.Status = "degraded"
	}

	// 返回响应
	statusCode := http.StatusOK
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// HandleLiveness 处理存活检查请求
// GET /health/live
// 简单的存活检查，只要服务能响应就返回 200
func (h *HealthHandler) HandleLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// HandleReadiness 处理就绪检查请求
// GET /health/ready
// 检查服务是否准备好接收流量
func (h *HealthHandler) HandleReadiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	ready := true
	components := make(map[string]string)

	// 检查关键组件
	if h.checkMilvus != nil {
		if err := h.checkMilvus(ctx); err != nil {
			ready = false
			components["milvus"] = "not ready"
		} else {
			components["milvus"] = "ready"
		}
	}

	if h.checkDB != nil {
		if err := h.checkDB(ctx); err != nil {
			ready = false
			components["database"] = "not ready"
		} else {
			components["database"] = "ready"
		}
	}

	statusCode := http.StatusOK
	status := "ready"
	if !ready {
		statusCode = http.StatusServiceUnavailable
		status = "not ready"
	}

	c.JSON(statusCode, gin.H{
		"status":     status,
		"timestamp":  time.Now().Format(time.RFC3339),
		"components": components,
	})
}
