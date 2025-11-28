package handler

import (
	"fmt"
	"net/http"

	"eino-qa/internal/adapter/http/middleware"
	"eino-qa/internal/usecase/vector"

	"github.com/gin-gonic/gin"
)

// VectorHandler 向量管理处理器
type VectorHandler struct {
	vectorUseCase vector.VectorUseCaseInterface
}

// NewVectorHandler 创建向量管理处理器
func NewVectorHandler(vectorUseCase vector.VectorUseCaseInterface) *VectorHandler {
	return &VectorHandler{
		vectorUseCase: vectorUseCase,
	}
}

// AddVectorRequestDTO 添加向量请求 DTO
type AddVectorRequestDTO struct {
	Texts    []string       `json:"texts" binding:"required"`
	TenantID string         `json:"tenant_id"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// AddVectorResponseDTO 添加向量响应 DTO
type AddVectorResponseDTO struct {
	Success     bool     `json:"success"`
	DocumentIDs []string `json:"document_ids"`
	Count       int      `json:"count"`
	Message     string   `json:"message"`
}

// DeleteVectorRequestDTO 删除向量请求 DTO
type DeleteVectorRequestDTO struct {
	IDs      []string `json:"ids" binding:"required"`
	TenantID string   `json:"tenant_id"`
}

// DeleteVectorResponseDTO 删除向量响应 DTO
type DeleteVectorResponseDTO struct {
	Success      bool   `json:"success"`
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// HandleAddVectors 处理添加向量请求
// POST /api/v1/vectors/items
// 需求: 9.1, 9.2, 9.3, 9.5
func (h *VectorHandler) HandleAddVectors(c *gin.Context) {
	var req AddVectorRequestDTO

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewBadRequestError(fmt.Sprintf("invalid request: %s", err.Error())))
		return
	}

	// 从 context 获取租户 ID（由中间件设置）
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok && req.TenantID == "" {
			req.TenantID = tid
		}
	}

	// 如果仍然没有租户 ID，使用默认值
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	// 验证请求
	if len(req.Texts) == 0 {
		c.Error(middleware.NewBadRequestError("texts cannot be empty"))
		return
	}

	// 构建用例请求
	useCaseReq := &vector.AddVectorRequest{
		Texts:    req.Texts,
		TenantID: req.TenantID,
		Metadata: req.Metadata,
	}

	// 执行添加向量用例
	resp, err := h.vectorUseCase.AddVectors(c.Request.Context(), useCaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为 DTO
	dto := &AddVectorResponseDTO{
		Success:     resp.Success,
		DocumentIDs: resp.DocumentIDs,
		Count:       resp.Count,
		Message:     resp.Message,
	}

	// 返回响应
	c.JSON(http.StatusOK, dto)
}

// HandleDeleteVectors 处理删除向量请求
// DELETE /api/v1/vectors/items
// 需求: 9.1, 9.4, 9.5
func (h *VectorHandler) HandleDeleteVectors(c *gin.Context) {
	var req DeleteVectorRequestDTO

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewBadRequestError(fmt.Sprintf("invalid request: %s", err.Error())))
		return
	}

	// 从 context 获取租户 ID（由中间件设置）
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok && req.TenantID == "" {
			req.TenantID = tid
		}
	}

	// 如果仍然没有租户 ID，使用默认值
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	// 验证请求
	if len(req.IDs) == 0 {
		c.Error(middleware.NewBadRequestError("ids cannot be empty"))
		return
	}

	// 构建用例请求
	useCaseReq := &vector.DeleteVectorRequest{
		IDs:      req.IDs,
		TenantID: req.TenantID,
	}

	// 执行删除向量用例
	resp, err := h.vectorUseCase.DeleteVectors(c.Request.Context(), useCaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为 DTO
	dto := &DeleteVectorResponseDTO{
		Success:      resp.Success,
		DeletedCount: resp.DeletedCount,
		Message:      resp.Message,
	}

	// 返回响应
	c.JSON(http.StatusOK, dto)
}

// HandleGetVectorCount 处理获取向量数量请求
// GET /api/v1/vectors/count
func (h *VectorHandler) HandleGetVectorCount(c *gin.Context) {
	// 从 context 获取租户 ID
	tenantID := "default"
	if tid, exists := c.Get("tenant_id"); exists {
		if id, ok := tid.(string); ok {
			tenantID = id
		}
	}

	// 获取向量数量
	count, err := h.vectorUseCase.GetVectorCount(c.Request.Context(), tenantID)
	if err != nil {
		c.Error(err)
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"tenant_id": tenantID,
		"count":     count,
	})
}

// HandleGetVector 处理获取向量请求
// GET /api/v1/vectors/items/:id
func (h *VectorHandler) HandleGetVector(c *gin.Context) {
	// 获取文档 ID
	id := c.Param("id")
	if id == "" {
		c.Error(middleware.NewBadRequestError("id is required"))
		return
	}

	// 从 context 获取租户 ID
	tenantID := "default"
	if tid, exists := c.Get("tenant_id"); exists {
		if t, ok := tid.(string); ok {
			tenantID = t
		}
	}

	// 获取向量
	doc, err := h.vectorUseCase.GetVectorByID(c.Request.Context(), id, tenantID)
	if err != nil {
		c.Error(middleware.NewNotFoundError(fmt.Sprintf("vector not found: %s", id)))
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"document": gin.H{
			"id":       doc.ID,
			"content":  doc.Content,
			"metadata": doc.Metadata,
			"score":    doc.Score,
		},
	})
}
