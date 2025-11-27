/**
 * 用户管理 API 模块
 * 
 * 这个文件包含了与用户管理相关的所有 API 函数，包括：
 * - 获取用户列表：分页查询用户信息
 * - 获取用户详情：根据 ID 获取单个用户的完整信息
 * - 创建用户：添加新用户到系统
 * - 更新用户：修改现有用户的信息
 * - 删除用户：从系统中移除用户
 * 
 * 这些函数都依赖于统一的 HTTP 请求工具(request.js)来处理网络请求
 * 和错误处理，统一处理认证头和错误响应
 * 
 * ⚠️ 重要说明：
 * 响应拦截器（在 request.js 中）已经提取了响应数据中的 data 字段，
 * 因此所有 API 函数直接返回 response 即可，无需再访问 response.data
 */

// 导入统一的 HTTP 请求工具
import request from '../utils/request'

/**
 * 获取用户列表
 * 
 * 分页获取用户列表，支持搜索和排序功能
 * 这是管理后台中最常用的功能之一，用于显示所有用户的概览信息
 * 
 * @param {Object} params - 查询参数
 * @param {number} [params.page=1] - 当前页码，从1开始
 * @param {number} [params.pageSize=10] - 每页显示的条目数量
 * @param {string} [params.keyword] - 搜索关键词（用户名、邮箱等）
 * @param {string} [params.status] - 用户状态筛选（active、inactive等）
 * @param {string} [params.sortBy] - 排序字段（createdAt、updatedAt、name等）
 * @param {string} [params.sortOrder] - 排序方向（asc、desc）
 * 
 * @returns {Promise<Object>} 包含用户列表和分页信息的响应数据
 * 
 * @example
 * ```javascript
 * // 获取第一页，每页10条记录
 * const response = await getUserList({ page: 1, pageSize: 10 })
 * 
 * // 搜索用户并按创建时间倒序排列
 * const response = await getUserList({
 *   keyword: 'admin',
 *   sortBy: 'createdAt',
 *   sortOrder: 'desc'
 * })
 * ```
 */
export const getUserList = async (params = {}) => {
    try {
        /**
         * 构建查询参数
         * 
         * 只传递有值的参数，避免出现 undefined 值
         * 这有助于构建干净的 URL 查询字符串
         */
        const queryParams = {
            page: params.page || 1,
            pageSize: params.pageSize || 10,
            ...(params.keyword && { keyword: params.keyword }),
            ...(params.status && { status: params.status }),
            ...(params.sortBy && { sortBy: params.sortBy }),
            ...(params.sortOrder && { sortOrder: params.sortOrder })
        }

        /**
         * 发送 GET 请求到用户列表端点
         * 
         * 使用 query 参数传递查询条件，
         * request.get() 会自动将参数拼接到 URL 后面
         */
        const response = await request.get('/users', { params: queryParams })

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   data: {
         *     users: [
         *       {
         *         id: 1,
         *         username: 'admin',
         *         name: '管理员',
         *         email: 'admin@example.com',
         *         status: 'active',
         *         role: 'admin',
         *         createdAt: '2024-01-15T10:30:00Z',
         *         updatedAt: '2024-01-15T10:30:00Z'
         *       }
         *     ],
         *     pagination: {
         *       currentPage: 1,
         *       pageSize: 10,
         *       total: 50,
         *       totalPages: 5
         *     }
         *   }
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 获取用户列表失败可能的原因：
         * 1. 访问令牌无效或已过期
         * 2. 当前用户没有查看用户列表的权限
         * 3. 后端数据库查询异常
         * 4. 网络连接问题
         */
        console.error('获取用户列表失败:', error)
        throw error
    }
}

/**
 * 获取用户详情
 * 
 * 根据用户ID获取单个用户的完整详细信息
 * 通常在用户详情页面或编辑用户时使用
 * 
 * @param {number|string} userId - 用户ID（通常是数字）
 * 
 * @returns {Promise<Object>} 指定用户的详细信息
 * 
 * @example
 * ```javascript
 * try {
 *   const user = await getUserDetail(123)
 *   console.log('用户详情:', user.name, user.email)
 * } catch (error) {
 *   console.error('获取用户详情失败:', error)
 * }
 * ```
 */
