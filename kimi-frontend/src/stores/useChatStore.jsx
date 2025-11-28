/**
 * 聊天状态 Store
 * 
 * 管理聊天会话、消息和加载状态
 */
import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { v4 as uuidv4 } from 'uuid'

/**
 * 创建聊天状态 Store
 */
const useChatStore = create(
    persist(
        (set, get) => ({
            // 状态
            currentSession: null,    // 当前会话
            sessions: [],            // 会话列表
            isLoading: false,        // 加载状态
            isConnected: false,      // 连接状态
            error: null,             // 错误信息
            threadId: null,          // 当前线程 ID

            // 设置当前会话
            setCurrentSession: (session) => set({ currentSession: session }),

            // 添加消息
            addMessage: (message) => {
                const { currentSession, sessions } = get()
                if (!currentSession) return

                const updatedSession = {
                    ...currentSession,
                    messages: [...currentSession.messages, message],
                    updatedAt: new Date().toISOString(),
                }

                const updatedSessions = sessions.map(session =>
                    session.id === currentSession.id ? updatedSession : session
                )

                set({
                    currentSession: updatedSession,
                    sessions: updatedSessions,
                })
            },

            // 更新消息
            updateMessage: (messageId, updates) => {
                const { currentSession, sessions } = get()
                if (!currentSession) return

                const updatedMessages = currentSession.messages.map(msg =>
                    msg.id === messageId ? { ...msg, ...updates } : msg
                )

                const updatedSession = {
                    ...currentSession,
                    messages: updatedMessages,
                    updatedAt: new Date().toISOString(),
                }

                const updatedSessions = sessions.map(session =>
                    session.id === currentSession.id ? updatedSession : session
                )

                set({
                    currentSession: updatedSession,
                    sessions: updatedSessions,
                })
            },

            // 创建新会话
            createSession: (title = '新对话') => {
                const newSession = {
                    id: uuidv4(),
                    title,
                    messages: [],
                    createdAt: new Date().toISOString(),
                    updatedAt: new Date().toISOString(),
                    tenantId: localStorage.getItem('tenantId') || 'default',
                }

                const newThreadId = uuidv4()

                set(state => ({
                    currentSession: newSession,
                    sessions: [newSession, ...state.sessions],
                    threadId: newThreadId,
                }))

                return newSession
            },

            // 删除会话
            deleteSession: (sessionId) => {
                const { sessions, currentSession } = get()
                const updatedSessions = sessions.filter(session => session.id !== sessionId)

                if (currentSession?.id === sessionId) {
                    set({
                        currentSession: updatedSessions[0] || null,
                        sessions: updatedSessions,
                        threadId: updatedSessions[0] ? get().threadId : uuidv4(),
                    })
                } else {
                    set({ sessions: updatedSessions })
                }
            },

            // 切换会话
            switchSession: (sessionId) => {
                const { sessions } = get()
                const session = sessions.find(s => s.id === sessionId)
                if (session) {
                    set({
                        currentSession: session,
                        threadId: uuidv4(), // 为新会话生成新的 thread ID
                    })
                }
            },

            // 设置加载状态
            setLoading: (loading) => set({ isLoading: loading }),

            // 设置连接状态
            setConnected: (connected) => set({ isConnected: connected }),

            // 设置错误
            setError: (error) => set({ error }),

            // 设置线程 ID
            setThreadId: (threadId) => set({ threadId }),

            // 清空当前会话
            clearCurrentSession: () => {
                const { createSession } = get()
                createSession('新对话')
            },
        }),
        {
            name: 'kimi-chat-store',
            storage: createJSONStorage(() => ({
                getItem: (name) => {
                    try {
                        return localStorage.getItem(name)
                    } catch {
                        return null
                    }
                },
                setItem: (name, value) => {
                    try {
                        localStorage.setItem(name, value)
                    } catch (e) {
                        // 存储满了，清除旧数据
                        console.warn('localStorage 已满，清除聊天历史')
                        localStorage.removeItem(name)
                    }
                },
                removeItem: (name) => {
                    try {
                        localStorage.removeItem(name)
                    } catch {
                        // ignore
                    }
                },
            })),
            partialize: (state) => {
                // 清理消息：移除图片，限制消息数量
                const cleanMessages = (messages) => 
                    messages.slice(-20).map(msg => ({ // 每个会话只保留最近20条
                        ...msg,
                        images: undefined,
                        sources: undefined, // 来源数据也可能很大
                    }))
                
                const cleanSessions = state.sessions.slice(0, 5).map(session => ({ // 只保留5个会话
                    ...session,
                    messages: cleanMessages(session.messages),
                }))
                
                const cleanCurrentSession = state.currentSession ? {
                    ...state.currentSession,
                    messages: cleanMessages(state.currentSession.messages),
                } : null
                
                return {
                    sessions: cleanSessions,
                    currentSession: cleanCurrentSession,
                    threadId: state.threadId,
                }
            },
        }
    )
)

export default useChatStore
