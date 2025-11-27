/**
 * 消息列表组件
 * 
 * 渲染消息数组
 */
import React from 'react'
import MessageItem from './MessageItem'

/**
 * 消息列表组件
 * 
 * @param {Object} props
 * @param {Array} props.messages - 消息数组
 * @param {Function} props.onSuggestionClick - 建议点击回调
 */
const MessageList = ({ messages = [], onSuggestionClick }) => {
    if (!messages || messages.length === 0) {
        return (
            <div className="flex items-center justify-center h-full text-gray-400">
                暂无消息
            </div>
        )
    }

    return (
        <div className="flex flex-col">
            {messages.map((message) => (
                <MessageItem
                    key={message.id}
                    message={message}
                    onSuggestionClick={onSuggestionClick}
                />
            ))}
        </div>
    )
}

export default MessageList
