/**
 * 个人信息页面组件
 * 
 * 这是用户查看和编辑个人信息的页面
 * 
 * 主要功能：
 * 1. 查看个人基本信息
 * 2. 编辑个人资料
 * 3. 修改密码
 * 4. 上传头像
 * 5. 查看登录日志
 * 
 * 技术特点：
 * - 使用 Tabs 组件实现多选项卡
 * - 使用 Form 组件进行信息编辑
 * - 使用 Upload 组件上传头像
 * - 实时显示用户信息
 */

// 导入 React 和 Hooks
import { useState, useEffect } from 'react'

// 导入 Ant Design 组件
import {
    Card,           // 卡片组件
    Tabs,           // 选项卡组件
    Form,           // 表单组件
    Input,          // 输入框组件
    Button,         // 按钮组件
    Space,          // 间距组件
    message,        // 消息提示
    Avatar,         // 头像组件
    Upload,         // 上传组件
    Descriptions,   // 描述列表
    Table,          // 表格组件
    Tag,            // 标签组件
    Typography,     // 文本组件
    Divider,        // 分割线
    Row,            // 行布局
    Col,            // 列布局
    Alert           // 警告提示
} from 'antd'

// 导入图标
import {
    UserOutlined,       // 用户图标
    LockOutlined,       // 锁图标
    SaveOutlined,       // 保存图标
    UploadOutlined,     // 上传图标
    HistoryOutlined,    // 历史图标
    CameraOutlined      // 相机图标
} from '@ant-design/icons'

// 导入认证上下文
import { useAuth } from '../context/AuthContext'

// 导入 API
import { updateUser, changePassword } from '../api/user'

const { Title, Text, Paragraph } = Typography

/**
 * 个人信息页面组件
 * 
 * @component
 */
