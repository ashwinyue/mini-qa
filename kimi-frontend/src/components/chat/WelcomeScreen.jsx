/**
 * 欢迎屏幕组件
 * 
 * 显示在没有消息时的初始界面
 */
import { Card, Space, Tag } from 'antd'

const WelcomeScreen = ({ onQuickAction }) => {
    const quickTags = [
        { label: '推荐', icon: '⭐' },
        { label: '网页应用', icon: '🌐' },
        { label: '移动应用', icon: '📱' },
        { label: '数据分析', icon: '📊' },
        { label: 'PPT', icon: '📊' },
        { label: '录屏', icon: '🎥' },
    ]

    const recommendCards = [
        {
            title: 'KIMI × MANNER 合作限时上线！',
            image: '/api/placeholder/200/120',
            tag: '月限时活动'
        },
        {
            title: 'PPT',
            subtitle: 'Cologne Cathedral',
            image: '/api/placeholder/200/120',
            tag: 'PPT'
        },
        {
            title: 'AUDIOBOOK SHERLOCK HOLMES',
            image: '/api/placeholder/200/120',
        }
    ]

    return (
        <div className="flex flex-col items-center justify-center h-full px-4 py-8">
            {/* Logo */}
            <div className="mb-6">
                <img 
                    src="/logo.svg" 
                    alt="Kimi AI" 
                    className="w-24 h-24"
                />
            </div>

            {/* 标题 */}
            <h1 className="text-5xl font-bold mb-12 text-gray-800" style={{ letterSpacing: '0.2em' }}>
                KIMI
            </h1>

            {/* 快捷标签 */}
            <div className="mb-8">
                <Space size={[8, 8]} wrap>
                    {quickTags.map((tag, index) => (
                        <Tag
                            key={index}
                            className="cursor-pointer px-4 py-1 text-sm"
                            onClick={() => onQuickAction?.(tag.label)}
                        >
                            <span className="mr-1">{tag.icon}</span>
                            {tag.label}
                        </Tag>
                    ))}
                </Space>
            </div>

            {/* 推荐卡片 */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 w-full max-w-4xl mt-8">
                {recommendCards.map((card, index) => (
                    <Card
                        key={index}
                        hoverable
                        cover={
                            <div className="h-32 bg-gradient-to-br from-gray-800 to-gray-600 flex items-center justify-center text-white">
                                {card.tag && (
                                    <div className="absolute top-2 left-2 bg-blue-500 text-white px-2 py-1 rounded text-xs">
                                        {card.tag}
                                    </div>
                                )}
                                <div className="text-center p-4">
                                    <div className="font-bold text-lg mb-1">{card.title}</div>
                                    {card.subtitle && (
                                        <div className="text-sm opacity-80">{card.subtitle}</div>
                                    )}
                                </div>
                            </div>
                        }
                        bodyStyle={{ padding: 0 }}
                    />
                ))}
            </div>
        </div>
    )
}

export default WelcomeScreen
