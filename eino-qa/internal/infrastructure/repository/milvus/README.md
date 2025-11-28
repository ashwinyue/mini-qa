# Milvus 向量数据库集成

本模块实现了基于 Milvus 的向量数据库集成，提供高性能的向量相似度搜索和多租户支持。

## 功能特性

- ✅ Milvus 客户端连接管理
- ✅ Collection 自动创建和管理
- ✅ 多租户 Collection 隔离
- ✅ 向量插入、搜索、删除操作
- ✅ HNSW 索引优化
- ✅ 自动刷新和持久化

## 架构组件

### 1. Client (client.go)
Milvus 客户端封装，负责连接管理和健康检查。

```go
client, err := milvus.NewClient(milvus.ClientConfig{
    Host:     "localhost",
    Port:     19530,
    Username: "",
    Password: "",
    Timeout:  10 * time.Second,
}, logger)
```

### 2. CollectionManager (collection.go)
管理 Milvus Collection 的创建、删除和配置。

**Schema 定义:**
- `id` (VarChar): 主键，文档 ID
- `vector` (FloatVector): 向量数据，维度可配置
- `content` (VarChar): 文档内容
- `metadata` (JSON): 元数据
- `tenant_id` (VarChar): 租户 ID
- `created_at` (Int64): 创建时间戳

**索引配置:**
- 类型: HNSW (Hierarchical Navigable Small World)
- 距离度量: L2 (欧氏距离)
- 参数: M=16, efConstruction=256

### 3. TenantManager (tenant_manager.go)
管理多租户的 Collection 映射和自动创建。

**命名规则:**
- 默认租户: `kb_default`
- 其他租户: `kb_{tenant_id}`

**特性:**
- 自动创建租户 Collection
- 缓存租户映射关系
- 线程安全的并发访问
- 双重检查锁定防止重复创建

### 4. VectorRepository (vector_repository.go)
实现 `domain.VectorRepository` 接口，提供向量操作。

**主要方法:**
- `Search`: 向量相似度搜索
- `Insert`: 批量插入文档
- `Delete`: 批量删除文档
- `GetByID`: 根据 ID 获取文档
- `Count`: 获取文档总数

### 5. Factory (factory.go)
工厂模式，简化组件初始化。

```go
factory, err := milvus.NewFactory(config.MilvusConfig{
    Host: "localhost",
    Port: 19530,
}, 1536, logger)

vectorRepo := factory.CreateVectorRepository()
```

## 使用示例

### 基本使用

```go
// 1. 创建工厂
factory, err := milvus.NewFactory(milvusConfig, 1536, logger)
if err != nil {
    log.Fatal(err)
}
defer factory.Close()

// 2. 创建仓储
vectorRepo := factory.CreateVectorRepository()

// 3. 设置租户上下文
ctx := context.WithValue(context.Background(), "tenant_id", "my_tenant")

// 4. 插入文档
docs := []*entity.Document{
    {
        ID:       "doc_001",
        Content:  "Python 编程语言",
        Vector:   embeddings, // 从嵌入模型获取
        Metadata: map[string]any{"category": "tech"},
        TenantID: "my_tenant",
    },
}
err = vectorRepo.Insert(ctx, docs)

// 5. 搜索相似文档
queryVector := getQueryEmbedding("Python 教程")
results, err := vectorRepo.Search(ctx, queryVector, 5)
for _, doc := range results {
    fmt.Printf("Score: %.4f, Content: %s\n", doc.Score, doc.Content)
}
```

### 多租户使用

```go
// 租户 A
ctxA := context.WithValue(context.Background(), "tenant_id", "tenant_a")
vectorRepo.Insert(ctxA, docsA)
resultsA, _ := vectorRepo.Search(ctxA, queryVector, 5)

// 租户 B
ctxB := context.WithValue(context.Background(), "tenant_id", "tenant_b")
vectorRepo.Insert(ctxB, docsB)
resultsB, _ := vectorRepo.Search(ctxB, queryVector, 5)

// 租户 A 和 B 的数据完全隔离
```

### 租户管理

```go
tenantManager := factory.GetTenantManager()

// 检查租户 Collection 是否存在
exists, err := tenantManager.CollectionExists(ctx, "tenant_a")

// 获取租户 Collection 名称
collectionName, err := tenantManager.GetCollection(ctx, "tenant_a")

// 删除租户 Collection
err = tenantManager.DropTenantCollection(ctx, "tenant_a")

// 列出所有租户
tenants := tenantManager.GetAllTenants()
```

## 配置说明

### Milvus 配置 (config.yaml)

```yaml
milvus:
  host: localhost
  port: 19530
  username: ""
  password: ""
  timeout: 10s
```

### 向量维度

向量维度需要与嵌入模型匹配：
- DashScope text-embedding-v2: 1536
- DashScope text-embedding-v3: 1024
- OpenAI text-embedding-ada-002: 1536

## 性能优化

### 1. 索引选择
使用 HNSW 索引提供高性能搜索：
- 查询速度快
- 内存占用适中
- 适合大规模数据

### 2. 批量操作
建议批量插入和删除以提高性能：
```go
// 好的做法：批量插入
vectorRepo.Insert(ctx, docs) // docs 包含多个文档

// 避免：逐个插入
for _, doc := range docs {
    vectorRepo.Insert(ctx, []*entity.Document{doc})
}
```

### 3. 刷新策略
- 插入和删除后自动刷新
- 确保数据持久化
- 可能影响性能，根据需求调整

## 错误处理

常见错误和处理方式：

### 连接失败
```go
client, err := milvus.NewClient(config, logger)
if err != nil {
    // 检查 Milvus 服务是否运行
    // 检查网络连接
    // 检查配置是否正确
}
```

### Collection 不存在
自动创建机制会处理，无需手动干预。

### 向量维度不匹配
```go
// 确保向量维度与 Collection Schema 一致
if len(vector) != expectedDimension {
    return fmt.Errorf("vector dimension mismatch")
}
```

## 测试

运行测试：
```bash
cd eino-qa
go test ./internal/infrastructure/repository/milvus/...
```

运行示例：
```bash
# 确保 Milvus 服务运行
docker run -d --name milvus -p 19530:19530 milvusdb/milvus:latest

# 运行示例
go run examples/milvus_usage.go
```

## 依赖

- `github.com/milvus-io/milvus-sdk-go/v2`: Milvus Go SDK
- `github.com/sirupsen/logrus`: 日志库

## 相关文档

- [Milvus 官方文档](https://milvus.io/docs)
- [Milvus Go SDK](https://github.com/milvus-io/milvus-sdk-go)
- [HNSW 算法](https://arxiv.org/abs/1603.09320)

## 需求映射

本实现满足以下需求：
- **需求 3.2**: 在 Milvus 向量数据库中执行相似度搜索
- **需求 5.3**: 使用租户对应的 Milvus Collection
- **需求 5.4**: 创建新的 Collection 并初始化 Schema
- **需求 9.3**: 将向量插入到 Milvus Collection 中
- **需求 9.4**: 从 Milvus Collection 中删除向量
