/**
 * 私有路由组件
 * 
 * 这是一个路由保护组件，用于确保只有已认证的用户才能访问受保护的页面
 * 
 * 工作原理：
 * 1. 检查用户是否已登录（从 AuthContext 获取认证状态）
 * 2. 如果已登录，渲染被保护的子组件
 * 3. 如果未登录，重定向到登录页面
 * 4. 如果正在加载认证状态，显示加载界面
 * 
 * 这种组件模式被称为"路由守卫"或"页面保护"，
 * 在需要身份验证的应用中非常常见
 */

// 从 React Router 导入 Navigate 组件，用于实现重定向
// Navigate 组件会改变当前路由，将用户导航到指定路径
import { Navigate } from 'react-router-dom'
// 导入 useAuth Hook，用于获取认证状态和方法
import { useAuth } from '../context/AuthContext'
// 导入 Ant Design 的加载指示器组件
import { Spin } from 'antd'

/**
 * 私有路由组件
 * 
 * @param {Object} props 组件属性
 * @param {React.ReactNode} props.children 受保护的子组件
 * 
 * 该组件接收一个 children 属性，这是一个 React 元素（组件）
 * 只有当用户通过认证时，这些子组件才会被渲染
 */
const PrivateRoute = ({ children }) => {
    /**
     * 从 AuthContext 获取认证相关状态
     * 
     * - isAuthenticated: 布尔值，表示用户是否已登录
     * - loading: 布尔值，表示认证状态是否正在检查中
     */
    const { isAuthenticated, loading } = useAuth()

    /**
     * 加载状态处理
     * 
     * 当正在检查用户认证状态时（通常是应用启动时的首次检查），
     * 显示一个加载指示器，避免在认证状态确定前就跳转或渲染内容
     */
    if (loading) {
        return (
            <div className="flex justify-center items-center h-screen">
                {/* Spin 组件：Ant Design 的加载指示器 */}
                <Spin size="large" />
                {/* size="large" 设置为大型加载指示器 */}
            </div>
        )
    }

    /**
     * 认证状态判断和渲染逻辑
     * 
     * 使用三元运算符简化条件渲染：
     * - 如果 isAuthenticated 为 true，渲染子组件
     * - 如果 isAuthenticated 为 false，重定向到登录页面
     */
    return isAuthenticated ?
        // 已认证：渲染受保护的子组件
        children :
        // 未认证：重定向到登录页面
        <Navigate to="/login" replace />
}

export default PrivateRoute
