# 实现计划

- [x] 1. 项目初始化和基础设施搭建
  - 创建 Go 模块和项目目录结构（Clean Architecture）
  - 配置 go.mod 依赖（Gin、GORM、Eino、Milvus SDK）
  - 创建配置文件结构（config.yaml）
  - 设置环境变量管理（.env 文件）
  - _需求: 1.1_

- [x] 2. Domain Layer 实现
  - 定义核心实体（Message、Intent、Document、Order）
  - 定义仓储接口（VectorRepository、OrderRepository、SessionRepository）
  - 定义值对象和业务规则
  - _需求: 所有需求的基础_

- [x] 3. Infrastructure Layer - 配置和日志
  - 实现配置加载（infrastructure/config）
  - 实现结构化日志（infrastructure/logger）
  - 实现 DashScope 客户端初始化
  - _需求: 1.1, 7.1, 7.2_

- [x] 4. Infrastructure Layer - Milvus 集成
  - 实现 Milvus 连接管理
  - 实现 VectorRepository（Milvus）
  - 实现 Collection 创建和 Schema 定义
  - 实现多租户 Collection 管理
  - _需求: 3.2, 5.3, 5.4, 9.3, 9.4_

- [x] 5. Infrastructure Layer - SQLite 集成
  - 实现 GORM 数据库连接
  - 定义 GORM 模型（Order、Session、MissedQuery）
  - 实现 OrderRepository（SQLite）
  - 实现 SessionRepository（SQLite）
  - 实现多租户数据库文件管理
  - _需求: 4.4, 5.5, 10.2, 10.3, 10.4_

- [x] 6. Infrastructure Layer - Eino AI 集成
  - 实现 DashScope ChatModel 初始化
  - 实现 DashScope Embedding 模型初始化
  - 实现意图识别器（IntentRecognizer）
  - 实现 RAG 检索器（RAGRetriever）
  - 实现订单查询器（OrderQuerier）
  - 实现响应生成器（ResponseGenerator）
  - _需求: 1.1, 2.1, 3.1, 4.1, 4.2_

- [x] 7. Use Case Layer - Chat 用例
  - 实现 ChatUseCase 核心逻辑
  - 实现意图路由逻辑
  - 实现会话历史管理
  - 集成 Eino ADK Sequential Agent
  - 集成 Eino ADK Parallel Agent（信息收集）
  - _需求: 1.2, 2.1, 2.3, 2.4, 2.5, 10.1, 10.2, 10.3, 10.4_

- [x] 8. Use Case Layer - Vector 管理用例
  - 实现 VectorManagementUseCase
  - 实现向量添加逻辑
  - 实现向量删除逻辑
  - 实现批量操作支持
  - _需求: 9.2, 9.3, 9.4, 9.5_

- [x] 9. Interface Adapter Layer - HTTP 中间件
  - 实现租户识别中间件（TenantMiddleware）
  - 实现安全脱敏中间件（SecurityMiddleware）
  - 实现日志记录中间件（LoggingMiddleware）
  - 实现错误处理中间件（ErrorHandler）
  - 实现 API Key 验证中间件（AuthMiddleware）
  - _需求: 5.1, 5.2, 7.1, 8.1, 8.2, 8.3, 8.4, 9.1_

