# 语音转文字回显功能

## 功能描述

用户发送语音消息后，语音识别的文本会自动回显在对话框中，替换原来的"[语音消息 - 时长: X秒]"占位符。

## 实现方式

### 后端修改 (work_v3/app.py)

#### 在响应中添加语音识别文本

**修改前**:
```python
# 返回路由、答案与来源
return {"route": route, "answer": answer, "sources": sources}
```

**修改后**:
```python
# 返回路由、答案与来源，如果有语音识别文本也返回
response = {"route": route, "answer": answer, "sources": sources}
if audio_text:
    response["audio_text"] = audio_text
return response
```

### 前端修改 (ChatInterface.jsx)

#### 1. 添加消息 ID 跟踪

```javascript
// 添加语音消息（先显示占位符）
const audioMessageId = uuidv4()
const audioMessage = {
    id: audioMessageId,
    role: 'user',
    content: `[语音消息 - 时长: ${Math.round(duration)}秒]`,
    timestamp: new Date().toISOString(),
    isAudio: true,
}
addMessage(audioMessage)
```

#### 2. 收到响应后更新消息

```javascript
.then(response => {
    // 如果后端返回了语音识别文本，更新用户消息
    if (response.audio_text) {
        updateMessage(audioMessageId, {
            content: response.audio_text,
            audioText: response.audio_text,
        })
    }
    
    // ... 添加助手消息
})
```

## 工作流程

```
用户点击录音按钮
    ↓
自动开始录音
    ↓
用户说话
    ↓
点击"停止录音"
    ↓
前端：添加占位符消息
    "[语音消息 - 时长: 5秒]"
    ↓
发送音频数据到后端
    ↓
后端：语音识别
    audio_text = "你好，今天天气怎么样？"
    ↓
后端：处理查询并返回
    {
        "route": "...",
        "answer": "...",
        "sources": [...],
        "audio_text": "你好，今天天气怎么样？"
    }
    ↓
前端：更新用户消息
    "[语音消息 - 时长: 5秒]"
    ↓ 更新为
    "你好，今天天气怎么样？"
    ↓
前端：显示助手回复
```

## 用户体验

### 发送语音前
```
┌─────────────────────────┐
│ 用户: [输入框]          │
└─────────────────────────┘
```

### 录音中
```
┌─────────────────────────┐
│ ● 00:05                 │
│ ████████░░░░░░░░░░      │
│ [停止录音] [取消]       │
│ 正在录音...             │
└─────────────────────────┘
```

### 发送后（识别前）
```
┌─────────────────────────┐
│ 用户: [语音消息 - 时长: 5秒] │
│                         │
│ 助手: 正在处理...       │
└─────────────────────────┘
```

### 识别完成后
```
┌─────────────────────────┐
│ 用户: 你好，今天天气怎么样？ │
│                         │
│ 助手: 今天天气晴朗...   │
└─────────────────────────┘
```

## 数据流

### 请求数据
```javascript
{
    query: '',  // 语音消息时为空
    thread_id: 'xxx',
    images: undefined,
    audio: 'data:audio/wav;base64,UklGR...'  // Base64 音频数据
}
```

### 响应数据
```javascript
{
    route: 'general',
    answer: '今天天气晴朗，温度适宜...',
    sources: [...],
    audio_text: '你好，今天天气怎么样？'  // 新增字段
}
```

## 消息对象结构

### 初始语音消息
```javascript
{
    id: 'uuid-xxx',
    role: 'user',
    content: '[语音消息 - 时长: 5秒]',
    timestamp: '2024-01-01T12:00:00.000Z',
    isAudio: true
}
```

### 更新后的消息
```javascript
{
    id: 'uuid-xxx',
    role: 'user',
    content: '你好，今天天气怎么样？',  // 更新为识别文本
    timestamp: '2024-01-01T12:00:00.000Z',
    isAudio: true,
    audioText: '你好，今天天气怎么样？'  // 保存识别文本
}
```

## 状态管理

使用 Zustand store 的 `updateMessage` 方法：

