# API 文档

## 概述

Eino QA System 提供 RESTful API 接口，支持智能对话、向量管理、健康检查等功能。

**基础 URL**: `http://localhost:8080`

**认证方式**:
- 对话接口：无需认证
- 向量管理接口：需要 API Key（通过 `X-API-Key` 请求头）

**内容类型**: `application/json`

## 目录

- [对话接口](#对话接口)
- [向量管理接口](#向量管理接口)
- [健康检查接口](#健康检查接口)
- [错误处理](#错误处理)
- [多租户支持](#多租户支持)

---

## 对话接口

### POST /chat

发送用户查询，获取智能回复。

#### 请求

**Headers**:
```
Content-Type: application/json
X-Tenant-ID: <tenant_id>  (可选)
```

**Body**:
```json
{
  "query": "Python 课程包含哪些内容？",
  "tenant_id": "default",
  "session_id": "session-123",
  "stream": false
}
```

**参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| query | string | 是 | 用户查询内容 |
| tenant_id | string | 否 | 租户 ID，默认为 "default" |
| session_id | string | 否 | 会话 ID，用于多轮对话上下文 |
| stream | boolean | 否 | 是否使用流式响应，默认 false |

#### 响应

**成功响应** (200 OK):

```json
{
  "answer": "Python 课程包含以下内容：\n1. 基础语法\n2. 数据结构\n3. 面向对象编程\n4. 常用库的使用",
  "route": "course",
  "session_id": "session-123",
  "sources": [
    {
      "content": "Python 课程包含基础语法、数据结构、面向对象编程等内容",
      "score": 0.95,
      "metadata": {
        "doc_id": "doc-001"
      }
    }
  ],
  "metadata": {
    "intent": "course",
    "confidence": 0.92,
    "duration_ms": 234,
    "timestamp": "2024-11-29T10:00:00Z"
  }
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| answer | string | 系统生成的回答 |
| route | string | 路由类型：course（课程咨询）、order（订单查询）、direct（直接回答）、handoff（人工转接） |
| session_id | string | 会话 ID |
| sources | array | 检索到的相关文档（仅 course 路由） |
| metadata | object | 元数据信息 |

**流式响应** (stream=true):

使用 Server-Sent Events (SSE) 协议：

```
Content-Type: text/event-stream

data: {"type":"start","session_id":"session-123"}

data: {"type":"chunk","content":"Python"}

data: {"type":"chunk","content":" 课程"}

data: {"type":"end","metadata":{"duration_ms":234}}
```

#### 示例

**课程咨询**:
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Python 课程包含哪些内容？",
    "tenant_id": "default"
  }'
```

**订单查询**:
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "查询订单 #20251114001",
    "tenant_id": "default"
  }'
```

**多轮对话**:
```bash
# 第一轮
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Python 课程多少钱？",
    "session_id": "session-001"
  }'

# 第二轮（带上下文）
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "有优惠吗？",
    "session_id": "session-001"
  }'
```

---

## 向量管理接口

### POST /api/v1/vectors/items

添加文档到向量数据库。

#### 请求

**Headers**:
```
Content-Type: application/json
X-API-Key: <your_api_key>
X-Tenant-ID: <tenant_id>  (可选)
```

**Body**:
```json
{
  "texts": [
    "Python 课程包含基础语法、数据结构、面向对象编程等内容",
    "Go 语言课程涵盖并发编程、网络编程、微服务开发等主题"
  ],
  "tenant_id": "default",
  "metadata": {
    "category": "course",
    "author": "admin"
  }
}
```

**参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| texts | array | 是 | 要添加的文档内容列表 |
| tenant_id | string | 否 | 租户 ID，默认为 "default" |
| metadata | object | 否 | 文档元数据 |

#### 响应

**成功响应** (200 OK):

```json
{
  "success": true,
  "inserted_count": 2,
  "document_ids": [
    "doc-uuid-001",
    "doc-uuid-002"
  ],
  "message": "Successfully inserted 2 documents"
}
```

#### 示例

```bash
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key_here" \
  -d '{
    "texts": [
      "Python 课程包含基础语法、数据结构、面向对象编程等内容"
    ],
    "tenant_id": "default"
  }'
```

---

### DELETE /api/v1/vectors/items

从向量数据库删除文档。

#### 请求

**Headers**:
```
Content-Type: application/json
X-API-Key: <your_api_key>
X-Tenant-ID: <tenant_id>  (可选)
```

**Body**:
```json
{
  "ids": [
    "doc-uuid-001",
    "doc-uuid-002"
  ],
  "tenant_id": "default"
}
```

**参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ids | array | 是 | 要删除的文档 ID 列表 |
| tenant_id | string | 否 | 租户 ID，默认为 "default" |

#### 响应

**成功响应** (200 OK):

```json
{
  "success": true,
  "deleted_count": 2,
  "message": "Successfully deleted 2 documents"
}
```

#### 示例

```bash
curl -X DELETE http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key_here" \
  -d '{
    "ids": ["doc-uuid-001", "doc-uuid-002"],
    "tenant_id": "default"
  }'
