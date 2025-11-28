/**
 * 聊天界面组件
 * 
 * 集成消息列表、输入框、欢迎屏幕
 */
import { useState, useRef, useEffect } from 'react'
import { Alert } from 'antd'
import { v4 as uuidv4 } from 'uuid'
import { useChatStore } from '../../stores/index.jsx'
import { sendMessage, getGreeting, getSuggestions } from '../../services/chatApi.jsx'
import WelcomeScreen from './WelcomeScreen.jsx'
import MessageList from './MessageList.jsx'
import ChatInput from './ChatInput.jsx'

const ChatInterface = () => {
    const {
        currentSession,
        isLoading,
        threadId,
        addMessage,
        updateMessage,
        setLoading,
        createSession,
    } = useChatStore()

    const eventSourceRef = useRef(null)
    const [error, setError] = useState(null)

    // 初始化会话
    useEffect(() => {
        if (!currentSession) {
            createSession('AI助手对话')
        }
    }, [currentSession, createSession])

    // 不自动加载欢迎语，等用户发送第一条消息后再显示对话界面

    // 清理 SSE 连接
    useEffect(() => {
        return () => {
            if (eventSourceRef.current) {
                eventSourceRef.current.close()
            }
        }
    }, [])

    // 移除自动加载欢迎语的函数

    const handleSendMessage = async (content, images) => {
        if (!threadId) {
            setError('会话初始化失败，请刷新页面重试')
            return
        }

        // 添加用户消息
        const userMessage = {
            id: uuidv4(),
            role: 'user',
            content,
            timestamp: new Date().toISOString(),
            images: images && images.length > 0 ? images : undefined,
        }

        addMessage(userMessage)
        setLoading(true)
        setError(null)

        try {
            // 发送消息到后端
            const response = await sendMessage(
                content,
                threadId,
                images,
                undefined
            )

            // 处理命令响应
            if (response.commands || response.history || response.reset !== undefined) {
                let commandResult = ''
                if (response.commands) {
                    commandResult = '可用命令：\n' + response.commands.map(cmd =>
                        `${cmd.cmd}: ${cmd.desc}`
                    ).join('\n')
                } else if (response.history) {
                    commandResult = '历史消息：\n' + response.history.map(msg =>
                        `${msg.role}: ${msg.content}`
                    ).join('\n')
                } else if (response.reset) {
                    commandResult = '会话已重置'
                }

                const systemMessage = {
                    id: uuidv4(),
                    role: 'system',
                    content: commandResult,
                    timestamp: new Date().toISOString(),
                }
                addMessage(systemMessage)
            } else {
                // 添加助手响应
                const assistantMessage = {
                    id: uuidv4(),
                    role: 'assistant',
                    content: response.answer || '',
                    timestamp: new Date().toISOString(),
                    sources: response.sources,
                    route: response.route,
                }
                addMessage(assistantMessage)

                // 启动建议流
                startSuggestionStream(threadId, assistantMessage.id)
            }

        } catch (error) {
            console.error('发送消息失败:', error)
            const errorMessage = {
                id: uuidv4(),
                role: 'system',
                content: error.message || '消息发送失败，请重试',
                timestamp: new Date().toISOString(),
                error: true,
            }
            addMessage(errorMessage)
            setError(error.message)
        } finally {
            setLoading(false)
        }
    }

    const startSuggestionStream = (threadId, messageId) => {
        // 关闭之前的连接
        if (eventSourceRef.current) {
            eventSourceRef.current.close()
        }

        // 延迟启动 SSE
        setTimeout(() => {
            eventSourceRef.current = getSuggestions(
                threadId,
                (data) => {
                    if (data.event === 'react' && data.suggestions) {
                        updateMessage(messageId, { suggestions: data.suggestions })
                    }
                },
                (error) => {
                    console.error('建议流错误:', error)
                }
            )
        }, 500)
    }

    const handleSuggestionClick = (suggestion) => {
        handleSendMessage(suggestion, [])
    }

    // 只有当有用户消息时才显示对话界面（排除系统欢迎语）
    const hasMessages = currentSession && currentSession.messages.some(msg => msg.role === 'user')

    return (
        <div className="flex flex-col h-full">
            {hasMessages ? (
                <>
                    {/* 消息区域 */}
                    <div className="flex-1 overflow-y-auto bg-gray-50">
                        <div className="max-w-4xl mx-auto p-4">
                            <MessageList
                                messages={currentSession.messages}
                                onSuggestionClick={handleSuggestionClick}
                            />
                        </div>
                    </div>

                    {/* 错误提示 */}
                    {error && (
                        <div className="mx-4 mb-2">
                            <Alert
                                description={error}
                                type="error"
                                closable
                                afterClose={() => setError(null)}
                            />
                        </div>
                    )}

                    {/* 输入区域 - 固定在底部 */}
                    <div className="border-t bg-white">
                        <div className="max-w-4xl mx-auto">
                            <ChatInput
                                onSend={handleSendMessage}
                                disabled={isLoading}
                                placeholder="问我问..."
                            />
                        </div>
                    </div>
                </>
            ) : (
                <>
                    {/* 欢迎屏幕 - 占据大部分空间 */}
                    <div className="flex-1 overflow-y-auto">
                        <WelcomeScreen onQuickAction={handleSuggestionClick} />
                    </div>

                    {/* 输入区域 - 居中显示 */}
                    <div className="pb-8">
                        <div className="max-w-3xl mx-auto px-4">
                            <ChatInput
                                onSend={handleSendMessage}
                                disabled={isLoading}
                                placeholder="问我问..."
                                centered={true}
                            />
                        </div>
                    </div>
                </>
            )}
        </div>
    )
}

export default ChatInterface
