import React from 'react';
import { Card, Button, Input, Tag, Space } from 'antd';
import { SearchOutlined, PlusOutlined } from '@ant-design/icons';

/**
 * Ant Design + Tailwind CSS 混合使用示例
 * 
 * 使用原则：
 * 1. Ant Design - 用于复杂组件（Card, Button, Input, Modal 等）
 * 2. Tailwind CSS - 用于布局、间距、响应式设计
 */
const TailwindExample: React.FC = () => {
  return (
    <div className="p-6 bg-gray-50 min-h-screen">
      {/* 使用 Tailwind 做响应式网格布局 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        {/* Ant Design Card + Tailwind 样式 */}
        <Card 
          title="课程统计" 
          className="hover:shadow-lg transition-shadow duration-300"
        >
          <div className="flex flex-col gap-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600">总课程数</span>
              <span className="text-2xl font-bold text-primary">128</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600">在线学员</span>
              <span className="text-2xl font-bold text-success">1,234</span>
            </div>
          </div>
        </Card>

        <Card 
          title="订单统计" 
          className="hover:shadow-lg transition-shadow duration-300"
        >
          <div className="flex flex-col gap-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600">今日订单</span>
              <span className="text-2xl font-bold text-primary">45</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600">总收入</span>
              <span className="text-2xl font-bold text-success">¥89,234</span>
            </div>
          </div>
        </Card>

        <Card 
          title="系统状态" 
          className="hover:shadow-lg transition-shadow duration-300"
        >
          <div className="flex flex-col gap-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600">知识库</span>
              <Tag color="success">正常</Tag>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600">AI 模型</span>
              <Tag color="processing">运行中</Tag>
            </div>
          </div>
        </Card>
      </div>

      {/* 搜索栏 - Ant Design 组件 + Tailwind 布局 */}
      <Card className="mb-6">
        <div className="flex flex-col md:flex-row gap-4">
          <Input
            placeholder="搜索课程、订单..."
            prefix={<SearchOutlined />}
            className="flex-1"
            size="large"
          />
          <Space>
            <Button type="primary" size="large" icon={<PlusOutlined />}>
              新建
            </Button>
            <Button size="large">导出</Button>
          </Space>
        </div>
      </Card>

      {/* 课程列表 - 响应式卡片 */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {[1, 2, 3, 4, 5, 6, 7, 8].map((item) => (
          <Card
            key={item}
            hoverable
            className="group"
            cover={
              <div className="h-40 bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
                <span className="text-white text-4xl font-bold">课程 {item}</span>
              </div>
            }
          >
            <div className="space-y-2">
              <h3 className="text-lg font-semibold truncate group-hover:text-primary transition-colors">
                Python 基础课程 {item}
              </h3>
              <div className="flex items-center justify-between text-sm text-gray-500">
                <span>40 小时</span>
                <Tag color="blue">初级</Tag>
              </div>
              <div className="flex items-center justify-between pt-2 border-t">
                <span className="text-primary font-bold">¥299</span>
                <Button type="link" size="small">查看详情</Button>
              </div>
            </div>
          </Card>
        ))}
      </div>

      {/* 底部信息 - Tailwind 布局 */}
      <div className="mt-8 p-4 bg-white rounded-lg shadow-sm">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <div className="text-gray-600">
            © 2024 AI 智能客服系统. All rights reserved.
          </div>
          <Space>
            <Button type="link">关于我们</Button>
            <Button type="link">帮助中心</Button>
            <Button type="link">联系我们</Button>
          </Space>
        </div>
      </div>
    </div>
  );
};

export default TailwindExample;
