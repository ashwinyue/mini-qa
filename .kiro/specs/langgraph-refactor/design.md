# Design Document

## Overview

本设计文档描述了将现有 Python 智能客服系统重构为符合 LangChain/LangGraph 最佳实践的详细方案。重构采用分层架构，将系统划分为 API 层、服务层、Agent 层、工具层和数据层，每层职责清晰且相互解耦。

核心设计原则：
- **关注点分离**：不同功能模块独立组织，降低耦合
- **依赖注入**：通过参数传递依赖，提升可测试性
- **类型安全**：使用类型注解和 Pydantic 模型
- **配置驱动**：支持多环境和多租户配置
- **向后兼容**：保持现有 API 和配置接口不变

## Architecture

### 目录结构

```
work_v3/
├── src/
│   └── edu_agent/              # 主应用包
│       ├── __init__.py
│       ├── main.py             # FastAPI 应用入口
│       ├── api/                # API 路由层
│       │   ├── __init__.py
│       │   ├── deps.py         # 依赖注入
│       │   ├── chat.py         # 对话相关端点
│       │   ├── auth.py         # 认证相关端点
│       │   ├── admin.py        # 管理相关端点
│       │   └── models.py       # API 请求/响应模型
│       ├── services/           # 服务层
│       │   ├── __init__.py
│       │   ├── chat_service.py # 对话服务
│       │   ├── auth_service.py # 认证服务
│       │   └── vector_service.py # 向量服务
│       ├── agents/             # Agent 定义
│       │   ├── __init__.py
│       │   ├── base.py         # Agent 基类
│       │   ├── customer_service/ # 客服 Agent
│       │   │   ├── __init__.py
│       │   │   ├── agent.py    # Agent 构建
│       │   │   ├── nodes.py    # 节点实现
│       │   │   ├── state.py    # 状态定义
│       │   │   └── routes.py   # 路由逻辑
│       │   └── suggestion/     # 建议生成 Agent
│       │       ├── __init__.py
│       │       ├── agent.py
│       │       └── state.py
│       ├── tools/              # LangChain 工具
│       │   ├── __init__.py
│       │   ├── kb_tools.py     # 知识库工具
│       │   ├── order_tools.py  # 订单工具
│       │   ├── sql_tools.py    # SQL 工具
│       │   └── handoff_tools.py # 转人工工具
│       ├── core/               # 核心组件
│       │   ├── __init__.py
│       │   ├── config.py       # 配置管理
│       │   ├── llm.py          # LLM 工厂
│       │   ├── embeddings.py   # Embedding 工厂
│       │   ├── vector_store.py # 向量存储管理
│       │   └── checkpointer.py # 检查点管理
│       ├── db/                 # 数据库访问
│       │   ├── __init__.py
│       │   ├── orders.py       # 订单数据库
│       │   ├── support.py      # 支持数据库
│       │   └── auth.py         # 认证数据库
│       ├── middleware/         # 中间件
│       │   ├── __init__.py
│       │   ├── security.py     # 安全中间件
│       │   ├── logging.py      # 日志中间件
│       │   └── metrics.py      # 指标中间件
│       ├── schemas/            # 数据模型
│       │   ├── __init__.py
│       │   ├── chat.py         # 对话模型
│       │   ├── auth.py         # 认证模型
│       │   └── common.py       # 通用模型
│       ├── utils/              # 工具函数
│       │   ├── __init__.py
│       │   ├── text.py         # 文本处理
│       │   ├── time.py         # 时间处理
│       │   └── validators.py   # 验证器
│       └── prompts/            # 提示词模板
│           ├── __init__.py
│           ├── intent.py       # 意图识别
│           ├── rag.py          # RAG 提示词
│           └── order.py        # 订单提示词
├── tests/                      # 测试目录
│   ├── __init__.py
│   ├── conftest.py             # pytest 配置
│   ├── unit/                   # 单元测试
│   │   ├── test_tools.py
│   │   ├── test_nodes.py
│   │   └── test_services.py
│   ├── integration/            # 集成测试
│   │   ├── test_agents.py
│   │   └── test_api.py
│   └── fixtures/               # 测试数据
│       ├── __init__.py
│       └── factories.py
├── docs/                       # 文档
│   ├── architecture.md
│   ├── api.md
│   └── development.md
├── examples/                   # 示例代码
│   ├── basic_chat.py
│   └── custom_agent.py
├── pyproject.toml              # 项目配置
├── requirements.txt            # 生产依赖
├── requirements-dev.txt        # 开发依赖
└── README.md
```

