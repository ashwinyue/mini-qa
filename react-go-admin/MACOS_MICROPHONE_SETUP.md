# macOS 麦克风权限设置指南

## 问题症状

错误信息：`NotFoundError: Requested device not found`

这表示浏览器无法找到麦克风设备，通常是因为 macOS 系统级别的隐私设置未授权。

## 解决方案

### 方法 1: 系统设置中授权（推荐）

#### 对于 macOS Ventura (13.0) 及更新版本：

1. **打开系统设置**
   - 点击屏幕左上角的  菜单
   - 选择"系统设置..."

2. **进入隐私与安全性**
   - 在左侧边栏找到"隐私与安全性"
   - 点击进入

3. **找到麦克风设置**
   - 向下滚动找到"麦克风"
   - 点击"麦克风"

4. **授权浏览器**
   - 找到你使用的浏览器（Chrome、Safari、Firefox 等）
   - 打开开关，允许访问麦克风
   - 如果列表中没有浏览器，说明浏览器还没有请求过权限

5. **重启浏览器**
   - 完全退出浏览器（Cmd+Q）
   - 重新打开浏览器
   - 访问应用并尝试录音

#### 对于 macOS Monterey (12.0) 及更早版本：

1. **打开系统偏好设置**
   - 点击屏幕左上角的  菜单
   - 选择"系统偏好设置..."

2. **进入安全性与隐私**
   - 点击"安全性与隐私"图标

3. **选择隐私标签**
   - 点击顶部的"隐私"标签

4. **找到麦克风**
   - 在左侧列表中找到"麦克风"
   - 点击选中

5. **授权浏览器**
   - 点击左下角的锁图标解锁（需要输入密码）
   - 勾选你使用的浏览器
   - 点击锁图标重新锁定

6. **重启浏览器**
   - 完全退出浏览器（Cmd+Q）
   - 重新打开浏览器

### 方法 2: 浏览器级别授权

#### Chrome/Edge:

1. **清除网站权限**
   - 在地址栏左侧点击锁图标或信息图标
   - 找到"麦克风"设置
   - 如果显示"已阻止"，点击并选择"允许"
   - 刷新页面

2. **检查浏览器设置**
   - 打开 `chrome://settings/content/microphone`
   - 确保"网站可以请求使用您的麦克风"已启用
   - 检查"不允许使用麦克风"列表，移除当前网站

#### Safari:

1. **网站设置**
   - Safari → 设置 → 网站
   - 在左侧找到"麦克风"
   - 为当前网站选择"允许"

2. **重置权限**
   - Safari → 设置 → 隐私
   - 点击"管理网站数据..."
   - 找到并删除当前网站的数据
   - 刷新页面重新请求权限

#### Firefox:

1. **网站权限**
   - 点击地址栏左侧的图标
   - 找到"使用麦克风"
   - 选择"允许"

2. **设置页面**
   - 打开 `about:preferences#privacy`
   - 向下滚动到"权限"部分
   - 点击"麦克风"旁的"设置..."
   - 检查当前网站的权限

### 方法 3: 检查麦克风硬件

1. **测试麦克风**
   - 打开"系统设置" → "声音"
   - 选择"输入"标签
   - 对着麦克风说话，观察输入电平是否有变化
   - 如果没有变化，说明麦克风可能有问题

2. **选择正确的输入设备**
   - 在"输入"标签中
   - 确保选择了正确的麦克风设备
   - 如果使用外接麦克风，确保已连接

3. **检查麦克风音量**
   - 确保"输入音量"滑块不是最小
   - 建议设置在 50-80% 之间

## 常见问题

### Q1: 系统设置中没有浏览器选项
**A**: 这是因为浏览器还没有请求过麦克风权限。解决方法：
1. 在浏览器中访问应用
2. 点击录音按钮
3. 应该会弹出系统级权限请求
4. 点击"好"允许访问

