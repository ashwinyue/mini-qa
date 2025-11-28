# 粘贴图片功能

## 功能描述

在聊天输入框中支持直接粘贴图片，无需通过上传按钮。

## 实现方式

### 修改的文件
`src/components/ai/ChatInterface.jsx`

### 新增功能

#### 1. 粘贴事件处理函数

```javascript
const handlePaste = async (e) => {
    const items = e.clipboardData?.items
    if (!items) return

    const imageFiles = []
    for (let i = 0; i < items.length; i++) {
        const item = items[i]
        if (item.type.indexOf('image') !== -1) {
            const file = item.getAsFile()
            if (file) {
                imageFiles.push(file)
            }
        }
    }

    if (imageFiles.length > 0) {
        e.preventDefault()
        
        // 转换为 Base64
        const base64Images = await Promise.all(
            imageFiles.map(file => {
                return new Promise((resolve, reject) => {
                    const reader = new FileReader()
                    reader.onload = (e) => resolve(e.target.result)
                    reader.onerror = reject
                    reader.readAsDataURL(file)
                })
            })
        )

        setUploadedImages(prev => [...prev, ...base64Images])
        antMessage.success(`已粘贴 ${imageFiles.length} 张图片`)
    }
}
```

#### 2. 绑定到输入框

```jsx
<TextArea
    value={inputValue}
    onChange={(e) => setInputValue(e.target.value)}
    onKeyDown={handleKeyDown}
    onPaste={handlePaste}  // 新增粘贴事件
    placeholder="输入消息，按 Enter 发送，Shift+Enter 换行，支持粘贴图片..."
    autoSize={{ minRows: 1, maxRows: 4 }}
    disabled={isLoading}
    className="flex-1"
/>
```

## 功能特性

### 1. 多图支持
- 可以一次粘贴多张图片
- 自动识别剪贴板中的所有图片

### 2. 格式支持
- 支持所有浏览器支持的图片格式（PNG, JPEG, GIF, WebP 等）
- 自动转换为 Base64 编码

### 3. 用户体验
- 粘贴后立即显示预览
- 显示成功提示消息
- 与现有上传功能无缝集成

### 4. 智能处理
- 只处理图片类型的剪贴板内容
- 文本粘贴不受影响
- 异步处理，不阻塞 UI

## 使用方式

### 方式一：截图粘贴
1. 使用系统截图工具（如 macOS 的 Cmd+Shift+4）
2. 截图后自动复制到剪贴板
3. 在聊天输入框中按 `Cmd+V` (macOS) 或 `Ctrl+V` (Windows)
4. 图片自动添加到预览区域

### 方式二：复制图片文件
1. 在文件管理器中选择图片文件
2. 右键选择"复制"或按 `Cmd+C` / `Ctrl+C`
3. 在聊天输入框中粘贴
4. 图片自动添加到预览区域

### 方式三：从网页复制图片
1. 在网页上右键点击图片
2. 选择"复制图片"
3. 在聊天输入框中粘贴
4. 图片自动添加到预览区域

## 工作流程

```
用户粘贴 (Cmd+V / Ctrl+V)
    ↓
检测剪贴板内容
    ↓
识别图片文件
    ↓
转换为 Base64
    ↓
添加到 uploadedImages 状态
    ↓
显示预览缩略图
    ↓
用户点击发送
    ↓
图片随消息一起发送
    ↓
在对话框中显示
```

## 与现有功能的集成

### 1. 图片预览
- 粘贴的图片会显示在输入框上方
- 与上传按钮添加的图片显示方式一致
- 可以点击 × 删除单张图片

### 2. 消息发送
- 粘贴的图片与上传的图片使用相同的发送逻辑
- 支持图片 + 文本同时发送
- 支持只发送图片（无文本）

### 3. 消息显示
- 使用 MessageItem 组件的图片显示功能
- 在用户消息气泡中正确显示

## 测试步骤

### 测试 1: 截图粘贴
1. 使用系统截图工具截取屏幕
2. 在聊天输入框中粘贴
3. 验证：
   - ✅ 图片出现在预览区域
   - ✅ 显示成功提示
   - ✅ 可以删除图片
   - ✅ 可以发送消息

### 测试 2: 多图粘贴
1. 复制多张图片（如果浏览器支持）
2. 在输入框中粘贴
3. 验证：
   - ✅ 所有图片都显示
   - ✅ 提示显示正确数量

### 测试 3: 文本粘贴不受影响
1. 复制一段文本
2. 在输入框中粘贴
3. 验证：
   - ✅ 文本正常粘贴
   - ✅ 不触发图片处理

### 测试 4: 混合使用
1. 通过上传按钮添加一张图片
2. 通过粘贴添加另一张图片
3. 输入文本
4. 发送消息
5. 验证：
   - ✅ 两张图片都显示在预览区
   - ✅ 发送后都显示在消息中

## 技术细节

### 剪贴板 API
使用标准的 `ClipboardEvent` API：
- `e.clipboardData.items` - 获取剪贴板项目
- `item.type` - 检查 MIME 类型
- `item.getAsFile()` - 获取文件对象

### 文件读取
使用 `FileReader` API：
- `readAsDataURL()` - 转换为 Base64
- 异步处理，使用 Promise

### 状态管理
- 使用 `setUploadedImages` 更新状态
- 使用展开运算符合并新旧图片
- 保持与上传功能的一致性

## 浏览器兼容性

- ✅ Chrome/Edge (Chromium)
- ✅ Firefox
- ✅ Safari
- ⚠️ 部分移动浏览器可能不支持

## 注意事项

1. **性能考虑**
   - 大图片会增加 Base64 编码时间
   - 建议在实际应用中添加图片压缩

2. **用户提示**
   - 已在 placeholder 中添加"支持粘贴图片"提示
   - 粘贴成功后显示 Toast 提示

3. **错误处理**
   - 文件读取失败会被 Promise 捕获
   - 不会影响正常的文本粘贴

## 未来改进

- [ ] 添加图片大小限制
- [ ] 添加图片格式验证
- [ ] 添加图片压缩功能
- [ ] 支持拖拽上传
- [ ] 添加粘贴动画效果
