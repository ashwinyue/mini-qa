import React, { useState } from 'react';
import { Card, Form, Input, Button, message, Tabs } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { authService } from '../services/authService';

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('login');
  const navigate = useNavigate();

  const handleLogin = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      const result = await authService.login(values.username, values.password);
      if (result.code === 0) {
        message.success('登录成功！');
        navigate('/');
      } else {
        message.error(result.message || '登录失败');
      }
    } catch (error: any) {
      message.error(error.response?.data?.detail || '登录失败，请检查网络连接');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (values: { username: string; password: string; nickname: string; email?: string }) => {
    setLoading(true);
    try {
      const result = await authService.register(
        values.username,
        values.password,
        values.nickname,
        values.email
      );
      if (result.code === 0) {
        message.success('注册成功！请登录');
        setActiveTab('login');
      } else {
        message.error(result.message || '注册失败');
      }
    } catch (error: any) {
      message.error(error.response?.data?.detail || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Logo 和标题 */}
        <div className="text-center mb-8">
          <div className="inline-block p-4 bg-white rounded-full shadow-lg mb-4">
            <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full flex items-center justify-center">
              <span className="text-3xl font-bold text-white">AI</span>
            </div>
          </div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">AI 智能客服系统</h1>
          <p className="text-gray-600">欢迎使用智能对话助手</p>
        </div>

        {/* 登录/注册卡片 */}
        <Card className="shadow-2xl">
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            centered
            items={[
              {
                key: 'login',
                label: '登录',
                children: (
                  <Form
                    name="login"
                    onFinish={handleLogin}
                    autoComplete="off"
                    size="large"
                    initialValues={{
                      username: 'admin',
                      password: 'admin123',
                    }}
                  >
                    <Form.Item
                      name="username"
                      rules={[{ required: true, message: '请输入用户名' }]}
                    >
                      <Input
                        prefix={<UserOutlined className="text-gray-400" />}
                        placeholder="用户名"
                      />
                    </Form.Item>

                    <Form.Item
                      name="password"
                      rules={[{ required: true, message: '请输入密码' }]}
                    >
                      <Input.Password
                        prefix={<LockOutlined className="text-gray-400" />}
                        placeholder="密码"
                      />
                    </Form.Item>

                    <Form.Item>
                      <Button
                        type="primary"
                        htmlType="submit"
                        loading={loading}
                        block
                        size="large"
                        className="h-12"
                      >
                        登录
                      </Button>
                    </Form.Item>

                    <div className="text-center text-sm space-y-2">
                      <div className="bg-blue-50 p-3 rounded-lg">
                        <p className="text-blue-600 font-medium mb-1">💡 演示账号（已预填）</p>
                        <p className="text-gray-600">管理员：admin / admin123</p>
                        <p className="text-gray-600">普通用户：demo / demo123</p>
                      </div>
                    </div>
                  </Form>
                ),
              },
              {
                key: 'register',
                label: '注册',
                children: (
                  <Form
                    name="register"
                    onFinish={handleRegister}
                    autoComplete="off"
                    size="large"
                  >
                    <Form.Item
                      name="username"
                      rules={[
                        { required: true, message: '请输入用户名' },
                        { min: 3, message: '用户名至少3个字符' },
                        { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含字母、数字和下划线' },
                      ]}
                    >
                      <Input
                        prefix={<UserOutlined className="text-gray-400" />}
                        placeholder="用户名"
                      />
                    </Form.Item>

                    <Form.Item
                      name="nickname"
                      rules={[{ required: true, message: '请输入昵称' }]}
                    >
                      <Input
                        prefix={<UserOutlined className="text-gray-400" />}
                        placeholder="昵称"
                      />
                    </Form.Item>

                    <Form.Item
                      name="email"
                      rules={[
                        { type: 'email', message: '请输入有效的邮箱地址' },
                      ]}
                    >
                      <Input
                        prefix={<MailOutlined className="text-gray-400" />}
                        placeholder="邮箱（可选）"
                      />
                    </Form.Item>

                    <Form.Item
                      name="password"
                      rules={[
                        { required: true, message: '请输入密码' },
                        { min: 6, message: '密码至少6个字符' },
                      ]}
                    >
                      <Input.Password
                        prefix={<LockOutlined className="text-gray-400" />}
                        placeholder="密码"
                      />
                    </Form.Item>

                    <Form.Item
                      name="confirm"
                      dependencies={['password']}
                      rules={[
                        { required: true, message: '请确认密码' },
                        ({ getFieldValue }) => ({
                          validator(_, value) {
                            if (!value || getFieldValue('password') === value) {
                              return Promise.resolve();
                            }
                            return Promise.reject(new Error('两次输入的密码不一致'));
                          },
                        }),
                      ]}
                    >
                      <Input.Password
                        prefix={<LockOutlined className="text-gray-400" />}
                        placeholder="确认密码"
                      />
                    </Form.Item>

                    <Form.Item>
                      <Button
                        type="primary"
                        htmlType="submit"
                        loading={loading}
                        block
                        size="large"
                        className="h-12"
                      >
                        注册
                      </Button>
                    </Form.Item>
                  </Form>
                ),
              },
            ]}
          />
        </Card>

        {/* 底部信息 */}
        <div className="text-center mt-6 text-sm text-gray-500">
          <p>© 2024 AI 智能客服系统. All rights reserved.</p>
        </div>
      </div>
    </div>
  );
};

export default Login;
