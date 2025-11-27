/**
 * 对话 API 模块
 * 
 * 调用 Python 后端的对话接口
 */
import apiService from '../services/api'

/**
 * 发送聊天消息
 * 
 * @param {string} query - 用户消息
 * @param {string} threadId - 会话线程 ID
 * @param {string[]} images - 图片数组 (Base64)
 * @param {string} audio - 音频数据 (Base64)
 * @returns {Promise<Object>} 对话响应
 */
export const sendMessage = async (query, threadId, images, audio) => {
    try {
        const response = await apiService.post('/chat', {
            query,
            thread_id: threadId,
            images,
            audio,
        })

        const data = response.data

        // 处理命令响应 (/help, /history, /reset)
        if (data.commands || data.history || data.reset !== undefined) {
            return data
        }

        // 标准对话响应: { route, answer, sources }
        return {
            route: data.route,
            answer: data.answer,
            sources: data.sources || [],
        }
    } catch (error) {
        console.error('发送消息失败:', error)
        throw error
    }
}

/**
 * 获取欢迎语
 * 
 * @returns {Promise<Object>} 欢迎语和选项
 */
export const getGreeting = async () => {
    try {
        const response = await apiService.get('/greet')
        return response.data
    } catch (error) {
        console.error('获取欢迎语失败:', error)
        // 返回默认欢迎语
        return {
            message: '您好，请问有什么可以帮您？',
            options: [],
        }
    }
}

/**
 * 建立建议问题 SSE 连接
 * 
 * @param {string} threadId - 会话线程 ID
 * @param {Function} onMessage - 消息回调
 * @param {Function} onError - 错误回调
 * @returns {EventSource} SSE 连接
 */
export const getSuggestions = (threadId, onMessage, onError) => {
    const baseURL = apiService.getBaseURL()
    const eventSource = new EventSource(`${baseURL}/suggest/${threadId}`)

    eventSource.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data)
            onMessage && onMessage(data)
        } catch (error) {
            console.error('解析建议数据失败:', error)
        }
    }

    eventSource.addEventListener('react_start', (event) => {
        try {
            const data = JSON.parse(event.data)
            onMessage && onMessage({ ...data, event: 'react_start' })
        } catch (error) {
            console.error('解析 react_start 数据失败:', error)
        }
    })

    eventSource.addEventListener('react', (event) => {
        try {
            const data = JSON.parse(event.data)
            onMessage && onMessage({ ...data, event: 'react' })
            // 收到最终建议后关闭连接
            if (data.final) {
                eventSource.close()
            }
        } catch (error) {
            console.error('解析 react 数据失败:', error)
        }
    })

    eventSource.addEventListener('error', (event) => {
        try {
            const data = JSON.parse(event.data)
            onError && onError(data)
        } catch {
            onError && onError({ message: '建议生成失败' })
        }
        eventSource.close()
    })

    eventSource.onerror = (error) => {
        console.error('SSE 连接错误:', error)
        onError && onError({ message: 'SSE 连接失败' })
        eventSource.close()
    }

    return eventSource
}

/**
 * 执行命令 (/help, /history, /reset)
 * 
 * @param {string} command - 命令
 * @param {string} threadId - 会话线程 ID
 * @returns {Promise<Object>} 命令响应
 */
export const executeCommand = async (command, threadId) => {
    return sendMessage(command, threadId)
}

export default {
    sendMessage,
    getGreeting,
    getSuggestions,
    executeCommand,
}
