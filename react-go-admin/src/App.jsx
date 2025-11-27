// 导入 React Router 的核心组件
// BrowserRouter: 提供路由功能的根组件，使用 HTML5 的 History API
import { BrowserRouter } from 'react-router-dom'
// 导入认证上下文提供器，用于全局管理用户认证状态
import { AuthProvider } from './context/AuthContext'
// 导入路由配置组件，管理所有页面路由
import AppRouter from './router'

/**
 * 应用根组件
 * 
 * 这是整个应用的最顶层组件，负责：
 * 1. 设置路由（BrowserRouter）
 * 2. 提供认证上下文（AuthProvider）
 * 3. 渲染路由配置（AppRouter）
 * 
 * 组件层次结构：
 * <BrowserRouter> // 提供路由功能
 *   <AuthProvider> // 提供全局认证状态
 *     <AppRouter /> // 管理所有页面路由
 *   </AuthProvider>
 * </BrowserRouter>
 */
function App() {
  return (
    // BrowserRouter 包装整个应用，启用客户端路由功能
    // 客户端路由的优势：
    // - 页面切换更快（无需重新加载整个页面）
    // - 支持浏览器前进/后退按钮
    // - 提供更好的用户体验
    <BrowserRouter>
      {/* 
        AuthProvider 提供全局的用户认证状态管理
        所有需要访问用户信息的子组件都可以通过 useAuth Hook 获取认证状态
        认证状态包括：
        - user: 当前用户信息
        - loading: 加载状态
        - login: 登录方法
        - logout: 登出方法
        - isAuthenticated: 是否已认证
      */}
      <AuthProvider>
        {/* 
          AppRouter 管理应用的所有路由配置
          包括：
          - 公开路由（登录页面）
          - 受保护的路由（需要登录才能访问）
          - 嵌套路由（用户管理、角色管理等）
        */}
        <AppRouter />
      </AuthProvider>
    </BrowserRouter>
  )
}

// 导出默认组件，供 main.jsx 引用
export default App