### 分层架构

```
┌─────────────────────────────────────────┐
│          API Layer (FastAPI)            │
│  - 路由定义                              │
│  - 请求验证                              │
│  - 响应格式化                            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Service Layer                   │
│  - 业务逻辑编排                          │
│  - Agent 调用封装                        │
│  - 事务管理                              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Agent Layer (LangGraph)         │
│  - 状态机定义                            │
│  - 节点实现                              │
│  - 路由逻辑                              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Tool Layer (LangChain)          │
│  - 知识库检索                            │
│  - 数据库查询                            │
│  - 外部 API 调用                         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Data Layer                      │
│  - 向量存储                              │
│  - 关系数据库                            │
│  - 文件系统                              │
└─────────────────────────────────────────┘
```

## Components and Interfaces

### 1. API Layer

**职责**：处理 HTTP 请求，验证输入，调用服务层，格式化响应

**主要组件**：
- `api/chat.py`: 对话相关端点（/chat, /suggest, /greet）
- `api/auth.py`: 认证端点（/login, /register, /logout）
- `api/admin.py`: 管理端点（/users, /roles, /vectors）
- `api/deps.py`: 依赖注入函数（获取当前用户、租户等）
- `api/models.py`: Pydantic 请求/响应模型

**接口示例**：
```python
# api/chat.py
from fastapi import APIRouter, Depends
from ..services.chat_service import ChatService
from ..api.deps import get_chat_service, get_current_tenant
from ..api.models import ChatRequest, ChatResponse

router = APIRouter(prefix="/chat", tags=["chat"])

@router.post("/", response_model=ChatResponse)
async def chat(
    request: ChatRequest,
    service: ChatService = Depends(get_chat_service),
    tenant_id: str = Depends(get_current_tenant)
):
    return await service.process_message(
        query=request.query,
        thread_id=request.thread_id,
        tenant_id=tenant_id,
        images=request.images,
        audio=request.audio
    )
```

### 2. Service Layer

**职责**：封装业务逻辑，协调 Agent 和工具调用，处理事务

**主要组件**：
- `ChatService`: 处理对话请求，管理会话状态
- `AuthService`: 处理认证和授权
- `VectorService`: 管理向量存储的增删改查

**接口示例**：
```python
# services/chat_service.py
from typing import Optional, List
from ..agents.customer_service import CustomerServiceAgent
from ..core.config import Settings

class ChatService:
    def __init__(self, settings: Settings):
        self.settings = settings
        self.agent = CustomerServiceAgent(settings)
    
    async def process_message(
        self,
        query: str,
        thread_id: str,
        tenant_id: str,
        images: Optional[List[str]] = None,
        audio: Optional[str] = None
    ) -> dict:
        # 处理音频转文本
        if audio:
            query = await self._transcribe_audio(audio)
        
        # 构建状态
        state = {
            "query": query,
            "tenant_id": tenant_id,
            "images": images or []
        }
        
        # 调用 Agent
        result = await self.agent.invoke(state, thread_id)
        
        # 格式化响应
        return self._format_response(result)
```

### 3. Agent Layer

**职责**：定义状态机，实现节点逻辑，管理状态转换

**主要组件**：
- `CustomerServiceAgent`: 客服对话 Agent
- `SuggestionAgent`: 建议生成 Agent

