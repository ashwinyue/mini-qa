/**
 * 认证上下文模块
 * 
 * 该模块使用 React Context API 和 useState/useEffect Hooks 实现全局用户认证状态管理
 * 
 * 主要功能：
 * 1. 管理用户登录状态
 * 2. 提供登录/登出方法
 * 3. 持久化用户认证信息（使用 localStorage）
 * 4. 应用启动时自动检查登录状态
 * 
 * Context 模式的优势：
 * - 避免了 prop drilling（逐层传递 props）
 * - 任何嵌套组件都可以直接访问认证状态
 * - 状态变化会自动触发相关组件的重新渲染
 */

// 从 React 导入必要的 Hook 和 Context API
// createContext: 创建 Context 对象的工厂函数
// useContext: 在组件中访问 Context 对象的 Hook
// useState: 管理组件本地状态的 Hook
// useEffect: 处理副作用的 Hook（如数据获取、DOM 操作、定时器等）
import { createContext, useContext, useEffect } from 'react'
// 导入获取当前用户信息的 API 函数
import { getCurrentUser } from '../api/auth'
import { useAuthStore } from '../stores'


/**
 * 创建认证 Context
 * 
 * createContext() 返回一个 Context 对象，用于在组件树中共享认证状态
 * 初始值为 null，后续通过 AuthProvider 提供实际值
 * 
 * Context 对象包含：
 * - Provider 组件：用于提供 Context 值
 * - Consumer 组件：用于消费 Context 值（新版本推荐使用 useContext Hook）
 */
const AuthContext = createContext(null)

/**
 * 认证上下文提供器组件
 * 
 * 该组件：
 * 1. 包装需要访问认证状态的子组件
 * 2. 提供认证相关的状态和方法
 * 3. 管理登录状态的持久化
 * 
 * @param {Object} props 组件属性
 * @param {React.ReactNode} props.children 子组件
 */
export const AuthProvider = ({ children }) => {
  /**
   * Use Zustand store hooks to access state and actions
   * This replaces the local useState and ensures the store is the single source of truth
   */
  const user = useAuthStore((state) => state.user)
  const loading = useAuthStore((state) => state.loading)
  const login = useAuthStore((state) => state.login)
  const logout = useAuthStore((state) => state.logout)
  const setUser = useAuthStore((state) => state.setUser)
  const setLoading = useAuthStore((state) => state.setLoading)

  /**
   * 组件挂载时检查登录状态
   * 
   * useEffect Hook 用于处理副作用操作：
   * - 组件首次渲染时自动执行
   * - 依赖数组为空，表示只在挂载时执行一次
   * 
   * 这里的作用是检查 localStorage 中是否有 token，
   * 如果有，则获取用户信息并更新状态
   */
  useEffect(() => {
    // 从 localStorage 中获取 token
    const token = localStorage.getItem('token')
    if (token) {
      // 如果存在 token，说明用户之前已经登录过
      // 尝试获取当前用户信息
      fetchUserInfo()
    } else {
      // 没有 token，直接结束加载状态
      setLoading(false)
    }
  }, []) // 空依赖数组表示只执行一次

  /**
   * 获取用户信息的异步函数
   * 
   * 使用 async/await 语法处理异步操作：
   * - try: 尝试执行用户信息获取逻辑
   * - catch: 处理可能出现的错误
   * - finally: 无论成功或失败都会执行，用于清理工作
   */
  const fetchUserInfo = async () => {
    try {
      // 调用 API 获取用户信息
      const userData = await getCurrentUser()
      // 设置用户信息状态
      setUser(userData)
    } catch (error) {
      // 获取用户信息失败，可能 token 已过期
      console.error('获取用户信息失败:', error)
      // 清除过期的 token
      localStorage.removeItem('token')
      logout() // Use store logout action which clears state
    } finally {
      // 无论成功或失败，都要结束加载状态
      setLoading(false)
    }
  }

  /**
   * 构建 Context 值对象
   * 
   * 该对象将传递给 Provider 的 value 属性，
   * 所有被 Provider 包装的子组件都可以通过 useAuth() 访问这些值和方法
   */
  const value = {
    // user: 当前用户信息（为 null 时表示未登录）
    user,
    // loading: 加载状态（true 时显示加载界面）
    loading,
    // login: 登录方法
    login,
    // logout: 登出方法
    logout,
    // isAuthenticated: 便捷的属性，直接返回布尔值表示是否已登录
    isAuthenticated: !!user, // !!user 将 user 转换为布尔值，null 为 false，非 null 为 true
  }

  /**
   * 返回 Context Provider 组件
   * 
   * Provider 组件：
   * - 接收一个 value 属性，包含要共享的状态和方法
   * - 包装其 children 组件，使它们可以访问 Context 值
   * - 当 value 发生变化时，所有使用该 Context 的组件都会重新渲染
   */
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

/**
 * 自定义 Hook：useAuth
 * 
 * 该 Hook：
 * 1. 封装了 useContext(AuthContext) 的调用
 * 2. 提供错误处理，防止在不正确的 Context 中使用
 * 3. 简化了组件中使用认证状态的方式
 * 
 * 使用方式：
 * const { user, loading, login, logout, isAuthenticated } = useAuth()
 * 
 * @returns {Object} 认证上下文的对象，包含状态和方法
 * @throws {Error} 如果在 AuthProvider 外部使用，抛出错误
 */
export const useAuth = () => {
  // 使用 useContext Hook 获取 Context 值
  const context = useContext(AuthContext)

  // 如果在 AuthProvider 外部使用，context 为 null，此时抛出错误
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }

  // 返回 Context 的值
  return context
}