- [x] 10. Interface Adapter Layer - HTTP Handler
  - 实现 ChatHandler（/chat 端点）
  - 实现 VectorHandler（/api/v1/vectors/* 端点）
  - 实现 HealthHandler（/health 端点）
  - 实现 ModelHandler（/models/* 端点）
  - 实现流式响应支持（SSE）
  - _需求: 6.1, 6.2, 6.3, 6.4, 6.5, 7.5_

- [x] 11. Interface Adapter Layer - Router 和服务器
  - 实现 Gin 路由配置
  - 实现服务器启动和优雅关闭
  - 集成所有中间件和 Handler
  - 实现健康检查逻辑
  - _需求: 6.1, 7.5_

- [x] 12. 多租户管理器实现
  - 实现 TenantManager
  - 实现租户 Collection 映射
  - 实现租户数据库映射
  - 实现租户资源自动创建
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 13. 日志和指标系统
  - 实现结构化日志输出
  - 实现请求日志记录
  - 实现错误日志记录
  - 实现指标收集（请求计数、响应时间）
  - 实现健康检查指标
  - _需求: 7.1, 7.2, 7.3, 7.4, 7.5_

- [x] 14. 错误处理和重试机制
  - 实现统一错误类型定义
  - 实现错误处理中间件
  - 实现重试策略（指数退避）
  - 实现降级策略
  - _需求: 1.4, 1.5_

- [x] 15. 应用入口和依赖注入
  - 实现 main.go 入口
  - 实现依赖注入容器
  - 实现组件初始化顺序
  - 实现优雅关闭
  - _需求: 1.1, 6.1_

- [x] 16. 配置和文档
  - 创建 config.yaml 示例
  - 创建 .env.example 文件
  - 编写 README.md
  - 编写 API 文档
  - 编写部署指南
  - _需求: 所有需求_

- [ ] 17. Docker 和部署
  - 创建 Dockerfile
  - 创建 docker-compose.yml（包含 Milvus、SQLite）
  - 配置健康检查
  - 测试容器化部署
  - _需求: 所有需求_

- [ ] 18. 检查点 - 基础功能验证
  - 确保系统可以启动并响应基本请求，如有问题请询问用户

- [ ]* 19. 单元测试 - Domain 层
  - 测试实体创建和验证逻辑
  - 测试值对象的不变性
  - _需求: 2.1, 2.2_

- [ ]* 20. 单元测试 - Infrastructure 层
  - 测试 Milvus Repository（向量插入、删除、搜索）
  - 测试 SQLite Repository（订单查询、会话存储）
  - 测试租户隔离
  - _需求: 3.2, 4.4, 5.3, 5.5, 9.3, 9.4, 10.2, 10.3_

- [ ]* 21. 单元测试 - AI 组件
  - 测试意图识别（课程、订单、直接、人工）
  - 测试 RAG 检索
  - 测试订单查询
  - _需求: 2.1, 2.4, 2.5, 3.1, 4.1_

- [ ]* 22. 单元测试 - Use Case 层
  - 测试 Chat 用例（各种意图流程）
  - 测试 Vector 管理用例
  - _需求: 2.4, 2.5, 9.2, 9.3, 9.4_

- [ ]* 23. 单元测试 - HTTP 层
  - 测试中间件（租户、安全、日志、认证）
  - 测试 Handler（/chat、/vectors、/health）
  - _需求: 5.1, 6.2, 6.3, 7.5, 8.1, 8.2, 8.3_

- [ ]* 24. 单元测试 - 多租户和日志
  - 测试多租户管理器
  - 测试日志系统
  - 测试错误处理
  - _需求: 5.3, 5.4, 7.1, 7.2, 7.3, 1.5_

- [ ]* 25. 属性测试 - 意图识别
  - **Property 1: 意图分类完整性**
  - **Property 2: 意图识别结果结构完整性**
  - **验证需求: 2.1, 2.2**

- [ ]* 26. 属性测试 - RAG 系统
  - **Property 5: RAG 向量生成**
  - **Property 6: RAG 搜索执行**
  - **Property 7: RAG 结果数量限制**
  - **Property 8: RAG 文档传递给 LLM**
  - **验证需求: 3.1, 3.2, 3.3, 3.4**

- [ ]* 27. 属性测试 - 订单查询
  - **Property 9: 订单 ID 提取**
  - **Property 10: SQL 安全性验证**
  - **Property 11: SQL 注入防护**
  - **Property 12: 订单查询结果格式化**
  - **验证需求: 4.1, 4.2, 4.3, 4.5, 8.5**

- [ ]* 28. 属性测试 - 多租户系统
  - **Property 13: 租户 ID 提取**
  - **Property 14: 租户 Collection 映射**
  - **Property 15: 租户数据库隔离**
  - **验证需求: 5.1, 5.3, 5.5**

- [ ]* 29. 属性测试 - HTTP 接口
  - **Property 16: HTTP 请求解析**
  - **Property 17: HTTP 响应格式**
  - **Property 18: HTTP 错误响应**
  - **验证需求: 6.2, 6.3, 6.5**

- [ ]* 30. 属性测试 - 日志和指标
  - **Property 19: 日志记录完整性**
  - **Property 20: 外部调用日志**
  - **Property 21: 指标统计**
  - **验证需求: 7.1, 7.3, 7.4**

- [ ]* 31. 属性测试 - 安全系统
  - **Property 22: 敏感信息脱敏**
  - **Property 23: 输入脱敏**
  - **Property 24: API Key 验证**
  - **验证需求: 8.1, 8.2, 8.3, 9.1**

- [ ]* 32. 属性测试 - 向量管理
  - **Property 25: 向量生成**
  - **Property 26: 向量插入**
  - **Property 27: 向量删除**
  - **Property 28: 向量操作结果**
  - **验证需求: 9.2, 9.3, 9.4, 9.5**

- [ ]* 33. 属性测试 - 会话管理
  - **Property 29: 会话 ID 唯一性**
  - **Property 30: 消息历史追加**
  - **Property 31: 上下文加载**
  - **验证需求: 10.1, 10.2, 10.3, 10.4**

- [ ]* 34. 集成测试
  - 编写端到端对话流程测试
  - 编写多租户隔离测试
  - 编写向量管理流程测试
  - _需求: 所有需求_

- [ ]* 35. 性能测试
  - 编写 /chat 端点性能测试
  - 编写 RAG 检索性能测试
  - 编写并发请求测试
  - 验证 P95 响应时间 < 500ms
  - _需求: 所有需求_

- [ ] 36. 最终检查点
  - 确保所有测试通过，如有问题请询问用户