```

---

## 健康检查接口

### GET /health

获取系统健康状态。

#### 响应

**成功响应** (200 OK):

```json
{
  "status": "healthy",
  "timestamp": "2024-11-29T10:00:00Z",
  "version": "1.0.0",
  "components": {
    "milvus": {
      "status": "healthy",
      "latency_ms": 5
    },
    "database": {
      "status": "healthy",
      "latency_ms": 2
    },
    "dashscope": {
      "status": "healthy"
    }
  },
  "metrics": {
    "total_requests": 1234,
    "avg_response_time_ms": 234,
    "error_rate": 0.01
  }
}
```

#### 示例

```bash
curl http://localhost:8080/health
```

---

### GET /health/live

存活检查（Liveness Probe）。

#### 响应

**成功响应** (200 OK):

```json
{
  "status": "alive"
}
```

#### 示例

```bash
curl http://localhost:8080/health/live
```

---

### GET /health/ready

就绪检查（Readiness Probe）。

#### 响应

**成功响应** (200 OK):

```json
{
  "status": "ready",
  "components": {
    "milvus": "ready",
    "database": "ready"
  }
}
```

**未就绪响应** (503 Service Unavailable):

```json
{
  "status": "not_ready",
  "components": {
    "milvus": "not_ready",
    "database": "ready"
  }
}
```

#### 示例

```bash
curl http://localhost:8080/health/ready
```

---

## 错误处理

### 错误响应格式

所有错误响应遵循统一格式：

```json
{
  "code": 400,
  "message": "Invalid request: missing required field 'query'",
  "details": {
    "field": "query",
    "reason": "required field is missing"
  },
  "trace_id": "trace-uuid-123"
}
```

### HTTP 状态码

| 状态码 | 说明 | 示例 |
|--------|------|------|
| 200 | 成功 | 请求处理成功 |
| 400 | 请求错误 | 参数格式错误、缺少必填字段 |
| 401 | 未授权 | API Key 验证失败 |
| 404 | 资源不存在 | 订单 ID 不存在 |
| 429 | 请求过多 | 超过速率限制 |
| 500 | 服务器错误 | 内部错误 |
| 502 | 网关错误 | 外部服务调用失败 |
| 503 | 服务不可用 | 系统过载或维护中 |

### 常见错误

#### 1. 缺少必填参数

**请求**:
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{}'
```

**响应** (400):
```json
{
  "code": 400,
  "message": "Invalid request: missing required field 'query'",
  "trace_id": "trace-uuid-123"
}
```

#### 2. API Key 验证失败

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -d '{"texts": ["test"]}'
```

**响应** (401):
```json
{
  "code": 401,
  "message": "Unauthorized: invalid or missing API key",
  "trace_id": "trace-uuid-123"
}
```

#### 3. 订单不存在

**请求**:
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "查询订单 #99999999"
  }'
```

**响应** (200，但在 answer 中说明):
```json
{
  "answer": "抱歉，未找到订单 #99999999，请检查订单号是否正确。",
  "route": "order",
  "metadata": {
    "intent": "order",
    "order_found": false
  }
}
```

#### 4. 外部服务不可用

**响应** (502):
```json
{
  "code": 502,
  "message": "External service error: failed to connect to Milvus",
  "details": {
    "service": "milvus",
    "error": "connection refused"
  },
  "trace_id": "trace-uuid-123"
}
```

---

## 多租户支持

### 租户识别

系统支持两种方式指定租户：

1. **请求头** (推荐):
```bash
curl -X POST http://localhost:8080/chat \
  -H "X-Tenant-ID: tenant1" \
  -H "Content-Type: application/json" \
  -d '{"query": "test"}'
```

2. **请求体**:
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "test",
    "tenant_id": "tenant1"
  }'
