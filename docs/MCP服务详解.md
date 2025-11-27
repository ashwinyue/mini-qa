# MCP 服务详解

## 什么是 MCP？

**MCP (Model Context Protocol)** 是一个开放协议，用于在 AI 应用和外部工具/数据源之间建立标准化的通信接口。

### 核心概念

MCP 类似于 API 网关，但专门为 AI Agent 设计：
- **工具 (Tools)**：AI 可以调用的函数
- **资源 (Resources)**：AI 可以访问的数据
- **提示词 (Prompts)**：预定义的提示模板

### 为什么需要 MCP？

在 work_v2 中引入 MCP 服务器是为了：

1. **标准化接口**：提供统一的工具调用规范
2. **外部集成**：让其他系统（如 Claude Desktop、Cursor 等）能调用智能客服的能力
3. **工具复用**：将内部功能封装为可复用的工具
4. **协议兼容**：支持多种传输方式（stdio、HTTP、SSE）

---

## Work V2 的 MCP 服务架构

### 整体架构

```
┌─────────────────────────────────────────┐
│         外部 MCP 客户端                  │
│  (Claude Desktop / Cursor / 自定义)     │
└────────────────┬────────────────────────┘
                 │ MCP Protocol
                 │ (stdio/HTTP/SSE)
┌────────────────▼────────────────────────┐
│         work_v2/mcp_server.py           │
│                                          │
│  ┌────────────────────────────────┐    │
│  │  FastMCP 框架                   │    │
│  │  - 工具注册                     │    │
│  │  - 资源管理                     │    │
│  │  - 协议转换                     │    │
│  └────────────────────────────────┘    │
│                                          │
│  ┌──────────┐  ┌──────────┐           │
│  │  Tools   │  │Resources │           │
│  │  (4个)   │  │  (2个)   │           │
│  └──────────┘  └──────────┘           │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      智能客服核心功能                    │
│  - tools.py (知识库、订单查询)          │
│  - graph.py (对话流程)                  │
│  - config.py (配置管理)                 │
└─────────────────────────────────────────┘
```

---

## MCP 服务提供的功能

### 1. 工具 (Tools)

#### 1.1 知识库检索 `kb_search`

**功能：** 在 FAISS 向量数据库中检索相关文档

**签名：**
```python
@mcp.tool()
def kb_search(query: str, k: int = 2) -> Dict[str, Any]:
    """
    在知识库中搜索相关内容
    
    参数:
        query: 搜索查询文本
        k: 返回的文档数量，默认 2
    
    返回:
        {
            "context": "检索到的文档内容",
            "sources": [{"source": "文件路径"}]
        }
    """
```

**使用场景：**
- 外部系统需要查询课程信息
- AI Agent 需要获取知识库上下文
- 自动化脚本批量检索

**示例：**
```python
result = kb_search("Python 课程适合零基础吗？", k=3)
print(result["context"])  # 相关文档内容
print(result["sources"])  # 来源列表
```

#### 1.2 订单查询 `order_lookup`

**功能：** 查询订单状态和详情

**签名：**
```python
@mcp.tool()
def order_lookup(text: str) -> Dict[str, Any]:
    """
    查询订单信息
    
    参数:
        text: 包含订单号的文本（如 "#20251114001"）
    
    返回:
        {
            "order_id": "订单号",
            "status": "订单状态",
            "amount": 金额,
            "updated_at": "更新时间",
            "timeline": ["时间线"]
        }
    """
```

**使用场景：**
- 客服系统查询订单
- 自动化订单状态监控
- 第三方系统集成

**示例：**
```python
order = order_lookup("#20251114001")
print(f"订单状态: {order['status']}")
print(f"金额: {order['amount']}")
```

#### 1.3 课程目录 `course_catalog`

**功能：** 获取课程列表和分类

**签名：**
```python
@mcp.tool()
def course_catalog(limit: int = 20) -> Dict[str, Any]:
    """
    获取课程目录
    
    参数:
        limit: 返回的课程数量，默认 20
    
    返回:
        {
            "sections": ["分类1", "分类2"],
            "items": [
                {"section": "分类", "q": "问题", "a": "答案"}
            ]
        }
    """
```

**使用场景：**
- 展示课程列表
- 生成课程导航
- 课程推荐系统

**示例：**
```python
catalog = course_catalog(limit=10)
for item in catalog["items"]:
    print(f"{item['section']}: {item['q']}")
```

#### 1.4 对话接口 `chat`

**功能：** 完整的智能客服对话流程

**签名：**
```python
@mcp.tool()
def chat(query: str, thread_id: Optional[str] = None) -> Dict[str, Any]:
    """
    智能客服对话
    
    参数:
        query: 用户问题
        thread_id: 会话 ID（可选）
    
    返回:
        {
            "route": "路由类型",
            "answer": "回答内容",
            "sources": [来源列表]
        }
    """
```

