/**
 * 主布局组件
 * 
 * 这是应用的主要布局结构，包含：
 * - 顶部导航栏：显示用户信息和操作按钮
 * - 侧边栏：应用的主要导航菜单
 * - 内容区域：渲染当前路由对应的页面内容
 * 
 * 布局特点：
 * 1. 使用 Ant Design 的 Layout 组件作为基础
 * 2. 响应式设计，支持移动端适配
 * 3. 侧边栏支持折叠/展开功能
 * 4. 集成了用户认证相关的操作（用户信息、登出等）
 */

// 导入 React 相关的 Hook
import React, { useState } from 'react'
// 导入 React Router 的 Outlet 组件，用于渲染子路由内容
import { Outlet, useLocation } from 'react-router-dom'
// 导入 Ant Design 组件
import {
    Layout,      // 布局容器组件
    Menu,        // 菜单组件
    Button,      // 按钮组件
    Avatar,      // 头像组件
    Dropdown,    // 下拉菜单组件
    Typography,  // 文本组件
    Space        // 间距组件
} from 'antd'
// 导入 Ant Design 图标
import {
    MenuFoldOutlined,    // 折叠菜单图标
    MenuUnfoldOutlined,  // 展开菜单图标
    UserOutlined,        // 用户图标
    LogoutOutlined,      // 登出图标
    DashboardOutlined,   // 仪表盘图标
    TeamOutlined,        // 团队图标
    SettingOutlined      // 设置图标
} from '@ant-design/icons'

// 导入认证相关的 Hook 和工具函数
import { useAuth } from '../context/AuthContext'
import { useNavigate } from 'react-router-dom'

// 解构 Layout 组件的子组件
const { Header, Sider, Content } = Layout
// 获取文本组件
const { Text } = Typography

/**
 * 主布局组件
 * 
 * 这是应用的根布局组件，包裹了应用的所有页面
 * 它负责渲染公共的导航结构和页面容器
 */
