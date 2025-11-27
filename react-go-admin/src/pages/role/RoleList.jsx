/**
 * 角色管理页面组件
 * 
 * 这是一个完整的 CRUD (增删改查) 功能页面，用于管理系统角色
 * 
 * 主要功能：
 * 1. 角色列表展示（分页）
 * 2. 新增角色
 * 3. 编辑角色
 * 4. 删除角色（带确认）
 * 
 * 技术要点：
 * - 使用 Ant Design 的 Table 组件展示数据
 * - 使用 Modal 组件实现弹窗表单
 * - 使用 Form 组件进行表单验证
 * - 使用 Popconfirm 组件实现删除确认
 * - 使用 useState 管理组件状态
 * - 使用 useEffect 在组件加载时获取数据
 */

// 导入 React Hooks
import { useState, useEffect } from 'react'
// useState: 用于管理组件的状态（如加载状态、数据列表等）
// useEffect: 用于处理副作用（如数据加载、订阅等）

// 导入 Ant Design 组件
import { Table, Button, Space, Modal, Form, Input, message, Popconfirm } from 'antd'
// Table: 表格组件，用于展示角色列表
// Button: 按钮组件，用于触发操作
// Space: 间距组件，用于控制元素之间的间距
// Modal: 对话框组件，用于显示新增/编辑表单
// Form: 表单组件，用于数据输入和验证
// Input: 输入框组件，用于文本输入
// message: 消息提示组件，用于显示操作结果
// Popconfirm: 气泡确认框，用于删除前的二次确认

// 导入 Ant Design 图标
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
// PlusOutlined: 加号图标，用于新增按钮
// EditOutlined: 编辑图标，用于编辑按钮
// DeleteOutlined: 删除图标，用于删除按钮

// 导入角色相关的 API 接口
import { getRoleList, createRole, updateRole, deleteRole } from '../../api/role'
// getRoleList: 获取角色列表（支持分页）
// createRole: 创建新角色
// updateRole: 更新角色信息
// deleteRole: 删除角色

/**
 * 角色列表组件
 * 
 * 这是一个完整的角色管理页面，包含列表展示和 CRUD 操作
 * 
 * @component
 */