**Agent 结构**：
```python
# agents/customer_service/agent.py
from langgraph.graph import StateGraph, START, END
from .state import CustomerServiceState
from .nodes import (
    intent_node,
    kb_node,
    order_node,
    direct_node,
    handoff_node
)
from .routes import decide_after_intent, decide_after_kb

class CustomerServiceAgent:
    def __init__(self, settings):
        self.settings = settings
        self.graph = self._build_graph()
    
    def _build_graph(self) -> StateGraph:
        graph = StateGraph(CustomerServiceState)
        
        # 添加节点
        graph.add_node("intent", intent_node)
        graph.add_node("kb", kb_node)
        graph.add_node("order", order_node)
        graph.add_node("direct", direct_node)
        graph.add_node("handoff", handoff_node)
        
        # 添加边
        graph.add_edge(START, "intent")
        graph.add_conditional_edges(
            "intent",
            decide_after_intent,
            {
                "kb": "kb",
                "order": "order",
                "direct": "direct",
                "handoff": "handoff"
            }
        )
        graph.add_conditional_edges(
            "kb",
            decide_after_kb,
            {"has_answer": END, "no_answer": "handoff"}
        )
        graph.add_edge("order", END)
        graph.add_edge("direct", END)
        graph.add_edge("handoff", END)
        
        # 编译
        checkpointer = self.settings.get_checkpointer()
        return graph.compile(checkpointer=checkpointer)
    
    async def invoke(self, state: dict, thread_id: str) -> dict:
        config = {"configurable": {"thread_id": thread_id}}
        return await self.graph.ainvoke(state, config)
```

**节点实现**：
```python
# agents/customer_service/nodes.py
from typing import Dict, Any
from .state import CustomerServiceState
from ...tools.kb_tools import retrieve_kb
from ...tools.order_tools import query_order
from ...core.llm import get_llm

async def intent_node(state: CustomerServiceState) -> Dict[str, Any]:
    """意图识别节点"""
    query = state["query"]
    
    # 关键词快速匹配
    if any(k in query for k in ["人工", "客服"]):
        return {"intent": "handoff"}
    if any(k in query for k in ["订单", "支付"]):
        return {"intent": "order"}
    
    # LLM 路由
    llm = get_llm()
    router = llm.with_structured_output(Route)
    result = await router.ainvoke(f"分类用户意图：{query}")
    
    return {"intent": result.step}

async def kb_node(state: CustomerServiceState) -> Dict[str, Any]:
    """知识库检索节点"""
    query = state["query"]
    tenant_id = state["tenant_id"]
    
    # 检索
    docs = await retrieve_kb(query, tenant_id)
    
    if not docs:
        return {"kb_answer": "", "sources": []}
    
    # 生成答案
    llm = get_llm()
    context = "\n\n".join(doc.page_content for doc in docs)
    prompt = f"基于以下资料回答：\n{context}\n\n问题：{query}"
    answer = await llm.ainvoke(prompt)
    
    sources = [
        {
            "content": doc.page_content[:200],
            "metadata": doc.metadata
        }
        for doc in docs
    ]
    
    return {"kb_answer": answer.content, "sources": sources}
```

### 4. Tool Layer

**职责**：封装可复用的操作，供 Agent 调用

**主要工具**：
- `retrieve_kb`: 知识库检索
- `query_order`: 订单查询
- `execute_sql`: SQL 执行
- `handoff_to_human`: 转人工

**工具实现**：
```python
# tools/kb_tools.py
from langchain_core.tools import tool
from typing import List
from ..core.vector_store import VectorStoreManager

@tool
async def retrieve_kb(query: str, tenant_id: str, k: int = 3) -> List[dict]:
    """从知识库检索相关文档
    
    Args:
        query: 查询文本
        tenant_id: 租户 ID
        k: 返回文档数量
    
    Returns:
        文档列表，每个文档包含 content 和 metadata
    """
    manager = VectorStoreManager()
    vector_store = manager.get_store(tenant_id)
    
    if not vector_store:
        return []
    
    docs = await vector_store.asimilarity_search(query, k=k)
    
    return [
        {
            "content": doc.page_content,
            "metadata": doc.metadata
        }
        for doc in docs
    ]
```

### 5. Core Components

**配置管理**：
```python
# core/config.py
from pydantic_settings import BaseSettings
from functools import lru_cache

class Settings(BaseSettings):
    # 模型配置
    model_name: str = "qwen-turbo"
    embedding_model: str = "text-embedding-v4"
    dashscope_api_key: str
    
    # 数据库配置
    orders_db_path: str = "db/orders.sqlite"
    support_db_path: str = "support.db"
    
    # 向量存储配置
    kb_index_dir: str = "faiss_index"
    tenants_base_dir: str = "tenants"
    
    # Redis 配置
    redis_url: str = "redis://127.0.0.1:6379/0"
    
    # 认证配置
    secret_key: str
    token_expire_days: int = 7
    
    class Config:
        env_file = ".env"
        case_sensitive = False

@lru_cache()
def get_settings() -> Settings:
    return Settings()
```

