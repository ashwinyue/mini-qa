/**
 * 消息项组件
 * 
 * 显示单条消息，支持用户/助手/系统消息的不同样式
 */
import ReactMarkdown from 'react-markdown'
import { Avatar, Tag, Button, Space } from 'antd'
import { UserOutlined, RobotOutlined, WarningOutlined } from '@ant-design/icons'
import { formatTime } from '../../utils/helpers.jsx'

const MessageItem = ({ message, onSuggestionClick }) => {
    const { role, content, timestamp, suggestions, sources, error } = message

    // 根据角色确定样式
    const isUser = role === 'user'
    const isAssistant = role === 'assistant'
    const isSystem = role === 'system'

    return (
        <div className={`flex gap-3 mb-6 ${isUser ? 'flex-row-reverse' : ''}`}>
            {/* 头像 */}
            <Avatar
                size={32}
                icon={isUser ? <UserOutlined /> : isSystem ? <WarningOutlined /> : <RobotOutlined />}
                style={{
                    backgroundColor: isUser ? '#4F46E5' : isSystem ? '#faad14' : '#1890ff',
                    color: '#fff'
                }}
            />

            {/* 消息内容 */}
            <div className={`flex-1 ${isUser ? 'flex justify-end' : ''}`}>
                <div className={`max-w-[85%] ${isUser ? 'text-right' : ''}`}>
                    {/* 消息气泡 */}
                    <div className={`inline-block px-4 py-2.5 rounded-2xl ${
                        isUser ? 'bg-gray-100 text-gray-800' : 
                        error ? 'bg-red-50 text-red-800 border border-red-200' :
                        'bg-white text-gray-800 shadow-sm'
                    }`}>
                        {isAssistant ? (
                            <div className="markdown-content text-sm leading-relaxed">
                                <ReactMarkdown>{content}</ReactMarkdown>
                            </div>
                        ) : (
                            <div className="whitespace-pre-wrap text-sm">{content}</div>
                        )}
                    </div>

                    {/* 操作按钮 - 仅助手消息显示 */}
                    {isAssistant && !error && (
                        <div className={`mt-2 flex gap-2 ${isUser ? 'justify-end' : ''}`}>
                            <Button type="text" size="small" icon={<span>📋</span>} title="复制" />
                            <Button type="text" size="small" icon={<span>🔄</span>} title="重新生成" />
                            <Button type="text" size="small" icon={<span>👍</span>} title="点赞" />
                            <Button type="text" size="small" icon={<span>👎</span>} title="点踩" />
                            <Button type="text" size="small" icon={<span>💬</span>} title="分享" />
                        </div>
                    )}

                    {/* 来源 */}
                    {sources && sources.length > 0 && (
                        <div className="mt-2">
                            <Space size={[0, 8]} wrap>
                                {sources.map((source, index) => (
                                    <Tag key={index} color="blue" className="text-xs">
                                        {source.title}
                                    </Tag>
                                ))}
                            </Space>
                        </div>
                    )}

                    {/* 建议问题 */}
                    {suggestions && suggestions.length > 0 && (
                        <div className="mt-3 space-y-2">
                            <Space direction="vertical" className="w-full">
                                {suggestions.map((suggestion, index) => (
                                    <Button
                                        key={index}
                                        block
                                        size="small"
                                        onClick={() => onSuggestionClick?.(suggestion)}
                                        className="text-left text-xs"
                                        style={{ height: 'auto', padding: '6px 12px' }}
                                    >
                                        {suggestion}
                                    </Button>
                                ))}
                            </Space>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

export default MessageItem
