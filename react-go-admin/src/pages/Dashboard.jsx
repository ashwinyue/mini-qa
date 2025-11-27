/**
 * 仪表盘页面组件
 * 
 * 这是应用的主页，显示系统的关键指标和概览信息
 * 仪表盘通常包含：
 * - 关键数据卡片：显示重要指标的当前状态
 * - 图表和图形：展示数据趋势和分布
 * - 快捷操作：提供常用功能的快速入口
 * - 最近活动：显示系统的最新动态
 * 
 * 这个组件使用了 Ant Design 的统计卡片和图表组件
 * 来创建一个美观且功能丰富的管理仪表板
 */

// 导入 React 相关的 Hook
import React, { useState, useEffect } from 'react'

// 导入 Ant Design 组件
import {
    Card,           // 卡片组件：用于显示统计信息
    Row,            // 行布局组件：用于网格布局
    Col,            // 列布局组件：用于网格布局
    Statistic,      // 统计组件：显示数字和趋势
    Progress,       // 进度条组件：显示进度信息
    Space,          // 间距组件：控制元素间距
    Typography,     // 文本组件：显示标题和文本
    Divider,        // 分割线组件：分隔内容
    Table,          // 表格组件：显示数据表格
    Tag,            // 标签组件：显示状态标签
    Avatar,         // 头像组件：显示用户头像
    List,           // 列表组件：显示列表数据
    Calendar,       // 日历组件：显示日历视图
    Badge,          // 徽标组件：显示状态标记
    Tooltip,        // 提示组件：显示悬停提示
    Button          // 按钮组件：提供操作按钮
} from 'antd'

// 导入 Ant Design 图标
import {
    UserOutlined,       // 用户图标
    TeamOutlined,       // 团队图标
    KeyOutlined,        // 权限图标
    EyeOutlined,        // 访问图标
    ArrowUpOutlined,    // 上升箭头图标
    ArrowDownOutlined,  // 下降箭头图标
    CalendarOutlined,   // 日历图标
    BarChartOutlined,   // 图表图标
    LineChartOutlined,  // 折线图图标
    PieChartOutlined,   // 饼图图标
    DotChartOutlined    // 点图图标
} from '@ant-design/icons'

// 导入图表组件（如果项目中安装了图表库）
// import { LineChart, BarChart, PieChart } from 'recharts'

// 解构文本组件
const { Title, Text } = Typography

/**
 * 仪表盘页面组件
 * 
 * 这是应用的主页，显示系统的关键指标和概览信息
 * 
 * 功能特点：
 * 1. 显示关键业务指标（用户数、角色数、权限数等）
 * 2. 展示数据趋势和变化
 * 3. 提供快捷操作入口
 * 4. 显示最近活动信息
 * 
 * @component
 */
