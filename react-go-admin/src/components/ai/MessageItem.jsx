/**
 * 消息项组件
 * 
 * 根据消息角色渲染不同样式的消息气泡
 */
import React from 'react'
import { Avatar, Tag, Typography } from 'antd'
import { UserOutlined, RobotOutlined, InfoCircleOutlined } from '@ant-design/icons'
import ReactMarkdown from 'react-markdown'
import SuggestionList from './SuggestionList'
import dayjs from 'dayjs'

const { Text, Paragraph } = Typography

/**
 * 消息项组件
 * 
 * @param {Object} props
 * @param {Object} props.message - 消息对象
 * @param {Function} props.onSuggestionClick - 建议点击回调
 */
const MessageItem = ({ message, onSuggestionClick }) => {
    const { role, content: rawContent, timestamp, sources, suggestions, error, route } = message
    
    // 确保 content 是字符串
    const content = typeof rawContent === 'string' 
        ? rawContent 
        : (rawContent ? JSON.stringify(rawContent) : '')

    // 根据角色确定样式
    const isUser = role === 'user'
    const isAssistant = role === 'assistant'
    const isSystem = role === 'system'

    // 格式化时间
    const timeStr = timestamp ? dayjs(timestamp).format('HH:mm') : ''

    // 用户消息
    if (isUser) {
        return (
            <div className="flex justify-end mb-4">
                <div className="flex items-start gap-2 max-w-[80%]">
                    <div className="flex flex-col items-end">
                        <div 
                            className="bg-blue-500 text-white px-4 py-2 rounded-lg rounded-tr-none"
                            style={{ wordBreak: 'break-word' }}
                        >
                            {content}
                        </div>
                        {timeStr && (
                            <Text type="secondary" className="text-xs mt-1">
                                {timeStr}
                            </Text>
                        )}
                    </div>
                    <Avatar 
                        icon={<UserOutlined />} 
                        className="bg-blue-500 flex-shrink-0"
                    />
                </div>
            </div>
        )
    }

    // 助手消息
    if (isAssistant) {
        return (
            <div className="flex justify-start mb-4">
                <div className="flex items-start gap-2 max-w-[80%]">
                    <Avatar 
                        icon={<RobotOutlined />} 
                        className="bg-green-500 flex-shrink-0"
                    />
                    <div className="flex flex-col">
                        <div 
                            className="bg-gray-100 px-4 py-2 rounded-lg rounded-tl-none"
                            style={{ wordBreak: 'break-word' }}
                        >
                            {/* 使用 Markdown 渲染 */}
                            <div className="markdown-content">
                                <ReactMarkdown>{content}</ReactMarkdown>
                            </div>
                            
                            {/* 路由标签 */}
                            {route && (
                                <Tag color="blue" className="mt-2">
                                    {route}
                                </Tag>
                            )}
                            
                            {/* 来源引用 */}
                            {sources && sources.length > 0 && (
                                <div className="mt-2 pt-2 border-t border-gray-200">
                                    <Text type="secondary" className="text-xs">
                                        来源：
                                    </Text>
                                    {sources.map((source, index) => (
                                        <div key={index} className="text-xs text-gray-500 mt-1">
                                            {source.source || source.title}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                        
                        {/* 建议问题 */}
                        {suggestions && suggestions.length > 0 && (
                            <SuggestionList 
                                suggestions={suggestions} 
                                onClick={onSuggestionClick}
                            />
                        )}
                        
                        {timeStr && (
                            <Text type="secondary" className="text-xs mt-1">
                                {timeStr}
                            </Text>
                        )}
                    </div>
                </div>
            </div>
        )
    }

    // 系统消息
    if (isSystem) {
        return (
            <div className="flex justify-center mb-4">
                <div 
                    className={`px-4 py-2 rounded-lg max-w-[80%] ${
                        error ? 'bg-red-50 text-red-600' : 'bg-gray-50 text-gray-600'
                    }`}
                >
                    <div className="flex items-center gap-2">
                        <InfoCircleOutlined />
                        <span style={{ whiteSpace: 'pre-wrap' }}>{content}</span>
                    </div>
                </div>
            </div>
        )
    }

    return null
}

export default MessageItem
