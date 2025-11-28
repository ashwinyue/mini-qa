/**
 * 音频录制组件
 * 
 * 使用 Web Audio API 录制音频
 */
import React, { useState, useRef, useEffect } from 'react'
import { Button, Progress, Space } from 'antd'
import { AudioOutlined, StopOutlined, CloseOutlined } from '@ant-design/icons'

/**
 * 音频录制组件
 * 
 * @param {Object} props
 * @param {Function} props.onRecord - 录制完成回调，参数为 (audioData, duration)
 * @param {Function} props.onCancel - 取消回调
 * @param {boolean} props.autoStart - 是否自动开始录音
 */
const AudioRecorder = ({ onRecord, onCancel, autoStart = false }) => {
    const [isRecording, setIsRecording] = useState(false)
    const [duration, setDuration] = useState(0)
    const [error, setError] = useState(null)

    const mediaRecorderRef = useRef(null)
    const audioChunksRef = useRef([])
    const timerRef = useRef(null)
    const startTimeRef = useRef(null)

    // 清理定时器
    useEffect(() => {
        return () => {
            if (timerRef.current) {
                clearInterval(timerRef.current)
            }
        }
    }, [])

    // 自动开始录音
    useEffect(() => {
        if (autoStart) {
            startRecording()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [autoStart])

    // 开始录制
    const startRecording = async () => {
        try {
            setError(null)
            audioChunksRef.current = []

            // 请求麦克风权限
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true })

            // 创建 MediaRecorder
            const mediaRecorder = new MediaRecorder(stream)
            mediaRecorderRef.current = mediaRecorder

            mediaRecorder.ondataavailable = (event) => {
                if (event.data.size > 0) {
                    audioChunksRef.current.push(event.data)
                }
            }

            mediaRecorder.onstop = () => {
                // 停止所有音轨
                stream.getTracks().forEach(track => track.stop())

                // 创建音频 Blob
                const audioBlob = new Blob(audioChunksRef.current, { type: 'audio/wav' })

                // 转换为 Base64
                const reader = new FileReader()
                reader.onloadend = () => {
                    const base64Audio = reader.result
                    const finalDuration = (Date.now() - startTimeRef.current) / 1000
                    onRecord && onRecord(base64Audio, finalDuration)
                }
                reader.readAsDataURL(audioBlob)
            }

            // 开始录制
            mediaRecorder.start()
            setIsRecording(true)
            startTimeRef.current = Date.now()

            // 开始计时
            timerRef.current = setInterval(() => {
                setDuration((Date.now() - startTimeRef.current) / 1000)
            }, 100)

        } catch (err) {
            console.error('录音失败:', err)
            console.error('错误类型:', err.name)
            console.error('错误信息:', err.message)
            
            let errorMessage = '无法访问麦克风'
            
            if (err.name === 'NotAllowedError') {
                errorMessage = '麦克风权限被拒绝。请在浏览器地址栏左侧点击锁图标，允许麦克风访问'
            } else if (err.name === 'NotFoundError') {
                // macOS 特殊处理
                const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0
                if (isMac) {
                    errorMessage = '未找到麦克风设备。\n\nmacOS 用户请检查：\n1. 打开"系统设置" → "隐私与安全性" → "麦克风"\n2. 确保浏览器已被授权\n3. 重启浏览器后重试'
                } else {
                    errorMessage = '未找到麦克风设备。请检查麦克风是否已连接'
                }
            } else if (err.name === 'NotReadableError') {
                errorMessage = '麦克风被其他应用占用。请关闭其他使用麦克风的应用'
            } else if (err.name === 'OverconstrainedError') {
                errorMessage = '麦克风不支持所需的配置'
            } else if (err.name === 'SecurityError') {
                errorMessage = '安全限制：请确保使用 HTTPS 或 localhost 访问'
            }
            
            setError(errorMessage)
        }
    }

    // 停止录制
    const stopRecording = () => {
        if (mediaRecorderRef.current && isRecording) {
            mediaRecorderRef.current.stop()
            setIsRecording(false)

            if (timerRef.current) {
                clearInterval(timerRef.current)
                timerRef.current = null
            }
        }
    }

    // 取消录制
    const cancelRecording = () => {
        if (mediaRecorderRef.current && isRecording) {
            mediaRecorderRef.current.stop()
            setIsRecording(false)

            if (timerRef.current) {
                clearInterval(timerRef.current)
                timerRef.current = null
            }
        }
        setDuration(0)
        onCancel && onCancel()
    }

    // 格式化时长
    const formatDuration = (seconds) => {
        const mins = Math.floor(seconds / 60)
        const secs = Math.floor(seconds % 60)
        return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
    }

    return (
        <div className="p-4 bg-gray-50 rounded-lg mt-2">
            {error ? (
                <div className="flex flex-col items-center gap-3">
                    <div className="text-red-500 text-center whitespace-pre-line text-sm max-w-md">
                        {error}
                    </div>
                    <Space>
                        <Button
                            type="primary"
                            icon={<AudioOutlined />}
                            onClick={startRecording}
                        >
                            重试
                        </Button>
                        <Button
                            icon={<CloseOutlined />}
                            onClick={onCancel}
                        >
                            关闭
                        </Button>
                    </Space>
                </div>
            ) : (
                <div className="flex flex-col items-center gap-3">
                    {/* 录制状态指示 */}
                    <div className="flex items-center gap-2">
                        {isRecording && (
                            <span className="w-3 h-3 bg-red-500 rounded-full animate-pulse" />
                        )}
                        <span className="text-lg font-mono">
                            {formatDuration(duration)}
                        </span>
                    </div>

                    {/* 进度条 */}
                    {isRecording && (
                        <Progress 
                            percent={Math.min((duration / 60) * 100, 100)} 
                            showInfo={false}
                            strokeColor="#1890ff"
                            className="w-full max-w-xs"
                        />
                    )}

                    {/* 控制按钮 */}
                    <Space>
                        {isRecording ? (
                            <Button
                                type="primary"
                                danger
                                icon={<StopOutlined />}
                                onClick={stopRecording}
                            >
                                停止录音
                            </Button>
                        ) : (
                            !autoStart && (
                                <Button
                                    type="primary"
                                    icon={<AudioOutlined />}
                                    onClick={startRecording}
                                >
                                    开始录音
                                </Button>
                            )
                        )}
                        <Button
                            icon={<CloseOutlined />}
                            onClick={cancelRecording}
                        >
                            取消
                        </Button>
                    </Space>

                    {/* 提示文字 */}
                    <div className="text-gray-400 text-xs">
                        {isRecording ? '正在录音...' : '点击开始录音，最长 60 秒'}
                    </div>
                </div>
            )}
        </div>
    )
}

export default AudioRecorder