export const getUserDetail = async (userId) => {
    try {
        /**
         * 参数验证
         * 
         * 确保传入的 userId 是有效的
         * 这有助于在后端抛出有意义的错误信息
         */
        if (!userId) {
            throw new Error('用户ID不能为空')
        }

        /**
         * 发送 GET 请求到用户详情端点
         * 
         * 路径参数直接拼接到 URL 中，
         * request.get() 支持动态路径参数
         */
        const response = await request.get(`/users/${userId}`)

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   data: {
         *     id: 123,
         *     username: 'john_doe',
         *     name: '张三',
         *     email: 'zhangsan@example.com',
         *     phone: '+86 138-0000-0000',
         *     status: 'active',
         *     role: 'user',
         *     avatar: 'https://example.com/avatar.jpg',
         *     bio: '这是一个用户简介',
         *     createdAt: '2024-01-15T10:30:00Z',
         *     updatedAt: '2024-01-20T14:22:00Z',
         *     lastLoginAt: '2024-01-21T09:15:00Z',
         *     permissions: ['read:users', 'write:own_profile']
         *   }
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 获取用户详情失败可能的原因：
         * 1. 指定的用户ID不存在
         * 2. 当前用户没有查看该用户详情的权限
         * 3. 访问令牌无效或已过期
         * 
         * 这些错误需要不同的处理方式
         */
        console.error(`获取用户详情失败 (ID: ${userId}):`, error)
        throw error
    }
}

/**
 * 创建用户
 * 
 * 向系统添加一个新用户
 * 通常在用户管理页面中由管理员使用
 * 
 * @param {Object} userData - 新用户的数据
 * @param {string} userData.username - 用户名（唯一）
 * @param {string} userData.realname - 用户真实姓名
 * @param {string} userData.email - 用户邮箱（唯一）
 * @param {string} userData.password - 用户密码
 * @param {string} [userData.role] - 用户角色（默认：user）
 * @param {string} [userData.status] - 用户状态（默认：active）
 * @param {string} [userData.phone] - 用户电话
 * @param {string} [userData.bio] - 用户简介
 * 
 * @returns {Promise<Object>} 创建成功后的用户信息
 * 
 * @example
 * ```javascript
 * try {
 *   const newUser = await createUser({
 *     username: 'new_user',
 *     realname: '新用户',
 *     email: 'newuser@example.com',
 *     password: 'password123',
 *     role: 'user'
 *   })
 *   console.log('用户创建成功:', newUser.id)
 * } catch (error) {
 *   console.error('创建用户失败:', error)
 * }
 * ```
 */
export const createUser = async (userData) => {
    try {
        /**
         * 参数验证
         * 
         * 确保必需字段都存在且有效
         * 这里只做基本的验证，详细的验证逻辑通常在后端处理
         * 
         * 注意：后端使用 'realname' 字段而不是 'name'
         */
        const requiredFields = ['username', 'realname', 'email', 'password']
        const missingFields = requiredFields.filter(field => !userData[field])

        if (missingFields.length > 0) {
            throw new Error(`缺少必需字段: ${missingFields.join(', ')}`)
        }

        /**
         * 发送 POST 请求到用户创建端点
         * 
         * userData 会作为请求体发送到后端
         * request.post() 会自动处理 JSON 序列化
         */
        const response = await request.post('/users', userData)

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   message: '用户创建成功',
         *   data: {
         *     id: 124,
         *     username: 'new_user',
         *     realname: '新用户',
         *     email: 'newuser@example.com',
         *     status: 'active',
         *     role: 'user',
         *     createdAt: '2024-01-22T11:30:00Z'
         *   }
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 创建用户失败可能的原因：
         * 1. 用户名或邮箱已存在
         * 2. 密码不符合安全要求
         * 3. 当前用户没有创建用户的权限
         * 4. 输入数据格式不正确
         * 
         * 这些错误信息通常来自后端验证
         */
        console.error('创建用户失败:', error)
        throw error
    }
}

/**
 * 更新用户
 * 
 * 修改现有用户的信息
 * 可以更新用户的基本信息、角色、状态等
 * 
 * @param {number|string} userId - 要更新的用户ID
 * @param {Object} updateData - 要更新的数据
 * @param {string} [updateData.name] - 用户姓名
 * @param {string} [updateData.email] - 用户邮箱
 * @param {string} [updateData.role] - 用户角色
 * @param {string} [updateData.status] - 用户状态
 * @param {string} [updateData.phone] - 用户电话
 * @param {string} [updateData.bio] - 用户简介
 * @param {string} [updateData.avatar] - 头像URL
 * 
 * @returns {Promise<Object>} 更新后的用户信息
 * 
 * @example
 * ```javascript
 * try {
 *   const updatedUser = await updateUser(123, {
 *     name: '张三（已更新）',
 *     status: 'inactive'
 *   })
 *   console.log('用户更新成功:', updatedUser.name)
 * } catch (error) {
 *   console.error('更新用户失败:', error)
 * }
 * ```
 */
