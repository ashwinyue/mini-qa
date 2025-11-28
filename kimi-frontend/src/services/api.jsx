/**
 * API 服务层
 * 
 * 配置 axios 实例，处理认证、租户头和统一错误处理
 */
import axios from 'axios'

// API 基础配置 - 使用 /api 前缀，通过 vite 代理转发到后端
const BASE_URL = '/api'
const TIMEOUT = 30000

// 创建 axios 实例
const apiClient = axios.create({
    baseURL: BASE_URL,
    timeout: TIMEOUT,
    headers: {
        'Content-Type': 'application/json',
    },
})

// 不需要租户头的简单 GET 路径
const SIMPLE_GET_PATHS = ['/health', '/models/list', '/greet']

/**
 * 请求拦截器
 * - 添加 Authorization 头
 * - 添加 X-Tenant-ID 头
 * - 添加 X-API-Key 头（向量库操作）
 */
apiClient.interceptors.request.use(
    (config) => {
        const method = (config.method || 'get').toLowerCase()
        const pathname = config.url || '/'
        const isSimpleGet = method === 'get' && SIMPLE_GET_PATHS.includes(pathname)

        // 添加认证 token
        const token = localStorage.getItem('auth_token')
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }

        // 非简单 GET 请求添加租户和 API Key
        if (!isSimpleGet) {
            const tenantId = localStorage.getItem('tenantId') || 'default'
            config.headers['X-Tenant-ID'] = tenantId

            const apiKey = localStorage.getItem('apiKey')
            if (apiKey) {
                config.headers['X-API-Key'] = apiKey
            }
        }

        // GET 请求不需要 Content-Type
        if (method === 'get') {
            delete config.headers['Content-Type']
        }

        return config
    },
    (error) => {
        console.error('[API] 请求拦截错误', error)
        return Promise.reject(error)
    }
)

/**
 * 响应拦截器
 * - 处理 401 未授权
 * - 统一错误处理
 */
apiClient.interceptors.response.use(
    (response) => {
        return response
    },
    (error) => {
        const status = error.response?.status

        // 401 未授权 - 清除认证
        if (status === 401) {
            localStorage.removeItem('auth_token')
            localStorage.removeItem('user_info')
            localStorage.removeItem('apiKey')
        }

        // 构造错误消息
        let errorMessage = '请求失败'
        if (status === 404) {
            errorMessage = '请求的资源不存在'
        } else if (status >= 500) {
            errorMessage = '服务器错误，请稍后重试'
        } else if (!error.response) {
            errorMessage = '网络连接失败，请检查网络设置'
        } else if (error.response?.data?.message) {
            errorMessage = error.response.data.message
        }

        // 将错误消息附加到 error 对象
        error.message = errorMessage

        return Promise.reject(error)
    }
)

/**
 * API 服务类
 */
class ApiService {
    /**
     * GET 请求
     */
    async get(url, config = {}) {
        return apiClient.get(url, config)
    }

    /**
     * POST 请求
     */
    async post(url, data = {}, config = {}) {
        return apiClient.post(url, data, config)
    }

    /**
     * PUT 请求
     */
    async put(url, data = {}, config = {}) {
        return apiClient.put(url, data, config)
    }

    /**
     * DELETE 请求
     */
    async delete(url, config = {}) {
        return apiClient.delete(url, config)
    }

    /**
     * 获取基础 URL（用于 SSE 等需要完整 URL 的场景）
     */
    getBaseURL() {
        return BASE_URL
    }
}

// 导出单例实例
export const apiService = new ApiService()
export default apiService
