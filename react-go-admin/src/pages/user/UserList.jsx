/**
 * 用户管理页面组件
 * 
 * 这个组件是系统的核心功能页面之一，负责管理所有用户账户
 * 提供了完整的用户生命周期管理功能
 * 
 * 主要功能：
 * 1. 用户列表展示：显示所有用户的详细信息
 * 2. 分页查询：支持大数据量的分页展示
 * 3. 新增用户：创建新的用户账户
 * 4. 编辑用户：修改现有用户信息
 * 5. 删除用户：移除用户账户（物理删除）
 * 6. 状态管理：用户账户的启用/禁用状态
 * 
 * 技术特点：
 * - 使用 React Hooks 进行状态管理
 * - 集成 Ant Design 组件库
 * - 响应式设计，适配不同屏幕尺寸
 * - 表单验证和错误处理
 * - 异步数据加载和错误处理
 * 
 * @component
 */

// 导入 React 和 Hooks
import { useState, useEffect } from 'react'

// 导入 Ant Design 组件
import {
    Table,           // 表格组件：展示用户列表数据
    Button,          // 按钮组件：操作按钮
    Space,           // 间距组件：控制元素间距
    Modal,           // 模态框组件：用户新增/编辑对话框
    Form,            // 表单组件：表单验证和提交
    Input,           // 输入框组件：表单字段输入
    message,         // 消息提示：操作反馈
    Popconfirm,      // 确认弹框：删除确认对话框
    Tag              // 标签组件：用户状态显示
} from 'antd'

// 导入 Ant Design 图标
import {
    PlusOutlined,        // 新增图标
    EditOutlined,        // 编辑图标
    DeleteOutlined       // 删除图标
} from '@ant-design/icons'

// 导入用户相关的 API 函数
import {
    getUserList,     // 获取用户列表 API
    createUser,      // 创建用户 API
    updateUser,      // 更新用户 API
    deleteUser       // 删除用户 API
} from '../../api/user'

/**
 * 用户列表页面组件
 * 
 * 这是用户管理的核心页面，提供完整的用户管理功能
 * 包含数据展示、新增、编辑、删除等操作
 * 
 * 组件状态：
 * - loading: 数据加载状态
 * - dataSource: 表格数据源
 * - pagination: 分页配置
 * - modalVisible: 模态框显示状态
 * - editingUser: 当前编辑的用户对象
 * - form: 表单实例
 */
