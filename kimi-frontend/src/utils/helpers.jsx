/**
 * 工具函数
 */

/**
 * 格式化时间戳
 * @param {string} timestamp - ISO 时间戳
 * @returns {string} 格式化后的时间
 */
export function formatTime(timestamp) {
    const date = new Date(timestamp)
    const now = new Date()
    const diff = now - date
    
    // 小于1分钟
    if (diff < 60000) {
        return '刚刚'
    }
    
    // 小于1小时
    if (diff < 3600000) {
        return `${Math.floor(diff / 60000)}分钟前`
    }
    
    // 小于24小时
    if (diff < 86400000) {
        return `${Math.floor(diff / 3600000)}小时前`
    }
    
    // 今天
    if (date.toDateString() === now.toDateString()) {
        return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }
    
    // 昨天
    const yesterday = new Date(now)
    yesterday.setDate(yesterday.getDate() - 1)
    if (date.toDateString() === yesterday.toDateString()) {
        return '昨天 ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }
    
    // 其他
    return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

/**
 * 手机号脱敏
 * @param {string} phone - 手机号
 * @returns {string} 脱敏后的手机号
 */
export function maskPhone(phone) {
    if (!phone || phone.length < 11) {
        return phone
    }
    return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/**
 * 生成随机 ID
 * @returns {string} 随机 ID
 */
export function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
export function debounce(func, wait) {
    let timeout
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout)
            func(...args)
        }
        clearTimeout(timeout)
        timeout = setTimeout(later, wait)
    }
}

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} limit - 时间限制（毫秒）
 * @returns {Function} 节流后的函数
 */
export function throttle(func, limit) {
    let inThrottle
    return function executedFunction(...args) {
        if (!inThrottle) {
            func(...args)
            inThrottle = true
            setTimeout(() => inThrottle = false, limit)
        }
    }
}

/**
 * 复制文本到剪贴板
 * @param {string} text - 要复制的文本
 * @returns {Promise<boolean>} 是否成功
 */
export async function copyToClipboard(text) {
    try {
        await navigator.clipboard.writeText(text)
        return true
    } catch (err) {
        console.error('复制失败:', err)
        return false
    }
}