const MainLayout = () => {
    /**
     * 状态管理
     * 
     * - collapsed: 控制侧边栏的折叠/展开状态
     * - selectedKeys: 当前选中的菜单项（高亮显示）
     * - openKeys: 当前展开的子菜单项（支持多级菜单）
     */
    const [collapsed, setCollapsed] = useState(false)
    const [selectedKeys, setSelectedKeys] = useState(['dashboard'])
    const [openKeys, setOpenKeys] = useState([])

    /**
     * 获取认证相关信息
     * 
     * - user: 当前登录的用户信息
     * - logout: 登出方法
     */
    const { user, logout } = useAuth()
    // 获取当前路由信息，用于同步菜单状态
    const location = useLocation()
    const navigate = useNavigate()

    /**
     * 根据当前路径更新菜单选中状态
     * 
     * 当路由改变时，自动更新侧边栏中对应的菜单项为选中状态
     */
    React.useEffect(() => {
        const path = location.pathname
        // 根据不同路径设置对应的菜单项
        if (path.startsWith('/dashboard')) {
            setSelectedKeys(['dashboard'])
            setOpenKeys([])
        } else if (path.startsWith('/system/user')) {
            setSelectedKeys(['user-list'])
            setOpenKeys(['system'])
        } else if (path.startsWith('/system/role')) {
            setSelectedKeys(['role-list'])
            setOpenKeys(['system'])
        } else if (path.startsWith('/settings')) {
            setSelectedKeys(['settings'])
            setOpenKeys([])
        }
    }, [location])

    /**
     * 侧边栏菜单项配置
     * 
     * 定义了应用的主要导航菜单，包括：
     * - 图标
     * - 菜单文本
     * - 对应的路由路径
     * - 可选的子菜单项
     */
    const menuItems = [
        {
            key: 'dashboard',              // 菜单项的唯一标识
            icon: <DashboardOutlined />,   // 菜单项图标
            label: '仪表盘',                // 菜单项文本
            onClick: () => navigate('/dashboard')  // 点击时的导航操作
        },
        {
            key: 'system',
            icon: <TeamOutlined />,
            label: '系统管理',
            children: [
                {
                    key: 'user-list',
                    label: '用户管理',
                    onClick: () => navigate('/system/user')
                },
                {
                    key: 'role-list',
                    label: '角色管理',
                    onClick: () => navigate('/system/role')
                }
            ]
        },
        {
            key: 'settings',
            icon: <SettingOutlined />,
            label: '系统设置',
            onClick: () => navigate('/settings')
        }
    ]

    /**
     * 用户下拉菜单配置
     * 
     * 定义了用户头像下拉菜单中的选项：
     * - 查看个人信息
     * - 登出
     */
    const userMenuItems = [
        {
            key: 'profile',
            icon: <UserOutlined />,
            label: '个人信息',
            onClick: () => navigate('/profile')
        },
        {
            key: 'logout',
            icon: <LogoutOutlined />,
            label: '退出登录',
            onClick: logout  // 直接调用登出方法
        }
    ]

    /**
     * 渲染主布局
     * 
     * 使用 Ant Design 的 Layout 组件构建整体结构：
     * - Sider：侧边栏（可折叠）
     * - Header：顶部导航栏
     * - Content：内容区域
     * - Outlet：渲染当前路由对应的页面组件
     */
    return (
        <Layout className="min-h-screen">
            {/* 侧边栏 */}
            <Sider
                trigger={null}              // 不显示默认的折叠按钮
                collapsible                // 允许折叠
                collapsed={collapsed}       // 当前折叠状态
                theme="dark"               // 使用深色主题
                width={240}                // 展开时的宽度
                collapsedWidth={80}        // 折叠时的宽度
            >
                {/* Logo/标题区域 - 隐藏，因为标题已移到顶部 */}
                <div className={`h-16 px-4 flex items-center ${collapsed ? 'justify-center' : 'justify-start'}`} style={{ background: '#001529' }}>
                    {/* 根据折叠状态显示不同内容 */}
                    {collapsed ? (
                        <div className="w-8 h-8 rounded-md flex items-center justify-center text-white font-bold" style={{ background: '#1677ff' }}>
                            R
                        </div>
                    ) : (
                        <div style={{ height: '64px' }}></div>
                    )}
                </div>

                {/* 菜单组件 */}
                <Menu
                    theme="dark"                          // 使用深色主题
                    mode="inline"                         // 内联模式（垂直菜单）
                    selectedKeys={selectedKeys}            // 当前选中的菜单项
                    defaultOpenKeys={openKeys}            // 默认展开的菜单项
                    items={menuItems}                      // 菜单项数据
                    className="border-none"             // 移除边框样式
                />
            </Sider>

            {/* 主要内容区域 */}
            <Layout>
                {/* 顶部导航栏 */}
                <Header 
                    className="flex items-center justify-between" 
                    style={{ 
                        height: '64px', 
                        lineHeight: '64px', 
                        padding: '0 24px',
                        background: '#ffffff',
                        boxShadow: '0 1px 4px rgba(0,21,41,0.08)'
                    }}
                >
                    {/* 左侧：折叠按钮和系统名称 */}
                    <div className="flex items-center">
                        <Button
                            type="text"
                            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                            onClick={() => setCollapsed(!collapsed)}
                            style={{ fontSize: '18px', width: '48px', height: '48px', color: '#001529' }}
                        />
                        <Text style={{ fontSize: '18px', fontWeight: 'bold', marginLeft: '16px', color: '#001529' }}>
                            React管理系统
                        </Text>
                    </div>

                    {/* 右侧：用户信息和操作 */}
                    <Dropdown
                        menu={{ items: userMenuItems }}
                        placement="bottomRight"
                        arrow
                    >
                        <Button type="text" style={{ height: '48px' }}>
                            <Space>
                                <Avatar size="small" icon={<UserOutlined />} />
                                <Text style={{ color: '#001529' }}>{user?.realname || user?.username || '用户'}</Text>
                            </Space>
                        </Button>
                    </Dropdown>
                </Header>

                {/* 页面内容区域 */}
                <Content className="m-6 p-6 bg-white rounded-lg shadow-sm" style={{ minHeight: 'calc(100vh - 112px)' }}>
                    {/* Outlet 组件：渲染当前路由对应的页面组件 */}
                    {/* 当路由改变时，Outlet 会自动渲染对应的页面内容 */}
                    <Outlet />
                </Content>
            </Layout>
        </Layout>
    )
}

export default MainLayout
