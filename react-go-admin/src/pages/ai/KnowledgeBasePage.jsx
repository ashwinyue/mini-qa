/**
 * 知识库查询页面
 */
import React, { useState } from 'react'
import { Input, Button, Card, List, Empty, Spin, message } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import apiService from '../../services/api'

const { Search } = Input

const KnowledgeBasePage = () => {
    const [loading, setLoading] = useState(false)
    const [results, setResults] = useState([])
    const [searched, setSearched] = useState(false)

    const handleSearch = async (value) => {
        if (!value.trim()) {
            message.warning('请输入搜索内容')
            return
        }

        setLoading(true)
        setSearched(true)

        try {
            // 通过对话接口搜索知识库
            const response = await apiService.post('/chat', {
                query: value,
                thread_id: 'kb-search-' + Date.now(),
            })

            const data = response.data

            // 提取来源作为搜索结果
            if (data.sources && data.sources.length > 0) {
                setResults(data.sources.map((source, index) => ({
                    id: index,
                    title: source.source || source.title || `结果 ${index + 1}`,
                    content: source.content || '',
                    metadata: source.metadata,
                })))
            } else if (data.answer) {
                // 如果没有来源但有答案，显示答案
                setResults([{
                    id: 0,
                    title: '搜索结果',
                    content: data.answer,
                }])
            } else {
                setResults([])
            }
        } catch (error) {
            console.error('搜索失败:', error)
            message.error('搜索失败，请重试')
            setResults([])
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="p-4">
            <Card title="知识库搜索" className="mb-4">
                <Search
                    placeholder="输入关键词搜索知识库..."
                    allowClear
                    enterButton={<><SearchOutlined /> 搜索</>}
                    size="large"
                    onSearch={handleSearch}
                    loading={loading}
                />
            </Card>

            <Spin spinning={loading}>
                {searched && (
                    <Card title={`搜索结果 (${results.length})`}>
                        {results.length > 0 ? (
                            <List
                                itemLayout="vertical"
                                dataSource={results}
                                renderItem={(item) => (
                                    <List.Item key={item.id}>
                                        <List.Item.Meta
                                            title={item.title}
                                        />
                                        <div className="text-gray-600 whitespace-pre-wrap">
                                            {item.content}
                                        </div>
                                    </List.Item>
                                )}
                            />
                        ) : (
                            <Empty description="未找到相关内容" />
                        )}
                    </Card>
                )}
            </Spin>
        </div>
    )
}

export default KnowledgeBasePage
