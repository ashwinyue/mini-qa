/**
 * 系统 API 模块
 * 
 * 调用 Python 后端的系统接口
 */
import apiService from '../services/api'

/**
 * 健康检查
 * 
 * @returns {Promise<Object>} 系统健康状态
 */
export const getHealth = async () => {
    try {
        const response = await apiService.get('/health')
        return response.data
    } catch (error) {
        console.error('健康检查失败:', error)
        throw error
    }
}

/**
 * 获取模型列表
 * 
 * @returns {Promise<Object>} 模型信息
 */
export const getModels = async () => {
    try {
        const response = await apiService.get('/models/list')
        const data = response.data

        // 返回格式: { code: 0, message: 'OK', data: { current, models } }
        if (data.code === 0 && data.data) {
            return {
                current: data.data.current,
                models: data.data.models || [],
            }
        }

        return {
            current: '',
            models: [],
        }
    } catch (error) {
        console.error('获取模型列表失败:', error)
        throw error
    }
}

/**
 * 切换模型
 * 
 * @param {string} modelName - 模型名称
 * @returns {Promise<Object>} 切换结果
 */
export const switchModel = async (modelName) => {
    try {
        const response = await apiService.post('/models/switch', {
            name: modelName,
        })
        const data = response.data

        if (data.code === 0 && data.data) {
            return {
                success: true,
                current: data.data.current,
                models: data.data.models || [],
            }
        }

        throw new Error(data.message || '切换模型失败')
    } catch (error) {
        console.error('切换模型失败:', error)
        throw error
    }
}

export default {
    getHealth,
    getModels,
    switchModel,
}