const Dashboard = () => {
    /**
     * 状态管理
     * 
     * - loading: 页面加载状态
     * - statistics: 统计数据
     * - recentActivities: 最近活动数据
     * - quickActions: 快捷操作列表
     */
    const [loading, setLoading] = useState(false)
    const [statistics, setStatistics] = useState({
        totalUsers: 1248,
        totalRoles: 12,
        totalPermissions: 156,
        todayVisits: 3421,
        userGrowth: '+12.5%',
        visitTrend: '+8.2%',
        roleGrowth: '+2',
        permissionGrowth: '+15'
    })
    const [recentActivities, setRecentActivities] = useState([
        {
            id: 1,
            user: '张三',
            action: '登录系统',
            time: '2分钟前',
            type: 'login',
            avatar: 'https://example.com/avatar1.jpg'
        },
        {
            id: 2,
            user: '李四',
            action: '创建了用户',
            time: '5分钟前',
            type: 'create',
            avatar: 'https://example.com/avatar2.jpg'
        },
        {
            id: 3,
            user: '王五',
            action: '更新了权限',
            time: '10分钟前',
            type: 'update',
            avatar: 'https://example.com/avatar3.jpg'
        },
        {
            id: 4,
            user: '赵六',
            action: '删除了用户',
            time: '15分钟前',
            type: 'delete',
            avatar: 'https://example.com/avatar4.jpg'
        },
        {
            id: 5,
            user: '钱七',
            action: '修改了密码',
            time: '20分钟前',
            type: 'password',
            avatar: 'https://example.com/avatar5.jpg'
        }
    ])

    /**
     * 快捷操作配置
     * 
     * 定义了仪表盘上的快捷操作按钮
     * 每个操作包含图标、标题、描述和点击处理函数
     */
    const quickActions = [
        {
            key: 'add-user',
            title: '添加用户',
            description: '创建新用户账户',
            icon: <UserOutlined />,
            color: '#1677ff',
            onClick: () => {
                // 导航到用户创建页面
                console.log('跳转到添加用户页面')
            }
        },
        {
            key: 'manage-roles',
            title: '管理角色',
            description: '配置用户角色和权限',
            icon: <KeyOutlined />,
            color: '#52c41a',
            onClick: () => {
                // 导航到角色管理页面
                console.log('跳转到角色管理页面')
            }
        },
        {
            key: 'view-logs',
            title: '查看日志',
            description: '查看系统操作日志',
            icon: <BarChartOutlined />,
            color: '#faad14',
            onClick: () => {
                // 导航到日志页面
                console.log('跳转到日志页面')
            }
        },
        {
            key: 'system-settings',
            title: '系统设置',
            description: '配置系统参数',
            icon: <LineChartOutlined />,
            color: '#722ed1',
            onClick: () => {
                // 导航到系统设置页面
                console.log('跳转到系统设置页面')
            }
        }
    ]

    /**
     * 获取活动类型对应的颜色和标签
     * 
     * @param {string} type - 活动类型
     * @returns {Object} 包含颜色和标签的对象
     */
    const getActivityConfig = (type) => {
        const configs = {
            login: { color: 'blue', text: '登录' },
            create: { color: 'green', text: '创建' },
            update: { color: 'orange', text: '更新' },
            delete: { color: 'red', text: '删除' },
            password: { color: 'purple', text: '密码' }
        }
        return configs[type] || { color: 'default', text: '操作' }
    }

    /**
     * 模拟数据加载
     * 
     * 在实际应用中，这里会调用 API 获取真实数据
     * 这里使用模拟数据来展示组件功能
     */
    useEffect(() => {
        // 模拟加载状态
        setLoading(true)

        // 模拟异步数据获取
        const timer = setTimeout(() => {
            setLoading(false)
            console.log('仪表盘数据加载完成')
        }, 1000)

        // 清理定时器
        return () => clearTimeout(timer)
    }, [])

    /**
     * 渲染统计卡片
     * 
     * 显示关键业务指标的卡片组件
     * 包含数值、趋势和图标
     */
    const renderStatisticCard = (title, value, prefix, suffix, trend, color) => (
        <Card hoverable className="h-full">
            <Statistic
                title={title}
                value={value}
                prefix={prefix}
                suffix={suffix}
                valueStyle={{ color: color }}
            />
            {trend && (
                <div className="mt-2">
                    <Text type={trend.startsWith('+') ? 'success' : 'danger'}>
                        {trend.startsWith('+') ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                        {trend}
                    </Text>
                    <Text className="ml-2 text-xs">较昨日</Text>
                </div>
            )}
        </Card>
    )

    /**
     * 渲染快捷操作卡片
     * 
     * 显示常用功能的快速入口
     * 每个操作都是可点击的卡片
     */
    const renderQuickActionCard = (action) => (
        <Card
            key={action.key}
            hoverable
            onClick={action.onClick}
            className="h-full text-center cursor-pointer"
        >
            <div className="text-4xl mb-4" style={{ color: action.color }}>
                {action.icon}
            </div>
            <Title level={4} style={{ marginBottom: '8px' }}>
                {action.title}
            </Title>
            <Text type="secondary" className="text-xs">
                {action.description}
            </Text>
        </Card>
    )

    /**
     * 渲染最近活动项
     * 
     * 显示系统最近的操作记录
     * 包含用户头像、操作描述和时间
     */
    const renderActivityItem = (activity) => {
        const config = getActivityConfig(activity.type)

        return (
            <List.Item key={activity.id}>
                <List.Item.Meta
                    avatar={<Avatar src={activity.avatar} icon={<UserOutlined />} />}
                    title={
                        <Space>
                            <Text>{activity.user}</Text>
                            <Tag color={config.color}>{config.text}</Tag>
                        </Space>
                    }
                    description={
                        <Space>
                            <Text>{activity.action}</Text>
                            <Text type="secondary" className="text-xs">
                                {activity.time}
                            </Text>
                        </Space>
                    }
                />
            </List.Item>
        )
    }

    /**
     * 主渲染函数
     * 
     * 使用 Ant Design 的网格系统布局页面
     * 包含统计卡片、快捷操作、最近活动等模块
     */
    return (
        <div>
            {/* 页面标题 */}
            <div className="mb-6">
                <Title level={2}>
                    仪表盘
                </Title>
                <Text type="secondary">
                    欢迎使用 React 管理系统，以下是系统的概览信息
                </Text>
            </div>

            {/* 统计卡片区域 */}
            <Row gutter={[16, 16]} className="mb-6">
                <Col xs={24} sm={12} md={6}>
                    {renderStatisticCard(
                        '用户总数',
                        statistics.totalUsers,
                        <UserOutlined />,
                        '人',
                        statistics.userGrowth,
                        '#1677ff'
                    )}
                </Col>
                <Col xs={24} sm={12} md={6}>
                    {renderStatisticCard(
                        '角色数量',
                        statistics.totalRoles,
                        <TeamOutlined />,
                        '个',
                        statistics.roleGrowth,
                        '#52c41a'
                    )}
                </Col>
                <Col xs={24} sm={12} md={6}>
                    {renderStatisticCard(
                        '权限数量',
                        statistics.totalPermissions,
                        <KeyOutlined />,
                        '项',
                        statistics.permissionGrowth,
                        '#faad14'
                    )}
                </Col>
                <Col xs={24} sm={12} md={6}>
                    {renderStatisticCard(
                        '今日访问量',
                        statistics.todayVisits,
                        <EyeOutlined />,
                        '次',
                        statistics.visitTrend,
                        '#722ed1'
                    )}
                </Col>
            </Row>

            {/* 快捷操作区域 */}
            <Row gutter={[16, 16]} className="mb-6">
                <Col span={24}>
                    <Card title="快捷操作">
                        <Row gutter={[16, 16]}>
                            {quickActions.map(action => (
                                <Col key={action.key} xs={24} sm={12} md={6}>
                                    {renderQuickActionCard(action)}
                                </Col>
                            ))}
                        </Row>
                    </Card>
                </Col>
            </Row>

            {/* 最近活动和系统信息 */}
            <Row gutter={[16, 16]}>
                {/* 最近活动 */}
                <Col xs={24} md={12}>
                    <Card
                        title="最近活动"
                        extra={<Button type="link">查看全部</Button>}
                    >
                        <List
                            dataSource={recentActivities}
                            renderItem={renderActivityItem}
                            loading={loading}
                        />
                    </Card>
                </Col>

                {/* 系统信息 */}
                <Col xs={24} md={12}>
                    <Card title="系统信息">
                        <div className="mb-4">
                            <Text strong>系统状态：</Text>
                            <Badge status="success" text="正常运行" />
                        </div>
                        <div className="mb-4">
                            <Text strong>系统版本：</Text>
                            <Text>v1.0.0</Text>
                        </div>
                        <div className="mb-4">
                            <Text strong>最后更新：</Text>
                            <Text>2024-01-21 14:30:00</Text>
                        </div>
                        <div className="mb-4">
                            <Text strong>服务器时间：</Text>
                            <Text>{new Date().toLocaleString()}</Text>
                        </div>
                        <Divider />
                        <div>
                            <Text strong>今日日历</Text>
                            <div className="mt-4">
                                {/* 这里可以集成完整的日历组件 */}
                                <Calendar
                                    fullscreen={false}
                                    headerRender={({ value, type, onChange, onTypeChange }) => (
                                        <div className="p-2">
                                            <Text strong>{value.format('YYYY年MM月')}</Text>
                                        </div>
                                    )}
                                />
                            </div>
                        </div>
                    </Card>
                </Col>
            </Row>
        </div>
    )
}

export default Dashboard
