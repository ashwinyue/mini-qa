# Kimi Style Frontend

基于 Kimi AI 风格的现代化前端应用，采用 React + Vite + Tailwind CSS 构建。

## 功能特性

- ✅ AI 对话界面
- ✅ 会话管理（创建、切换、删除）
- ✅ 系统设置（主题、语言）
- ✅ 响应式设计（支持移动端）
- ✅ Markdown 消息渲染
- ✅ 建议问题流式推送
- ✅ 图片上传支持

## 技术栈

- **框架**: React 18 + JSX
- **构建工具**: Vite
- **样式**: Tailwind CSS 4.x
- **状态管理**: Zustand
- **路由**: React Router v6
- **HTTP 客户端**: Axios
- **图标**: Lucide React
- **Markdown**: react-markdown

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

应用将在 http://localhost:5199 启动

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
kimi-frontend/
├── public/              # 静态资源
├── src/
│   ├── components/      # 组件
│   │   ├── layout/      # 布局组件
│   │   ├── chat/        # 聊天组件
│   │   ├── settings/    # 设置组件
│   │   └── common/      # 通用组件
│   ├── pages/           # 页面
│   ├── stores/          # Zustand 状态管理
│   ├── services/        # API 服务
│   ├── utils/           # 工具函数
│   ├── App.jsx          # 主应用组件
│   ├── main.jsx         # 入口文件
│   ├── index.css        # 全局样式
│   └── App.css          # 应用样式
├── index.html           # HTML 模板
├── vite.config.js       # Vite 配置
├── tailwind.config.js   # Tailwind 配置
└── package.json         # 项目配置
```

## API 配置

后端 API 默认代理到 `http://localhost:8000`，可在 `vite.config.js` 中修改：

```javascript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```

## 主要组件

### 布局组件
- `MainLayout`: 主布局，包含侧边栏和内容区
- `Sidebar`: 侧边栏，包含会话列表和用户菜单

### 聊天组件
- `ChatInterface`: 聊天界面主组件
- `ChatInput`: 消息输入框
- `MessageList`: 消息列表
- `MessageItem`: 单条消息
- `SessionList`: 会话列表
- `WelcomeScreen`: 欢迎屏幕

### 设置组件
- `SettingsPanel`: 设置面板
- `UserMenu`: 用户菜单

## 状态管理

### useChatStore
管理聊天会话和消息：
- `currentSession`: 当前会话
- `sessions`: 会话列表
- `createSession()`: 创建新会话
- `addMessage()`: 添加消息
- `switchSession()`: 切换会话

### useSettingsStore
管理用户设置：
- `theme`: 主题（light/dark/system）
- `language`: 语言（zh/en）
- `setTheme()`: 设置主题
- `setLanguage()`: 设置语言

## 样式定制

主要颜色变量定义在 `src/index.css`：

```css
:root {
  --color-sidebar-dark: #1a1a1a;
  --color-content-light: #ffffff;
  --color-primary: #4F46E5;
  --color-primary-hover: #4338CA;
}
```

## 浏览器支持

- Chrome (最新版)
- Firefox (最新版)
- Safari (最新版)
- Edge (最新版)

## License

MIT