**LLM 工厂**：
```python
# core/llm.py
from langchain_community.chat_models import ChatTongyi
from functools import lru_cache
from .config import get_settings

@lru_cache()
def get_llm(model_name: str = None):
    settings = get_settings()
    name = model_name or settings.model_name
    return ChatTongyi(
        model=name,
        dashscope_api_key=settings.dashscope_api_key
    )
```

## Data Models

### State 定义

```python
# agents/customer_service/state.py
from typing import TypedDict, Optional, List, Literal

class CustomerServiceState(TypedDict, total=False):
    """客服 Agent 状态"""
    # 输入
    query: str                          # 用户查询
    tenant_id: str                      # 租户 ID
    images: List[str]                   # 图片列表（base64）
    history: str                        # 对话历史
    
    # 中间状态
    intent: Literal["kb", "order", "direct", "handoff"]  # 意图
    
    # 输出
    route: str                          # 最终路由
    kb_answer: str                      # 知识库答案
    order_summary: str                  # 订单摘要
    human_handoff: dict                 # 转人工信息
    sources: List[dict]                 # 来源列表
```

### API 模型

```python
# api/models.py
from pydantic import BaseModel, Field
from typing import Optional, List

class ChatRequest(BaseModel):
    """对话请求"""
    query: Optional[str] = None
    thread_id: Optional[str] = None
    user_id: Optional[str] = None
    images: Optional[List[str]] = None
    audio: Optional[str] = None
    asr_language: Optional[str] = "zh"
    asr_itn: Optional[bool] = True

class ChatResponse(BaseModel):
    """对话响应"""
    route: str = Field(..., description="路由类型")
    answer: str = Field(..., description="回答内容")
    sources: List[dict] = Field(default_factory=list, description="来源列表")
    audio_text: Optional[str] = Field(None, description="语音识别文本")
```

## C
orrectness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Agent StateGraph consistency

*For any* Agent implementation, the agent MUST use StateGraph to build its state machine and define all state transitions explicitly through add_node and add_edge/add_conditional_edges methods.

**Validates: Requirements 2.1**

### Property 2: Node function signature consistency

*For any* node function in an Agent, the function MUST accept a State parameter and return a dictionary containing state updates.

**Validates: Requirements 2.2**

### Property 3: Routing function validity

*For any* conditional routing function, the function MUST return a string that corresponds to a valid node name defined in the graph.

**Validates: Requirements 2.3**

### Property 4: Checkpointer configuration

*For any* Agent that requires state persistence, the agent MUST be compiled with a Checkpointer instance to enable conversation history and state recovery.

**Validates: Requirements 2.4**

### Property 5: LLM dependency injection

*For any* Agent or node that uses an LLM, the LLM instance MUST be obtained through dependency injection (passed as parameter or retrieved from config) rather than being instantiated directly in the code.

**Validates: Requirements 2.5**

### Property 6: Tool standardization

*For any* Tool implementation, the tool MUST either use the @tool decorator or inherit from the LangChain Tool class to ensure standard tool interface.

**Validates: Requirements 3.1**

### Property 7: Tool dependency injection

*For any* Tool that accesses external resources (databases, APIs, file systems), dependencies MUST be passed as function parameters rather than accessed through global variables or module-level imports.

**Validates: Requirements 3.2**

### Property 8: Tool error handling

*For any* Tool that performs operations that may fail, the tool MUST catch exceptions and return structured error information rather than allowing exceptions to propagate.

**Validates: Requirements 3.3**

### Property 9: Tool configuration usage

*For any* Tool that requires configuration parameters, the tool MUST retrieve parameters from the config module to support multi-tenant and environment-specific configurations.

**Validates: Requirements 3.4**

### Property 10: State type definition

*For any* State definition, the state MUST be defined using either TypedDict or Pydantic BaseModel to ensure type safety and validation.

