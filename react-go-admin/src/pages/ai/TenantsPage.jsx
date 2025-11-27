/**
 * 租户管理页面
 */
import React, { useState, useEffect } from 'react'
import { Card, Input, Button, List, Tag, message, Space, Modal } from 'antd'
import { PlusOutlined, CheckOutlined, DeleteOutlined } from '@ant-design/icons'
import { useSystemStore } from '../../stores'

const TenantsPage = () => {
    const { config, setTenantId, getTenantId } = useSystemStore()
    const [tenants, setTenants] = useState(['default'])
    const [newTenant, setNewTenant] = useState('')
    const [currentTenant, setCurrentTenant] = useState('default')

    useEffect(() => {
        // 从 localStorage 加载租户列表
        const savedTenants = localStorage.getItem('tenantList')
        if (savedTenants) {
            try {
                const parsed = JSON.parse(savedTenants)
                setTenants(parsed)
            } catch {
                setTenants(['default'])
            }
        }

        // 获取当前租户
        const current = getTenantId()
        setCurrentTenant(current)
    }, [])

    const saveTenants = (newList) => {
        localStorage.setItem('tenantList', JSON.stringify(newList))
        setTenants(newList)
    }

    const handleAddTenant = () => {
        const trimmed = newTenant.trim()
        if (!trimmed) {
            message.warning('请输入租户 ID')
            return
        }

        if (tenants.includes(trimmed)) {
            message.warning('租户已存在')
            return
        }

        const newList = [...tenants, trimmed]
        saveTenants(newList)
        setNewTenant('')
        message.success('租户添加成功')
    }

    const handleSwitchTenant = (tenantId) => {
        setTenantId(tenantId)
        setCurrentTenant(tenantId)
        message.success(`已切换到租户: ${tenantId}`)
    }

    const handleDeleteTenant = (tenantId) => {
        if (tenantId === 'default') {
            message.warning('默认租户不能删除')
            return
        }

        Modal.confirm({
            title: '确认删除',
            content: `确定要删除租户 "${tenantId}" 吗？`,
            onOk: () => {
                const newList = tenants.filter(t => t !== tenantId)
                saveTenants(newList)

                // 如果删除的是当前租户，切换到 default
                if (currentTenant === tenantId) {
                    handleSwitchTenant('default')
                }

                message.success('租户删除成功')
            },
        })
    }

    return (
        <div className="p-4">
            <Card title="当前租户" className="mb-4">
                <div className="flex items-center gap-4">
                    <span className="text-gray-500">当前使用的租户：</span>
                    <Tag color="blue" className="text-lg px-4 py-1">
                        {currentTenant}
                    </Tag>
                </div>
            </Card>

            <Card title="添加租户" className="mb-4">
                <div className="flex gap-2">
                    <Input
                        placeholder="输入新租户 ID..."
                        value={newTenant}
                        onChange={(e) => setNewTenant(e.target.value)}
                        onPressEnter={handleAddTenant}
                        className="flex-1"
                    />
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={handleAddTenant}
                    >
                        添加
                    </Button>
                </div>
            </Card>

            <Card title="租户列表">
                <List
                    dataSource={tenants}
                    renderItem={(tenant) => (
                        <List.Item
                            actions={[
                                tenant === currentTenant ? (
                                    <Tag color="green" key="current">
                                        <CheckOutlined /> 当前
                                    </Tag>
                                ) : (
                                    <Button
                                        key="switch"
                                        type="link"
                                        onClick={() => handleSwitchTenant(tenant)}
                                    >
                                        切换
                                    </Button>
                                ),
                                tenant !== 'default' && (
                                    <Button
                                        key="delete"
                                        type="link"
                                        danger
                                        icon={<DeleteOutlined />}
                                        onClick={() => handleDeleteTenant(tenant)}
                                    >
                                        删除
                                    </Button>
                                ),
                            ].filter(Boolean)}
                        >
                            <List.Item.Meta
                                title={
                                    <span className="flex items-center gap-2">
                                        {tenant}
                                        {tenant === 'default' && (
                                            <Tag color="default">默认</Tag>
                                        )}
                                    </span>
                                }
                            />
                        </List.Item>
                    )}
                />
            </Card>
        </div>
    )
}

export default TenantsPage