export const updateUser = async (userId, updateData) => {
    try {
        /**
         * 参数验证
         * 
         * 确保用户ID和更新数据都存在
         */
        if (!userId) {
            throw new Error('用户ID不能为空')
        }

        if (!updateData || typeof updateData !== 'object') {
            throw new Error('更新数据必须是有效对象')
        }

        /**
         * 发送 PUT/PATCH 请求到用户更新端点
         * 
         * PUT 通常用于完整替换资源，
         * PATCH 用于部分更新资源
         * 
         * 这里使用 PUT 来替换整个用户对象
         */
        const response = await request.put(`/users/${userId}`, updateData)

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   message: '用户信息更新成功',
         *   data: {
         *     id: 123,
         *     username: 'john_doe',
         *     name: '张三（已更新）',
         *     email: 'zhangsan@example.com',
         *     status: 'inactive',
         *     role: 'user',
         *     updatedAt: '2024-01-22T15:45:00Z'
         *   }
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 更新用户失败可能的原因：
         * 1. 用户ID不存在
         * 2. 邮箱已被其他用户使用
         * 3. 当前用户没有更新该用户的权限
         * 4. 更新数据格式不正确
         */
        console.error(`更新用户失败 (ID: ${userId}):`, error)
        throw error
    }
}

/**
 * 删除用户
 * 
 * 从系统中永久删除用户
 * 这是一个危险操作，通常需要特殊权限和确认步骤
 * 
 * @param {number|string} userId - 要删除的用户ID
 * 
 * @returns {Promise<Object>} 删除操作的结果
 * 
 * @example
 * ```javascript
 * try {
 *   const result = await deleteUser(123)
 *   console.log('用户删除成功:', result.message)
 * } catch (error) {
 *   console.error('删除用户失败:', error)
 * }
 * ```
 */
export const deleteUser = async (userId) => {
    try {
        /**
         * 参数验证
         */
        if (!userId) {
            throw new Error('用户ID不能为空')
        }

        /**
         * 安全检查：防止误删重要用户
         * 
         * 在实际应用中，可能需要：
         * 1. 检查当前用户权限
         * 2. 确认用户不在关键系统中
         * 3. 提供用户确认对话框
         */
        if (userId === 1) {
            throw new Error('不能删除系统管理员账户')
        }

        /**
         * 发送 DELETE 请求到用户删除端点
         */
        const response = await request.delete(`/users/${userId}`)

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   message: '用户删除成功'
         * }
         */
        return response
    } catch (error) {
        /**
         * 错误处理
         * 
         * 删除用户失败可能的原因：
         * 1. 用户ID不存在
         * 2. 当前用户没有删除权限
         * 3. 用户被其他数据关联，无法删除（外键约束）
         * 4. 系统保护的用户（如管理员）
         */
        console.error(`删除用户失败 (ID: ${userId}):`, error)
        throw error
    }
}

/**
 * 批量删除用户
 * 
 * 一次性删除多个用户
 * 这是一个高效的批量操作，但同样需要谨慎使用
 * 
 * @param {Array<number|string>} userIds - 要删除的用户ID数组
 * 
 * @returns {Promise<Object>} 批量删除操作的结果
 * 
 * @example
 * ```javascript
 * try {
 *   const result = await batchDeleteUsers([123, 124, 125])
 *   console.log(`成功删除 ${result.deletedCount} 个用户`)
 * } catch (error) {
 *   console.error('批量删除用户失败:', error)
 * }
 * ```
 */
export const batchDeleteUsers = async (userIds) => {
    try {
        /**
         * 参数验证
         */
        if (!Array.isArray(userIds) || userIds.length === 0) {
            throw new Error('用户ID列表不能为空')
        }

        /**
         * 安全检查：过滤掉受保护的用户ID
         */
        const filteredIds = userIds.filter(id => id !== 1)

        if (filteredIds.length === 0) {
            throw new Error('不能删除系统管理员账户')
        }

        /**
         * 发送 POST 请求到批量删除端点
         * 
         * 使用 POST 而不是 DELETE 是因为：
         * 1. DELETE 不支持请求体
         * 2. 需要传递用户ID数组
         */
        const response = await request.post('/users/batch-delete', {
            userIds: filteredIds
        })

        /**
         * 响应数据格式（示例）：
         * {
         *   success: true,
         *   message: '批量删除成功',
         *   data: {
         *     deletedCount: 3,
         *     failedCount: 0,
         *     totalRequested: 4,
         *     skippedIds: [1] // 被跳过的受保护用户
         *   }
         * }
         */
        return response
    } catch (error) {
        console.error('批量删除用户失败:', error)
        throw error
    }
}

/**
 * 修改密码
 * 
 * @param {string} oldPassword - 原密码
 * @param {string} newPassword - 新密码
 * @returns {Promise<Object>}
 */
export const changePassword = async (oldPassword, newPassword) => {
    try {
        const response = await request.post('/auth/change-password', {
            oldPassword,
            newPassword
        })
        return response
    } catch (error) {
        console.error('修改密码失败:', error)
        throw error
    }
}

export default {
    getUserList,
    getUserDetail,
    createUser,
    updateUser,
    deleteUser,
    batchDeleteUsers,
    changePassword
}
