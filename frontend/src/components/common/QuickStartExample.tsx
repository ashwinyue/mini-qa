import React from 'react';
import { Card, Button, Input, Tag } from 'antd';
import { SearchOutlined, HeartOutlined, ShoppingCartOutlined } from '@ant-design/icons';

/**
 * 快速开始示例 - 展示 Ant Design + Tailwind CSS 的基本用法
 */
const QuickStartExample: React.FC = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50 p-8">
      <div className="max-w-6xl mx-auto">
        {/* 页面标题 */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-800 mb-4">
            Ant Design + Tailwind CSS
          </h1>
          <p className="text-lg text-gray-600">
            完美结合，打造现代化的用户界面
          </p>
        </div>

        {/* 搜索栏 */}
        <Card className="mb-8 shadow-lg">
          <div className="flex gap-4">
            <Input
              size="large"
              placeholder="搜索课程、文章、视频..."
              prefix={<SearchOutlined className="text-gray-400" />}
              className="flex-1"
            />
            <Button type="primary" size="large">
              搜索
            </Button>
          </div>
        </Card>

        {/* 特色课程 */}
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">特色课程</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
              { title: 'Python 全栈开发', price: 299, level: '初级', students: 1234 },
              { title: 'React 进阶实战', price: 399, level: '中级', students: 856 },
              { title: 'AI 机器学习', price: 599, level: '高级', students: 432 },
            ].map((course, index) => (
              <Card
                key={index}
                hoverable
                className="group overflow-hidden"
                cover={
                  <div className="h-48 bg-gradient-to-br from-blue-400 via-purple-400 to-pink-400 flex items-center justify-center transform group-hover:scale-105 transition-transform duration-300">
                    <span className="text-white text-5xl font-bold opacity-20">
                      {index + 1}
                    </span>
                  </div>
                }
              >
                <div className="space-y-3">
                  <div className="flex items-start justify-between">
                    <h3 className="text-lg font-semibold text-gray-800 group-hover:text-primary transition-colors">
                      {course.title}
                    </h3>
                    <HeartOutlined className="text-gray-400 hover:text-red-500 cursor-pointer transition-colors" />
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <Tag color="blue">{course.level}</Tag>
                    <span className="text-sm text-gray-500">
                      {course.students.toLocaleString()} 人学习
                    </span>
                  </div>

                  <div className="flex items-center justify-between pt-3 border-t border-gray-100">
                    <div>
                      <span className="text-2xl font-bold text-primary">
                        ¥{course.price}
                      </span>
                    </div>
                    <Button 
                      type="primary" 
                      icon={<ShoppingCartOutlined />}
                      className="shadow-md"
                    >
                      立即购买
                    </Button>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </div>

        {/* 统计信息 */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          {[
            { label: '在线课程', value: '128+', color: 'from-blue-500 to-blue-600' },
            { label: '注册学员', value: '12K+', color: 'from-green-500 to-green-600' },
            { label: '讲师团队', value: '50+', color: 'from-purple-500 to-purple-600' },
            { label: '好评率', value: '98%', color: 'from-pink-500 to-pink-600' },
          ].map((stat, index) => (
            <Card 
              key={index}
              className="text-center hover:shadow-lg transition-shadow"
            >
              <div className={`inline-block px-6 py-3 rounded-lg bg-gradient-to-r ${stat.color} mb-3`}>
                <span className="text-3xl font-bold text-white">
                  {stat.value}
                </span>
              </div>
              <p className="text-gray-600 font-medium">{stat.label}</p>
            </Card>
          ))}
        </div>

        {/* 底部信息 */}
        <Card className="bg-gradient-to-r from-blue-500 to-indigo-600 border-0">
          <div className="text-center text-white">
            <h3 className="text-2xl font-bold mb-2">开始你的学习之旅</h3>
            <p className="mb-6 opacity-90">
              加入我们，掌握最新的技术技能
            </p>
            <div className="flex gap-4 justify-center">
              <Button size="large" className="bg-white text-primary hover:bg-gray-100">
                免费试学
              </Button>
              <Button size="large" type="default" ghost>
                了解更多
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default QuickStartExample;