**Validates: Requirements 4.1**

### Property 11: State field type annotations

*For any* State field, the field MUST have a type annotation to enable static type checking and IDE support.

**Validates: Requirements 4.2**

### Property 12: State type consistency

*For any* state passed between nodes, the state MUST maintain type consistency as defined in the State schema, verifiable through static type checking tools like mypy.

**Validates: Requirements 4.3**

### Property 13: Optional state fields

*For any* State field that may be absent, the field MUST be marked as Optional in the type annotation or have a default value specified.

**Validates: Requirements 4.4**

### Property 14: State inheritance

*For any* Agent-specific State that shares fields with other agents, the state MUST inherit from a base State type to promote code reuse and consistency.

**Validates: Requirements 4.5**

### Property 15: Configuration default values

*For any* configuration parameter, if the parameter is not provided through environment variables, the system MUST use a reasonable default value that allows the system to start successfully.

**Validates: Requirements 5.2**

### Property 16: No hardcoded secrets

*For any* sensitive configuration (API keys, passwords, tokens), the value MUST be loaded from environment variables and MUST NOT be hardcoded in source code.

**Validates: Requirements 5.3**

### Property 17: Configuration validation

*For any* required configuration parameter, if the parameter is missing or invalid, the system MUST raise a clear error message during startup indicating which configuration is problematic.

**Validates: Requirements 5.4**

### Property 18: Tenant-specific configuration

*For any* tenant in a multi-tenant system, the system MUST support tenant-specific configuration overrides that allow different tenants to use different models or data sources.

**Validates: Requirements 5.5**

### Property 19: API route separation

*For any* API route handler, the handler MUST be a thin layer that delegates business logic to service layer functions, containing only request parsing and response formatting logic.

**Validates: Requirements 6.1**

### Property 20: Service layer encapsulation

*For any* API route that needs to invoke an Agent, the route MUST call a service layer method rather than directly importing and invoking LangGraph components.

**Validates: Requirements 6.2**

### Property 21: Pydantic request validation

*For any* API endpoint, the endpoint MUST use Pydantic models for request and response schemas to enable automatic validation and documentation.

**Validates: Requirements 6.3**

### Property 22: Unified error response format

*For any* API error response, the response MUST follow a consistent error format (code, message, data) across all endpoints.

**Validates: Requirements 6.4**

### Property 23: Authentication dependency injection

*For any* API endpoint requiring authentication, authentication MUST be implemented as a FastAPI dependency function rather than inline authentication logic.

**Validates: Requirements 6.5**

### Property 24: Module docstrings

*For any* module containing complex logic, the module MUST include a module-level docstring explaining the design intent and key concepts.

**Validates: Requirements 8.2**

### Property 25: Function parameter documentation

*For any* function with parameters, the function MUST include a docstring with parameter descriptions and type information.

**Validates: Requirements 8.3**

### Property 26: Dependency version constraints

*For any* dependency in requirements.txt or pyproject.toml, major dependencies MUST specify version constraints to prevent incompatible updates.

**Validates: Requirements 9.2**

### Property 27: API endpoint backward compatibility

*For any* existing API endpoint from the original system, the refactored system MUST maintain the same endpoint path, HTTP method, and request/response structure.

**Validates: Requirements 10.1**

### Property 28: Environment variable backward compatibility

*For any* environment variable used in the original system, the refactored system MUST continue to recognize and use the same variable name and format.

**Validates: Requirements 10.2**

### Property 29: Data format backward compatibility

*For any* database schema or file format used in the original system, the refactored system MUST be able to read and write data in the same format without requiring migration.

**Validates: Requirements 10.3**

### Property 30: Feature flag control

*For any* new feature introduced during refactoring, the feature MUST be controlled by a feature flag that is disabled by default to maintain backward compatibility.

**Validates: Requirements 10.5**

## Error Handling

### Error Categories

1. **Configuration Errors**: Missing or invalid configuration parameters
   - Handled at startup with clear error messages
   - System fails fast if critical config is missing

2. **Validation Errors**: Invalid request data or state
   - Caught by Pydantic validators
   - Returned as 400 Bad Request with detailed field errors

