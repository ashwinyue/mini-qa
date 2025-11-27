/**
 * 订单 API 模块
 * 
 * 调用 Python 后端的订单查询接口
 */
import apiService from '../services/api'

/**
 * 查询订单详情
 * 
 * @param {string} orderId - 订单 ID
 * @returns {Promise<Object>} 订单信息
 */
export const getOrder = async (orderId) => {
    try {
        const response = await apiService.get(`/api/orders/${orderId}`)
        const data = response.data

        // 返回订单信息
        return {
            orderId: data.order_id,
            status: data.status,
            amount: data.amount,
            updatedAt: data.updated_at,
            enrollTime: data.enroll_time,
            startTime: data.start_time,
        }
    } catch (error) {
        console.error('查询订单失败:', error)
        
        // 处理 404 错误
        if (error.response?.status === 404) {
            throw new Error('订单不存在')
        }
        
        throw error
    }
}

/**
 * 订单状态映射
 */
export const ORDER_STATUS = {
    PENDING: { label: '待支付', color: 'orange' },
    PAID: { label: '已支付', color: 'green' },
    CANCELLED: { label: '已取消', color: 'red' },
    REFUNDED: { label: '已退款', color: 'gray' },
    COMPLETED: { label: '已完成', color: 'blue' },
}

/**
 * 获取订单状态显示信息
 * 
 * @param {string} status - 订单状态
 * @returns {Object} 状态显示信息
 */
export const getOrderStatusInfo = (status) => {
    return ORDER_STATUS[status] || { label: status, color: 'default' }
}

export default {
    getOrder,
    ORDER_STATUS,
    getOrderStatusInfo,
}
