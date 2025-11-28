# 欢迎语功能测试文档

## 功能概述

react-go-admin 前端已成功接入 Python 后端的 `/greet` 接口，实现了欢迎语功能。

## 技术架构

### 后端 (Python - work_v3/app.py)
- **接口**: `GET /greet`
- **端口**: 8000
- **返回数据**:
```json
{
  "message": "您好，请问有什么可以帮您？",
  "options": [
    {"key": "course", "title": "课程咨询", "desc": "显示课程目录和详细信息"},
    {"key": "order", "title": "订单查询", "desc": "验证用户身份后显示订单状态"},
    {"key": "human", "title": "人工转接", "desc": "直接转人工客服"}
  ]
}
```

### 前端 (React - react-go-admin)
- **端口**: 5175
- **代理配置**: vite.config.js 中已配置 `/greet` 代理到 `http://localhost:8000`
- **API 调用**: `src/api/chat.js` 中的 `getGreeting()` 函数
- **UI 组件**: 
  - `ChatInterface.jsx` - 主聊天界面，负责加载欢迎语
  - `MessageItem.jsx` - 消息项组件，显示欢迎语和建议
  - `SuggestionList.jsx` - 建议问题列表组件

## 工作流程

1. **会话初始化**: 用户打开聊天界面或创建新会话
2. **加载欢迎语**: `ChatInterface` 检测到空会话时自动调用 `loadGreeting()`
3. **API 请求**: 通过 vite 代理向 Python 后端请求 `/greet` 接口
4. **数据处理**: 将返回的 `options` 数组转换为 `suggestions` 数组
5. **UI 渲染**: 
   - 欢迎消息显示在聊天界面
   - 建议问题显示为可点击的按钮
6. **用户交互**: 点击建议按钮自动填充到输入框

## 测试步骤

### 1. 启动后端服务
```bash
cd work_v3
python app.py
# 确保服务运行在 http://localhost:8000
```

### 2. 启动前端服务
```bash
cd react-go-admin
npm run dev
# 前端运行在 http://localhost:5175
```

### 3. 测试接口连通性
```bash
# 直接测试后端
curl http://localhost:8000/greet

# 测试前端代理
curl http://localhost:5175/greet
```

### 4. 浏览器测试
1. 打开 http://localhost:5175
2. 登录系统
3. 进入 AI 聊天页面
4. 观察是否显示欢迎语和三个建议按钮：
   - 课程咨询
   - 订单查询
   - 人工转接
5. 点击任意建议按钮，确认文本填充到输入框

## 代码关键点

### API 配置 (src/services/api.js)
```javascript
// /greet 被标记为简单 GET 请求，不需要租户头
const SIMPLE_GET_PATHS = ['/health', '/models/list', '/greet']
```

### 欢迎语加载 (src/components/ai/ChatInterface.jsx)
```javascript
const loadGreeting = async () => {
    try {
        const greeting = await getGreeting()
        const greetingMessage = {
            id: 'greeting-' + Date.now(),
            role: 'assistant',
            content: greeting.message || '您好，请问有什么可以帮您？',
            timestamp: new Date().toISOString(),
            suggestions: greeting.options?.map(opt => opt.title) || [],
        }
        addMessage(greetingMessage)
    } catch (error) {
        console.error('加载欢迎语失败:', error)
    }
}
```

### 建议显示 (src/components/ai/MessageItem.jsx)
```javascript
{suggestions && suggestions.length > 0 && (
    <SuggestionList 
        suggestions={suggestions} 
        onClick={onSuggestionClick}
    />
)}
```

## 故障排查

### 问题1: 欢迎语不显示
- 检查 Python 后端是否运行在 8000 端口
- 检查浏览器控制台是否有 API 错误
- 确认会话是否为空（只有空会话才会加载欢迎语）

### 问题2: 建议按钮不显示
- 检查 API 返回的 `options` 数组是否正确
- 确认 `options` 中的 `title` 字段存在
- 查看 `SuggestionList` 组件是否正确渲染

### 问题3: 代理不工作
- 检查 vite.config.js 中的代理配置
- 确认前端和后端端口正确
- 重启前端开发服务器

## 当前状态

✅ **已完成**:
- Python 后端 `/greet` 接口实现
- 前端 API 调用封装
- Vite 代理配置
- UI 组件完整实现
- 自动加载欢迎语逻辑
- 建议问题点击交互

✅ **测试通过**:
- 后端接口正常响应
- 前端代理正常工作
- 所有组件代码完整

## 下一步

如需自定义欢迎语内容，修改 `work_v3/app.py` 中的 `greet()` 函数即可。
