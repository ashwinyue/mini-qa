# 更新日志

本文档记录 Eino QA System 的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 计划中
- Kubernetes Helm Chart
- Redis 缓存支持
- 更多 LLM 模型支持
- 流式响应优化

## [1.0.0] - 2024-11-29

### 新增
- ✅ 基于 Eino 框架的智能对话系统
- ✅ 意图识别（课程咨询、订单查询、直接回答、人工转接）
- ✅ RAG 知识库检索（基于 Milvus）
- ✅ 订单查询功能（基于 SQLite）
- ✅ 多租户支持（独立的向量 Collection 和数据库）
- ✅ HTTP API 接口（Gin 框架）
- ✅ 向量管理接口（添加/删除文档）
- ✅ 健康检查接口
- ✅ 流式响应支持（SSE）
- ✅ 结构化日志系统
- ✅ 指标收集和监控
- ✅ 敏感信息脱敏
- ✅ API Key 认证
- ✅ SQL 注入防护
- ✅ 错误处理和重试机制
- ✅ 优雅关闭
- ✅ Docker 支持
- ✅ Docker Compose 配置
- ✅ 完整的文档（API、部署、架构）

### 架构
- Clean Architecture 分层设计
- Domain Layer（领域层）
- Use Case Layer（用例层）
- Interface Adapter Layer（接口适配层）
- Infrastructure Layer（基础设施层）

### 技术栈
- Go 1.23+
- CloudWeGo Eino (ADK + Compose)
- Gin Web Framework
- GORM ORM
- SQLite 数据库
- Milvus 向量数据库
- DashScope (通义千问)

### 文档
- README.md - 项目概述和快速开始
- QUICKSTART.md - 5 分钟快速上手指南
- STARTUP_GUIDE.md - 详细启动指南
- PROJECT_STRUCTURE.md - 项目结构说明
- API_DOCUMENTATION.md - 完整 API 文档
- DEPLOYMENT_GUIDE.md - 部署指南
- 设计文档 - 架构设计和技术选型
- 需求文档 - 功能需求和验收标准

### 配置
- config.yaml - 主配置文件
- .env.example - 环境变量模板
- Dockerfile - 容器化配置
- docker-compose.yml - 服务编排配置
- Makefile - 构建和开发命令

## [0.1.0] - 2024-11-15

### 新增
- 项目初始化
- 基础目录结构
- Go 模块配置
- 基础依赖安装

---

## 版本说明

### 版本号格式

版本号格式：`主版本号.次版本号.修订号`

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

### 变更类型

- **新增**：新功能
- **变更**：现有功能的变更
- **废弃**：即将移除的功能
- **移除**：已移除的功能
- **修复**：Bug 修复
- **安全**：安全相关的修复

---

## 贡献指南

如需贡献代码或报告问题，请：
1. 查看 [CONTRIBUTING.md](CONTRIBUTING.md)
2. 提交 [Issue](https://github.com/your-repo/issues)
3. 创建 [Pull Request](https://github.com/your-repo/pulls)

---

**维护者**: Eino QA Team  
**最后更新**: 2024-11-29
