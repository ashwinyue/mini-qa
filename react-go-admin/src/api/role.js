/**
 * 角色管理 API 模块
 * 
 * 这个模块包含了与角色（Role）相关的所有 API 调用函数
 * 角色是权限管理中的核心概念，用于组织和管理用户的权限集合
 * 
 * 角色管理通常包括：
 * - 角色的创建、读取、更新和删除（CRUD 操作）
 * - 角色权限的分配和管理
 * - 用户与角色的关联管理
 * - 角色层级和继承关系
 * 
 * 每个角色通常包含：
 * - 角色 ID：唯一标识符
 * - 角色名称：显示名称
 * - 角色描述：详细说明
 * - 权限列表：角色包含的所有权限
 * - 创建时间和更新时间
 * - 创建人和最后修改人
 * 
 * 在使用这些函数时，通常需要配合 React 组件或状态管理库
 * 来处理异步请求和用户交互
 */

import request from '../utils/request'

/**
 * 获取角色列表
 * 
 * 从服务器获取所有角色的列表，支持分页和过滤
 * 
 * @param {Object} params - 查询参数
 * @param {number} params.page - 当前页码，默认为 1
 * @param {number} params.pageSize - 每页显示数量，默认为 10
 * @param {string} params.keyword - 关键词搜索，支持角色名称模糊搜索
 * @param {string} params.status - 角色状态过滤，active/inactive
 * 
 * @returns {Promise} 返回包含角色列表数据的 Promise 对象
 *         成功时返回 { data: { list: [], total: number, page: number, pageSize: number } }
 *         失败时返回错误对象
 * 
 * @example
 * // 获取第一页，10条记录
 * const response = await getRoles({ page: 1, pageSize: 10 })
 * 
 * @example
 * // 搜索包含"admin"的角色
 * const response = await getRoles({ keyword: 'admin' })
 */
export const getRoles = (params) => {
    return request({
        url: '/roles',
        method: 'get',
        params,
    })
}

/**
 * 获取角色列表（兼容旧版本）
 * 
 * 这个函数与 getRoles 功能相同，主要是为了向后兼容
 * 在新项目中建议使用 getRoles 函数
 * 
 * @param {Object} params - 查询参数，与 getRoles 相同
 * @returns {Promise} 返回与 getRoles 相同格式的数据
 */
export const getRoleList = (params) => {
    return request({
        url: '/roles',
        method: 'get',
        params,
    })
}

/**
 * 获取角色详情
 * 
 * 根据角色 ID 获取特定角色的详细信息
 * 包括角色的基本信息、权限列表、关联用户等
 * 
 * @param {string|number} id - 角色 ID，必须是有效的角色标识符
 * 
 * @returns {Promise} 返回角色详细信息的 Promise 对象
 *         成功时返回 { data: { 
 *           id: number, 
 *           name: string, 
 *           description: string, 
 *           permissions: Array, 
 *           userCount: number,
 *           createdAt: string,
 *           updatedAt: string 
 *         } }
 *         失败时返回错误对象
 * 
 * @example
 * // 获取 ID 为 1 的角色详情
 * const response = await getRoleDetail(1)
 * console.log(response.data.name) // 输出角色名称
 * 
 * @throws {Error} 当角色 ID 不存在时，服务器会返回 404 错误
 */
export const getRoleDetail = (id) => {
    return request({
        url: `/roles/${id}`,
        method: 'get',
    })
}

