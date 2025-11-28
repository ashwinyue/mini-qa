/**
 * 聊天界面组件
 * 
 * 集成消息列表、输入框、语音录制、图片上传
 */
import { useState, useRef, useEffect } from 'react'
import { Input, Button, message as antMessage } from 'antd'
import { SendOutlined, AudioOutlined, PictureOutlined, LoadingOutlined, PlusOutlined } from '@ant-design/icons'
import { useChatStore } from '../../stores'
import { sendMessage, getGreeting, getSuggestions } from '../../api/chat'
import MessageList from './MessageList'
import AudioRecorder from './AudioRecorder'
import ImageUploader from './ImageUploader'
import { v4 as uuidv4 } from 'uuid'

const { TextArea } = Input

/**
 * 聊天界面组件
 */
const ChatInterface = () => {
    const [inputValue, setInputValue] = useState('')
    const [uploadedImages, setUploadedImages] = useState([])
    const [showAudioRecorder, setShowAudioRecorder] = useState(false)
    const [showImageUploader, setShowImageUploader] = useState(false)

    const {
        currentSession,
        isLoading,
        threadId,
        addMessage,
        updateMessage,
        setLoading,
        createSession,
    } = useChatStore()

    const messagesEndRef = useRef(null)
    const eventSourceRef = useRef(null)

    // 初始化会话
    useEffect(() => {
        if (!currentSession) {
            createSession('AI助手对话')
        }
    }, [currentSession, createSession])

    // 加载欢迎语
    useEffect(() => {
        if (currentSession && currentSession.messages.length === 0) {
            loadGreeting()
        }
    }, [currentSession?.id])

    // 滚动到底部
    useEffect(() => {
        scrollToBottom()
    }, [currentSession?.messages])

    // 清理 SSE 连接
    useEffect(() => {
        return () => {
            if (eventSourceRef.current) {
                eventSourceRef.current.close()
            }
        }
    }, [])

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }

    const loadGreeting = async () => {
        try {
            const greeting = await getGreeting()
            const greetingMessage = {
                id: 'greeting-' + Date.now(),
                role: 'assistant',
                content: greeting.message || '您好，请问有什么可以帮您？',
                timestamp: new Date().toISOString(),
                suggestions: greeting.options?.map(opt => opt.title) || [],
            }
            addMessage(greetingMessage)
        } catch (error) {
            console.error('加载欢迎语失败:', error)
        }
    }

    const handleSendMessage = async () => {
        if (!inputValue.trim() && uploadedImages.length === 0) {
            antMessage.warning('请输入消息或上传图片')
            return
        }

        if (!threadId) {
            antMessage.error('会话初始化失败，请刷新页面重试')
            return
        }

        // 添加用户消息
        const userMessage = {
            id: uuidv4(),
            role: 'user',
            content: inputValue.trim(),
            timestamp: new Date().toISOString(),
            images: uploadedImages.length > 0 ? uploadedImages : undefined,
        }

        addMessage(userMessage)
        const query = inputValue.trim()
        const imagesToSend = uploadedImages.length > 0 ? [...uploadedImages] : undefined
        setInputValue('')
        setUploadedImages([])
        setLoading(true)

        try {
            // 发送消息到后端
            const response = await sendMessage(
                query,
                threadId,
                imagesToSend,
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

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault()
            handleSendMessage()
        }
    }

    const handlePaste = async (e) => {
        const items = e.clipboardData?.items
        if (!items) return

        const imageFiles = []
        for (let i = 0; i < items.length; i++) {
            const item = items[i]
            if (item.type.indexOf('image') !== -1) {
                const file = item.getAsFile()
                if (file) {
                    imageFiles.push(file)
                }
            }
        }

        if (imageFiles.length > 0) {
            e.preventDefault()
            
            // 转换为 Base64
            const base64Images = await Promise.all(
                imageFiles.map(file => {
                    return new Promise((resolve, reject) => {
                        const reader = new FileReader()
                        reader.onload = (e) => resolve(e.target.result)
                        reader.onerror = reject
                        reader.readAsDataURL(file)
                    })
                })
            )

            setUploadedImages(prev => [...prev, ...base64Images])
            antMessage.success(`已粘贴 ${imageFiles.length} 张图片`)
        }
    }

    const handleAudioRecorded = (audioData, duration) => {
        setShowAudioRecorder(false)

        if (!threadId) {
            antMessage.error('会话初始化失败，请刷新页面重试')
            return
        }

        // 添加语音消息
        const audioMessage = {
            id: uuidv4(),
            role: 'user',
            content: `[语音消息 - 时长: ${Math.round(duration)}秒]`,
            timestamp: new Date().toISOString(),
        }
        addMessage(audioMessage)
        setLoading(true)

        // 发送语音消息 - query 为空，让后端使用语音识别结果
        sendMessage('', threadId, undefined, audioData)
            .then(response => {
                const assistantMessage = {
                    id: uuidv4(),
                    role: 'assistant',
                    content: response.answer || '',
                    timestamp: new Date().toISOString(),
                    sources: response.sources,
                    route: response.route,
                }
                addMessage(assistantMessage)
                startSuggestionStream(threadId, assistantMessage.id)
            })
            .catch(error => {
                console.error('语音消息发送失败:', error)
                const errorMessage = {
                    id: uuidv4(),
                    role: 'system',
                    content: '语音消息发送失败，请重试',
                    timestamp: new Date().toISOString(),
                    error: true,
                }
                addMessage(errorMessage)
            })
            .finally(() => {
                setLoading(false)
            })
    }

    const handleImagesUploaded = (images) => {
        setUploadedImages(images)
        setShowImageUploader(false)
        antMessage.success(`已上传 ${images.length} 张图片`)
    }

    const handleSuggestionClick = (suggestion) => {
        setInputValue(suggestion)
    }

    const handleNewChat = () => {
        // 关闭 SSE 连接
        if (eventSourceRef.current) {
            eventSourceRef.current.close()
        }
        createSession('新对话')
    }

    return (
        <div className="flex flex-col h-full">
            {/* 顶部工具栏 */}
            <div className="flex items-center justify-between px-4 py-2 border-b bg-white">
                <span className="text-gray-600 text-sm">
                    {currentSession?.title || 'AI助手'}
                </span>
                <Button
                    type="text"
                    icon={<PlusOutlined />}
                    onClick={handleNewChat}
                    disabled={isLoading}
                >
                    新对话
                </Button>
            </div>

            {/* 消息列表 */}
            <div className="flex-1 overflow-y-auto p-4">
                {currentSession && (
                    <MessageList
                        messages={currentSession.messages}
                        onSuggestionClick={handleSuggestionClick}
                    />
                )}
                <div ref={messagesEndRef} />
            </div>

            {/* 输入区域 */}
            <div className="border-t p-4 bg-white">
                {/* 已上传图片预览 */}
                {uploadedImages.length > 0 && (
                    <div className="flex gap-2 mb-2 flex-wrap">
                        {uploadedImages.map((image, index) => (
                            <div key={index} className="relative">
                                <img
                                    src={image}
                                    alt={`上传图片 ${index + 1}`}
                                    className="w-16 h-16 object-cover rounded"
                                />
                                <button
                                    className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white rounded-full text-xs"
                                    onClick={() => setUploadedImages(prev => prev.filter((_, i) => i !== index))}
                                >
                                    ×
                                </button>
                            </div>
                        ))}
                    </div>
                )}

                {/* 输入框和按钮 */}
                <div className="flex items-end gap-2">
                    <div className="flex gap-1">
                        <Button
                            type="text"
                            icon={<AudioOutlined />}
                            onClick={() => setShowAudioRecorder(!showAudioRecorder)}
                            className={showAudioRecorder ? 'text-blue-500' : ''}
                            disabled={isLoading}
                        />
                        <Button
                            type="text"
                            icon={<PictureOutlined />}
                            onClick={() => setShowImageUploader(!showImageUploader)}
                            className={showImageUploader ? 'text-blue-500' : ''}
                            disabled={isLoading}
                        />
                    </div>

                    <TextArea
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        onKeyDown={handleKeyDown}
                        onPaste={handlePaste}
                        placeholder="输入消息，按 Enter 发送，Shift+Enter 换行，支持粘贴图片..."
                        autoSize={{ minRows: 1, maxRows: 4 }}
                        disabled={isLoading}
                        className="flex-1"
                    />

                    <Button
                        type="primary"
                        icon={isLoading ? <LoadingOutlined /> : <SendOutlined />}
                        onClick={handleSendMessage}
                        disabled={isLoading || (!inputValue.trim() && uploadedImages.length === 0)}
                    />
                </div>

                {/* 语音录制 */}
                {showAudioRecorder && (
                    <AudioRecorder
                        onRecord={handleAudioRecorded}
                        onCancel={() => setShowAudioRecorder(false)}
                    />
                )}

                {/* 图片上传 */}
                {showImageUploader && (
                    <ImageUploader
                        onUpload={handleImagesUploaded}
                        onCancel={() => setShowImageUploader(false)}
                    />
                )}
            </div>
        </div>
    )
}

export default ChatInterface