3. **Resource Errors**: Database, vector store, or external API failures
   - Caught in tool layer
   - Logged with context
   - Returned as structured error to caller
   - Retry logic for transient failures

4. **Agent Errors**: LLM failures, routing errors, state inconsistencies
   - Caught in node functions
   - Logged with state snapshot
   - Fallback to safe default behavior (e.g., handoff to human)

5. **Authentication Errors**: Invalid tokens, expired sessions
   - Caught by auth middleware
   - Returned as 401 Unauthorized

### Error Response Format

```python
{
    "code": "error_code",           # Machine-readable error code
    "message": "Human readable",    # User-friendly message
    "data": {                       # Optional additional context
        "field": "error detail"
    }
}
```

### Retry Strategy

- **Vector Store Operations**: 3 retries with exponential backoff
- **Database Queries**: 2 retries with 100ms delay
- **LLM Calls**: 2 retries with 500ms delay
- **External APIs**: No automatic retry (fail fast)

### Logging Strategy

- **Structured Logging**: JSON format with timestamp, level, module, message
- **Context Injection**: Request ID, tenant ID, user ID in all logs
- **Sensitive Data Redaction**: Automatic redaction of PII and secrets
- **Log Levels**:
  - DEBUG: Detailed state transitions, tool inputs/outputs
  - INFO: Request/response, agent invocations
  - WARNING: Retries, fallback behaviors
  - ERROR: Exceptions, failures

## Testing Strategy

### Unit Testing

**Scope**: Individual functions, tools, and utilities

**Approach**:
- Test each tool in isolation with mocked dependencies
- Test node functions with sample state inputs
- Test utility functions with edge cases
- Use pytest fixtures for common test data

**Example**:
```python
# tests/unit/test_tools.py
from unittest.mock import Mock, patch
from edu_agent.tools.kb_tools import retrieve_kb

@patch('edu_agent.core.vector_store.VectorStoreManager')
def test_retrieve_kb_returns_documents(mock_manager):
    # Arrange
    mock_store = Mock()
    mock_store.asimilarity_search.return_value = [
        Mock(page_content="doc1", metadata={"source": "test"})
    ]
    mock_manager.return_value.get_store.return_value = mock_store
    
    # Act
    result = await retrieve_kb("test query", "tenant1")
    
    # Assert
    assert len(result) == 1
    assert result[0]["content"] == "doc1"
```

### Property-Based Testing

**Scope**: Universal properties that should hold across all inputs

**Library**: Hypothesis (Python property-based testing library)

**Configuration**: Each property test runs minimum 100 iterations

**Key Properties to Test**:

1. **State Consistency**: For any valid state input to a node, the output must be a valid state update
2. **Routing Validity**: For any state, routing functions must return valid node names
3. **Type Safety**: For any state, all fields must match their type annotations
4. **Error Handling**: For any tool with invalid inputs, must return structured error
5. **Configuration Defaults**: For any missing config, system must use valid defaults

**Example**:
```python
# tests/property/test_nodes.py
from hypothesis import given, strategies as st
from edu_agent.agents.customer_service.nodes import intent_node
from edu_agent.agents.customer_service.state import CustomerServiceState

@given(st.text(min_size=1, max_size=1000))
async def test_intent_node_always_returns_valid_intent(query):
    """
    Feature: langgraph-refactor, Property 3: Routing function validity
    
    For any query string, intent_node must return a state update
    with an 'intent' field containing a valid intent value.
    """
    # Arrange
    state = CustomerServiceState(
        query=query,
        tenant_id="test",
        images=[],
        history=""
    )
    
    # Act
    result = await intent_node(state)
    
    # Assert
    assert "intent" in result
    assert result["intent"] in ["kb", "order", "direct", "handoff"]
```

### Integration Testing

**Scope**: Agent workflows, service layer, database interactions

**Approach**:
- Test complete agent execution from start to end
- Test service layer with real (test) databases
- Test API endpoints with TestClient
- Use test fixtures for database setup/teardown