/**
 * 创建新角色
 * 
 * 在系统中创建一个新的角色，需要提供角色基本信息
 * 创建成功后，系统会自动分配一个唯一的角色 ID
 * 
 * @param {Object} data - 角色数据对象
 * @param {string} data.name - 角色名称（必填），2-50 个字符
 * @param {string} data.description - 角色描述，可选，1-200 个字符
 * @param {Array<number>} data.permissionIds - 权限 ID 数组，角色包含的权限列表
 * @param {string} data.status - 角色状态，active 或 inactive
 * 
 * @returns {Promise} 返回创建结果的 Promise 对象
 *         成功时返回 { data: { id: number, name: string, ... } }
 *         失败时返回错误对象
 * 
 * @example
 * // 创建一个管理员角色
 * const response = await createRole({
 *     name: '管理员',
 *     description: '系统管理员，拥有所有权限',
 *     permissionIds: [1, 2, 3, 4, 5],
 *     status: 'active'
 * })
 * 
 * @throws {Error} 
 *         - 角色名称已存在（409 冲突）
 *         - 权限 ID 无效（400 请求错误）
 *         - 权限不足（403 禁止访问）
 */
export const createRole = (data) => {
    return request({
        url: '/roles',
        method: 'post',
        data,
    })
}

/**
 * 更新角色信息
 * 
 * 修改现有角色的基本信息，可以更新角色名称、描述、权限等
 * 更新时会记录操作日志和最后修改时间
 * 
 * @param {string|number} id - 角色 ID，必须是存在的角色
 * @param {Object} data - 更新的角色数据对象
 * @param {string} data.name - 新的角色名称，2-50 个字符
 * @param {string} data.description - 新的角色描述，1-200 个字符
 * @param {Array<number>} data.permissionIds - 新的权限 ID 数组
 * @param {string} data.status - 新的角色状态，active 或 inactive
 * 
 * @returns {Promise} 返回更新结果的 Promise 对象
 *         成功时返回 { data: { id: number, name: string, updatedAt: string } }
 *         失败时返回错误对象
 * 
 * @example
 * // 更新角色名称和权限
 * const response = await updateRole(1, {
 *     name: '超级管理员',
 *     permissionIds: [1, 2, 3, 4, 5, 6, 7]
 * })
 * 
 * @throws {Error}
 *         - 角色不存在（404 未找到）
 *         - 角色名称冲突（409 冲突）
 *         - 权限 ID 无效（400 请求错误）
 *         - 权限不足（403 禁止访问）
 */
export const updateRole = (id, data) => {
    return request({
        url: `/roles/${id}`,
        method: 'put',
        data,
    })
}

/**
 * 删除角色
 * 
 * 从系统中删除指定的角色，这是一个不可逆的操作
 * 删除前需要确保该角色没有关联的用户，否则可能需要先转移用户到其他角色
 * 
 * @param {string|number} id - 要删除的角色 ID
 * 
 * @returns {Promise} 返回删除结果的 Promise 对象
 *         成功时返回 { message: "角色删除成功" }
 *         失败时返回错误对象
 * 
 * @example
 * // 删除 ID 为 5 的角色
 * const response = await deleteRole(5)
 * console.log(response.message) // 输出成功消息
 * 
 * @throws {Error}
 *         - 角色不存在（404 未找到）
 *         - 角色有关联用户（409 冲突，需要先转移用户）
 *         - 系统保留角色不能删除（400 请求错误）
 *         - 权限不足（403 禁止访问）
 * 
 * @important 安全提示
 * 在实际应用中，删除操作通常需要二次确认，
 * 并且可能会提供"软删除"选项（将角色标记为已删除而不是物理删除）
 */
export const deleteRole = (id) => {
    return request({
        url: `/roles/${id}`,
        method: 'delete',
    })
}

/**
 * 角色管理最佳实践提示：
 * 
 * 1. 在创建角色前，应该先规划好权限体系
 * 2. 为角色命名时应该使用描述性的名称，如"内容管理员"、"财务审核员"等
 * 3. 避免创建过多细粒度的角色，可以考虑使用权限组合的方式
 * 4. 定期审查和清理不再使用的角色
 * 5. 对于重要的角色变更，建议保留操作日志
 * 6. 在删除角色前，应该考虑对现有用户的影响
 * 
 * 常见使用场景：
 * - 用户管理系统中的权限分配
 * - 企业组织架构中的岗位权限
 * - 多租户系统中的角色隔离
 * - 内容管理系统的编辑权限控制
 */