const Profile = () => {
    /**
     * 状态管理
     */
    const [activeTab, setActiveTab] = useState('info')
    const [loading, setLoading] = useState(false)
    const [avatarUrl, setAvatarUrl] = useState('')

    // 获取当前用户信息
    const { user } = useAuth()

    // 创建表单实例
    const [infoForm] = Form.useForm()
    const [passwordForm] = Form.useForm()

    /**
     * 组件挂载时初始化表单
     */
    useEffect(() => {
        if (user) {
            infoForm.setFieldsValue({
                username: user.username,
                realname: user.realname,
                email: user.email,
                phone: user.phone || '',
                avatar: user.avatar || ''
            })
            setAvatarUrl(user.avatar || '')
        }
    }, [user, infoForm])

    /**
     * 保存个人信息
     */
    const handleSaveInfo = async (values) => {
        setLoading(true)
        try {
            // 如果 avatarUrl 有变化，也添加到 values 中
            if (avatarUrl !== user.avatar) {
                values.avatar = avatarUrl
            }
            await updateUser(user.id, values)
            message.success('个人信息更新成功')
        } catch (error) {
            console.error('更新个人信息失败:', error)
            message.error('更新失败，请稍后重试')
        } finally {
            setLoading(false)
        }
    }

    /**
     * 修改密码
     */
    const handleChangePassword = async (values) => {
        setLoading(true)
        try {
            // 验证新密码和确认密码是否一致
            if (values.newPassword !== values.confirmPassword) {
                message.error('两次输入的密码不一致')
                return
            }

            // 调用修改密码 API
            await changePassword(values.oldPassword, values.newPassword)

            message.success('密码修改成功，请重新登录')
            passwordForm.resetFields()
        } catch (error) {
            console.error('修改密码失败:', error)
            message.error('修改密码失败，请检查原密码是否正确')
        } finally {
            setLoading(false)
        }
    }

    /**
     * 处理头像上传
     */
    const handleAvatarChange = (info) => {
        if (info.file.status === 'uploading') {
            setLoading(true)
            return
        }
        if (info.file.status === 'done') {
            // 获取上传后的图片 URL
            const url = info.file.response?.url || ''
            setAvatarUrl(url)
            message.success('头像上传成功')
            setLoading(false)
        }
        if (info.file.status === 'error') {
            message.error('头像上传失败')
            setLoading(false)
        }
    }

    /**
     * 上传前的验证
     */
    const beforeUpload = (file) => {
        const isImage = file.type.startsWith('image/')
        if (!isImage) {
            message.error('只能上传图片文件！')
            return false
        }
        const isLt2M = file.size / 1024 / 1024 < 2
        if (!isLt2M) {
            message.error('图片大小不能超过 2MB！')
            return false
        }
        return true
    }

    /**
     * 渲染基本信息选项卡
     */
    const renderInfoTab = () => (
        <Row gutter={24}>
            {/* 左侧：头像和基本信息 */}
            <Col xs={24} md={8}>
                <Card>
                    <div className="text-center">
                        <Avatar
                            size={120}
                            src={avatarUrl}
                            icon={<UserOutlined />}
                            className="mb-4"
                        />
                        <Title level={4}>{user?.realname || user?.username}</Title>
                        <Text type="secondary">{user?.email}</Text>

                        <Divider />

                        <div className="mb-4">
                            <Input
                                placeholder="输入头像 URL"
                                value={avatarUrl}
                                onChange={(e) => setAvatarUrl(e.target.value)}
                                prefix={<UploadOutlined />}
                            />
                        </div>

                        <Divider />

                        <Descriptions column={1} size="small">
                            <Descriptions.Item label="用户ID">
                                {user?.id}
                            </Descriptions.Item>
                            <Descriptions.Item label="注册时间">
                                {user?.created_at || '2024-01-01'}
                            </Descriptions.Item>
                            <Descriptions.Item label="最后登录">
                                {user?.last_login_at || '刚刚'}
                            </Descriptions.Item>
                        </Descriptions>
                    </div>
                </Card>
            </Col>

            {/* 右侧：编辑表单 */}
            <Col xs={24} md={16}>
                <Card title="编辑个人信息">
                    <Alert
                        message="提示"
                        description="用户名不可修改，如需修改请联系管理员"
                        type="info"
                        showIcon
                        className="mb-6"
                    />

                    <Form
                        form={infoForm}
                        layout="vertical"
                        onFinish={handleSaveInfo}
                    >
                        <Form.Item
                            name="username"
                            label="用户名"
                        >
                            <Input disabled placeholder="用户名不可修改" />
                        </Form.Item>

                        <Form.Item
                            name="realname"
                            label="真实姓名"
                            rules={[
                                { required: true, message: '请输入真实姓名' }
                            ]}
                        >
                            <Input placeholder="请输入真实姓名" />
                        </Form.Item>

                        <Form.Item
                            name="email"
                            label="邮箱地址"
                            rules={[
                                { required: true, message: '请输入邮箱地址' },
                                { type: 'email', message: '请输入有效的邮箱地址' }
                            ]}
                        >
                            <Input placeholder="请输入邮箱地址" />
                        </Form.Item>

                        <Form.Item
                            name="phone"
                            label="手机号码"
                            rules={[
                                { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' }
                            ]}
                        >
                            <Input placeholder="请输入手机号码" />
                        </Form.Item>

                        <Form.Item>
                            <Space>
                                <Button
                                    type="primary"
                                    htmlType="submit"
                                    icon={<SaveOutlined />}
                                    loading={loading}
                                >
                                    保存修改
                                </Button>
                                <Button onClick={() => infoForm.resetFields()}>
                                    重置
                                </Button>
                            </Space>
                        </Form.Item>
                    </Form>
                </Card>
            </Col>
        </Row>
    )

    /**
     * 渲染修改密码选项卡
     */
    const renderPasswordTab = () => (
        <Row justify="center">
            <Col xs={24} md={12}>
                <Card title="修改密码">
                    <Alert
                        message="安全提示"
                        description="为了您的账户安全，建议定期修改密码，密码长度至少 6 位，包含字母和数字"
                        type="warning"
                        showIcon
                        className="mb-6"
                    />

                    <Form
                        form={passwordForm}
                        layout="vertical"
                        onFinish={handleChangePassword}
                    >
                        <Form.Item
                            name="oldPassword"
                            label="原密码"
                            rules={[
                                { required: true, message: '请输入原密码' }
                            ]}
                        >
                            <Input.Password
                                prefix={<LockOutlined />}
                                placeholder="请输入原密码"
                            />
                        </Form.Item>

                        <Form.Item
                            name="newPassword"
                            label="新密码"
                            rules={[
                                { required: true, message: '请输入新密码' },
                                { min: 6, message: '密码长度至少 6 位' },
                                {
                                    pattern: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]/,
                                    message: '密码必须包含字母和数字'
                                }
                            ]}
                        >
                            <Input.Password
                                prefix={<LockOutlined />}
                                placeholder="请输入新密码"
                            />
                        </Form.Item>

                        <Form.Item
                            name="confirmPassword"
                            label="确认新密码"
                            dependencies={['newPassword']}
                            rules={[
                                { required: true, message: '请确认新密码' },
                                ({ getFieldValue }) => ({
                                    validator(_, value) {
                                        if (!value || getFieldValue('newPassword') === value) {
                                            return Promise.resolve()
                                        }
                                        return Promise.reject(new Error('两次输入的密码不一致'))
                                    }
                                })
                            ]}
                        >
                            <Input.Password
                                prefix={<LockOutlined />}
                                placeholder="请再次输入新密码"
                            />
                        </Form.Item>

                        <Form.Item>
                            <Space>
                                <Button
                                    type="primary"
                                    htmlType="submit"
                                    icon={<SaveOutlined />}
                                    loading={loading}
                                >
                                    修改密码
                                </Button>
                                <Button onClick={() => passwordForm.resetFields()}>
                                    重置
                                </Button>
                            </Space>
                        </Form.Item>
                    </Form>
                </Card>
            </Col>
        </Row>
    )

    /**
     * 渲染登录日志选项卡
     */
    const renderLoginLogTab = () => {
        // 模拟登录日志数据
        const loginLogs = [
            {
                id: 1,
                ip: '192.168.1.100',
                location: '北京市',
                device: 'Chrome 120 / Windows 10',
                time: '2024-01-25 11:30:00',
                status: 'success'
            },
            {
                id: 2,
                ip: '192.168.1.101',
                location: '上海市',
                device: 'Safari 17 / macOS',
                time: '2024-01-24 09:15:00',
                status: 'success'
            },
            {
                id: 3,
                ip: '192.168.1.102',
                location: '广州市',
                device: 'Firefox 121 / Ubuntu',
                time: '2024-01-23 14:20:00',
                status: 'failed'
            }
        ]

        const columns = [
            {
                title: 'IP 地址',
                dataIndex: 'ip',
                key: 'ip'
            },
            {
                title: '登录地点',
                dataIndex: 'location',
                key: 'location'
            },
            {
                title: '设备信息',
                dataIndex: 'device',
                key: 'device'
            },
            {
                title: '登录时间',
                dataIndex: 'time',
                key: 'time'
            },
            {
                title: '状态',
                dataIndex: 'status',
                key: 'status',
                render: (status) => (
                    <Tag color={status === 'success' ? 'success' : 'error'}>
                        {status === 'success' ? '成功' : '失败'}
                    </Tag>
                )
            }
        ]

        return (
            <Card title="登录日志">
                <Alert
                    message="安全提示"
                    description="如发现异常登录记录，请立即修改密码并联系管理员"
                    type="info"
                    showIcon
                    className="mb-4"
                />
                <Table
                    columns={columns}
                    dataSource={loginLogs}
                    rowKey="id"
                    pagination={{ pageSize: 10 }}
                />
            </Card>
        )
    }

    /**
     * 选项卡配置
     */
    const tabItems = [
        {
            key: 'info',
            label: (
                <span>
                    <UserOutlined />
                    基本信息
                </span>
            ),
            children: renderInfoTab()
        },
        {
            key: 'password',
            label: (
                <span>
                    <LockOutlined />
                    修改密码
                </span>
            ),
            children: renderPasswordTab()
        },
        {
            key: 'logs',
            label: (
                <span>
                    <HistoryOutlined />
                    登录日志
                </span>
            ),
            children: renderLoginLogTab()
        }
    ]

    /**
     * 主渲染函数
     */
    return (
        <div>
            <Title level={2}>个人中心</Title>
            <Paragraph type="secondary">
                查看和管理您的个人信息、修改密码、查看登录记录
            </Paragraph>

            <Tabs
                activeKey={activeTab}
                onChange={setActiveTab}
                items={tabItems}
                size="large"
                className="mt-6"
            />
        </div>
    )
}

// 导出组件
export default Profile
