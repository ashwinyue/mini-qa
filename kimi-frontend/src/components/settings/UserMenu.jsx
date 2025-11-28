/**
 * 用户菜单组件
 * 
 * 显示用户头像和下拉菜单
 */
import { Dropdown, Avatar, Typography } from 'antd'
import { SettingOutlined, UserOutlined, UpOutlined, DownOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useState } from 'react'

const { Text } = Typography

const UserMenu = () => {
    const [isOpen, setIsOpen] = useState(false)
    const navigate = useNavigate()

    // 模拟用户信息
    const user = {
        name: 'Rosas',
        phone: '178****6418',
        avatar: null
    }

    const menuItems = [
        {
            key: 'settings',
            icon: <SettingOutlined />,
            label: '设置',
            onClick: () => navigate('/settings')
        }
    ]

    return (
        <Dropdown
            menu={{ items: menuItems }}
            trigger={['click']}
            placement="topLeft"
            onOpenChange={setIsOpen}
        >
            <div className="w-full p-4 hover:bg-white/10 transition-colors flex items-center gap-3 cursor-pointer">
                {/* 头像 */}
                <Avatar
                    size={40}
                    icon={<UserOutlined />}
                    src={user.avatar}
                    style={{ backgroundColor: 'rgba(255, 255, 255, 0.2)' }}
                />

                {/* 用户信息 */}
                <div className="flex-1 text-left min-w-0">
                    <Text strong style={{ color: '#fff', display: 'block' }} ellipsis>
                        {user.name}
                    </Text>
                    <Text style={{ color: '#999', fontSize: 12 }}>
                        {user.phone}
                    </Text>
                </div>

                {/* 展开图标 */}
                {isOpen ? (
                    <UpOutlined style={{ color: '#999', fontSize: 12 }} />
                ) : (
                    <DownOutlined style={{ color: '#999', fontSize: 12 }} />
                )}
            </div>
        </Dropdown>
    )
}

export default UserMenu
