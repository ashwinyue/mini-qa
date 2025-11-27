/**
 * 认证 API 模块
 * 
 * 这个文件包含了与用户认证相关的所有 API 函数，包括：
 * - 登录：发送用户凭据并获取访问令牌
 * - 登出：使当前用户的访问令牌失效
 * - 获取当前用户信息：验证令牌并获取用户详情
 * 
 * 这些函数都依赖于统一的 HTTP 请求工具(request.js)来处理网络请求
 * 和错误处理
 * 
 * ⚠️ 重要说明：
 * 响应拦截器（在 request.js 中）已经提取了响应数据中的 data 字段，
 * 因此所有 API 函数直接返回 response 即可，无需再访问 response.data
 */

// 导入统一的 HTTP 请求工具
import request from '../utils/request'

/**
 * 用户登录
 * 
 * 向后端发送登录请求，验证用户的用户名和密码
 * 如果验证成功，返回包含访问令牌的用户信息
 * 
 * @param {Object} credentials - 登录凭据
 * @param {string} credentials.username - 用户名
 * @param {string} credentials.password - 密码
 * 
 * @returns {Promise<Object>} 包含用户信息和访问令牌的响应数据
 * 
 * @example
 * ```javascript
 * const response = await login({
 *   username: 'admin',
 *   password: 'password123'
 * })
 * console.log(response.token) // 访问令牌
 * console.log(response.user)  // 用户信息
 * ```
 */
export const login = async (credentials) => {
    try {
        /**
         * 发送 POST 请求到登录端点
         * 
         * request.post() 是对 axios.post() 的封装，包含了：
         * - 自动的请求头设置
         * - 请求/响应拦截器
         * - 统一的错误处理
         * - 请求/响应转换
         */
        const response = await request.post('/auth/login', credentials)

        /**
         * 响应数据格式（示例）：
         * {
         *   token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
         *   user: {
         *     id: 1,
         *     username: 'admin',
         *     name: '管理员',
         *     email: 'admin@example.com',
         *     role: 'admin'
         *   },
         *   expiresIn: 3600
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 如果登录失败，request.js 中的拦截器会处理通用错误，
         * 这里只需要记录日志或进行额外的处理
         */
        console.error('登录失败:', error)
        // 重新抛出错误，让调用者可以处理具体的登录失败原因
        throw error
    }
}

/**
 * 用户登出
 * 
 * 通知后端用户即将登出，使当前的访问令牌失效
 * 这是一个安全最佳实践，确保令牌不会被滥用
 * 
 * @returns {Promise<Object>} 登出操作的响应结果
 * 
 * @example
 * ```javascript
 * try {
 *   await logout()
 *   console.log('登出成功')
 * } catch (error) {
 *   console.error('登出失败:', error)
 * }
 * ```
 */
export const logout = async () => {
    try {
        /**
         * 发送 POST 请求到登出端点
         * 
         * 注意：即使登出请求失败，我们也应该清除本地的认证状态，
         * 因为用户的意图是明确的（要登出）
         */
        const response = await request.post('/auth/logout')

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   message: '登出成功'
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 即使后端登出失败，我们也应该记录日志，
         * 但这不应该阻止用户在前端登出
         */
        console.error('登出请求失败:', error)

        // 即使请求失败，也返回成功状态，
        // 因为用户的登出意图是明确的
        return { success: true, message: '本地登出完成' }
    }
}

/**
 * 获取当前用户信息
 * 
 * 向后端发送请求，验证当前访问令牌的有效性，
 * 并获取用户的详细信息
 * 
 * 这个函数通常在以下场景使用：
 * 1. 应用启动时验证用户是否已登录
 * 2. 刷新用户信息
 * 3. 验证令牌是否仍然有效
 * 
 * @returns {Promise<Object>} 当前用户的详细信息
 * 
 * @example
 * ```javascript
 * try {
 *   const userInfo = await getCurrentUser()
 *   console.log('当前用户:', userInfo.name)
 * } catch (error) {
 *   console.error('获取用户信息失败:', error)
 * }
 * ```
 */
export const getCurrentUser = async () => {
    try {
        /**
         * 发送 GET 请求到用户信息端点
         * 
         * request.get() 会自动在请求头中包含认证令牌
         * (在 request.js 的请求拦截器中设置)
         */
        const response = await request.get('/auth/profile')

        /**
         * 响应数据格式（示例）：
         * {
         *   id: 1,
         *   username: 'admin',
         *   name: '管理员',
         *   email: 'admin@example.com',
         *   role: 'admin',
         *   permissions: ['read:users', 'write:users'],
         *   lastLoginAt: '2024-01-15T10:30:00Z',
         *   avatar: 'https://example.com/avatar.jpg'
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 获取用户信息失败可能的原因：
         * 1. 访问令牌已过期
         * 2. 访问令牌无效
         * 3. 后端服务异常
         * 4. 网络连接问题
         * 
         * 这些错误通常表明用户需要重新登录
         */
        console.error('获取用户信息失败:', error)

        // 重新抛出错误，让调用者可以决定如何处理
        // (例如：重定向到登录页面)
        throw error
    }
}

/**
 * 刷新访问令牌
 * 
 * 使用刷新令牌获取新的访问令牌，延长用户会话
 * 
 * @param {string} refreshToken - 刷新令牌
 * 
 * @returns {Promise<Object>} 新的访问令牌和用户信息
 * 
 * @example
 * ```javascript
 * const response = await refreshToken(refreshToken)
 * console.log(response.token) // 新的访问令牌
 * ```
 */
export const refreshToken = async (refreshToken) => {
    try {
        /**
         * 发送 POST 请求到令牌刷新端点
         * 
         * 刷新令牌的目的是延长用户的会话时间，
         * 而不需要用户重新输入凭据
         */
        const response = await request.post('/auth/refresh', {
            refreshToken
        })

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   data: {
         *     token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
         *     refreshToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
         *     expiresIn: 3600
         *   }
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 刷新令牌失败通常意味着用户的会话已完全过期，
         * 需要重新登录
         */
        console.error('刷新令牌失败:', error)
        throw error
    }
}

export default {
    login,
    logout,
    getCurrentUser,
    refreshToken
}
