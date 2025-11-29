# 需求文档

## 简介

基于 CloudWeGo Eino 框架构建的 Go 语言智能客服系统，用于替代现有的 Python 技术栈（LangChain + LangGraph + FastAPI）实现。系统将提供多租户支持、RAG 知识库检索、订单查询、意图识别等核心功能，并充分利用 Go 语言的高性能和并发处理能力。

系统采用 **ADK 为主、Compose 为辅**的混合架构模式：
- 主对话流程使用 Eino ADK 的顺序、并行和循环模式进行高层编排
- 性能关键路径和细粒度控制使用 Eino Compose 进行底层图编排
- 充分发挥 ADK 的开发效率优势和 Compose 的性能优化能力

## 术语表

- **Eino**: CloudWeGo 开源的 LLM 应用开发框架
- **ADK**: Agent Development Kit，Eino 提供的智能体开发工具包
- **Compose**: Eino 的图编排组件，用于构建数据流处理管道
- **DashScope**: 阿里云通义千问模型服务平台
- **RAG**: Retrieval-Augmented Generation，检索增强生成
- **Milvus**: 开源向量数据库，用于高性能向量相似度搜索
- **Gin**: 高性能的 Go HTTP Web 框架
- **GORM**: Go 语言的 ORM 库，用于数据库操作
- **SQLite**: 轻量级嵌入式关系数据库
- **租户**: Tenant，系统中的独立业务单元，拥有独立的知识库和数据
- **意图**: Intent，用户查询的目的分类（如课程咨询、订单查询等）
- **会话**: Session，用户与系统的一次完整对话交互
- **检查点**: Checkpoint，对话状态的持久化快照
- **Collection**: Milvus 中的向量集合，类似于数据库中的表

## 需求

### 需求 1

**用户故事:** 作为系统架构师，我希望使用 Eino 框架构建核心对话引擎，以便获得更好的性能和可维护性。

#### 验收标准

1. WHEN 系统初始化 THEN Eino 系统 SHALL 成功加载 DashScope 聊天模型配置
2. WHEN 用户发送对话请求 THEN Eino 系统 SHALL 使用 Compose 图编排处理请求流程
3. WHEN 对话流程执行 THEN Eino 系统 SHALL 支持流式响应输出
4. WHEN 系统运行 THEN Eino 系统 SHALL 提供统一的错误处理和日志追踪机制
5. WHEN 模型调用失败 THEN Eino 系统 SHALL 执行重试策略并记录错误日志

### 需求 2

**用户故事:** 作为用户，我希望系统能够识别我的查询意图，以便获得准确的响应路由。

#### 验收标准

1. WHEN 用户输入查询文本 THEN 意图识别系统 SHALL 将查询分类为课程咨询、订单查询、直接回答或人工转接四种类型之一
2. WHEN 意图识别完成 THEN 意图识别系统 SHALL 返回意图类型和置信度分数
3. WHEN 置信度低于阈值 THEN 意图识别系统 SHALL 将查询路由到人工转接
4. WHEN 意图为课程咨询 THEN 意图识别系统 SHALL 触发 RAG 知识库检索流程
5. WHEN 意图为订单查询 THEN 意图识别系统 SHALL 触发订单数据库查询流程

### 需求 3

**用户故事:** 作为用户，我希望系统能够从知识库中检索相关信息，以便回答我的课程相关问题。

#### 验收标准

1. WHEN 系统接收到课程咨询查询 THEN RAG 系统 SHALL 使用 DashScope 嵌入模型生成查询向量
2. WHEN 查询向量生成完成 THEN RAG 系统 SHALL 在 Milvus 向量数据库中执行相似度搜索
3. WHEN 相似度搜索完成 THEN RAG 系统 SHALL 返回前 K 个最相关的文档片段
4. WHEN 检索到相关文档 THEN RAG 系统 SHALL 将文档内容和查询一起发送给 LLM 生成答案
5. WHEN 未检索到相关文档 THEN RAG 系统 SHALL 记录未命中事件并触发人工转接

### 需求 4

**用户故事:** 作为用户，我希望系统能够查询我的订单信息，以便了解订单状态和详情。

#### 验收标准

1. WHEN 用户请求订单查询 THEN 订单查询系统 SHALL 从用户输入中提取订单 ID
2. WHEN 订单 ID 提取完成 THEN 订单查询系统 SHALL 生成安全的 SQL 查询语句
3. WHEN SQL 查询生成 THEN 订单查询系统 SHALL 验证查询语句的安全性（防止 SQL 注入）
4. WHEN 查询验证通过 THEN 订单查询系统 SHALL 在 SQLite 数据库中执行查询
5. WHEN 查询结果返回 THEN 订单查询系统 SHALL 将结构化数据格式化为自然语言响应

