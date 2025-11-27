/**
 * 认证 API 模块
 * 
 * 调用 Python 后端的认证接口
 */
import apiService from '../services/api'

// Token 和用户信息的存储键
const TOKEN_KEY = 'auth_token'
const USER_KEY = 'user_info'

/**
 * 用户登录
 * 
 * @param {Object} credentials - 登录凭据
 * @param {string} credentials.username - 用户名
 * @param {string} credentials.password - 密码
 * @returns {Promise<Object>} 包含 token 和用户信息
 */
export const login = async (credentials) => {
    try {
        const response = await apiService.post('/api/auth/login', credentials)
        const data = response.data

        // Python 后端返回格式: { code: 0, message: 'OK', data: { token, user, expiresIn } }
        if (data.code === 0 && data.data) {
            const { token, user, expiresIn } = data.data
            
            // 存储 token 和用户信息
            localStorage.setItem(TOKEN_KEY, token)
            localStorage.setItem(USER_KEY, JSON.stringify(user))
            
            return {
                token,
                user,
                expiresIn,
            }
        }

        throw new Error(data.message || '登录失败')
    } catch (error) {
        console.error('登录失败:', error)
        throw error
    }
}

/**
 * 用户登出
 * 
 * @returns {Promise<Object>} 登出结果
 */
export const logout = async () => {
    try {
        await apiService.post('/api/auth/logout')
    } catch (error) {
        console.error('登出请求失败:', error)
    } finally {
        // 无论请求是否成功，都清除本地存储
        clearAuth()
    }
    
    return { success: true, message: '登出成功' }
}

/**
 * 获取当前用户信息
 * 
 * @returns {Promise<Object>} 用户信息
 */
export const getCurrentUser = async () => {
    try {
        const response = await apiService.get('/api/auth/me')
        const data = response.data

        // Python 后端返回格式: { code: 0, message: 'OK', data: { username, nickname, role, email } }
        if (data.code === 0 && data.data) {
            const user = data.data
            // 更新本地存储的用户信息
            localStorage.setItem(USER_KEY, JSON.stringify(user))
            return user
        }

        throw new Error(data.message || '获取用户信息失败')
    } catch (error) {
        console.error('获取用户信息失败:', error)
        // 清除可能过期的认证信息
        clearAuth()
        throw error
    }
}

/**
 * 获取存储的 token
 * 
 * @returns {string|null} token
 */
export const getToken = () => {
    return localStorage.getItem(TOKEN_KEY)
}

/**
 * 获取存储的用户信息
 * 
 * @returns {Object|null} 用户信息
 */
export const getUser = () => {
    const userStr = localStorage.getItem(USER_KEY)
    if (userStr) {
        try {
            return JSON.parse(userStr)
        } catch {
            return null
        }
    }
    return null
}

/**
 * 清除认证信息
 */
export const clearAuth = () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
}

/**
 * 检查是否已登录
 * 
 * @returns {boolean}
 */
export const isAuthenticated = () => {
    return !!getToken()
}

export default {
    login,
    logout,
    getCurrentUser,
    getToken,
    getUser,
    clearAuth,
    isAuthenticated,
}