**使用场景：**
- 外部聊天机器人集成
- 批量问答测试
- API 调用

**示例：**
```python
response = chat("我想了解 AI 课程", thread_id="user123")
print(f"路由: {response['route']}")
print(f"回答: {response['answer']}")
```

### 2. 资源 (Resources)

#### 2.1 知识库资源 `kb://{query}`

**功能：** 以资源形式访问知识库

**URI 格式：** `kb://Python课程介绍`

**返回：** 纯文本格式的检索结果

**使用场景：**
- 静态资源访问
- 文档预览
- 内容导出

#### 2.2 订单资源 `orders://{order_id}`

**功能：** 以资源形式访问订单

**URI 格式：** `orders://20251114001`

**返回：** JSON 格式的订单详情

**使用场景：**
- RESTful 风格访问
- 资源缓存
- 批量导出

---

## 在 Work V2 中的集成方式

### 1. FastAPI 应用集成

**文件：** `work_v2/app.py`

```python
# 导入 MCP 服务器
from mcp_server import mcp as _mcp

# 创建 SSE 应用
_mcp_app = _mcp.sse_app()

# 挂载到 FastAPI
app.mount("/mcp", _mcp_app)
```

**访问地址：** `http://localhost:8000/mcp`

### 2. 独立运行模式

**直接运行 MCP 服务器：**

```bash
# stdio 模式（标准输入输出）
python3 work_v2/mcp_server.py

# HTTP 模式
MCP_TRANSPORT=http MCP_PORT=6278 python3 work_v2/mcp_server.py

# SSE 模式
MCP_TRANSPORT=sse MCP_PORT=6278 python3 work_v2/mcp_server.py
```

---

## 传输协议详解

### 1. stdio 模式

**特点：**
- 通过标准输入输出通信
- 适合本地进程间通信
- Claude Desktop 默认使用此模式

**配置示例（Claude Desktop）：**
```json
{
  "mcpServers": {
    "edu-agent": {
      "command": "python3",
      "args": ["/path/to/work_v2/mcp_server.py"],
      "env": {
        "PYTHONPATH": "/path/to/project"
      }
    }
  }
}
```

### 2. HTTP 模式

**特点：**
- 标准 HTTP REST API
- 适合远程调用
- 易于调试和测试

**请求示例：**
```bash
curl -X POST http://localhost:6278/tools/kb_search \
  -H "Content-Type: application/json" \
  -d '{"query": "Python课程", "k": 2}'
```

### 3. SSE 模式

**特点：**
- Server-Sent Events 流式传输
- 支持实时推送
- 适合长连接场景

**集成到 FastAPI：**
```python
# 已在 app.py 中集成
# 访问: http://localhost:8000/mcp
```

---

## 实际使用场景

### 场景 1：Claude Desktop 集成

**目标：** 让 Claude 能够调用智能客服的知识库

**步骤：**

1. 配置 Claude Desktop 的 `claude_desktop_config.json`：
```json
{
  "mcpServers": {
    "edu-agent": {
      "command": "python3",
      "args": ["/Users/yourname/project/work_v2/mcp_server.py"]
    }
  }
}
```

2. 重启 Claude Desktop

3. 在对话中使用：
```
用户: 帮我查询一下 Python 课程的信息
Claude: [调用 kb_search 工具] 根据知识库，Python 课程...
```

### 场景 2：Cursor IDE 集成

**目标：** 在 Cursor 中使用智能客服能力

**配置：** 类似 Claude Desktop

**使用：**
```
# 在 Cursor 中
@edu-agent 这门课程适合零基础吗？
```

### 场景 3：自定义脚本调用

**目标：** 批量处理问题

**代码：**
```python
from mcp.client import Client

# 连接到 MCP 服务器
client = Client("http://localhost:6278")

# 批量查询
questions = [
    "Python 课程适合零基础吗？",
    "如何报名课程？",
    "课程费用是多少？"
]

for q in questions:
    result = client.call_tool("kb_search", {"query": q, "k": 1})
    print(f"Q: {q}")
    print(f"A: {result['context']}\n")
```

### 场景 4：第三方系统集成

**目标：** 将智能客服能力集成到现有系统

**架构：**
```
┌─────────────┐
│  现有系统    │
└──────┬──────┘
       │ HTTP
┌──────▼──────────────┐
│  MCP HTTP Server    │
│  (port 6278)        │
└──────┬──────────────┘
       │
┌──────▼──────────────┐
│  智能客服核心        │
└─────────────────────┘
```

**调用示例：**
```javascript
// JavaScript 调用
fetch('http://localhost:6278/tools/chat', {
  method: 'POST',
  headers: {'Content-Type': 'application/json'},
  body: JSON.stringify({
    query: '我想了解课程',
    thread_id: 'user123'
  })
})
.then(res => res.json())
.then(data => console.log(data.answer));
```