**Example**:
```python
# tests/integration/test_agents.py
import pytest
from edu_agent.agents.customer_service import CustomerServiceAgent
from edu_agent.core.config import Settings

@pytest.mark.asyncio
async def test_customer_service_agent_kb_flow():
    # Arrange
    settings = Settings(
        kb_index_dir="tests/fixtures/test_index",
        model_name="qwen-turbo"
    )
    agent = CustomerServiceAgent(settings)
    state = {
        "query": "课程费用是多少",
        "tenant_id": "test",
        "images": []
    }
    
    # Act
    result = await agent.invoke(state, "test-thread-123")
    
    # Assert
    assert result["route"] == "kb"
    assert "kb_answer" in result
    assert len(result["kb_answer"]) > 0
```

### End-to-End Testing

**Scope**: Complete API workflows from HTTP request to response

**Approach**:
- Use FastAPI TestClient
- Test authentication flows
- Test multi-step conversations
- Test error scenarios

**Example**:
```python
# tests/integration/test_api.py
from fastapi.testclient import TestClient
from edu_agent.main import app

client = TestClient(app)

def test_chat_endpoint_returns_answer():
    # Act
    response = client.post(
        "/chat",
        json={
            "query": "课程费用",
            "thread_id": "test-123"
        },
        headers={"X-Tenant-ID": "default"}
    )
    
    # Assert
    assert response.status_code == 200
    data = response.json()
    assert "answer" in data
    assert "route" in data
```

### Test Data Management

**Fixtures**:
- `conftest.py`: Shared fixtures for database, vector store, settings
- `factories.py`: Factory functions for generating test data

**Test Databases**:
- SQLite in-memory databases for unit tests
- Temporary FAISS indices for vector store tests
- Isolated test directories for file-based storage

**Mocking Strategy**:
- Mock external APIs (LLM, ASR) in unit tests
- Use real implementations in integration tests with test data
- Mock time-dependent operations for deterministic tests

## Migration Strategy

### Phase 1: Preparation (Week 1)

1. Create new directory structure
2. Set up testing framework
3. Document current API contracts
4. Create compatibility layer

### Phase 2: Core Components (Week 2)

1. Migrate configuration management
2. Migrate LLM and embedding factories
3. Migrate vector store management
4. Add tests for core components

### Phase 3: Tools and State (Week 3)

1. Migrate and refactor tools
2. Define state schemas
3. Add tool tests
4. Verify backward compatibility

### Phase 4: Agents (Week 4)

1. Migrate customer service agent
2. Migrate suggestion agent
3. Add agent tests
4. Verify state transitions

### Phase 5: API Layer (Week 5)

1. Migrate API routes
2. Implement service layer
3. Add API tests
4. Verify endpoint compatibility

### Phase 6: Integration and Cleanup (Week 6)

1. Integration testing
2. Performance testing
3. Documentation updates
4. Remove old code
5. Final compatibility verification

### Rollback Plan

- Keep old code in separate branch
- Feature flag for new vs old implementation
- Gradual traffic migration (10% -> 50% -> 100%)
- Monitoring and alerting for errors
- Quick rollback capability if issues detected

## Performance Considerations

### Optimization Targets

- **API Response Time**: < 2 seconds for 95th percentile
- **Vector Search**: < 500ms for similarity search
- **Database Queries**: < 100ms for order lookups
- **LLM Calls**: < 1.5 seconds (dependent on model)

### Caching Strategy

- **Vector Store**: Cache loaded indices in memory
- **Configuration**: Cache settings with TTL
- **Session History**: Redis cache with 30-minute TTL
- **LLM Responses**: Optional caching for common queries

### Resource Management

- **Connection Pooling**: Database connection pools
- **Async Operations**: Use async/await throughout
- **Lazy Loading**: Load resources on-demand
- **Graceful Degradation**: Fallback to simpler operations on timeout

## Security Considerations

### Authentication

- JWT tokens with expiration
- Token refresh mechanism
- Secure password hashing (bcrypt)
- Rate limiting on auth endpoints

### Authorization

- Role-based access control (RBAC)
- Tenant isolation
- API key validation for admin endpoints

### Data Protection

- PII redaction in logs
- Encrypted storage for sensitive data
- Secure environment variable handling
- Input validation and sanitization

### API Security

- CORS configuration
- Request size limits
- SQL injection prevention (parameterized queries)
- XSS prevention (output encoding)
