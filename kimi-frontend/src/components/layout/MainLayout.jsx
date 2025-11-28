/**
 * 主布局组件
 * 
 * 包含侧边栏和内容区域
 */
import { useState } from 'react'
import { Button } from 'antd'
import { MenuOutlined, CloseOutlined } from '@ant-design/icons'
import Sidebar from './Sidebar'

const MainLayout = ({ children }) => {
    const [sidebarOpen, setSidebarOpen] = useState(false)

    const toggleSidebar = () => {
        setSidebarOpen(!sidebarOpen)
    }

    return (
        <div className="app-container">
            {/* 移动端遮罩层 */}
            {sidebarOpen && (
                <div 
                    className="sidebar-overlay md:hidden"
                    onClick={() => setSidebarOpen(false)}
                />
            )}

            {/* 侧边栏 */}
            <Sidebar 
                isOpen={sidebarOpen} 
                onClose={() => setSidebarOpen(false)}
            />

            {/* 内容区域 */}
            <div className="content-area">
                {/* 移动端顶部栏 */}
                <div className="md:hidden flex items-center justify-between px-4 py-3 border-b bg-white">
                    <Button
                        type="text"
                        icon={sidebarOpen ? <CloseOutlined /> : <MenuOutlined />}
                        onClick={toggleSidebar}
                    />
                    <span className="font-semibold text-lg">Kimi AI</span>
                    <div className="w-10" /> {/* 占位，保持标题居中 */}
                </div>

                {/* 页面内容 */}
                <div className="flex-1 overflow-hidden">
                    {children}
                </div>
            </div>
        </div>
    )
}

export default MainLayout