---

## MCP vs 直接 API 调用

### 对比表

| 维度 | MCP 服务 | 直接 API |
|------|---------|---------|
| **标准化** | ✅ 遵循 MCP 协议 | ❌ 自定义格式 |
| **工具发现** | ✅ 自动发现工具 | ❌ 需要文档 |
| **AI 集成** | ✅ 原生支持 | ⚠️ 需要适配 |
| **类型安全** | ✅ 自动验证 | ⚠️ 手动验证 |
| **多传输** | ✅ stdio/HTTP/SSE | ❌ 仅 HTTP |
| **调试难度** | ⚠️ 需要 MCP 客户端 | ✅ 简单 |

### 何时使用 MCP？

**推荐使用 MCP：**
- ✅ 需要与 AI 工具（Claude、Cursor）集成
- ✅ 需要标准化的工具接口
- ✅ 需要工具自动发现
- ✅ 需要多种传输方式

**推荐直接 API：**
- ✅ 简单的 HTTP 调用
- ✅ 不需要 AI 集成
- ✅ 已有成熟的 API 体系

---

## Work V3 的 MCP 增强

### 主要改进

#### 1. 租户支持

所有工具都支持 `tenant_id` 参数：

```python
@mcp.tool()
def kb_search(query: str, k: int = 2, tenant_id: Optional[str] = None):
    """支持多租户的知识库检索"""
    serialized, docs = retrieve_kb(query, tenant_id)
    # ...
```

#### 2. 租户级资源

新增租户级资源 URI：

```python
@mcp.resource("kb://{tenant_id}/{query}")
def kb_resource(tenant_id: str, query: str) -> str:
    """租户级知识库资源"""
    serialized, _ = retrieve_kb(query, tenant_id)
    return serialized

@mcp.resource("orders://{tenant_id}/{order_id}")
def order_resource(tenant_id: str, order_id: str) -> Dict[str, Any]:
    """租户级订单资源"""
    # ...
```

#### 3. 使用示例

```python
# V2: 单租户
result = kb_search("Python课程")

# V3: 多租户
result = kb_search("Python课程", tenant_id="t1")
```

---

## 调试与测试

### 1. 测试工具调用

```bash
# 使用 MCP Inspector（官方工具）
npx @modelcontextprotocol/inspector python3 work_v2/mcp_server.py
```

### 2. 手动测试

```python
# 直接导入测试
from work_v2.mcp_server import kb_search, order_lookup

# 测试知识库检索
result = kb_search("Python课程", k=2)
print(result)

# 测试订单查询
order = order_lookup("#20251114001")
print(order)
```

### 3. HTTP 模式测试

```bash
# 启动 HTTP 服务器
MCP_TRANSPORT=http MCP_PORT=6278 python3 work_v2/mcp_server.py

# 测试工具列表
curl http://localhost:6278/tools

# 测试工具调用
curl -X POST http://localhost:6278/tools/kb_search \
  -H "Content-Type: application/json" \
  -d '{"query":"Python课程","k":2}'
```

---

## 常见问题

### Q1: MCP 服务器启动失败？

**原因：** 缺少依赖

**解决：**
```bash
pip install mcp-server-fastmcp
```

### Q2: Claude Desktop 无法连接？

**检查清单：**
- ✅ 配置文件路径正确
- ✅ Python 路径正确
- ✅ 环境变量设置正确
- ✅ 重启 Claude Desktop

### Q3: 工具调用超时？

**原因：** 知识库索引加载慢

**解决：**
- 预热索引（首次调用会慢）
- 优化索引大小
- 增加超时时间

### Q4: 如何添加新工具？

**步骤：**
```python
# 在 mcp_server.py 中添加
@mcp.tool()
def your_new_tool(param: str) -> Dict[str, Any]:
    """工具描述"""
    # 实现逻辑
    return {"result": "..."}
```

---

## 总结

### MCP 服务的价值

1. **标准化**：统一的工具接口规范
2. **可集成**：轻松接入 AI 工具生态
3. **可扩展**：方便添加新工具
4. **多协议**：支持多种传输方式

### 在 Work V2 中的作用

- 🔌 **外部集成桥梁**：连接智能客服与外部系统
- 🛠️ **工具标准化**：将内部功能封装为标准工具
- 🤖 **AI 原生支持**：让 AI 助手能直接调用
- 📦 **功能复用**：一次封装，多处使用

### 未来展望

- 支持更多工具类型
- 增强安全认证
- 提供工具组合能力
- 支持流式响应

---

**文档版本：** v1.0  
**创建日期：** 2025-11-27  
**作者：** Kiro AI Assistant  
**适用版本：** work_v2, work_v3