```

优先级：请求头 > 请求体 > 默认值 ("default")

### 租户隔离

每个租户拥有独立的：
- **Milvus Collection**: `kb_{tenant_id}`
- **SQLite 数据库**: `./data/db/{tenant_id}.db`

### 租户自动创建

首次使用时，系统会自动创建租户资源：
- 创建 Milvus Collection
- 创建 SQLite 数据库文件
- 初始化数据表结构

---

## 速率限制

### 限制规则

| 接口 | 限制 | 说明 |
|------|------|------|
| /chat | 100 请求/分钟/租户 | 对话接口 |
| /api/v1/vectors/* | 50 请求/分钟/API Key | 向量管理接口 |

### 超限响应

**响应** (429):
```json
{
  "code": 429,
  "message": "Rate limit exceeded",
  "details": {
    "limit": 100,
    "window": "1m",
    "retry_after": 30
  },
  "trace_id": "trace-uuid-123"
}
```

**Headers**:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1701234567
Retry-After: 30
```

---

## 最佳实践

### 1. 使用会话 ID

对于多轮对话，始终使用相同的 `session_id`：

```bash
# 生成唯一的 session_id
SESSION_ID=$(uuidgen)

# 第一轮对话
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"query\": \"Python 课程多少钱？\",
    \"session_id\": \"$SESSION_ID\"
  }"

# 第二轮对话（带上下文）
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"query\": \"有优惠吗？\",
    \"session_id\": \"$SESSION_ID\"
  }"
```

### 2. 处理流式响应

使用流式响应提升用户体验：

```javascript
const eventSource = new EventSource('/chat?stream=true');

eventSource.addEventListener('message', (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'chunk') {
    // 追加内容到 UI
    appendToChat(data.content);
  } else if (data.type === 'end') {
    // 对话结束
    eventSource.close();
  }
});
```

### 3. 错误重试

对于临时性错误（5xx），实现指数退避重试：

```python
import time
import requests

def chat_with_retry(query, max_retries=3):
    for attempt in range(max_retries):
        try:
            response = requests.post(
                'http://localhost:8080/chat',
                json={'query': query}
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if e.response.status_code >= 500 and attempt < max_retries - 1:
                wait_time = 2 ** attempt
                time.sleep(wait_time)
                continue
            raise
```

### 4. 批量添加向量

批量添加向量以提高效率：

```bash
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "texts": [
      "文档1内容",
      "文档2内容",
      "文档3内容"
    ]
  }'
```

---

## SDK 示例

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type ChatRequest struct {
    Query     string `json:"query"`
    TenantID  string `json:"tenant_id,omitempty"`
    SessionID string `json:"session_id,omitempty"`
}

type ChatResponse struct {
    Answer    string                 `json:"answer"`
    Route     string                 `json:"route"`
    SessionID string                 `json:"session_id"`
    Metadata  map[string]interface{} `json:"metadata"`
}

func Chat(query string) (*ChatResponse, error) {
    req := ChatRequest{
        Query:    query,
        TenantID: "default",
    }
    
    body, _ := json.Marshal(req)
    resp, err := http.Post(
        "http://localhost:8080/chat",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result ChatResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

### Python

```python
import requests

class EinoQAClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
    
    def chat(self, query, tenant_id="default", session_id=None):
        response = requests.post(
            f"{self.base_url}/chat",
            json={
                "query": query,
                "tenant_id": tenant_id,
                "session_id": session_id
            }
        )
        response.raise_for_status()
        return response.json()
    
    def add_vectors(self, texts, api_key, tenant_id="default"):
        response = requests.post(
            f"{self.base_url}/api/v1/vectors/items",
            headers={"X-API-Key": api_key},
            json={
                "texts": texts,
                "tenant_id": tenant_id
            }
        )
        response.raise_for_status()
        return response.json()

# 使用示例
client = EinoQAClient()
result = client.chat("Python 课程包含哪些内容？")
print(result["answer"])
```

---

## 附录

### A. 意图类型说明

| 意图类型 | 说明 | 示例查询 |
|---------|------|---------|
| course | 课程咨询 | "Python 课程包含哪些内容？" |
| order | 订单查询 | "查询订单 #20251114001" |
| direct | 直接回答 | "你好"、"谢谢" |
| handoff | 人工转接 | 复杂问题或低置信度查询 |

### B. 元数据字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| intent | string | 识别的意图类型 |
| confidence | float | 意图识别置信度 (0-1) |
| duration_ms | int | 处理耗时（毫秒） |
| timestamp | string | 响应时间戳 (ISO 8601) |
| order_found | bool | 订单是否找到（仅 order 路由） |
| sources_count | int | 检索到的文档数量（仅 course 路由） |

### C. 配置参数参考

详细配置参数请参考 `config/config.yaml` 文件。

---

## 更新日志

### v1.0.0 (2024-11-29)

- 初始版本发布
- 支持对话接口
- 支持向量管理接口
- 支持多租户隔离
- 支持健康检查

---

## 联系我们

如有问题或建议，请：
- 提交 Issue
- 查看项目文档
- 联系技术支持
