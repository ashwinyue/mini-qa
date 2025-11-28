package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestModelHandler_HandleListModels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	req := httptest.NewRequest(http.MethodGet, "/models", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleListModels(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	models := response["models"].([]any)
	assert.NotEmpty(t, models)

	current := response["current"].(map[string]any)
	assert.Equal(t, "qwen-turbo", current["chat"])
	assert.Equal(t, "text-embedding-v2", current["embedding"])
}

func TestModelHandler_HandleGetCurrentModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-plus", "text-embedding-v3")

	req := httptest.NewRequest(http.MethodGet, "/models/current", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetCurrentModel(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	current := response["current"].(map[string]any)
	assert.Equal(t, "qwen-plus", current["chat"])
	assert.Equal(t, "text-embedding-v3", current["embedding"])
}

func TestModelHandler_HandleSwitchModel_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	requestBody := SwitchModelRequest{
		Type:  "chat",
		Model: "qwen-plus",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/models/switch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleSwitchModel(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "chat", response["type"])
	assert.Equal(t, "qwen-plus", response["model"])

	// 验证模型已切换
	assert.Equal(t, "qwen-plus", handler.currentChatModel)
}

func TestModelHandler_HandleSwitchModel_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	requestBody := SwitchModelRequest{
		Type:  "invalid",
		Model: "qwen-plus",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/models/switch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleSwitchModel(c)

	// 验证有错误记录
	assert.NotEmpty(t, c.Errors)
}

func TestModelHandler_HandleSwitchModel_InvalidModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	requestBody := SwitchModelRequest{
		Type:  "chat",
		Model: "invalid-model",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/models/switch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleSwitchModel(c)

	// 验证有错误记录
	assert.NotEmpty(t, c.Errors)
}

func TestModelHandler_HandleGetModelInfo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	req := httptest.NewRequest(http.MethodGet, "/models/chat/qwen-turbo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "type", Value: "chat"},
		{Key: "name", Value: "qwen-turbo"},
	}

	handler.HandleGetModelInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	model := response["model"].(map[string]any)
	assert.Equal(t, "chat", model["type"])
	assert.Equal(t, "qwen-turbo", model["name"])
	assert.True(t, model["current"].(bool))
}

func TestModelHandler_HandleGetModelInfo_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	req := httptest.NewRequest(http.MethodGet, "/models/chat/invalid-model", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "type", Value: "chat"},
		{Key: "name", Value: "invalid-model"},
	}

	handler.HandleGetModelInfo(c)

	// 验证有错误记录
	assert.NotEmpty(t, c.Errors)
}

func TestModelHandler_HandleSwitchModel_EmbeddingModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewModelHandler("qwen-turbo", "text-embedding-v2")

	requestBody := SwitchModelRequest{
		Type:  "embedding",
		Model: "text-embedding-v3",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/models/switch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleSwitchModel(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// 验证嵌入模型已切换
	assert.Equal(t, "text-embedding-v3", handler.currentEmbedModel)
	// 聊天模型应该保持不变
	assert.Equal(t, "qwen-turbo", handler.currentChatModel)
}
