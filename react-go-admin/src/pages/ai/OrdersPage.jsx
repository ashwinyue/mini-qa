/**
 * 订单查询页面
 */
import React, { useState } from 'react'
import { Input, Button, Card, Descriptions, Empty, Spin, message, Tag } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { getOrder, getOrderStatusInfo } from '../../api/orders'
import dayjs from 'dayjs'

const OrdersPage = () => {
    const [orderId, setOrderId] = useState('')
    const [loading, setLoading] = useState(false)
    const [order, setOrder] = useState(null)
    const [error, setError] = useState(null)
    const [searched, setSearched] = useState(false)

    const handleSearch = async () => {
        if (!orderId.trim()) {
            message.warning('请输入订单号')
            return
        }

        setLoading(true)
        setSearched(true)
        setError(null)
        setOrder(null)

        try {
            const orderData = await getOrder(orderId.trim())
            setOrder(orderData)
        } catch (err) {
            console.error('查询订单失败:', err)
            setError(err.message || '查询失败')
        } finally {
            setLoading(false)
        }
    }

    const handleKeyPress = (e) => {
        if (e.key === 'Enter') {
            handleSearch()
        }
    }

    const formatDate = (dateStr) => {
        if (!dateStr) return '-'
        return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
    }

    const formatAmount = (amount) => {
        if (amount === null || amount === undefined) return '-'
        return `¥ ${amount.toFixed(2)}`
    }

    return (
        <div className="p-4">
            <Card title="订单查询" className="mb-4">
                <div className="flex gap-2">
                    <Input
                        placeholder="请输入订单号..."
                        value={orderId}
                        onChange={(e) => setOrderId(e.target.value)}
                        onKeyPress={handleKeyPress}
                        size="large"
                        className="flex-1"
                        allowClear
                    />
                    <Button
                        type="primary"
                        icon={<SearchOutlined />}
                        size="large"
                        onClick={handleSearch}
                        loading={loading}
                    >
                        查询
                    </Button>
                </div>
            </Card>

            <Spin spinning={loading}>
                {searched && (
                    <Card title="查询结果">
                        {error ? (
                            <Empty
                                description={error}
                                image={Empty.PRESENTED_IMAGE_SIMPLE}
                            />
                        ) : order ? (
                            <Descriptions bordered column={2}>
                                <Descriptions.Item label="订单号" span={2}>
                                    {order.orderId}
                                </Descriptions.Item>
                                <Descriptions.Item label="订单状态">
                                    <Tag color={getOrderStatusInfo(order.status).color}>
                                        {getOrderStatusInfo(order.status).label}
                                    </Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="订单金额">
                                    <span className="text-lg font-semibold text-red-500">
                                        {formatAmount(order.amount)}
                                    </span>
                                </Descriptions.Item>
                                <Descriptions.Item label="更新时间">
                                    {formatDate(order.updatedAt)}
                                </Descriptions.Item>
                                <Descriptions.Item label="开始时间">
                                    {formatDate(order.startTime)}
                                </Descriptions.Item>
                                {order.enrollTime && (
                                    <Descriptions.Item label="报名时间" span={2}>
                                        {formatDate(order.enrollTime)}
                                    </Descriptions.Item>
                                )}
                            </Descriptions>
                        ) : (
                            <Empty description="请输入订单号进行查询" />
                        )}
                    </Card>
                )}
            </Spin>
        </div>
    )
}

export default OrdersPage
