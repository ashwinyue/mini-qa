# Ant Design + Tailwind CSS 混合使用指南

## 配置说明

### 关键配置项

**tailwind.config.js**
```javascript
corePlugins: {
  preflight: false, // 禁用 Tailwind 的样式重置，避免与 Ant Design 冲突
}
```

这个配置非常重要，可以避免 Tailwind 的基础样式覆盖 Ant Design 的样式。

## 使用原则

### ✅ 推荐使用场景

| 场景 | 使用框架 | 示例 |
|------|---------|------|
| 复杂组件 | Ant Design | `<Table>`, `<Form>`, `<Modal>`, `<Drawer>` |
| 布局 | Tailwind CSS | `grid`, `flex`, `gap-4`, `p-6` |
| 间距 | Tailwind CSS | `mt-4`, `mb-6`, `px-4`, `py-2` |
| 响应式 | Tailwind CSS | `md:grid-cols-2`, `lg:flex-row` |
| 颜色 | Tailwind CSS | `bg-gray-50`, `text-primary` |
| 动画 | Tailwind CSS | `hover:shadow-lg`, `transition-all` |
| 表单输入 | Ant Design | `<Input>`, `<Select>`, `<DatePicker>` |
| 按钮 | Ant Design | `<Button type="primary">` |

### ❌ 避免的做法

```tsx
// ❌ 不要用 Tailwind 重写 Ant Design 组件的核心样式
<Button className="bg-blue-500 text-white px-4 py-2 rounded">
  按钮
</Button>

// ✅ 使用 Ant Design 的 Button，用 Tailwind 做外层布局
<div className="flex gap-2 mt-4">
  <Button type="primary">确定</Button>
  <Button>取消</Button>
</div>
```

## 常用模式

### 1. 响应式网格布局

```tsx
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  <Card>内容1</Card>
  <Card>内容2</Card>
  <Card>内容3</Card>
</div>
```

### 2. Flexbox 布局

```tsx
<div className="flex flex-col md:flex-row gap-4 items-center justify-between">
  <Input placeholder="搜索" className="flex-1" />
  <Button type="primary">搜索</Button>
</div>
```

### 3. 卡片悬停效果

```tsx
<Card className="hover:shadow-lg transition-shadow duration-300">
  内容
</Card>
```

### 4. 间距和内边距

```tsx
<div className="p-6 mb-4">
  <Card className="mt-2">
    <div className="space-y-4">
      <div>项目1</div>
      <div>项目2</div>
    </div>
  </Card>
</div>
```

### 5. 文本样式

```tsx
<div>
  <h1 className="text-2xl font-bold text-gray-800 mb-4">标题</h1>
  <p className="text-gray-600 leading-relaxed">正文内容</p>
</div>
```

### 6. 响应式显示/隐藏

```tsx
<div className="hidden md:block">
  桌面端显示
</div>
<div className="block md:hidden">
  移动端显示
</div>
```

## 自定义主题色

在 `tailwind.config.js` 中已配置与 Ant Design 一致的主题色：

```javascript
colors: {
  primary: {
    DEFAULT: '#1890ff',
    hover: '#40a9ff',
    active: '#096dd9',
  },
  success: '#52c41a',
  warning: '#faad14',
  error: '#ff4d4f',
}
```

使用示例：
```tsx
<div className="bg-primary text-white p-4">
  主色背景
</div>
<span className="text-success">成功文本</span>
<Tag className="bg-warning">警告标签</Tag>
```

## 实际应用示例

### 聊天界面布局

```tsx
<div className="flex flex-col h-screen">
  {/* 头部 */}
  <div className="flex items-center justify-between p-4 border-b">
    <h1 className="text-xl font-bold">AI 助手</h1>
    <Button type="primary">新对话</Button>
  </div>

  {/* 消息区域 */}
  <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
    {messages.map((msg) => (
      <Card key={msg.id} className="max-w-2xl">
        {msg.content}
      </Card>
    ))}
  </div>

  {/* 输入区域 */}
  <div className="p-4 border-t bg-white">
    <div className="flex gap-2">
      <Input.TextArea 
        placeholder="输入消息..." 
        className="flex-1"
        autoSize={{ minRows: 1, maxRows: 4 }}
      />
      <Button type="primary" icon={<SendOutlined />}>
        发送
      </Button>
    </div>
  </div>
</div>
```

### 仪表板布局

```tsx
<div className="p-6 bg-gray-50 min-h-screen">
  {/* 统计卡片 */}
  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
    <Card className="hover:shadow-lg transition-shadow">
      <Statistic title="总用户" value={1234} />
    </Card>
    {/* 更多卡片... */}
  </div>

  {/* 图表区域 */}
  <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
    <Card title="用户增长">
      {/* 图表组件 */}
    </Card>
    <Card title="订单统计">
      {/* 图表组件 */}
    </Card>
  </div>
</div>
```

## 常用 Tailwind 类名速查

### 间距
- `p-4` = padding: 1rem
- `m-4` = margin: 1rem
- `px-4` = padding-left & padding-right: 1rem
- `gap-4` = gap: 1rem
- `space-y-4` = 子元素垂直间距

### 布局
- `flex` = display: flex
- `grid` = display: grid
- `grid-cols-3` = 3列网格
- `items-center` = align-items: center
- `justify-between` = justify-content: space-between

### 尺寸
- `w-full` = width: 100%
- `h-screen` = height: 100vh
- `max-w-2xl` = max-width: 42rem
- `min-h-screen` = min-height: 100vh

### 颜色
- `bg-gray-50` = 背景色
- `text-gray-600` = 文字颜色
- `border-gray-200` = 边框颜色

### 响应式
- `md:grid-cols-2` = 中等屏幕2列
- `lg:flex-row` = 大屏幕横向排列
- `sm:text-lg` = 小屏幕大字体

## 调试技巧

### 1. 使用浏览器开发工具
查看元素的实际应用的 Tailwind 类名

### 2. Tailwind CSS IntelliSense
安装 VS Code 插件获得自动补全

### 3. 样式冲突检查
如果样式不生效，检查是否被 Ant Design 样式覆盖

## 性能优化

### 1. 生产环境自动清除未使用的样式
Tailwind 会自动清除未使用的 CSS，保持打包体积小

### 2. 避免过度使用内联样式
优先使用 Tailwind 类名，而不是 `style` 属性

### 3. 复用常用组合
对于频繁使用的样式组合，可以创建自定义组件

## 参考资源

- [Tailwind CSS 官方文档](https://tailwindcss.com/docs)
- [Ant Design 官方文档](https://ant.design/)
- [Tailwind + Ant Design 最佳实践](https://github.com/ant-design/ant-design/discussions)
