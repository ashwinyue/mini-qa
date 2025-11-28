/**
 * 聊天输入组件
 * 
 * 支持文本、语音、图片输入
 */
import { useState } from 'react'
import { Input, Button, Upload, Space, Image } from 'antd'
import { SendOutlined, AudioOutlined, PictureOutlined, LoadingOutlined, CloseCircleOutlined } from '@ant-design/icons'

const { TextArea } = Input

const ChatInput = ({ onSend, disabled, placeholder = '输入消息...', centered = false }) => {
    const [inputValue, setInputValue] = useState('')
    const [uploadedImages, setUploadedImages] = useState([])

    const handleSend = () => {
        if (!inputValue.trim() && uploadedImages.length === 0) {
            return
        }

        onSend(inputValue.trim(), uploadedImages)
        setInputValue('')
        setUploadedImages([])
    }

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault()
            handleSend()
        }
    }

    const handleImageUpload = (file) => {
        const reader = new FileReader()
        reader.onload = (event) => {
            setUploadedImages(prev => [...prev, event.target.result])
        }
        reader.readAsDataURL(file)
        return false // 阻止自动上传
    }

    const removeImage = (index) => {
        setUploadedImages(prev => prev.filter((_, i) => i !== index))
    }

    return (
        <div className={`bg-white p-4 ${!centered ? 'border-t' : ''}`}>
            {/* 已上传图片预览 */}
            {uploadedImages.length > 0 && (
                <div className="flex gap-2 mb-3 flex-wrap">
                    {uploadedImages.map((image, index) => (
                        <div key={index} className="relative">
                            <Image
                                src={image}
                                alt={`上传图片 ${index + 1}`}
                                width={64}
                                height={64}
                                className="object-cover rounded"
                                preview={true}
                            />
                            <Button
                                type="text"
                                danger
                                size="small"
                                icon={<CloseCircleOutlined />}
                                onClick={() => removeImage(index)}
                                className="absolute -top-2 -right-2"
                                style={{ padding: 0, minWidth: 20, height: 20 }}
                            />
                        </div>
                    ))}
                </div>
            )}

            {/* 输入区域 */}
            <div className={`flex items-center gap-2 ${centered ? 'border border-gray-300 rounded-full px-4 py-2 shadow-sm' : ''}`}>
                {/* 工具按钮 */}
                <Space size="small">
                    <Button
                        type="text"
                        icon={<span style={{ fontSize: 18 }}>🎯</span>}
                        disabled={disabled}
                        size="small"
                    />
                    <Upload
                        accept="image/*"
                        multiple
                        showUploadList={false}
                        beforeUpload={handleImageUpload}
                        disabled={disabled}
                    >
                        <Button 
                            type="text"
                            icon={<PictureOutlined />} 
                            disabled={disabled}
                            size="small"
                        />
                    </Upload>
                    
                    <Button
                        type="text"
                        icon={<AudioOutlined />}
                        disabled={disabled}
                        title="语音输入（开发中）"
                        size="small"
                    />
                </Space>

                {/* 文本输入框 */}
                <TextArea
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder={placeholder}
                    disabled={disabled}
                    autoSize={{ minRows: 1, maxRows: 4 }}
                    className="flex-1"
                    bordered={!centered}
                    style={centered ? { border: 'none', boxShadow: 'none' } : {}}
                />

                {/* 右侧按钮组 */}
                <Space size="small">
                    {centered && (
                        <>
                            <Button
                                type="text"
                                size="small"
                                disabled={disabled}
                            >
                                K2
                            </Button>
                            <Button
                                type="text"
                                icon={<span>📎</span>}
                                disabled={disabled}
                                size="small"
                            />
                            <Button
                                type="text"
                                icon={<span>⚙️</span>}
                                disabled={disabled}
                                size="small"
                            />
                        </>
                    )}
                    
                    {/* 发送按钮 */}
                    <Button
                        type={centered ? "default" : "primary"}
                        shape="circle"
                        icon={disabled ? <LoadingOutlined /> : <SendOutlined />}
                        onClick={handleSend}
                        disabled={disabled || (!inputValue.trim() && uploadedImages.length === 0)}
                        loading={disabled}
                        size="small"
                    />
                </Space>
            </div>

            {/* 提示文本 - 只在非居中模式显示 */}
            {!centered && (
                <div className="text-xs text-gray-400 mt-2 text-center">
                    Kimi 可能会出错，请核查重要信息
                </div>
            )}
        </div>
    )
}

export default ChatInput