### Q2: 授权后仍然无法使用
**A**: 尝试以下步骤：
1. 完全退出浏览器（Cmd+Q）
2. 重新打开浏览器
3. 清除浏览器缓存和 Cookie
4. 重新访问应用

### Q3: 使用外接麦克风
**A**: 
1. 确保麦克风已正确连接
2. 在系统设置 → 声音 → 输入中选择外接麦克风
3. 重启浏览器

### Q4: 虚拟机或远程桌面
**A**: 
- 虚拟机可能无法正确访问主机麦克风
- 需要在虚拟机设置中启用麦克风直通
- 远程桌面可能不支持麦克风

## 验证步骤

### 1. 测试麦克风是否工作
```bash
# 在终端中运行（需要安装 sox）
rec -r 16000 -c 1 test.wav
# 说话几秒后按 Ctrl+C 停止
# 播放录音
play test.wav
```

### 2. 检查浏览器权限
在浏览器控制台运行：
```javascript
navigator.mediaDevices.enumerateDevices()
  .then(devices => {
    const audioInputs = devices.filter(d => d.kind === 'audioinput')
    console.log('可用麦克风:', audioInputs)
  })
```

如果返回空数组或没有权限，说明系统级权限未授予。

### 3. 测试权限请求
在浏览器控制台运行：
```javascript
navigator.mediaDevices.getUserMedia({ audio: true })
  .then(stream => {
    console.log('✅ 麦克风访问成功')
    stream.getTracks().forEach(track => track.stop())
  })
  .catch(err => {
    console.error('❌ 麦克风访问失败:', err.name, err.message)
  })
```

## 不同浏览器的系统权限名称

| 浏览器 | 系统设置中的名称 |
|--------|------------------|
| Chrome | Google Chrome |
| Edge | Microsoft Edge |
| Safari | Safari |
| Firefox | Firefox |
| Brave | Brave Browser |
| Arc | Arc |

## 截图指南

### macOS Ventura 系统设置路径：
```
系统设置
  └─ 隐私与安全性
      └─ 麦克风
          └─ [浏览器名称] ✓
```

### macOS Monterey 系统偏好设置路径：
```
系统偏好设置
  └─ 安全性与隐私
      └─ 隐私
          └─ 麦克风
              └─ [浏览器名称] ✓
```

## 快速诊断命令

在终端运行以下命令检查权限：

```bash
# 检查浏览器是否有麦克风权限
tccutil reset Microphone com.google.Chrome  # Chrome
tccutil reset Microphone com.apple.Safari   # Safari
tccutil reset Microphone org.mozilla.firefox # Firefox

# 注意：这会重置权限，需要重新授权
```

## 应用内提示改进

当检测到 `NotFoundError` 时，应用会显示：

```
未找到麦克风设备

可能的原因：
1. macOS 系统未授权浏览器访问麦克风
2. 没有可用的麦克风设备

解决方法：
1. 打开"系统设置" → "隐私与安全性" → "麦克风"
2. 确保浏览器已被授权
3. 重启浏览器后重试

[查看详细指南] [重试] [关闭]
```

## 开发者注意事项

1. **HTTPS 要求**
   - 麦克风 API 只能在 HTTPS 或 localhost 下使用
   - 确保开发环境使用 localhost

2. **权限请求时机**
   - 必须在用户交互（如点击）后请求权限
   - 页面加载时自动请求可能被阻止

3. **错误处理**
   - `NotFoundError`: 系统级权限未授予或无设备
   - `NotAllowedError`: 用户拒绝或浏览器级权限被阻止
   - `NotReadableError`: 设备被占用

4. **用户体验**
   - 提供清晰的错误提示
   - 提供系统设置指引链接
   - 允许用户重试

## 相关链接

- [Chrome 麦克风权限说明](https://support.google.com/chrome/answer/2693767)
- [Safari 网站权限管理](https://support.apple.com/guide/safari/websites-ibrwe2159f50/mac)
- [macOS 隐私设置](https://support.apple.com/guide/mac-help/control-access-to-the-microphone-mchla1b1e1fe/mac)