const UserList = () => {
    /**
     * 状态管理
     * 
     * loading: 控制表格加载状态，显示加载动画
     * dataSource: 存储用户列表数据，在表格中展示
     * pagination: 分页配置，包含当前页、每页条数、总数等
     * modalVisible: 控制新增/编辑用户对话框的显示/隐藏
     * editingUser: 存储当前编辑的用户对象，null 表示新增模式
     * form: Ant Design 表单实例，用于表单验证和提交
     */
    const [loading, setLoading] = useState(false)
    const [dataSource, setDataSource] = useState([])
    const [pagination, setPagination] = useState({
        current: 1,     // 当前页码
        pageSize: 10,   // 每页显示条数
        total: 0        // 总记录数
    })
    const [modalVisible, setModalVisible] = useState(false)           // 模态框显示状态
    const [editingUser, setEditingUser] = useState(null)             // 编辑的用户对象
    const [form] = Form.useForm()                                    // 表单实例

    /**
     * 加载用户列表数据
     * 
     * 异步函数，负责从服务器获取用户数据
     * 支持分页参数，可以根据页码和每页条数加载数据
     * 
     * @param {number} page - 页码，默认为 1
     * @param {number} pageSize - 每页条数，默认为 10
     * 
     * 流程：
     * 1. 设置加载状态为 true
     * 2. 调用 API 获取用户列表
     * 3. 更新数据源和分页状态
     * 4. 处理错误情况
     * 5. 设置加载状态为 false
     */
    const loadData = async (page = 1, pageSize = 10) => {
        setLoading(true)  // 开始加载，显示加载状态
        try {
            // 调用获取用户列表 API
            const data = await getUserList({ page, pageSize })

            // 更新数据源状态
            setDataSource(data.list || [])

            // 更新分页配置
            setPagination({
                current: page,
                pageSize,
                total: data.total || 0,
            })
        } catch (error) {
            // 错误处理：记录错误日志和显示错误消息
            console.error('加载用户列表失败:', error)
            message.error('加载用户列表失败，请稍后重试')
        } finally {
            // 无论成功或失败，都要关闭加载状态
            setLoading(false)
        }
    }

    /**
     * 组件挂载时初始化数据
     * 
     * 使用 useEffect Hook 在组件首次渲染时加载用户列表
     * 依赖数组为空，表示只在挂载时执行一次
     */
    useEffect(() => {
        loadData()
    }, [])

    /**
     * 表格列配置
     * 
     * 定义表格中每列的显示配置
     * 包括列标题、数据字段、渲染函数等
     * 
     * 每列配置说明：
     * - title: 列标题显示文本
     * - dataIndex: 对应数据对象的字段名
     * - key: React key 值，通常与 dataIndex 相同
     * - width: 列宽度
     * - render: 自定义渲染函数，用于复杂数据展示
     */
    const columns = [
        {
            title: 'ID',                    // 列标题
            dataIndex: 'id',               // 对应数据的字段名
            key: 'id',                     // React key 值
            width: 80,                     // 列宽度
        },
        {
            title: '用户名',
            dataIndex: 'username',
            key: 'username',
        },
        {
            title: '姓名',
            dataIndex: 'realname',
            key: 'realname',
        },
        {
            title: '邮箱',
            dataIndex: 'email',
            key: 'email',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            /**
             * 自定义渲染函数
             * 根据用户状态显示不同的标签
             * @param {number} status - 用户状态值
             * @returns {JSX.Element} 状态标签组件
             */
            render: (status) => (
                <Tag color={status === 1 ? 'success' : 'error'}>
                    {status === 1 ? '正常' : '禁用'}
                </Tag>
            ),
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
             * 操作列渲染函数
             * 为每行显示编辑和删除按钮
             * @param {Object} _ - 当前行数据（未使用）
             * @param {Object} record - 当前行数据记录
             * @returns {JSX.Element} 操作按钮组
             */
            render: (_, record) => (
                <Space>
                    {/* 编辑按钮 */}
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>

                    {/* 删除确认按钮 */}
                    <Popconfirm
                        title="确定要删除这个用户吗？"
                        onConfirm={() => handleDelete(record.id)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button
                            type="link"
                            size="small"
                            danger
                            icon={<DeleteOutlined />}
                        >
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ]

    /**
     * 处理表格分页变化
     * 
     * 当用户点击分页器或改变每页条数时触发
     * 重新加载对应页的数据
     * 
     * @param {Object} newPagination - 新的分页配置
     */
    const handleTableChange = (newPagination) => {
        loadData(newPagination.current, newPagination.pageSize)
    }

    /**
     * 打开新增用户对话框
     * 
     * 重置表单字段，清除编辑状态，显示新增对话框
     * 
     * 流程：
     * 1. 设置编辑用户为 null（新增模式）
     * 2. 重置表单字段
     * 3. 显示模态框
     */
    const handleAdd = () => {
        setEditingUser(null)         // 设置为新增模式
        form.resetFields()           // 清空表单字段
        setModalVisible(true)        // 显示模态框
    }

    /**
     * 打开编辑用户对话框
     * 
     * 将用户数据填充到表单中，设置为编辑模式
     * 
     * @param {Object} record - 当前要编辑的用户数据
     */
    const handleEdit = (record) => {
        setEditingUser(record)               // 设置为编辑模式，存储当前用户数据
        form.setFieldsValue(record)          // 将用户数据填充到表单字段
        setModalVisible(true)                // 显示模态框
    }

    /**
     * 删除用户
     * 
     * 异步删除用户操作，包含错误处理和成功反馈
     * 删除成功后刷新当前页数据
     * 
     * @param {string|number} id - 要删除的用户 ID
     */
    const handleDelete = async (id) => {
        try {
            // 调用删除用户 API
            await deleteUser(id)

            // 显示成功消息
            message.success('删除成功')

            // 重新加载当前页数据
            loadData(pagination.current, pagination.pageSize)
        } catch (error) {
            // 错误处理
            console.error('删除失败:', error)
            message.error('删除失败，请稍后重试')
        }
    }

    /**
     * 处理表单提交
     * 
     * 处理新增和编辑用户的表单提交逻辑
     * 包含表单验证、API 调用、错误处理和成功反馈
     * 
     * 流程：
     * 1. 表单字段验证
     * 2. 根据模式调用创建或更新 API
     * 3. 关闭模态框
     * 4. 刷新数据
     */
    const handleSubmit = async () => {
        try {
            // 表单验证，确保所有必填字段都已填写
            const values = await form.validateFields()

            if (editingUser) {
                // 编辑模式：调用更新用户 API
                await updateUser(editingUser.id, values)
                message.success('更新成功')
            } else {
                // 新增模式：调用创建用户 API
                await createUser(values)
                message.success('创建成功')
            }

            // 关闭模态框
            setModalVisible(false)

            // 重新加载当前页数据
            loadData(pagination.current, pagination.pageSize)
        } catch (error) {
            // 表单验证错误或 API 调用错误
            console.error('提交失败:', error)
            if (error.errorFields) {
                // 表单验证错误，不需要额外处理
                return
            }
            message.error('操作失败，请稍后重试')
        }
    }

    /**
     * 主渲染函数
     * 
     * 渲染完整的用户管理页面，包括：
     * - 页面标题和新增按钮
     * - 用户列表表格
     * - 新增/编辑用户模态框
     */
    return (
        <div>
            {/* 页面标题和操作区域 */}
            <div className="mb-4 flex justify-between items-center">
                <h1 className="text-2xl font-bold">用户管理</h1>
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={handleAdd}
                >
                    新增用户
                </Button>
            </div>

            {/* 用户列表表格 */}
            <Table
                columns={columns}                    // 表格列配置
                dataSource={dataSource}             // 表格数据源
                loading={loading}                   // 加载状态
                rowKey="id"                         // 行数据的 key 字段
                pagination={pagination}             // 分页配置
                onChange={handleTableChange}        // 分页变化处理
                // 其他表格配置
                scroll={{ x: 800 }}                 // 表格水平滚动
            />

            {/* 新增/编辑用户模态框 */}
            <Modal
                title={editingUser ? '编辑用户' : '新增用户'}  // 动态标题
                open={modalVisible}                 // 控制显示状态
                onOk={handleSubmit}                 // 确定按钮点击处理
                onCancel={() => setModalVisible(false)}  // 取消按钮点击处理
                width={600}                         // 模态框宽度
                destroyOnClose                      // 关闭时销毁内容
            >
                {/* 用户信息表单 */}
                <Form
                    form={form}
                    layout="vertical"               // 表单布局：垂直布局
                >
                    {/* 用户名字段 */}
                    <Form.Item
                        name="username"
                        label="用户名"
                        rules={[
                            {
                                required: true,
                                message: '请输入用户名'
                            },
                            {
                                min: 3,
                                max: 20,
                                message: '用户名长度在 3-20 个字符'
                            }
                        ]}
                    >
                        <Input placeholder="请输入用户名" />
                    </Form.Item>

                    {/* 密码字段（仅新增时显示） */}
                    {!editingUser && (
                        <Form.Item
                            name="password"
                            label="密码"
                            rules={[
                                {
                                    required: true,
                                    message: '请输入密码'
                                },
                                {
                                    min: 6,
                                    message: '密码长度不能少于 6 个字符'
                                }
                            ]}
                        >
                            <Input.Password placeholder="请输入密码" />
                        </Form.Item>
                    )}

                    {/* 真实姓名字段 */}
                    <Form.Item
                        name="realname"
                        label="姓名"
                        rules={[
                            {
                                required: true,
                                message: '请输入姓名'
                            }
                        ]}
                    >
                        <Input placeholder="请输入姓名" />
                    </Form.Item>

                    {/* 邮箱字段 */}
                    <Form.Item
                        name="email"
                        label="邮箱"
                        rules={[
                            {
                                required: true,
                                message: '请输入邮箱'
                            },
                            {
                                type: 'email',
                                message: '请输入有效的邮箱地址'
                            }
                        ]}
                    >
                        <Input placeholder="请输入邮箱" />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

/**
 * 组件导出
 * 
 * 将组件导出供其他模块使用
 */
export default UserList

/**
 * 用户管理最佳实践提示：
 * 
 * 1. 数据安全：
 *    - 密码字段不应该在编辑时显示和返回
 *    - 敏感信息需要加密存储
 *    - 用户删除应该使用软删除（标记为已删除状态）
 * 
 * 2. 用户体验：
 *    - 操作成功后应该有明确的反馈消息
 *    - 加载状态应该使用 loading 动画
 *    - 删除操作应该进行二次确认
 * 
 * 3. 性能优化：
 *    - 大数据量时应该使用虚拟滚动
 *    - 分页查询避免一次性加载所有数据
 *    - 搜索功能应该防抖处理
 * 
 * 4. 错误处理：
 *    - 网络错误应该提供重试机制
 *    - 表单验证应该在前端和后端双重验证
 *    - 敏感操作应该记录操作日志
 * 
 * 5. 数据权限：
 *    - 不同角色看到的用户数据可能不同
 *    - 普通用户可能只能查看个人信息
 *    - 管理员可以对所有用户进行操作
 */
