# React Go Admin

一个现代化的前后端分离的用户权限管理系统，基于 React + Vite + Ant Design 前端框架和 Go + Gin + GORM 后端框架开发。

## 🚀 项目特性

- **现代化前端**: React 18 + Vite + Ant Design，支持响应式设计
- **高性能后端**: Go 1.21 + Gin 框架，RESTful API 设计
- **权限管理**: 基于角色的访问控制（RBAC）
- **用户管理**: 完整的用户 CRUD 操作
- **安全认证**: JWT Token 认证机制
- **数据持久化**: SQLite 数据库，支持数据迁移
- **开发友好**: 热更新、代码分割、组件化开发

## 🛠️ 技术栈

### 前端技术
- **React 18** - 现代化前端框架
- **Vite** - 快速的构建工具
- **Ant Design** - 企业级 UI 组件库
- **React Router** - 客户端路由管理
- **Axios** - HTTP 客户端
- **@ant-design/icons** - Ant Design 图标库
- **dayjs** - 日期处理库

### 后端技术
- **Go 1.21** - 高性能后端语言
- **Gin Framework** - 轻量级 Web 框架
- **GORM** - ORM 框架
- **JWT** - JSON Web Token 认证
- **golang-jwt/jwt** - JWT 库
- **gin-contrib/cors** - CORS 中间件

### 数据库
- **SQLite** - 轻量级关系型数据库

## 📦 项目结构

```
react-go-admin/
├── backend/              # 后端代码
│   ├── api/             # API 控制器
│   ├── config/          # 配置管理
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   ├── utils/           # 工具函数
│   └── main.go          # 后端入口文件
├── src/                 # 前端代码
│   ├── api/             # API 接口封装
│   ├── components/      # 公共组件
│   ├── context/         # React 上下文
│   ├── layouts/         # 布局组件
│   ├── pages/           # 页面组件
│   ├── router/          # 路由配置
│   ├── utils/           # 工具函数
│   └── App.jsx          # 根组件
├── docs/                # 项目文档
└── package.json         # 前端依赖配置
```

## 🚦 快速开始

### 环境要求

- Node.js >= 16.0.0
- Go >= 1.21
- npm 或 yarn

### 安装依赖

1. **安装前端依赖**
   ```bash
   npm install
   ```

2. **安装后端依赖**
   ```bash
   cd backend
   go mod download
   ```

### 启动项目

1. **启动后端服务**
   ```bash
   cd backend
   go run main.go
   ```
   后端服务将在 `http://localhost:8000` 启动

2. **启动前端服务**
   ```bash
   npm run dev
   ```
   前端服务将在 `http://localhost:5173` 启动

3. **访问系统**
   打开浏览器访问 `http://localhost:5173`，使用默认管理员账号登录：
   - 用户名：admin
   - 密码：123456

### 构建项目

1. **构建前端**
   ```bash
   npm run build
   ```

2. **构建后端**
   ```bash
   cd backend
   go build -o admin-server main.go
   ```

## 📋 功能模块

### 认证模块
- 用户登录/登出
- JWT Token 认证
- 用户信息验证

### 用户管理
- 用户列表展示
- 用户创建、编辑、删除
- 用户状态管理
- 用户搜索和筛选

### 角色管理
- 角色列表展示
- 角色创建、编辑、删除
- 角色权限配置
- 基于角色的访问控制

### 系统功能
- 响应式布局设计
- 国际化支持（中文）
- 错误页面处理
- 权限验证

## 🔧 API 接口

系统提供完整的 RESTful API，主要接口包括：

### 认证接口
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/verify` - 验证 Token

### 用户管理接口
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### 角色管理接口
- `GET /api/roles` - 获取角色列表
- `POST /api/roles` - 创建角色
- `PUT /api/roles/:id` - 更新角色
- `DELETE /api/roles/:id` - 删除角色

详细的 API 文档请参考 [API 规范说明](docs/api_specification.md)

## 🏗️ 架构设计

系统采用前后端分离架构，详细的架构设计请参考 [架构设计文档](docs/architecture_design.md)

### 核心特性
- **模块化设计**: 前后端独立开发和部署
- **RESTful API**: 统一的接口规范
- **安全认证**: JWT Token 和角色权限控制
- **数据验证**: 前后端双重验证
- **错误处理**: 统一的错误响应格式

## 🔒 安全设计

### 认证机制
- JWT Token 认证，有效期 24 小时
- 支持 Token 刷新机制

### 授权机制
- 基于角色的访问控制（RBAC）
- 接口级别的权限验证

### 安全措施
- 密码使用 bcrypt 加密存储
- CORS 跨域保护
- 输入参数验证
- SQL 注入防护

## 🚀 部署指南

详细的部署指南请参考 [部署指南](docs/deployment_guide.md)

### 快速部署
1. 构建前端和后端
2. 配置环境变量
3. 启动后端服务
4. 部署前端静态文件

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目。

### 开发规范
- 遵循代码规范
- 编写单元测试
- 更新相关文档
- 提交清晰的 commit 信息

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

如果你在使用过程中遇到问题，可以：

1. 查看项目文档
2. 提交 Issue
3. 查看代码注释

## 🎯 后续计划

- [ ] 增加更多权限控制功能
- [ ] 支持多语言国际化
- [ ] 添加操作日志功能
- [ ] 支持文件上传功能
- [ ] 增加数据导入导出
- [ ] 优化性能和用户体验

---

**⭐ 如果这个项目对你有帮助，请给个 Star！**