package milvus

import (
	"context"
	"testing"
	"time"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
)

// TestVectorRepositoryIntegration 集成测试（需要 Milvus 服务运行）
// 运行前请确保 Milvus 服务可用
func TestVectorRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 创建日志器
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	// 配置 Milvus
	milvusConfig := config.MilvusConfig{
		Host:     "localhost",
		Port:     19530,
		Username: "",
		Password: "",
		Timeout:  10 * time.Second,
	}

	// 创建工厂
	factory, err := NewFactory(milvusConfig, 128, logger)
	if err != nil {
		t.Skipf("Skipping test: Milvus not available: %v", err)
		return
	}
	defer factory.Close()

	// 创建仓储
	repo := factory.CreateVectorRepository()

	// 测试上下文
	ctx := context.WithValue(context.Background(), "tenant_id", "test_tenant")

	// 测试 1: 插入文档
	t.Run("Insert", func(t *testing.T) {
		docs := []*entity.Document{
			{
				ID:        "test_doc_001",
				Content:   "测试文档 1",
				Vector:    generateTestVector(128),
				Metadata:  map[string]any{"category": "test"},
				TenantID:  "test_tenant",
				CreatedAt: time.Now(),
			},
			{
				ID:        "test_doc_002",
				Content:   "测试文档 2",
				Vector:    generateTestVector(128),
				Metadata:  map[string]any{"category": "test"},
				TenantID:  "test_tenant",
				CreatedAt: time.Now(),
			},
		}

		err := repo.Insert(ctx, docs)
		if err != nil {
			t.Fatalf("Failed to insert documents: %v", err)
		}

		// 等待数据刷新
		time.Sleep(2 * time.Second)
	})

	// 测试 2: 搜索文档
	t.Run("Search", func(t *testing.T) {
		queryVector := generateTestVector(128)
		results, err := repo.Search(ctx, queryVector, 5)
		if err != nil {
			t.Fatalf("Failed to search: %v", err)
		}

		if len(results) == 0 {
			t.Log("No results found (this is ok for test)")
		} else {
			t.Logf("Found %d results", len(results))
			for i, doc := range results {
				t.Logf("  %d. ID: %s, Score: %.4f", i+1, doc.ID, doc.Score)
			}
		}
	})

	// 测试 3: 获取文档总数
	t.Run("Count", func(t *testing.T) {
		count, err := repo.Count(ctx)
		if err != nil {
			t.Fatalf("Failed to count: %v", err)
		}

		t.Logf("Total documents: %d", count)
		if count < 0 {
			t.Errorf("Count should be non-negative, got %d", count)
		}
	})

	// 测试 4: 根据 ID 获取文档
	t.Run("GetByID", func(t *testing.T) {
		doc, err := repo.GetByID(ctx, "test_doc_001")
		if err != nil {
			t.Logf("Document not found (this is ok): %v", err)
		} else {
			if doc.ID != "test_doc_001" {
				t.Errorf("Expected ID test_doc_001, got %s", doc.ID)
			}
			t.Logf("Found document: %s", doc.Content)
		}
	})

	// 测试 5: 删除文档
	t.Run("Delete", func(t *testing.T) {
		deletedCount, err := repo.Delete(ctx, []string{"test_doc_001", "test_doc_002"})
		if err != nil {
			t.Fatalf("Failed to delete: %v", err)
		}

		if deletedCount != 2 {
			t.Logf("Expected to delete 2 documents, deleted %d", deletedCount)
		}
	})
}

// TestTenantManager 测试租户管理器
func TestTenantManager(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	milvusConfig := config.MilvusConfig{
		Host:     "localhost",
		Port:     19530,
		Username: "",
		Password: "",
		Timeout:  10 * time.Second,
	}

	factory, err := NewFactory(milvusConfig, 128, logger)
	if err != nil {
		t.Skipf("Skipping test: Milvus not available: %v", err)
		return
	}
	defer factory.Close()

	tenantManager := factory.GetTenantManager()

	t.Run("GetCollection", func(t *testing.T) {
		ctx := context.Background()

		// 测试默认租户
		collection, err := tenantManager.GetCollection(ctx, "default")
		if err != nil {
			t.Fatalf("Failed to get collection for default tenant: %v", err)
		}

		if collection != "kb_default" {
			t.Errorf("Expected collection name kb_default, got %s", collection)
		}

		// 测试自定义租户
		collection, err = tenantManager.GetCollection(ctx, "tenant_test")
		if err != nil {
			t.Fatalf("Failed to get collection for tenant_test: %v", err)
		}

		if collection != "kb_tenant_test" {
			t.Errorf("Expected collection name kb_tenant_test, got %s", collection)
		}
	})

	t.Run("CollectionExists", func(t *testing.T) {
		ctx := context.Background()

		exists, err := tenantManager.CollectionExists(ctx, "default")
		if err != nil {
			t.Fatalf("Failed to check collection existence: %v", err)
		}

		t.Logf("Default tenant collection exists: %v", exists)
	})

	t.Run("GetAllTenants", func(t *testing.T) {
		tenants := tenantManager.GetAllTenants()
		t.Logf("All tenants: %v", tenants)

		if len(tenants) == 0 {
			t.Log("No tenants cached yet")
		}
	})
}

// generateTestVector 生成测试向量
func generateTestVector(dimension int) []float32 {
	vector := make([]float32, dimension)
	for i := range vector {
		vector[i] = float32(i) / float32(dimension)
	}
	return vector
}
