/**
 * 侧边栏组件
 * 
 * 包含新建会话按钮、会话列表、用户菜单
 */
import { Button } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useChatStore } from '../../stores/index.jsx'
import SessionList from '../chat/SessionList.jsx'
import UserMenu from '../settings/UserMenu.jsx'

const Sidebar = ({ isOpen, onClose }) => {
    const { createSession } = useChatStore()

    const handleNewChat = () => {
        createSession('新对话')
        // 移动端关闭侧边栏
        if (window.innerWidth < 768) {
            onClose?.()
        }
    }

    return (
        <div className={`sidebar ${isOpen ? 'open' : ''}`}>
            {/* 顶部：新建会话按钮 */}
            <div className="p-4 border-b border-gray-700">
                <Button
                    type="default"
                    icon={<PlusOutlined />}
                    onClick={handleNewChat}
                    block
                    size="large"
                    style={{
                        backgroundColor: 'rgba(255, 255, 255, 0.1)',
                        borderColor: 'transparent',
                        color: '#fff'
                    }}
                    className="hover:bg-white/20"
                >
                    新建会话
                </Button>
            </div>

            {/* 中间：会话列表 */}
            <div className="flex-1 overflow-y-auto">
                <SessionList onSessionClick={onClose} />
            </div>

            {/* 底部：用户菜单 */}
            <div className="border-t border-gray-700">
                <UserMenu />
            </div>
        </div>
    )
}

export default Sidebar
