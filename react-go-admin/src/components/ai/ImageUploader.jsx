/**
 * 图片上传组件
 * 
 * 支持点击或拖拽上传图片
 */
import React, { useState, useRef } from 'react'
import { Button, Space, Image } from 'antd'
import { PlusOutlined, DeleteOutlined, CloseOutlined, UploadOutlined } from '@ant-design/icons'

/**
 * 图片上传组件
 * 
 * @param {Object} props
 * @param {Function} props.onUpload - 上传完成回调，参数为图片数组 (Base64)
 * @param {Function} props.onCancel - 取消回调
 * @param {number} props.maxCount - 最大上传数量，默认 3
 */
const ImageUploader = ({ onUpload, onCancel, maxCount = 3 }) => {
    const [images, setImages] = useState([])
    const [isDragging, setIsDragging] = useState(false)
    const fileInputRef = useRef(null)

    // 处理文件选择
    const handleFileSelect = (files) => {
        const fileList = Array.from(files)
        const remainingSlots = maxCount - images.length

        if (remainingSlots <= 0) {
            return
        }

        const filesToProcess = fileList.slice(0, remainingSlots)

        filesToProcess.forEach(file => {
            if (!file.type.startsWith('image/')) {
                return
            }

            const reader = new FileReader()
            reader.onloadend = () => {
                setImages(prev => [...prev, reader.result])
            }
            reader.readAsDataURL(file)
        })
    }

    // 点击上传
    const handleClick = () => {
        fileInputRef.current?.click()
    }

    // 文件输入变化
    const handleInputChange = (e) => {
        if (e.target.files) {
            handleFileSelect(e.target.files)
        }
        // 清空 input 以便重复选择同一文件
        e.target.value = ''
    }

    // 拖拽事件
    const handleDragOver = (e) => {
        e.preventDefault()
        setIsDragging(true)
    }

    const handleDragLeave = (e) => {
        e.preventDefault()
        setIsDragging(false)
    }

    const handleDrop = (e) => {
        e.preventDefault()
        setIsDragging(false)
        if (e.dataTransfer.files) {
            handleFileSelect(e.dataTransfer.files)
        }
    }

    // 删除图片
    const handleDelete = (index) => {
        setImages(prev => prev.filter((_, i) => i !== index))
    }

    // 确认上传
    const handleConfirm = () => {
        onUpload && onUpload(images)
    }

    // 取消
    const handleCancel = () => {
        setImages([])
        onCancel && onCancel()
    }

    return (
        <div className="p-4 bg-gray-50 rounded-lg mt-2">
            {/* 隐藏的文件输入 */}
            <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                multiple
                onChange={handleInputChange}
                className="hidden"
            />

            {/* 拖拽区域 */}
            <div
                className={`border-2 border-dashed rounded-lg p-4 text-center cursor-pointer transition-colors ${
                    isDragging 
                        ? 'border-blue-500 bg-blue-50' 
                        : 'border-gray-300 hover:border-blue-400'
                }`}
                onClick={handleClick}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
            >
                <UploadOutlined className="text-2xl text-gray-400" />
                <div className="mt-2 text-gray-500">
                    点击或拖拽图片到此处上传
                </div>
                <div className="text-xs text-gray-400 mt-1">
                    最多上传 {maxCount} 张图片
                </div>
            </div>

            {/* 图片预览 */}
            {images.length > 0 && (
                <div className="mt-4">
                    <div className="flex flex-wrap gap-2">
                        {images.map((image, index) => (
                            <div key={index} className="relative group">
                                <Image
                                    src={image}
                                    alt={`上传图片 ${index + 1}`}
                                    width={80}
                                    height={80}
                                    className="object-cover rounded"
                                />
                                <button
                                    className="absolute -top-2 -right-2 w-5 h-5 bg-red-500 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                                    onClick={(e) => {
                                        e.stopPropagation()
                                        handleDelete(index)
                                    }}
                                >
                                    <DeleteOutlined className="text-xs" />
                                </button>
                            </div>
                        ))}
                        
                        {/* 添加更多按钮 */}
                        {images.length < maxCount && (
                            <div
                                className="w-20 h-20 border-2 border-dashed border-gray-300 rounded flex items-center justify-center cursor-pointer hover:border-blue-400"
                                onClick={handleClick}
                            >
                                <PlusOutlined className="text-gray-400" />
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* 操作按钮 */}
            <div className="mt-4 flex justify-end">
                <Space>
                    <Button onClick={handleCancel}>
                        取消
                    </Button>
                    <Button 
                        type="primary" 
                        onClick={handleConfirm}
                        disabled={images.length === 0}
                    >
                        确认 ({images.length})
                    </Button>
                </Space>
            </div>
        </div>
    )
}

export default ImageUploader
