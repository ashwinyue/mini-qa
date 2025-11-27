/**
 * 登录页面组件
 * 
 * 该组件实现用户登录功能，包含：
 * 1. 登录表单（用户名和密码）
 * 2. 表单验证
 * 3. 登录状态管理
 * 4. 登录成功后的页面跳转
 * 5. 美观的登录界面设计
 */

// 从 React 导入 useState Hook，用于管理组件状态
// useState 是 React 中最常用的 Hook，用于在函数组件中添加状态
import { useState } from 'react'
// 从 React Router 导入 useNavigate Hook，用于编程式导航
// useNavigate 替代了老版本的 useHistory，可以实现页面跳转
import { useNavigate } from 'react-router-dom'
// 导入 Ant Design 组件
import { Form, Input, Button, Card, message } from 'antd'
// 导入 Ant Design 图标
import { UserOutlined, LockOutlined } from '@ant-design/icons'
// 导入登录 API 函数
import { login as loginApi } from '../api/auth'
// 导入 useAuth Hook，用于访问认证上下文
import { useAuth } from '../context/AuthContext'

/**
 * 登录组件
 * 
 * 这是一个函数组件，使用 React Hooks 管理状态和副作用
 * 组件特点：
 * - 使用 useState 管理加载状态
 * - 使用 useNavigate 实现登录后跳转
 * - 使用 useAuth 获取认证方法和状态
 * - 使用 Ant Design 的 Form 组件实现表单功能
 */
const Login = () => {
    /**
     * 组件状态
     * - loading: 控制登录按钮的加载状态，防止重复提交
     */
    const [loading, setLoading] = useState(false)

    // 使用 useNavigate Hook 获取导航函数
    const navigate = useNavigate()

    // 使用 useAuth Hook 获取认证上下文的方法
    const { login } = useAuth()

    /**
     * 表单提交处理函数
     * 
     * @param {Object} values 表单值对象，包含 username 和 password
     * 
     * 该函数：
     * 1. 设置加载状态
     * 2. 调用登录 API
     * 3. 保存认证信息到上下文
     * 4. 显示成功消息
     * 5. 跳转到主页
     * 6. 错误处理和状态清理
     */
    const onFinish = async (values) => {
        // 设置加载状态，显示加载动画并禁用按钮
        setLoading(true)

        try {
            // 调用登录 API，传入表单值
            // loginApi 是一个异步函数，返回包含 token 和用户信息的对象
            const data = await loginApi(values)

            /**
             * 登录成功后的处理：
             * 1. 将 token 和用户信息保存到认证上下文
             * 2. 认证上下文会自动更新状态并触发相关组件重新渲染
             */
            login(data.token, data.user)

            // 显示成功消息
            message.success('登录成功')

            /**
             * 使用 navigate 函数进行编程式导航
             * 参数 '/' 表示跳转到根路径（仪表盘页面）
             * 这会触发 React Router 的页面跳转
             */
            navigate('/')

        } catch (error) {
            /**
             * 登录失败处理
             * 
             * 可能的错误情况：
             * - 用户名或密码错误
             * - 网络连接问题
             * - 服务器错误
             * - API 返回业务错误
             */
            console.error('登录失败:', error)
            // 注意：这里没有显示错误消息，因为错误拦截器已经在 request.js 中处理了

        } finally {
            // 无论成功或失败，都要清除加载状态
            setLoading(false)
        }
    }

    return (
        /**
         * 登录页面容器
         * 
         * 使用内联样式创建居中的登录界面：
         * - display: flex：使用 flexbox 布局
         * - justifyContent: center：水平居中
         * - alignItems: center：垂直居中
         * - minHeight: 100vh：最小高度为视口高度
         * - background: 渐变背景色
         */
        <div className="flex justify-center items-center min-h-screen bg-gradient-to-br from-[#667eea] to-[#764ba2]">
            {/**
             * Ant Design 的 Card 组件
             * 
             * Card 组件特点：
             * - 提供卡片式容器
             * - 可以设置标题
             * - 自带阴影和边框
             * - 支持自定义样式
             */}
            <Card
                title="管理系统登录"
                className="w-[400px] shadow-lg"
            >
                {/**
                 * Ant Design 的 Form 组件
                 * 
                 * Form 组件功能：
                 * - 表单状态管理
                 * - 表单验证
                 * - 提交处理
                 * - 自动布局
                 * 
                 * 属性说明：
                 * - name: 表单名称
                 * - initialValues: 初始值（用于演示）
                 * - onFinish: 表单提交回调
                 * - autoComplete: 禁用自动完成
                 * - size: 组件大小
                 */}
                <Form
                    name="login"
                    // 演示用的初始值，生产环境中应该删除
                    initialValues={{ username: 'admin', password: '123456' }}
                    onFinish={onFinish}
                    autoComplete="off"
                    size="large"
                >
                    {/**
                     * 用户名输入项
                     * 
                     * Form.Item 组件特点：
                     * - 自动包装输入组件
                     * - 支持标签和验证规则
                     * - 自动布局管理
                     * 
                     * 属性说明：
                     * - name: 表单字段名
                     * - rules: 验证规则数组
                     */}
                    <Form.Item
                        name="username"
                        rules={[{ required: true, message: '请输入用户名' }]}
                    >
                        {/**
                         * Input 组件
                         * 
                         * 属性说明：
                         * - prefix: 输入框前缀图标
                         * - placeholder: 占位符文本
                         */}
                        <Input
                            prefix={<UserOutlined />}
                            placeholder="用户名"
                        />
                    </Form.Item>

                    {/**
                     * 密码输入项
                     * 
                     * Input.Password 是 Ant Design 专门用于密码输入的组件，
                     * 自动处理密码的显示/隐藏功能
                     */}
                    <Form.Item
                        name="password"
                        rules={[{ required: true, message: '请输入密码' }]}
                    >
                        <Input.Password
                            prefix={<LockOutlined />}
                            placeholder="密码"
                        />
                    </Form.Item>

                    {/**
                     * 提交按钮
                     * 
                     * Form.Item 不设置 name 属性时，不会作为表单字段处理
                     * 
                     * Button 组件属性：
                     * - type="primary": 主要按钮样式（蓝色）
                     * - htmlType="submit": 表单提交按钮
                     * - block: 占满整个宽度
                     * - loading: 加载状态
                     */}
                    <Form.Item>
                        <Button
                            type="primary"
                            htmlType="submit"
                            block
                            loading={loading}
                        >
                            登录
                        </Button>
                    </Form.Item>
                </Form>

                {/**
                 * 提示信息
                 * 
                 * 用于显示默认用户名和密码，
                 * 仅用于演示，生产环境中应该删除
                 */}
                <div className="text-center text-gray-500 text-xs">
                    默认账号：admin / 123456
                </div>
            </Card>
        </div>
    )
}

// 导出默认组件
export default Login
