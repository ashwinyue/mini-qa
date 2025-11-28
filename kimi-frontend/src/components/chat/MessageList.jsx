/**
 * 消息列表组件
 * 
 * 显示所有消息
 */
import { useEffect, useRef } from 'react'
import MessageItem from './MessageItem.jsx'

const MessageList = ({ messages, onSuggestionClick }) => {
    const messagesEndRef = useRef(null)

    // 自动滚动到底部
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, [messages])

    if (!messages || messages.length === 0) {
        return null
    }

    return (
        <div className="space-y-4">
            {messages.map((message) => (
                <MessageItem
                    key={message.id}
                    message={message}
                    onSuggestionClick={onSuggestionClick}
                />
            ))}
            <div ref={messagesEndRef} />
        </div>
    )
}

export default MessageList
