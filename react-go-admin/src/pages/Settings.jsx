/**
 * 系统设置页面组件
 * 
 * 这是一个综合性的系统设置页面，提供多个配置选项
 * 
 * 主要功能模块：
 * 1. 基本设置：系统名称、Logo、描述等基础信息
 * 2. 安全设置：密码策略、会话超时、登录限制等
 * 3. 邮件设置：SMTP 配置、邮件模板等
 * 4. 通知设置：系统通知、消息推送配置
 * 5. 外观设置：主题颜色、布局配置等
 * 
 * 技术特点：
 * - 使用 Tabs 组件实现多选项卡切换
 * - 使用 Form 组件进行配置管理
 * - 支持配置的保存和重置
 * - 实时预览配置效果
 */

// 导入 React 和 Hooks
import { useState } from 'react'

// 导入 Ant Design 组件
import {
    Tabs,           // 选项卡组件
    Card,           // 卡片组件
    Form,           // 表单组件
    Input,          // 输入框组件
    InputNumber,    // 数字输入框
    Switch,         // 开关组件
    Button,         // 按钮组件
    Space,          // 间距组件
    message,        // 消息提示
    Divider,        // 分割线
    Select,         // 下拉选择
    ColorPicker,    // 颜色选择器
    Upload,         // 上传组件
    Radio,          // 单选框
    Slider,         // 滑块组件
    Typography,     // 文本组件
    Alert,          // 警告提示
    Row,            // 行布局
    Col             // 列布局
} from 'antd'

// 导入图标
import {
    SaveOutlined,       // 保存图标
    ReloadOutlined,     // 重置图标
    UploadOutlined,     // 上传图标
    InfoCircleOutlined  // 信息图标
} from '@ant-design/icons'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

/**
 * 系统设置页面组件
 * 
 * @component
 */
