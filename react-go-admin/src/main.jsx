import React from 'react'
// ReactDOM 是 React 的 DOM 渲染器，用于将 React 组件渲染到真实的 DOM 节点中
import ReactDOM from 'react-dom/client'
// ConfigProvider 是 Ant Design 的配置组件，用于全局配置组件的国际化、主题等
import { ConfigProvider } from 'antd'
// zhCN 是 Ant Design 的中文语言包，实现界面中文化
import zhCN from 'antd/es/locale/zh_CN'
// 导入根组件 App
import App from './App.jsx'
// 导入全局样式文件
import './App.css'

/**
 * React 应用的入口点
 * React 18 使用新的 createRoot API 来创建应用的根节点
 * 
 * ReactDOM.createRoot() 方法：
 * - 接收一个 DOM 元素（通常是 index.html 中的 root 元素）
 * - 返回一个根容器对象，该对象有一个 render() 方法用于渲染 React 组件
 * 
 * React.StrictMode 组件：
 * - 用于在开发模式下检测潜在的副作用和过时的 API
 * - 不会在生产环境中生效
 * - 帮助发现不安全的生命周期方法、废弃的 API 等问题
 */
ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        {/* 
            ConfigProvider 是 Antd 5.0 的新特性，用于全局配置
            - locale: 设置国际化语言，这里使用中文
            - theme: 可以配置主题色彩、圆角等样式
            - 其他配置如组件大小、表单验证等
        */}
        <ConfigProvider locale={zhCN}>
            {/* 
                渲染根组件 App
                在 React 18 中，可以直接返回 JSX 而不需要用 React.Fragment
            */}
            <App />
        </ConfigProvider>
    </React.StrictMode>,
)
