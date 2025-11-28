# Vector Management Use Case

向量管理用例层，负责处理知识库向量的添加、删除和查询操作。

## 功能概述

### 核心功能

1. **向量添加** (`AddVectors`)
   - 批量添加文本到知识库
   - 自动生成向量嵌入
   - 支持自定义元数据
   - 多租户隔离

2. **向量删除** (`DeleteVectors`)
   - 批量删除向量
   - 根据文档 ID 删除
   - 多租户隔离

3. **向量查询** (`GetVectorByID`)
   - 根据 ID 获取文档
   - 支持多租户

4. **统计功能** (`GetVectorCount`)
   - 获取租户的向量总数

## 使用示例

### 初始化

```go
import (
    "eino-qa/internal/usecase/vector"
    "github.com/cloudwego/eino/components/embedding"
    "eino-qa/internal/domain/repository"
)

// 创建向量管理用例
uc := vector.NewVectorManagementUseCase(
    embedder,    // embedding.Embedder
    vectorRepo,  // repository.VectorRepository
    logger,      // *logrus.Logger
)
```

### 添加向量

```go
// 单个文本添加
req := &vector.AddVectorRequest{
    Texts:    []string{"Python 课程介绍"},
    TenantID: "tenant1",
}

resp, err := uc.AddVectors(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("添加成功，文档 ID: %v\n", resp.DocumentIDs)
```

### 批量添加向量

```go
// 批量添加
req := &vector.AddVectorRequest{
    Texts: []string{
        "Python 基础课程",
        "Go 语言入门",
        "数据结构与算法",
    },
    TenantID: "tenant1",
    Metadata: map[string]any{
        "category": "course",
        "level":    "beginner",
    },
}

resp, err := uc.AddVectors(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("成功添加 %d 个向量\n", resp.Count)
```

### 删除向量

```go
// 删除指定 ID 的向量
req := &vector.DeleteVectorRequest{
    IDs:      []string{"doc_001", "doc_002"},
    TenantID: "tenant1",
}

resp, err := uc.DeleteVectors(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("成功删除 %d 个向量\n", resp.DeletedCount)
```

### 查询向量

```go
// 根据 ID 获取文档
doc, err := uc.GetVectorByID(ctx, "doc_001", "tenant1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("文档内容: %s\n", doc.Content)
```

### 获取向量数量

```go
// 获取租户的向量总数
count, err := uc.GetVectorCount(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("向量总数: %d\n", count)
```

## 数据结构

### AddVectorRequest

```go
type AddVectorRequest struct {
    Texts    []string          `json:"texts" binding:"required"`  // 文本列表
    TenantID string            `json:"tenant_id"`                 // 租户 ID
    Metadata map[string]any    `json:"metadata,omitempty"`        // 元数据
}
```

### AddVectorResponse

```go
type AddVectorResponse struct {
    Success     bool     `json:"success"`       // 是否成功
    DocumentIDs []string `json:"document_ids"`  // 文档 ID 列表
    Count       int      `json:"count"`         // 添加的数量
    Message     string   `json:"message"`       // 消息
}
```

### DeleteVectorRequest

```go
type DeleteVectorRequest struct {
    IDs      []string `json:"ids" binding:"required"`  // 文档 ID 列表
    TenantID string   `json:"tenant_id"`               // 租户 ID
}
```

### DeleteVectorResponse

```go
type DeleteVectorResponse struct {
    Success      bool   `json:"success"`        // 是否成功
    DeletedCount int    `json:"deleted_count"`  // 删除的数量
    Message      string `json:"message"`        // 消息
}
```

## 多租户支持

所有操作都支持多租户隔离：

- 如果请求中未指定 `TenantID`，将使用默认租户 `"default"`
- 每个租户的向量存储在独立的 Milvus Collection 中
- 租户之间的数据完全隔离

## 错误处理

用例层会处理以下错误：

1. **参数验证错误**
   - 空文本列表
   - 空 ID 列表
   - 空文档 ID

2. **嵌入模型错误**
   - 向量生成失败
   - 向量维度不匹配

3. **向量仓储错误**
   - 插入失败
   - 删除失败
   - 查询失败

## 性能考虑

### 批量操作

- 支持批量添加和删除，提高性能
- 建议每批次不超过 100 个文档
- 大批量操作建议分批处理

### 向量生成

- 使用 Eino 嵌入模型生成向量
- 支持并发生成（由 Eino 内部处理）
- 记录生成耗时用于性能监控

## 日志记录

用例层会记录以下日志：

- **Info 级别**
  - 添加向量操作（租户 ID、数量）
  - 删除向量操作（租户 ID、数量）
  - 操作成功信息

- **Debug 级别**
  - 向量生成耗时
  - 详细的操作参数

- **Error 级别**
  - 向量生成失败
  - 插入/删除失败
  - 验证失败

## 需求映射

本用例实现了以下需求：

- **需求 9.2**: 使用 Eino 嵌入模型生成文本向量
- **需求 9.3**: 将向量插入到 Milvus Collection 中
- **需求 9.4**: 根据文档 ID 从 Milvus Collection 中删除向量
- **需求 9.5**: 返回操作结果和受影响的记录数

## 测试

运行单元测试：

```bash
go test -v ./internal/usecase/vector/
```

测试覆盖：

- 添加向量（单个、批量、带元数据）
- 删除向量（单个、批量）
- 查询向量
- 获取向量数量
- 错误处理（空参数、验证失败）
