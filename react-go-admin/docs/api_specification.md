# API 规范说明

## 总览
- 基础 URL：`http://localhost:8000/api`
- 认证方式：JWT Token (Bearer Token)
- 内容类型：`application/json`
- 状态码约定：
  - 200 成功
  - 401 未授权
  - 404 资源不存在
  - 400 参数错误
  - 500 服务内部错误
- 统一响应格式：`{code, message, data}`

## 1. 认证接口

### 用户登录
- `POST /auth/login`
- 功能：用户身份验证，返回JWT令牌
- 请求体：
  ```json
  {
    "username": "admin",
    "password": "123456"
  }
  ```
- 响应体：
  ```json
  {
    "code": 200,
    "message": "登录成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin"
      }
    }
  }
  ```

### 验证令牌
- `GET /auth/verify`
- 功能：验证JWT令牌有效性
- 请求头：`Authorization: Bearer <token>`
- 响应体：
  ```json
  {
    "code": 200,
    "message": "令牌有效",
    "data": {
      "user": {
        "id": 1,
        "username": "admin",
        "role": "admin"
      }
    }
  }
  ```

## 2. 用户管理接口

### 获取用户列表
- `GET /users`
- 功能：分页获取用户列表
- 查询参数：
  - `page`: 页码 (默认: 1)
  - `pageSize`: 每页条数 (默认: 10)
  - `keyword`: 搜索关键词 (可选)
  - `status`: 状态筛选 (可选: active, inactive)
- 响应体：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "list": [
        {
          "id": 1,
          "username": "admin",
          "email": "admin@example.com",
          "role": "admin",
          "status": "active",
          "created_at": "2024-01-01T00:00:00Z",
          "updated_at": "2024-01-01T00:00:00Z"
        }
      ],
      "total": 100,
      "page": 1,
      "pageSize": 10
    }
  }
  ```

### 创建用户
- `POST /users`
- 功能：创建新用户
- 请求体：
  ```json
  {
    "username": "newuser",
    "email": "user@example.com",
    "password": "123456",
    "role": "user",
    "status": "active"
  }
  ```
- 响应体：
  ```json
  {
    "code": 200,
    "message": "用户创建成功",
    "data": {
      "id": 2,
      "username": "newuser",
      "email": "user@example.com",
      "role": "user",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
  ```

### 更新用户
- `PUT /users/:id`
- 功能：更新用户信息
- 请求体：
  ```json
  {
    "username": "updateduser",
    "email": "updated@example.com",
    "role": "user",
    "status": "active"
  }
  ```
- 响应体：
  ```json
  {
    "code": 200,
    "message": "用户更新成功",
    "data": {
      "id": 2,
      "username": "updateduser",
      "email": "updated@example.com",
      "role": "user",
      "status": "active",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 删除用户
- `DELETE /users/:id`
- 功能：删除用户
- 响应体：
  ```json
  {
    "code": 200,
    "message": "用户删除成功",
    "data": null
  }
  ```

## 3. 角色管理接口

### 获取角色列表
- `GET /roles`
- 功能：获取所有角色列表
- 响应体：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": [
      {
        "id": 1,
        "name": "admin",
        "description": "管理员",
        "permissions": ["user.manage", "role.manage"],
        "created_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "name": "user",
        "description": "普通用户",
        "permissions": ["profile.view"],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
  ```

### 创建角色
- `POST /roles`
- 功能：创建新角色
- 请求体：
  ```json
  {
    "name": "editor",
    "description": "编辑",
    "permissions": ["content.manage"]
  }
  ```
- 响应体：
  ```json
  {
    "code": 200,
    "message": "角色创建成功",
    "data": {
      "id": 3,
      "name": "editor",
      "description": "编辑",
      "permissions": ["content.manage"],
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
  ```

### 更新角色
- `PUT /roles/:id`
- 功能：更新角色信息
- 请求体：
  ```json
  {
    "name": "senior_editor",
    "description": "高级编辑",
    "permissions": ["content.manage", "user.view"]
  }
  ```
- 响应体：
  ```json
  {
    "code": 200,
    "message": "角色更新成功",
    "data": {
      "id": 3,
      "name": "senior_editor",
      "description": "高级编辑",
      "permissions": ["content.manage", "user.view"],
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 删除角色
- `DELETE /roles/:id`
- 功能：删除角色
- 响应体：
  ```json
  {
    "code": 200,
    "message": "角色删除成功",
    "data": null
  }
  ```

## 4. 仪表盘接口

### 获取统计数据
- `GET /dashboard/stats`
- 功能：获取系统统计数据
- 响应体：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "totalUsers": 150,
      "activeUsers": 120,
      "totalRoles": 5,
      "todayNewUsers": 10
    }
  }
  ```

## 认证与授权
- 认证方式：JWT Bearer Token
- 获取Token：通过登录接口获取
- 使用方式：在请求头中添加 `Authorization: Bearer <token>`
- Token有效期：默认24小时

## 错误处理
- 统一错误响应格式：
  ```json
  {
    "code": 400,
    "message": "错误描述信息",
    "data": null
  }
  ```
- 常见错误码：
  - 400: 参数错误
  - 401: 未授权或Token无效
  - 403: 权限不足
  - 404: 资源不存在
  - 500: 服务器内部错误

## 代码示例

### JavaScript (React)
```javascript
// 用户登录
const login = async (username, password) => {
  const response = await fetch('http://localhost:8000/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });
  return await response.json();
};

// 获取用户列表
const getUsers = async (token) => {
  const response = await fetch('http://localhost:8000/api/users?page=1&pageSize=10', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  return await response.json();
};
```


## 关联文档
- 架构设计：参见《[architecture_design.md](./architecture_design.md)》
- 部署说明：参见《[deployment_guide.md](./deployment_guide.md)》