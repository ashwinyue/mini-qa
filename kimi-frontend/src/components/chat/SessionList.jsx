/**
 * 会话列表组件
 * 
 * 显示历史会话列表
 */
import { List, Typography, Button, Popconfirm, Empty } from 'antd'
import { MessageOutlined, DeleteOutlined } from '@ant-design/icons'
import { useChatStore } from '../../stores/index.jsx'
import { formatTime } from '../../utils/helpers.jsx'

const { Text, Paragraph } = Typography

const SessionList = ({ onSessionClick }) => {
    const { sessions, currentSession, switchSession, deleteSession } = useChatStore()

    const handleSessionClick = (sessionId) => {
        switchSession(sessionId)
        // 移动端点击后关闭侧边栏
        if (window.innerWidth < 768) {
            onSessionClick?.()
        }
    }

    const handleDeleteSession = (sessionId) => {
        deleteSession(sessionId)
    }

    if (sessions.length === 0) {
        return (
            <div className="p-4">
                <Empty 
                    description="暂无历史会话" 
                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                    style={{ color: '#999' }}
                />
            </div>
        )
    }

    return (
        <List
            dataSource={sessions}
            renderItem={(session) => {
                const isActive = currentSession?.id === session.id
                const lastMessage = session.messages[session.messages.length - 1]
                
                return (
                    <List.Item
                        key={session.id}
                        onClick={() => handleSessionClick(session.id)}
                        className={`cursor-pointer transition-colors group px-4 ${
                            isActive ? 'bg-white/20' : 'hover:bg-white/10'
                        }`}
                        style={{ borderBottom: 'none' }}
                        actions={[
                            <Popconfirm
                                title="确定要删除这个会话吗？"
                                onConfirm={(e) => {
                                    e?.stopPropagation()
                                    handleDeleteSession(session.id)
                                }}
                                okText="确定"
                                cancelText="取消"
                                key="delete"
                            >
                                <Button
                                    type="text"
                                    size="small"
                                    icon={<DeleteOutlined />}
                                    onClick={(e) => e.stopPropagation()}
                                    className="opacity-0 group-hover:opacity-100 transition-opacity"
                                    style={{ color: '#999' }}
                                />
                            </Popconfirm>
                        ]}
                    >
                        <List.Item.Meta
                            avatar={<MessageOutlined style={{ color: '#999', fontSize: 18 }} />}
                            title={
                                <Text strong style={{ color: '#fff', fontSize: 14 }}>
                                    {session.title}
                                </Text>
                            }
                            description={
                                <div>
                                    {lastMessage && (
                                        <Paragraph
                                            ellipsis={{ rows: 1 }}
                                            style={{ color: '#999', fontSize: 12, marginBottom: 4 }}
                                        >
                                            {lastMessage.content}
                                        </Paragraph>
                                    )}
                                    <Text style={{ color: '#666', fontSize: 11 }}>
                                        {formatTime(session.updatedAt)}
                                    </Text>
                                </div>
                            }
                        />
                    </List.Item>
                )
            }}
        />
    )
}

export default SessionList
