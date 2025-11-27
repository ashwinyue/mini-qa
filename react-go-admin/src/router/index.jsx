/**
 * 路由配置组件
 * 
 * 使用 React Router v6 的声明式路由配置
 * - Routes: 路由容器，包含所有的 Route 组件
 * - Route: 单个路由配置，定义 URL 路径和对应的组件
 * - Navigate: 导航组件，用于重定向到指定路径
 * - PrivateRoute: 自定义组件，用于保护需要登录的路由
 * 
 * 路由结构说明：
 * /login - 登录页面（公开访问）
 * / - 主应用（需要登录，包含侧边栏和顶部导航）
 *   ├── /dashboard - 仪表盘（默认重定向到这里）
 *   ├── /users - 用户管理页面
 *   ├── /roles - 角色管理页面
 *   └── * - 404 页面（匹配所有未定义的路径）
 */

// 导入 React Router 的核心组件
import { Routes, Route, Navigate } from 'react-router-dom'
// 导入私有路由组件，用于保护需要登录才能访问的页面
import PrivateRoute from '../components/PrivateRoute'
// 导入主布局组件，包含侧边栏和顶部导航
import MainLayout from '../layouts/MainLayout'
// 导入各个页面组件
import Login from '../pages/Login'
import Dashboard from '../pages/Dashboard'
import UserList from '../pages/user/UserList'
import RoleList from '../pages/role/RoleList'
import Settings from '../pages/Settings'
import Profile from '../pages/Profile'
import NotFound from '../pages/NotFound'

const AppRouter = () => {
    return (
        <Routes>
            {/* 
                登录页面路由
                - path="/login": 匹配 URL 中的 /login 路径
                - element={<Login />}: 渲染 Login 组件
                - 这是一个公开路由，无需登录即可访问
            */}
            <Route path="/login" element={<Login />} />

            {/* 
                主应用路由
                - path="/": 匹配根路径
                - PrivateRoute 包装 MainLayout，确保只有登录用户才能访问
                - 登录用户会看到包含侧边栏和顶部导航的完整界面
            */}
            <Route
                path="/"
                element={
                    // PrivateRoute 是一个自定义组件，用于检查用户是否已登录
                    // 如果未登录，会自动重定向到登录页面
                    <PrivateRoute>
                        {/* MainLayout 是主布局组件，包含侧边栏和顶部导航 */}
                        <MainLayout />
                    </PrivateRoute>
                }
            >
                {/* 
                    默认重定向到仪表盘
                    - index: 表示这是父路由的默认子路由
                    - Navigate to="/dashboard" replace: 重定向到 /dashboard 路径
                    - replace: 使用 replace 模式，不会在浏览器历史中留下记录
                */}
                <Route index element={<Navigate to="/dashboard" replace />} />

                {/* 
                    仪表盘页面
                    - path="dashboard": 匹配 /dashboard 路径
                    - element={<Dashboard />}: 渲染 Dashboard 组件
                    - 由于 Dashboard 是 MainLayout 的子路由，会在 MainLayout 的内容区域显示
                */}
                <Route path="dashboard" element={<Dashboard />} />

                {/* 
                    用户管理页面
                    - path="users": 匹配 /users 路径
                    - element={<UserList />}: 渲染 UserList 组件
                    - 包含用户列表、添加用户、编辑用户、删除用户等功能
                */}
                {/* 
                    系统管理模块
                    - /system/user: 用户管理
                    - /system/role: 角色管理
                */}
                <Route path="system">
                    <Route path="user" element={<UserList />} />
                    <Route path="role" element={<RoleList />} />
                </Route>

                {/* 
                    系统设置页面
                    - path="settings": 匹配 /settings 路径
                    - element={<Settings />}: 渲染 Settings 组件
                    - 包含基本设置、安全设置、邮件设置、通知设置、外观设置等
                */}
                <Route path="settings" element={<Settings />} />

                {/* 
                    个人信息页面
                    - path="profile": 匹配 /profile 路径
                    - element={<Profile />}: 渲染 Profile 组件
                */}
                <Route path="profile" element={<Profile />} />

            </Route>

            {/* 
                404 页面路由
                - path="*": 匹配所有未定义的路径（通配符）
                - element={<NotFound />}: 渲染 NotFound 组件
                - 当用户访问不存在的页面时会显示这个组件
            */}
            <Route path="*" element={<NotFound />} />
        </Routes>
    )
}

export default AppRouter
