# 部署说明文档

## 目标与范围
- 说明如何在开发/测试/生产环境部署React用户管理系统，完成前后端依赖安装、配置、初始化、启动与验证。
- 适用组件：React前端、Go后端、SQLite数据库、JWT认证系统。

## 环境要求

### 前端环境
- **Node.js**: ≥ 18.0.0
- **npm** 或 **yarn**: 包管理器
- **浏览器**: Chrome, Firefox, Safari, Edge (最新版本)

### 后端环境  
- **Go**: ≥ 1.21
- **Git**: 版本控制工具
- **SQLite**: 数据库（无需额外安装）

### 端口占用
- 前端开发服务器：`5173`
- 后端API服务器：`8000`
- 数据库文件：本地文件存储

## 依赖安装

### 前端依赖安装
```bash
# 进入项目目录
cd /Users/solariswu/WebstormProjects/react-go-admin

# 安装前端依赖
npm install

# 或者使用 yarn
yarn install
```

### 后端依赖安装
```bash
# 进入后端目录
cd backend

# 下载Go依赖
go mod download

# 整理依赖
go mod tidy
```

## 配置项

### 前端配置
- **Vite配置**: `vite.config.js`
- **代理配置**: 开发环境代理到后端API
- **环境变量**: `.env` 文件（可选）

### 后端配置
后端配置通过代码硬编码，主要配置项：
- **服务器端口**: `:8000`（可在main.go中修改）
- **CORS配置**: 允许`http://localhost:5173`和`http://localhost:5174`
- **数据库路径**: 自动创建SQLite数据库文件
- **JWT密钥**: 硬编码密钥（生产环境应使用环境变量）

## 初始化步骤

### 数据库初始化
后端服务启动时会自动创建数据库和表结构，无需手动初始化。

### 创建默认数据（可选）
启动后端服务后，可以通过API创建初始用户和角色数据。

## 启动服务

### 开发模式启动

#### 启动后端服务
```bash
# 进入后端目录
cd backend

# 运行后端服务
go run main.go

# 或者编译后运行
go build -o app
./app
```

#### 启动前端服务
```bash
# 在项目根目录
npm run dev

# 或者使用 yarn
yarn dev
```

### 生产模式启动

#### 构建前端
```bash
# 构建生产版本
npm run build

# 构建文件在dist目录
```

#### 构建后端
```bash
# 进入后端目录
cd backend

# 构建二进制文件
go build -o react-go-admin-backend

# 运行二进制文件
./react-go-admin-backend
```

## 验证方法

### 前端验证
- 访问 `http://localhost:5173`
- 应该能看到登录页面
- 使用默认凭据登录（需要先在数据库中创建用户）

### 后端验证
- 健康检查：`curl http://localhost:8000/health`
- 应该返回服务状态信息

### API验证
```bash
# 测试用户登录
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 测试获取用户列表（需要Token）
curl -X GET http://localhost:8000/api/users \
  -H "Authorization: Bearer <your_token_here>"
```

## 环境差异配置

### 开发环境
- 启用热重载
- 详细的错误信息
- 开发工具扩展
- 代理配置解决跨域

### 测试环境
- 独立的测试数据库
- 测试数据准备
- API集成测试
- 端到端测试

### 生产环境
- 静态文件优化
- 错误页面定制
- 日志级别调整
- 性能监控
- 安全加固

## 故障排查

### 前端常见问题
- **依赖安装失败**: 检查Node.js版本，清除npm缓存
- **端口占用**: 修改vite.config.js中的端口配置
- **构建失败**: 检查代码语法，查看构建错误信息

### 后端常见问题
- **依赖下载失败**: 检查网络连接，设置Go代理
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

- **数据库权限问题**: 检查SQLite文件权限
- **端口占用**: 修改main.go中的端口配置

### 跨域问题
- 确保后端CORS配置正确
- 检查前端代理配置
- 验证请求头设置

## 性能调优

### 前端优化
- 启用代码分割
- 图片压缩和优化
- 使用CDN加速
- 启用Gzip压缩

### 后端优化
- 数据库索引优化
- 查询性能优化
- 内存使用优化
- 并发处理优化

## 安全加固

### 前端安全
- 输入验证和清理
- XSS防护
- CSRF防护
- 敏感信息脱敏

### 后端安全
- JWT密钥管理
- 数据库连接安全
- API限流
- 错误信息脱敏

## 日志与监控

### 前端日志
- 浏览器控制台日志
- 错误边界处理
- 用户行为追踪
- 性能监控

### 后端日志
- 请求日志记录
- 错误日志记录
- 数据库操作日志
- 性能指标收集

## 备份与恢复

### 数据库备份
```bash
# SQLite数据库文件路径（根据实际路径调整）
cp backend/users.db backend/users.db.backup
```

### 配置备份
- 备份重要配置文件
- 版本控制配置变更
- 文档化配置项

## 升级与维护

### 前端升级
- 依赖包更新
- 构建工具升级
- 框架版本升级
- 兼容性测试

### 后端升级
- Go版本升级
- 依赖包更新
- 数据库迁移
- API版本管理

## 关联文档
- 架构设计：参见《[architecture_design.md](./architecture_design.md)》
- API规范：参见《[api_specification.md](./api_specification.md)》