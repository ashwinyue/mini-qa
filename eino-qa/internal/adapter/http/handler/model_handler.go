package handler

import (
	"fmt"
	"net/http"

	"eino-qa/internal/adapter/http/middleware"

	"github.com/gin-gonic/gin"
)

// ModelHandler 模型管理处理器
type ModelHandler struct {
	// 当前使用的模型配置
	currentChatModel  string
	currentEmbedModel string
	availableModels   map[string][]string
}

// NewModelHandler 创建模型管理处理器
func NewModelHandler(chatModel, embedModel string) *ModelHandler {
	return &ModelHandler{
		currentChatModel:  chatModel,
		currentEmbedModel: embedModel,
		availableModels: map[string][]string{
			"chat": {
				"qwen-turbo",
				"qwen-plus",
				"qwen-max",
				"qwen-max-longcontext",
			},
			"embedding": {
				"text-embedding-v1",
				"text-embedding-v2",
				"text-embedding-v3",
			},
		},
	}
}

// ModelInfo 模型信息
type ModelInfo struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Current bool   `json:"current"`
}

// SwitchModelRequest 切换模型请求
type SwitchModelRequest struct {
	Type  string `json:"type" binding:"required"`  // "chat" or "embedding"
	Model string `json:"model" binding:"required"` // 模型名称
}

// HandleListModels 处理列出可用模型请求
// GET /models
func (h *ModelHandler) HandleListModels(c *gin.Context) {
	models := make([]ModelInfo, 0)

	// 添加聊天模型
	for _, model := range h.availableModels["chat"] {
		models = append(models, ModelInfo{
			Type:    "chat",
			Name:    model,
			Current: model == h.currentChatModel,
		})
	}

	// 添加嵌入模型
	for _, model := range h.availableModels["embedding"] {
		models = append(models, ModelInfo{
			Type:    "embedding",
			Name:    model,
			Current: model == h.currentEmbedModel,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"models":  models,
		"current": gin.H{
			"chat":      h.currentChatModel,
			"embedding": h.currentEmbedModel,
		},
	})
}

// HandleGetCurrentModel 处理获取当前模型请求
// GET /models/current
func (h *ModelHandler) HandleGetCurrentModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"current": gin.H{
			"chat":      h.currentChatModel,
			"embedding": h.currentEmbedModel,
		},
	})
}

// HandleSwitchModel 处理切换模型请求
// POST /models/switch
// 注意：实际的模型切换需要重新初始化 Eino 组件，这里只是演示接口
func (h *ModelHandler) HandleSwitchModel(c *gin.Context) {
	var req SwitchModelRequest

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewBadRequestError(fmt.Sprintf("invalid request: %s", err.Error())))
		return
	}

	// 验证模型类型
	if req.Type != "chat" && req.Type != "embedding" {
		c.Error(middleware.NewBadRequestError("type must be 'chat' or 'embedding'"))
		return
	}

	// 验证模型是否可用
	availableModels, ok := h.availableModels[req.Type]
	if !ok {
		c.Error(middleware.NewBadRequestError(fmt.Sprintf("unknown model type: %s", req.Type)))
		return
	}

	modelExists := false
	for _, model := range availableModels {
		if model == req.Model {
			modelExists = true
			break
		}
	}

	if !modelExists {
		c.Error(middleware.NewBadRequestError(fmt.Sprintf("model '%s' not available for type '%s'", req.Model, req.Type)))
		return
	}

	// 更新当前模型
	// 注意：实际应用中，这里需要重新初始化相应的 Eino 组件
	// 这里只是更新配置，实际的模型切换需要在应用层面处理
	oldModel := ""
	switch req.Type {
	case "chat":
		oldModel = h.currentChatModel
		h.currentChatModel = req.Model
	case "embedding":
		oldModel = h.currentEmbedModel
		h.currentEmbedModel = req.Model
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("model switched from '%s' to '%s'", oldModel, req.Model),
		"type":    req.Type,
		"model":   req.Model,
		"note":    "Model switch will take effect on next request. Some components may need restart.",
	})
}

// HandleGetModelInfo 处理获取模型信息请求
// GET /models/:type/:name
func (h *ModelHandler) HandleGetModelInfo(c *gin.Context) {
	modelType := c.Param("type")
	modelName := c.Param("name")

	if modelType == "" || modelName == "" {
		c.Error(middleware.NewBadRequestError("type and name are required"))
		return
	}

	// 验证模型类型
	availableModels, ok := h.availableModels[modelType]
	if !ok {
		c.Error(middleware.NewNotFoundError(fmt.Sprintf("unknown model type: %s", modelType)))
		return
	}

	// 检查模型是否存在
	modelExists := false
	for _, model := range availableModels {
		if model == modelName {
			modelExists = true
			break
		}
	}

	if !modelExists {
		c.Error(middleware.NewNotFoundError(fmt.Sprintf("model '%s' not found for type '%s'", modelName, modelType)))
		return
	}

	// 返回模型信息
	isCurrent := false
	if modelType == "chat" {
		isCurrent = modelName == h.currentChatModel
	} else if modelType == "embedding" {
		isCurrent = modelName == h.currentEmbedModel
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"model": ModelInfo{
			Type:    modelType,
			Name:    modelName,
			Current: isCurrent,
		},
	})
}