```javascript
updateMessage(messageId, updates)
```

这个方法会：
1. 找到对应 ID 的消息
2. 合并更新的字段
3. 触发 UI 重新渲染

## 错误处理

### 场景 1: 语音识别失败
- 后端不返回 `audio_text` 字段
- 前端保持显示占位符文本
- 用户仍能看到助手的回复

### 场景 2: 网络错误
- 显示错误消息
- 语音消息保持占位符状态
- 用户可以重新发送

### 场景 3: 识别为空
- 后端返回空的 `audio_text`
- 前端不更新消息
- 保持显示占位符

## 测试步骤

### 测试 1: 正常语音识别
1. 点击录音按钮
2. 说话："你好，今天天气怎么样？"
3. 点击停止录音
4. 验证：
   - ✅ 先显示"[语音消息 - 时长: X秒]"
   - ✅ 几秒后更新为识别的文本
   - ✅ 助手回复正确显示

### 测试 2: 短语音
1. 录音 1 秒
2. 说："你好"
3. 停止录音
4. 验证：
   - ✅ 识别文本正确显示
   - ✅ 时长显示正确

### 测试 3: 长语音
1. 录音 30 秒
2. 说一段长文本
3. 停止录音
4. 验证：
   - ✅ 完整文本显示
   - ✅ 不会被截断

### 测试 4: 识别失败
1. 录音但不说话（静音）
2. 停止录音
3. 验证：
   - ✅ 保持显示占位符
   - ✅ 或显示"[无法识别]"

### 测试 5: 网络错误
1. 断开网络
2. 发送语音
3. 验证：
   - ✅ 显示错误提示
   - ✅ 消息保持占位符状态

## 后端语音识别

后端使用 `graph.transcribe_audio()` 方法进行语音识别：

```python
audio_text = graph.transcribe_audio(
    req.audio,           # Base64 音频数据
    req.asr_language,    # 识别语言（可选）
    True if req.asr_itn is None else bool(req.asr_itn)  # 是否启用 ITN
)
```

### 支持的参数
- `audio`: Base64 编码的音频数据
- `asr_language`: 识别语言（如 'zh', 'en'）
- `asr_itn`: 是否启用逆文本归一化（Inverse Text Normalization）

## 优化建议

### 已实现
- ✅ 语音识别文本回显
- ✅ 占位符到文本的平滑更新
- ✅ 保留时长信息
- ✅ 错误处理

### 未来改进
- [ ] 显示识别进度（"正在识别..."）
- [ ] 支持编辑识别文本
- [ ] 显示识别置信度
- [ ] 支持多语言识别
- [ ] 添加语音播放功能
- [ ] 缓存识别结果
- [ ] 支持实时语音识别（流式）

## 注意事项

1. **音频格式**
   - 前端录制为 WAV 格式
   - Base64 编码传输
   - 后端需要支持解码

2. **识别延迟**
   - 语音识别需要时间
   - 用户会先看到占位符
   - 几秒后更新为文本

3. **识别准确性**
   - 依赖后端 ASR 服务质量
   - 环境噪音会影响识别
   - 方言可能识别不准

4. **数据大小**
   - 长语音会产生大量数据
   - 建议限制录音时长
   - 考虑压缩音频

## 相关文件

- `work_v3/app.py` - 后端 API（已修改）
- `react-go-admin/src/components/ai/ChatInterface.jsx` - 聊天界面（已修改）
- `react-go-admin/src/components/ai/AudioRecorder.jsx` - 录音组件
- `react-go-admin/src/components/ai/MessageItem.jsx` - 消息显示组件
- `react-go-admin/src/stores/chatStore.js` - 聊天状态管理

## API 文档

### POST /chat

**请求**:
```json
{
    "query": "",
    "thread_id": "xxx",
    "audio": "data:audio/wav;base64,..."
}
```

**响应**:
```json
{
    "route": "general",
    "answer": "回复内容",
    "sources": [],
    "audio_text": "识别的文本"  // 新增
}
```