### 需求 5

**用户故事:** 作为租户管理员，我希望系统支持多租户隔离，以便不同业务单元拥有独立的数据和配置。

#### 验收标准

1. WHEN 请求包含租户标识 THEN 多租户系统 SHALL 从请求头或查询参数中提取租户 ID
2. WHEN 租户 ID 缺失 THEN 多租户系统 SHALL 使用默认租户标识
3. WHEN 租户 ID 确定 THEN 多租户系统 SHALL 使用该租户对应的 Milvus Collection
4. WHEN 租户 Collection 不存在 THEN 多租户系统 SHALL 创建新的 Collection 并初始化 Schema
5. WHEN 执行数据库操作 THEN 多租户系统 SHALL 使用 GORM 操作租户独立的 SQLite 数据库文件

### 需求 6

**用户故事:** 作为系统管理员，我希望系统提供 HTTP API 接口，以便前端和其他服务能够集成。

#### 验收标准

1. WHEN 系统启动 THEN Gin HTTP 服务器 SHALL 在配置的端口上监听请求
2. WHEN 接收到 POST /chat 请求 THEN Gin HTTP 服务器 SHALL 解析 JSON 请求体并提取查询参数
3. WHEN 对话处理完成 THEN Gin HTTP 服务器 SHALL 返回 JSON 格式的响应包含答案和元数据
4. WHEN 客户端请求流式响应 THEN Gin HTTP 服务器 SHALL 使用 SSE（Server-Sent Events）协议推送增量内容
5. WHEN 请求处理失败 THEN Gin HTTP 服务器 SHALL 返回适当的 HTTP 状态码和错误信息

### 需求 7

**用户故事:** 作为系统管理员，我希望系统记录详细的日志和指标，以便监控系统运行状态和排查问题。

#### 验收标准

1. WHEN 系统处理请求 THEN 日志系统 SHALL 记录请求 ID、租户 ID、查询内容和处理时长
2. WHEN 发生错误 THEN 日志系统 SHALL 记录错误堆栈和上下文信息
3. WHEN 调用外部服务 THEN 日志系统 SHALL 记录调用参数和响应时间
4. WHEN 系统运行 THEN 指标系统 SHALL 统计各类请求的数量和平均响应时间
5. WHEN 访问健康检查接口 THEN 指标系统 SHALL 返回系统状态和关键指标快照

### 需求 8

**用户故事:** 作为安全管理员，我希望系统能够保护敏感信息，以便符合数据安全和隐私保护要求。

#### 验收标准

1. WHEN 系统记录日志 THEN 安全系统 SHALL 对敏感字段（如密码、身份证号）进行脱敏处理
2. WHEN 用户输入包含敏感信息 THEN 安全系统 SHALL 在处理前替换为占位符
3. WHEN 访问管理接口 THEN 安全系统 SHALL 验证 API Key 的有效性
4. WHEN API Key 验证失败 THEN 安全系统 SHALL 拒绝请求并返回 401 未授权错误
5. WHEN 执行 SQL 查询 THEN 安全系统 SHALL 验证查询语句不包含危险操作（如 DROP、DELETE）

### 需求 9

**用户故事:** 作为开发者，我希望系统提供向量库管理接口，以便动态添加和删除知识库内容。

#### 验收标准

1. WHEN 接收到添加向量请求 THEN 向量管理系统 SHALL 验证请求的 API Key
2. WHEN API Key 验证通过 THEN 向量管理系统 SHALL 使用 Eino 嵌入模型生成文本向量
3. WHEN 向量生成完成 THEN 向量管理系统 SHALL 将向量插入到 Milvus Collection 中
4. WHEN 接收到删除向量请求 THEN 向量管理系统 SHALL 根据文档 ID 从 Milvus Collection 中删除向量
5. WHEN 向量操作完成 THEN 向量管理系统 SHALL 返回操作结果和受影响的记录数

### 需求 10

**用户故事:** 作为用户，我希望系统能够维护对话上下文，以便进行多轮对话交互。

#### 验收标准

1. WHEN 用户开始新对话 THEN 会话管理系统 SHALL 创建唯一的会话 ID
2. WHEN 用户发送消息 THEN 会话管理系统 SHALL 将消息添加到会话历史中
3. WHEN 系统生成响应 THEN 会话管理系统 SHALL 将响应添加到会话历史中
4. WHEN 处理后续消息 THEN 会话管理系统 SHALL 加载完整的会话历史作为上下文
5. WHEN 会话超时或结束 THEN 会话管理系统 SHALL 清理会话数据释放资源
