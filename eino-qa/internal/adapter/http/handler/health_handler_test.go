package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HandleHealth_AllHealthy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler().
		WithMilvusCheck(func(ctx context.Context) error { return nil }).
		WithDBCheck(func(ctx context.Context) error { return nil }).
		WithDashScopeCheck(func(ctx context.Context) error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleHealth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
	assert.NotEmpty(t, response.Timestamp)
	assert.Equal(t, "healthy", response.Components["milvus"].Status)
	assert.Equal(t, "healthy", response.Components["database"].Status)
	assert.Equal(t, "healthy", response.Components["dashscope"].Status)
}

func TestHealthHandler_HandleHealth_PartiallyUnhealthy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler().
		WithMilvusCheck(func(ctx context.Context) error { return errors.New("connection failed") }).
		WithDBCheck(func(ctx context.Context) error { return nil }).
		WithDashScopeCheck(func(ctx context.Context) error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleHealth(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "degraded", response.Status)
	assert.Equal(t, "unhealthy", response.Components["milvus"].Status)
	assert.Equal(t, "connection failed", response.Components["milvus"].Message)
	assert.Equal(t, "healthy", response.Components["database"].Status)
}

func TestHealthHandler_HandleHealth_NoChecks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleHealth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
	assert.Empty(t, response.Components)
}

func TestHealthHandler_HandleLiveness(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleLiveness(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "alive", response["status"])
	assert.NotEmpty(t, response["timestamp"])
}

func TestHealthHandler_HandleReadiness_Ready(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler().
		WithMilvusCheck(func(ctx context.Context) error { return nil }).
		WithDBCheck(func(ctx context.Context) error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleReadiness(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ready", response["status"])

	components := response["components"].(map[string]any)
	assert.Equal(t, "ready", components["milvus"])
	assert.Equal(t, "ready", components["database"])
}

func TestHealthHandler_HandleReadiness_NotReady(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler().
		WithMilvusCheck(func(ctx context.Context) error { return errors.New("not ready") }).
		WithDBCheck(func(ctx context.Context) error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleReadiness(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not ready", response["status"])

	components := response["components"].(map[string]any)
	assert.Equal(t, "not ready", components["milvus"])
	assert.Equal(t, "ready", components["database"])
}
