/**
 * 404 错误页面组件
 * 
 * 当用户访问不存在的路由时显示此页面
 * 
 * 功能特点：
 * 1. 显示友好的 404 错误提示
 * 2. 提供返回首页的按钮
 * 3. 使用 Ant Design 的 Result 组件提供统一的错误页面样式
 * 
 * 使用场景：
 * - 用户输入了错误的 URL
 * - 访问已删除或不存在的页面
 * - 路由配置中未定义的路径
 */

// 导入 Ant Design 组件
import { Button, Result } from 'antd'
// Button: 按钮组件，用于提供返回首页的操作
// Result: 结果页组件，提供统一的成功/失败/404等状态页面

// 导入 React Router 的导航 Hook
import { useNavigate } from 'react-router-dom'
// useNavigate: 用于编程式导航，可以在代码中跳转到其他页面

/**
 * 404 页面组件
 * 
 * 这是一个简单的错误页面，当路由匹配失败时显示
 * 
 * @component
 */
const NotFound = () => {
    /**
     * 获取导航函数
     * 
     * navigate 函数用于在代码中进行页面跳转
     * 类似于后端的重定向功能
     */
    const navigate = useNavigate()

    return (
        /**
         * Result 组件：Ant Design 提供的结果页组件
         * 
         * 属性说明：
         * - status: 页面状态类型，支持 'success' | 'error' | 'info' | 'warning' | '404' | '403' | '500'
         * - title: 主标题，通常显示状态码或简短描述
         * - subTitle: 副标题，显示详细的错误信息
         * - extra: 额外的操作按钮或内容，通常放置返回、重试等按钮
         */
        <Result
            status="404"                              // 显示 404 状态图标和样式
            title="404"                               // 主标题：显示错误代码
            subTitle="抱歉，您访问的页面不存在"        // 副标题：显示友好的错误提示
            extra={
                /**
                 * 返回首页按钮
                 * 
                 * type="primary": 使用主要按钮样式（蓝色背景）
                 * onClick: 点击时调用 navigate('/') 跳转到首页
                 * 
                 * navigate('/') 的工作原理：
                 * 1. 更新浏览器地址栏为根路径 '/'
                 * 2. 触发路由匹配，渲染对应的组件
                 * 3. 不会刷新整个页面，只更新内容区域（SPA 特性）
                 */
                <Button type="primary" onClick={() => navigate('/')}>
                    返回首页
                </Button>
            }
        />
    )
}

// 导出组件，供路由配置使用
export default NotFound
