/**
 * 设置面板组件
 * 
 * 显示用户设置选项
 */
import { Card, Avatar, Typography, Segmented, Space, Divider } from 'antd'
import { UserOutlined, SunOutlined, MoonOutlined, DesktopOutlined, GlobalOutlined } from '@ant-design/icons'
import { useSettingsStore } from '../../stores/index.jsx'
import { maskPhone } from '../../utils/helpers.jsx'

const { Title, Text, Paragraph } = Typography

const SettingsPanel = () => {
    const { theme, language, setTheme, setLanguage } = useSettingsStore()

    // 模拟用户信息
    const user = {
        name: 'Rosas',
        phone: '17800006418',
        avatar: null
    }

    const themeOptions = [
        { label: '浅色', value: 'light', icon: <SunOutlined /> },
        { label: '深色', value: 'dark', icon: <MoonOutlined /> },
        { label: '跟随系统', value: 'system', icon: <DesktopOutlined /> },
    ]

    const languageOptions = [
        { label: '中文', value: 'zh' },
        { label: 'English', value: 'en' },
    ]

    return (
        <div className="max-w-2xl mx-auto p-6">
            <Title level={2}>设置</Title>

            {/* 用户信息 */}
            <Card className="mb-6">
                <Space size="large">
                    <Avatar
                        size={64}
                        icon={<UserOutlined />}
                        src={user.avatar}
                    />
                    <div>
                        <Title level={4} style={{ marginBottom: 4 }}>{user.name}</Title>
                        <Text type="secondary">{maskPhone(user.phone)}</Text>
                    </div>
                </Space>
            </Card>

            {/* 通用设置 */}
            <Card title="通用" className="mb-6">
                {/* 界面主题 */}
                <div className="mb-6">
                    <Space className="mb-3">
                        <SunOutlined style={{ fontSize: 18 }} />
                        <Text strong>界面主题</Text>
                    </Space>
                    <Segmented
                        block
                        options={themeOptions}
                        value={theme}
                        onChange={setTheme}
                    />
                </div>

                <Divider />

                {/* 语言 */}
                <div>
                    <Space className="mb-3">
                        <GlobalOutlined style={{ fontSize: 18 }} />
                        <Text strong>Language</Text>
                    </Space>
                    <Segmented
                        block
                        options={languageOptions}
                        value={language}
                        onChange={setLanguage}
                    />
                </div>
            </Card>

            {/* 关于我们 */}
            <Card title="关于我们">
                <Paragraph>
                    <Text>Kimi AI Assistant v1.0.0</Text>
                </Paragraph>
                <Paragraph type="secondary">
                    © 2024 All rights reserved
                </Paragraph>
            </Card>
        </div>
    )
}

export default SettingsPanel
