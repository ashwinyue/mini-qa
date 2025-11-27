/**
 * 智能对话页面
 */
import React from 'react'
import { ChatInterface } from '../../components/ai'

const ChatPage = () => {
    return (
        <div className="h-full" style={{ height: 'calc(100vh - 180px)' }}>
            <ChatInterface />
        </div>
    )
}

export default ChatPage