const Settings = () => {
    /**
     * 状态管理
     * 
     * - activeTab: 当前激活的选项卡
     * - loading: 保存加载状态
     */
    const [activeTab, setActiveTab] = useState('basic')
    const [loading, setLoading] = useState(false)

    // 创建表单实例
    const [basicForm] = Form.useForm()
    const [securityForm] = Form.useForm()
    const [emailForm] = Form.useForm()
    const [notificationForm] = Form.useForm()
    const [appearanceForm] = Form.useForm()

    /**
     * 保存设置
     * 
     * @param {string} formType - 表单类型
     * @param {Object} values - 表单值
     */
    const handleSave = async (formType, values) => {
        setLoading(true)
        try {
            // 模拟 API 调用
            await new Promise(resolve => setTimeout(resolve, 1000))

            console.log(`保存${formType}设置:`, values)
            message.success('设置保存成功')
        } catch (error) {
            console.error('保存设置失败:', error)
            message.error('保存设置失败')
        } finally {
            setLoading(false)
        }
    }

    /**
     * 重置表单
     * 
     * @param {Object} form - 表单实例
     */
    const handleReset = (form) => {
        form.resetFields()
        message.info('已重置为默认值')
    }

    /**
     * 渲染基本设置选项卡
     */
    const renderBasicSettings = () => (
        <Card>
            <Form
                form={basicForm}
                layout="vertical"
                initialValues={{
                    systemName: 'React 管理系统',
                    systemDescription: '基于 React + Ant Design 的现代化管理系统',
                    companyName: '示例公司',
                    contactEmail: 'admin@example.com',
                    contactPhone: '400-123-4567',
                    address: '北京市朝阳区示例大厦',
                    icp: '京ICP备12345678号',
                    enableRegistration: false,
                    enableGuestAccess: false
                }}
                onFinish={(values) => handleSave('基本', values)}
            >
                <Title level={4}>系统信息</Title>
                <Paragraph type="secondary">
                    配置系统的基本信息，这些信息会显示在系统的各个位置
                </Paragraph>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="systemName"
                            label="系统名称"
                            rules={[{ required: true, message: '请输入系统名称' }]}
                        >
                            <Input placeholder="请输入系统名称" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="companyName"
                            label="公司名称"
                        >
                            <Input placeholder="请输入公司名称" />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="systemDescription"
                    label="系统描述"
                >
                    <TextArea
                        rows={3}
                        placeholder="请输入系统描述"
                        showCount
                        maxLength={200}
                    />
                </Form.Item>

                <Divider />

                <Title level={4}>联系信息</Title>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="contactEmail"
                            label="联系邮箱"
                            rules={[{ type: 'email', message: '请输入有效的邮箱地址' }]}
                        >
                            <Input placeholder="请输入联系邮箱" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="contactPhone"
                            label="联系电话"
                        >
                            <Input placeholder="请输入联系电话" />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="address"
                    label="公司地址"
                >
                    <Input placeholder="请输入公司地址" />
                </Form.Item>

                <Form.Item
                    name="icp"
                    label="ICP 备案号"
                >
                    <Input placeholder="请输入 ICP 备案号" />
                </Form.Item>

                <Divider />

                <Title level={4}>功能开关</Title>

                <Form.Item
                    name="enableRegistration"
                    label="开放注册"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Form.Item
                    name="enableGuestAccess"
                    label="允许游客访问"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Form.Item>
                    <Space>
                        <Button
                            type="primary"
                            htmlType="submit"
                            icon={<SaveOutlined />}
                            loading={loading}
                        >
                            保存设置
                        </Button>
                        <Button
                            icon={<ReloadOutlined />}
                            onClick={() => handleReset(basicForm)}
                        >
                            重置
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )

    /**
     * 渲染安全设置选项卡
     */
    const renderSecuritySettings = () => (
        <Card>
            <Form
                form={securityForm}
                layout="vertical"
                initialValues={{
                    passwordMinLength: 6,
                    passwordRequireUppercase: true,
                    passwordRequireNumber: true,
                    passwordRequireSpecialChar: false,
                    sessionTimeout: 30,
                    maxLoginAttempts: 5,
                    lockoutDuration: 15,
                    enableTwoFactor: false,
                    enableIpWhitelist: false
                }}
                onFinish={(values) => handleSave('安全', values)}
            >
                <Alert
                    message="安全提示"
                    description="合理的安全策略可以有效保护系统免受攻击，建议启用多重安全措施"
                    type="info"
                    showIcon
                    icon={<InfoCircleOutlined />}
                    className="mb-6"
                />

                <Title level={4}>密码策略</Title>

                <Form.Item
                    name="passwordMinLength"
                    label="密码最小长度"
                    extra="建议设置为 8 位以上"
                >
                    <Slider min={6} max={20} marks={{ 6: '6', 10: '10', 15: '15', 20: '20' }} />
                </Form.Item>

                <Form.Item
                    name="passwordRequireUppercase"
                    label="要求包含大写字母"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="是" unCheckedChildren="否" />
                </Form.Item>

                <Form.Item
                    name="passwordRequireNumber"
                    label="要求包含数字"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="是" unCheckedChildren="否" />
                </Form.Item>

                <Form.Item
                    name="passwordRequireSpecialChar"
                    label="要求包含特殊字符"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="是" unCheckedChildren="否" />
                </Form.Item>

                <Divider />

                <Title level={4}>会话管理</Title>

                <Form.Item
                    name="sessionTimeout"
                    label="会话超时时间（分钟）"
                    extra="用户无操作后自动登出的时间"
                >
                    <InputNumber min={5} max={120} className="w-full" />
                </Form.Item>

                <Divider />

                <Title level={4}>登录保护</Title>

                <Form.Item
                    name="maxLoginAttempts"
                    label="最大登录尝试次数"
                    extra="超过此次数将锁定账户"
                >
                    <InputNumber min={3} max={10} className="w-full" />
                </Form.Item>

                <Form.Item
                    name="lockoutDuration"
                    label="账户锁定时长（分钟）"
                >
                    <InputNumber min={5} max={60} className="w-full" />
                </Form.Item>

                <Divider />

                <Title level={4}>高级安全</Title>

                <Form.Item
                    name="enableTwoFactor"
                    label="启用双因素认证"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Form.Item
                    name="enableIpWhitelist"
                    label="启用 IP 白名单"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Form.Item>
                    <Space>
                        <Button
                            type="primary"
                            htmlType="submit"
                            icon={<SaveOutlined />}
                            loading={loading}
                        >
                            保存设置
                        </Button>
                        <Button
                            icon={<ReloadOutlined />}
                            onClick={() => handleReset(securityForm)}
                        >
                            重置
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )

    /**
     * 渲染邮件设置选项卡
     */
    const renderEmailSettings = () => (
        <Card>
            <Form
                form={emailForm}
                layout="vertical"
                initialValues={{
                    smtpHost: 'smtp.example.com',
                    smtpPort: 587,
                    smtpUser: 'noreply@example.com',
                    smtpEncryption: 'tls',
                    senderName: 'React 管理系统',
                    senderEmail: 'noreply@example.com'
                }}
                onFinish={(values) => handleSave('邮件', values)}
            >
                <Title level={4}>SMTP 配置</Title>
                <Paragraph type="secondary">
                    配置 SMTP 服务器信息，用于发送系统邮件
                </Paragraph>

                <Row gutter={16}>
                    <Col span={16}>
                        <Form.Item
                            name="smtpHost"
                            label="SMTP 服务器"
                            rules={[{ required: true, message: '请输入 SMTP 服务器地址' }]}
                        >
                            <Input placeholder="smtp.example.com" />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="smtpPort"
                            label="端口"
                            rules={[{ required: true, message: '请输入端口' }]}
                        >
                            <InputNumber min={1} max={65535} className="w-full" />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="smtpUser"
                    label="SMTP 用户名"
                    rules={[{ required: true, message: '请输入 SMTP 用户名' }]}
                >
                    <Input placeholder="请输入 SMTP 用户名" />
                </Form.Item>

                <Form.Item
                    name="smtpPassword"
                    label="SMTP 密码"
                    extra="密码将加密存储"
                >
                    <Input.Password placeholder="请输入 SMTP 密码" />
                </Form.Item>

                <Form.Item
                    name="smtpEncryption"
                    label="加密方式"
                >
                    <Radio.Group>
                        <Radio value="none">无</Radio>
                        <Radio value="ssl">SSL</Radio>
                        <Radio value="tls">TLS</Radio>
                    </Radio.Group>
                </Form.Item>

                <Divider />

                <Title level={4}>发件人信息</Title>

                <Form.Item
                    name="senderName"
                    label="发件人名称"
                >
                    <Input placeholder="请输入发件人名称" />
                </Form.Item>

                <Form.Item
                    name="senderEmail"
                    label="发件人邮箱"
                    rules={[{ type: 'email', message: '请输入有效的邮箱地址' }]}
                >
                    <Input placeholder="请输入发件人邮箱" />
                </Form.Item>

                <Form.Item>
                    <Space>
                        <Button
                            type="primary"
                            htmlType="submit"
                            icon={<SaveOutlined />}
                            loading={loading}
                        >
                            保存设置
                        </Button>
                        <Button
                            icon={<ReloadOutlined />}
                            onClick={() => handleReset(emailForm)}
                        >
                            重置
                        </Button>
                        <Button>
                            发送测试邮件
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )

    /**
     * 渲染通知设置选项卡
     */
    const renderNotificationSettings = () => (
        <Card>
            <Form
                form={notificationForm}
                layout="vertical"
                initialValues={{
                    enableEmailNotification: true,
                    enableSmsNotification: false,
                    enablePushNotification: true,
                    notifyOnUserLogin: true,
                    notifyOnUserCreate: true,
                    notifyOnUserDelete: false,
                    notifyOnSystemError: true
                }}
                onFinish={(values) => handleSave('通知', values)}
            >
                <Title level={4}>通知渠道</Title>

                <Form.Item
                    name="enableEmailNotification"
                    label="邮件通知"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Form.Item
                    name="enableSmsNotification"
                    label="短信通知"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Form.Item
                    name="enablePushNotification"
                    label="推送通知"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                </Form.Item>

                <Divider />

                <Title level={4}>通知事件</Title>
                <Paragraph type="secondary">
                    选择需要接收通知的系统事件
                </Paragraph>

                <Form.Item
                    name="notifyOnUserLogin"
                    label="用户登录"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="通知" unCheckedChildren="不通知" />
                </Form.Item>

                <Form.Item
                    name="notifyOnUserCreate"
                    label="用户创建"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="通知" unCheckedChildren="不通知" />
                </Form.Item>

                <Form.Item
                    name="notifyOnUserDelete"
                    label="用户删除"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="通知" unCheckedChildren="不通知" />
                </Form.Item>

                <Form.Item
                    name="notifyOnSystemError"
                    label="系统错误"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="通知" unCheckedChildren="不通知" />
                </Form.Item>

                <Form.Item>
                    <Space>
                        <Button
                            type="primary"
                            htmlType="submit"
                            icon={<SaveOutlined />}
                            loading={loading}
                        >
                            保存设置
                        </Button>
                        <Button
                            icon={<ReloadOutlined />}
                            onClick={() => handleReset(notificationForm)}
                        >
                            重置
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )

    /**
     * 渲染外观设置选项卡
     */
    const renderAppearanceSettings = () => (
        <Card>
            <Form
                form={appearanceForm}
                layout="vertical"
                initialValues={{
                    theme: 'light',
                    primaryColor: '#1890ff',
                    sidebarCollapsed: false,
                    showBreadcrumb: true,
                    showFooter: true,
                    pageSize: 10
                }}
                onFinish={(values) => handleSave('外观', values)}
            >
                <Title level={4}>主题设置</Title>

                <Form.Item
                    name="theme"
                    label="主题模式"
                >
                    <Radio.Group>
                        <Radio value="light">浅色</Radio>
                        <Radio value="dark">深色</Radio>
                        <Radio value="auto">跟随系统</Radio>
                    </Radio.Group>
                </Form.Item>

                <Form.Item
                    name="primaryColor"
                    label="主题色"
                    extra="选择系统的主题颜色"
                >
                    <ColorPicker showText />
                </Form.Item>

                <Divider />

                <Title level={4}>布局设置</Title>

                <Form.Item
                    name="sidebarCollapsed"
                    label="默认折叠侧边栏"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="是" unCheckedChildren="否" />
                </Form.Item>

                <Form.Item
                    name="showBreadcrumb"
                    label="显示面包屑导航"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="显示" unCheckedChildren="隐藏" />
                </Form.Item>

                <Form.Item
                    name="showFooter"
                    label="显示页脚"
                    valuePropName="checked"
                >
                    <Switch checkedChildren="显示" unCheckedChildren="隐藏" />
                </Form.Item>

                <Divider />

                <Title level={4}>数据展示</Title>

                <Form.Item
                    name="pageSize"
                    label="默认每页条数"
                >
                    <Select>
                        <Select.Option value={10}>10 条/页</Select.Option>
                        <Select.Option value={20}>20 条/页</Select.Option>
                        <Select.Option value={50}>50 条/页</Select.Option>
                        <Select.Option value={100}>100 条/页</Select.Option>
                    </Select>
                </Form.Item>

                <Form.Item>
                    <Space>
                        <Button
                            type="primary"
                            htmlType="submit"
                            icon={<SaveOutlined />}
                            loading={loading}
                        >
                            保存设置
                        </Button>
                        <Button
                            icon={<ReloadOutlined />}
                            onClick={() => handleReset(appearanceForm)}
                        >
                            重置
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )

    /**
     * 选项卡配置
     */
    const tabItems = [
        {
            key: 'basic',
            label: '基本设置',
            children: renderBasicSettings()
        },
        {
            key: 'security',
            label: '安全设置',
            children: renderSecuritySettings()
        },
        {
            key: 'email',
            label: '邮件设置',
            children: renderEmailSettings()
        },
        {
            key: 'notification',
            label: '通知设置',
            children: renderNotificationSettings()
        },
        {
            key: 'appearance',
            label: '外观设置',
            children: renderAppearanceSettings()
        }
    ]

    /**
     * 主渲染函数
     */
    return (
        <div>
            <Title level={2}>系统设置</Title>
            <Paragraph type="secondary">
                管理系统的各项配置，包括基本信息、安全策略、邮件服务、通知规则和外观主题
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
export default Settings