const RoleList = () => {
    /**
     * 状态管理
     * 
     * 使用 useState Hook 管理组件的各种状态
     * 每个状态都有对应的更新函数
     */

    // 加载状态：控制表格的 loading 效果
    const [loading, setLoading] = useState(false)

    // 数据源：存储角色列表数据
    const [dataSource, setDataSource] = useState([])

    // 分页信息：当前页、每页条数、总条数
    const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })

    // 对话框显示状态：控制新增/编辑对话框的显示/隐藏
    const [modalVisible, setModalVisible] = useState(false)

    // 正在编辑的角色：null 表示新增，有值表示编辑
    const [editingRole, setEditingRole] = useState(null)

    // 表单实例：用于获取和设置表单值
    // Form.useForm() 返回一个数组，第一个元素是表单实例
    const [form] = Form.useForm()

    /**
     * 加载角色列表数据
     * 
     * 这是一个异步函数，从后端 API 获取角色列表
     * 
     * @param {number} page - 当前页码，默认为 1
     * @param {number} pageSize - 每页显示的条数，默认为 10
     * 
     * 工作流程：
     * 1. 设置 loading 状态为 true（显示加载动画）
     * 2. 调用 API 获取数据
     * 3. 更新数据源和分页信息
     * 4. 处理错误（如果有）
     * 5. 设置 loading 状态为 false（隐藏加载动画）
     */
    const loadData = async (page = 1, pageSize = 10) => {
        setLoading(true)  // 开始加载，显示 loading 动画
        try {
            // 调用 API 获取角色列表
            // 传入分页参数：page（页码）和 pageSize（每页条数）
            const data = await getRoleList({ page, pageSize })

            // 更新数据源
            // data.list 是角色列表数组，如果为空则使用空数组
            setDataSource(data.list || [])

            // 更新分页信息
            setPagination({
                current: page,              // 当前页码
                pageSize,                   // 每页条数
                total: data.total || 0,     // 总条数（用于计算总页数）
            })
        } catch (error) {
            // 捕获并记录错误
            // 在实际应用中，这里可以显示错误提示
            console.error('加载角色列表失败:', error)
        } finally {
            // 无论成功还是失败，都要关闭 loading 状态
            setLoading(false)
        }
    }

    /**
     * 组件挂载时的副作用
     * 
     * useEffect Hook 用于处理副作用操作
     * 第二个参数是依赖数组，空数组表示只在组件挂载时执行一次
     * 
     * 这里的作用：页面加载时自动获取第一页数据
     */
    useEffect(() => {
        loadData()  // 加载第一页数据
    }, [])  // 空依赖数组，只在组件挂载时执行一次

    /**
     * 表格列配置
     * 
     * 定义表格的列结构，包括：
     * - title: 列标题
     * - dataIndex: 数据字段名（对应后端返回的字段）
     * - key: 列的唯一标识
     * - width: 列宽度（可选）
     * - render: 自定义渲染函数（可选）
     */
    const columns = [
        {
            title: 'ID',              // 列标题
            dataIndex: 'id',          // 对应数据中的 id 字段
            key: 'id',                // 列的唯一标识
            width: 80,                // 列宽度
        },
        {
            title: '角色名称',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '角色标识',
            dataIndex: 'code',
            key: 'code',
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
        },
        {
            title: '创建时间',
            dataIndex: 'created_at',
            key: 'created_at',
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            /**
             * render 函数：自定义列的渲染内容
             * 
             * @param {any} _ - 当前单元格的值（这里不需要，用 _ 表示忽略）
             * @param {Object} record - 当前行的完整数据对象
             * @returns {JSX.Element} 渲染的内容
             */
            render: (_, record) => (
                // Space 组件：控制按钮之间的间距
                <Space>
                    {/* 编辑按钮 */}
                    <Button
                        type="link"                      // 链接样式按钮
                        size="small"                     // 小尺寸
                        icon={<EditOutlined />}          // 编辑图标
                        onClick={() => handleEdit(record)}  // 点击时调用编辑函数
                    >
                        编辑
                    </Button>

                    {/* 删除按钮（带确认） */}
                    <Popconfirm
                        title="确定要删除这个角色吗?"     // 确认提示文本
                        onConfirm={() => handleDelete(record.id)}  // 确认后执行删除
                        okText="确定"                     // 确认按钮文本
                        cancelText="取消"                 // 取消按钮文本
                    >
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ]

    /**
     * 处理表格分页、排序、筛选变化
     * 
     * 当用户切换页码或改变每页条数时，Table 组件会调用这个函数
     * 
     * @param {Object} newPagination - 新的分页信息
     */
    const handleTableChange = (newPagination) => {
        // 重新加载数据，使用新的页码和每页条数
        loadData(newPagination.current, newPagination.pageSize)
    }

    /**
     * 打开新增对话框
     * 
     * 工作流程：
     * 1. 清空正在编辑的角色（表示这是新增操作）
     * 2. 重置表单字段（清空之前的输入）
     * 3. 显示对话框
     */
    const handleAdd = () => {
        setEditingRole(null)      // 设置为 null，表示新增模式
        form.resetFields()        // 清空表单
        setModalVisible(true)     // 显示对话框
    }

    /**
     * 打开编辑对话框
     * 
     * @param {Object} record - 要编辑的角色数据
     * 
     * 工作流程：
     * 1. 设置正在编辑的角色（用于判断是新增还是编辑）
     * 2. 将角色数据填充到表单中
     * 3. 显示对话框
     */
    const handleEdit = (record) => {
        setEditingRole(record)           // 保存正在编辑的角色
        form.setFieldsValue(record)      // 将角色数据填充到表单
        setModalVisible(true)            // 显示对话框
    }

    /**
     * 处理删除操作
     * 
     * @param {number} id - 要删除的角色 ID
     * 
     * 工作流程：
     * 1. 调用删除 API
     * 2. 显示成功提示
     * 3. 重新加载当前页数据
     * 4. 如果失败，记录错误
     */
    const handleDelete = async (id) => {
        try {
            // 调用删除 API
            await deleteRole(id)

            // 显示成功提示
            message.success('删除成功')

            // 重新加载当前页数据
            // 保持在当前页，避免跳回第一页
            loadData(pagination.current, pagination.pageSize)
        } catch (error) {
            // 记录错误
            console.error('删除失败:', error)
        }
    }

    /**
     * 处理表单提交
     * 
     * 根据 editingRole 的值判断是新增还是编辑
     * 
     * 工作流程：
     * 1. 验证表单字段
     * 2. 根据模式调用新增或更新 API
     * 3. 显示成功提示
     * 4. 关闭对话框
     * 5. 重新加载数据
     */
    const handleSubmit = async () => {
        try {
            // 验证表单并获取表单值
            // 如果验证失败，会抛出异常，不会继续执行
            const values = await form.validateFields()

            // 判断是编辑还是新增
            if (editingRole) {
                // 编辑模式：调用更新 API
                await updateRole(editingRole.id, values)
                message.success('更新成功')
            } else {
                // 新增模式：调用创建 API
                await createRole(values)
                message.success('创建成功')
            }

            // 关闭对话框
            setModalVisible(false)

            // 重新加载数据
            loadData(pagination.current, pagination.pageSize)
        } catch (error) {
            // 捕获错误（可能是验证失败或 API 调用失败）
            console.error('提交失败:', error)
        }
    }

    /**
     * 渲染组件
     * 
     * 页面结构：
     * 1. 顶部：标题和新增按钮
     * 2. 中间：角色列表表格
     * 3. 底部：新增/编辑对话框（默认隐藏）
     */
    return (
        <div>
            {/* 顶部操作栏 */}
            <div className="mb-4 flex justify-between items-center">
                <h1 className="text-2xl font-bold">角色管理</h1>
                {/* 新增按钮 */}
                <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                    新增角色
                </Button>
            </div>

            {/* 角色列表表格 */}
            <Table
                columns={columns}              // 列配置
                dataSource={dataSource}        // 数据源
                loading={loading}              // 加载状态
                rowKey="id"                    // 行的唯一标识（必须指定）
                pagination={pagination}        // 分页配置
                onChange={handleTableChange}   // 分页变化时的回调
            />

            {/* 新增/编辑对话框 */}
            <Modal
                title={editingRole ? '编辑角色' : '新增角色'}  // 根据模式显示不同标题
                open={modalVisible}                           // 控制显示/隐藏
                onOk={handleSubmit}                           // 点击确定时的回调
                onCancel={() => setModalVisible(false)}       // 点击取消时的回调
                width={600}                                   // 对话框宽度
            >
                {/* 表单 */}
                <Form form={form} layout="vertical">
                    {/* 角色名称字段 */}
                    <Form.Item
                        name="name"                           // 字段名（对应提交数据的 key）
                        label="角色名称"                       // 字段标签
                        rules={[{ required: true, message: '请输入角色名称' }]}  // 验证规则
                    >
                        <Input placeholder="请输入角色名称" />
                    </Form.Item>

                    {/* 角色标识字段 */}
                    <Form.Item
                        name="code"
                        label="角色标识"
                        rules={[{ required: true, message: '请输入角色标识' }]}
                    >
                        <Input placeholder="请输入角色标识，如：admin" />
                    </Form.Item>

                    {/* 描述字段（非必填） */}
                    <Form.Item name="description" label="描述">
                        <Input.TextArea rows={4} placeholder="请输入角色描述" />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

// 导出组件，供路由使用
export default RoleList
