/**
 * 系统设置页面组件
 */
import { useState } from 'react'
import {
    Tabs, Card, Form, Input, InputNumber, Switch, Button, Space,
    message, Select, Radio, Typography, Row, Col
} from 'antd'
import { SaveOutlined, ReloadOutlined } from '@ant-design/icons'

const { Title, Paragraph } = Typography

const Settings = () => {
    const [activeTab, setActiveTab] = useState('basic')
    const [loading, setLoading] = useState(false)
    const [basicForm] = Form.useForm()
    const [securityForm] = Form.useForm()
    const [appearanceForm] = Form.useForm()

    const handleSave = async (formType, values) => {
        setLoading(true)
        try {
            await new Promise(resolve => setTimeout(resolve, 500))
            console.log(`保存${formType}设置:`, values)
            message.success('设置保存成功')
        } catch {
            message.error('保存设置失败')
        } finally {
            setLoading(false)
        }
    }

    const handleReset = (form) => {
        form.resetFields()
        message.info('已重置为默认值')
    }

    const FormActions = ({ form }) => (
        <Form.Item>
            <Space>
                <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={loading}>
                    保存
                </Button>
                <Button icon={<ReloadOutlined />} onClick={() => handleReset(form)}>
                    重置
                </Button>
            </Space>
        </Form.Item>
    )

    const tabItems = [
        {
            key: 'basic',
            label: '基本设置',
            children: (
                <Card>
                    <Form
                        form={basicForm}
                        layout="vertical"
                        initialValues={{
                            systemName: 'AI 智能客服系统',
                            companyName: '示例公司',
                            contactEmail: 'admin@example.com',
                        }}
                        onFinish={(values) => handleSave('基本', values)}
                    >
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item name="systemName" label="系统名称" rules={[{ required: true }]}>
                                    <Input placeholder="请输入系统名称" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item name="companyName" label="公司名称">
                                    <Input placeholder="请输入公司名称" />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Form.Item name="contactEmail" label="联系邮箱">
                            <Input placeholder="请输入联系邮箱" />
                        </Form.Item>
                        <FormActions form={basicForm} />
                    </Form>
                </Card>
            )
        },
        {
            key: 'security',
            label: '安全设置',
            children: (
                <Card>
                    <Form
                        form={securityForm}
                        layout="vertical"
                        initialValues={{
                            sessionTimeout: 30,
                            maxLoginAttempts: 5,
                            enableTwoFactor: false,
                        }}
                        onFinish={(values) => handleSave('安全', values)}
                    >
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item name="sessionTimeout" label="会话超时（分钟）">
                                    <InputNumber min={5} max={120} className="w-full" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item name="maxLoginAttempts" label="最大登录尝试次数">
                                    <InputNumber min={3} max={10} className="w-full" />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Form.Item name="enableTwoFactor" label="双因素认证" valuePropName="checked">
                            <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                        </Form.Item>
                        <FormActions form={securityForm} />
                    </Form>
                </Card>
            )
        },
        {
            key: 'appearance',
            label: '外观设置',
            children: (
                <Card>
                    <Form
                        form={appearanceForm}
                        layout="vertical"
                        initialValues={{
                            theme: 'light',
                            sidebarCollapsed: false,
                            pageSize: 10,
                        }}
                        onFinish={(values) => handleSave('外观', values)}
                    >
                        <Form.Item name="theme" label="主题模式">
                            <Radio.Group>
                                <Radio value="light">浅色</Radio>
                                <Radio value="dark">深色</Radio>
                            </Radio.Group>
                        </Form.Item>
                        <Form.Item name="sidebarCollapsed" label="默认折叠侧边栏" valuePropName="checked">
                            <Switch checkedChildren="是" unCheckedChildren="否" />
                        </Form.Item>
                        <Form.Item name="pageSize" label="默认每页条数">
                            <Select>
                                <Select.Option value={10}>10 条/页</Select.Option>
                                <Select.Option value={20}>20 条/页</Select.Option>
                                <Select.Option value={50}>50 条/页</Select.Option>
                            </Select>
                        </Form.Item>
                        <FormActions form={appearanceForm} />
                    </Form>
                </Card>
            )
        }
    ]

    return (
        <div>
            <Title level={3}>系统设置</Title>
            <Paragraph type="secondary">管理系统配置</Paragraph>
            <Tabs activeKey={activeTab} onChange={setActiveTab} items={tabItems} />
        </div>
    )
}

export default Settings
