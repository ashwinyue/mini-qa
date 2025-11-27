/**
 * HTTP 请求工具模块
 * 
 * 基于 Axios 封装的统一请求工具，提供：
 * 1. 请求拦截器：自动添加认证 token
 * 2. 响应拦截器：统一处理响应数据和错误
 * 3. 错误处理：自动处理常见 HTTP 错误
 * 4. 消息提示：统一的用户提示信息
 * 5. 配置管理：统一的 baseURL 和 timeout 设置
 */

import axios from 'axios'
import { message } from 'antd'
import { useAuthStore } from '../stores'

/**
 * 创建 Axios 实例
 * 
 * 使用 axios.create() 创建自定义配置的 Axios 实例，
 * 而不是直接使用默认实例，这样可以避免污染全局配置
 * 
 * 配置说明：
 * - baseURL: API 的基础 URL，所有请求都会自动添加此前缀
 * - timeout: 请求超时时间（毫秒），超过此时间请求将被取消
 */
const request = axios.create({
    // 基础 URL，通常指向后端 API 服务器
    baseURL: '/api',
    // 请求超时时间设置为 10 秒
    timeout: 10000,
})

/**
 * 请求拦截器
 * 
 * 拦截器是一种特殊的函数，会在请求发送之前自动执行
 * 这里主要用于在每个请求中添加认证 token
 * 
 * request.interceptors.request.use() 接受两个参数：
 * - 成功回调：请求发送前执行
 * - 错误回调：请求配置错误时执行
 */

request.interceptors.request.use(
    (config) => {
        // Get token from Zustand store
        const token = useAuthStore.getState().token

        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }

        return config
    },
    /**
     * 错误回调函数
     * 
     * 当请求配置有误时（如无效的 URL、错误的请求方法等）会执行
     * 这种错误通常在发送请求之前就会发生
     * 
     * @param {Object} error 错误对象
     * @returns {Promise} 返回一个 rejected 的 Promise
     */
    (error) => {
        // 直接返回错误，调用者可以捕获到这个错误
        return Promise.reject(error)
    }
)

/**
 * 响应拦截器
 * 
 * 拦截器会在响应数据到达客户端之前对其进行处理
 * 这里主要用于：
 * 1. 统一处理业务状态码
 * 2. 提取实际数据
 * 3. 处理错误响应
 * 
 * request.interceptors.response.use() 接受两个参数：
 * - 成功回调：服务器正常响应时执行
 * - 错误回调：服务器响应错误时执行
 */
request.interceptors.response.use(
    /**
     * 成功回调函数
     * 
     * @param {Object} response 响应对象
     * @returns {any} 返回处理后的数据
     */
    (response) => {
        /**
         * 解构响应数据
         * 
         * 假设服务器返回的数据结构为：
         * {
         *   code: 200,           // 业务状态码
         *   data: {...},         // 实际数据
         *   msg: "操作成功"       // 提示消息
         * }
         */
        const { code, data, msg } = response.data

        /**
         * 根据业务状态码处理响应
         * 
         * 常见的业务状态码约定：
         * - 200: 成功
         * - 400: 请求参数错误
         * - 401: 未授权
         * - 403: 权限不足
         * - 404: 资源不存在
         * - 500: 服务器内部错误
         */
        if (code === 200) {
            // 成功时返回实际数据，调用者可以直接使用
            return data
        } else {
            // 非 200 状态码，视为业务错误
            // 显示错误消息
            message.error(msg || '请求失败')
            // 抛出错误，让调用者可以捕获
            return Promise.reject(new Error(msg || '请求失败'))
        }
    },
    /**
     * 错误回调函数
     * 
     * 当服务器返回错误状态码（如 404, 500 等）或网络错误时会执行
     * 这里根据不同的错误类型提供不同的处理逻辑
     * 
     * @param {Object} error 错误对象
     * @returns {Promise} 返回一个 rejected 的 Promise
     */
    (error) => {
        /**
         * HTTP 错误处理
         * 
         * Axios 错误对象结构：
         * - error.response: 服务器响应对象（存在时表示服务器返回了响应）
         * - error.request: 请求对象（存在但无响应时表示网络错误）
         * - error.message: 错误消息
         */
        if (error.response) {
            /**
             * 服务器返回了响应，但状态码不是 2xx
             * error.response 包含完整的响应信息
             */
            const { status, data } = error.response

            /**
             * 根据 HTTP 状态码进行不同处理
             * 
             * 常见 HTTP 状态码：
             * - 400: Bad Request 请求参数错误
             * - 401: Unauthorized 未授权（token 无效或过期）
             * - 403: Forbidden 权限不足
             * - 404: Not Found 资源不存在
             * - 500: Internal Server Error 服务器内部错误
             */
            switch (status) {
                case 401:
                    // 401 表示未授权，通常是 token 过期或无效
                    message.error('未授权，请重新登录')
                    // 清除本地存储的 token
                    localStorage.removeItem('token')
                    // 跳转到登录页面
                    window.location.href = '/login'
                    break

                case 403:
                    // 403 表示权限不足，用户没有访问资源的权限
                    message.error('拒绝访问')
                    break

                case 404:
                    // 404 表示请求的资源不存在
                    message.error('请求的资源不存在')
                    break

                case 500:
                    // 500 表示服务器内部错误
                    message.error('服务器错误')
                    break

                default:
                    // 其他状态码，使用服务器返回的错误消息
                    message.error(data?.msg || '请求失败')
            }
        } else if (error.request) {
            /**
             * 请求已发送，但没有收到响应
             * 通常表示网络错误，如：
             * - 服务器宕机
             * - 网络连接中断
             * - CORS 错误
             */
            message.error('网络错误，请检查网络连接')
        } else {
            /**
             * 请求配置错误
             * 这种错误在请求发送之前就会发生，如：
             * - 无效的 URL
             * - 错误的请求方法
             * - 无效的请求配置
             */
            message.error('请求配置错误')
        }

        /**
         * 返回 rejected 的 Promise
         * 
         * 这让调用者可以使用 .catch() 捕获错误并进行处理
         * 同时阻止了 JavaScript 未处理的 Promise 拒绝警告
         */
        return Promise.reject(error)
    }
)

// 导出配置好的 request 实例
// 使用者可以调用 request({ url, method, data }) 发起请求
export default request
