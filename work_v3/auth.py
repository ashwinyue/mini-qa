"""用户认证模块

提供用户登录、注册、token 验证等功能
"""
import hashlib
import secrets
import time
from typing import Optional, Dict, Any, List
from datetime import datetime, timedelta

# 用户 ID 计数器
_user_id_counter = 3

# 简单的内存存储（生产环境应使用数据库）
USERS_DB = {
    "admin": {
        "id": 1,
        "password": hashlib.sha256("admin123".encode()).hexdigest(),
        "username": "admin",
        "realname": "管理员",
        "role": "admin",
        "email": "admin@example.com",
        "status": 1,
        "created_at": "2024-01-01 00:00:00",
    },
    "demo": {
        "id": 2,
        "password": hashlib.sha256("demo123".encode()).hexdigest(),
        "username": "demo",
        "realname": "演示用户",
        "role": "user",
        "email": "demo@example.com",
        "status": 1,
        "created_at": "2024-01-01 00:00:00",
    },
}

# 角色 ID 计数器
_role_id_counter = 3

# 角色存储
ROLES_DB = {
    1: {
        "id": 1,
        "name": "管理员",
        "code": "admin",
        "description": "系统管理员，拥有所有权限",
        "created_at": "2024-01-01 00:00:00",
    },
    2: {
        "id": 2,
        "name": "普通用户",
        "code": "user",
        "description": "普通用户，基本权限",
        "created_at": "2024-01-01 00:00:00",
    },
}

# Token 存储（生产环境应使用 Redis）
TOKENS = {}

# Token 过期时间（7天）
TOKEN_EXPIRE_DAYS = 7


def hash_password(password: str) -> str:
    """密码哈希"""
    return hashlib.sha256(password.encode()).hexdigest()


def generate_token() -> str:
    """生成随机 token"""
    return secrets.token_urlsafe(32)


def create_user(username: str, password: str, nickname: str, email: str) -> Dict[str, Any]:
    """创建新用户"""
    if username in USERS_DB:
        return {"ok": False, "error": "用户名已存在"}
    
    USERS_DB[username] = {
        "password": hash_password(password),
        "username": username,
        "nickname": nickname,
        "role": "user",
        "email": email,
    }
    
    return {"ok": True, "message": "注册成功"}


def verify_user(username: str, password: str) -> Optional[Dict[str, Any]]:
    """验证用户名和密码"""
    user = USERS_DB.get(username)
    if not user:
        return None
    
    if user["password"] != hash_password(password):
        return None
    
    return {
        "username": user["username"],
        "realname": user.get("realname", user.get("nickname", "")),
        "role": user["role"],
        "email": user.get("email", ""),
    }


def create_token(username: str) -> str:
    """创建 token"""
    token = generate_token()
    expire_at = datetime.now() + timedelta(days=TOKEN_EXPIRE_DAYS)
    
    TOKENS[token] = {
        "username": username,
        "expire_at": expire_at,
        "created_at": datetime.now(),
    }
    
    return token


def verify_token(token: str) -> Optional[Dict[str, Any]]:
    """验证 token"""
    token_data = TOKENS.get(token)
    if not token_data:
        return None
    
    # 检查是否过期
    if datetime.now() > token_data["expire_at"]:
        del TOKENS[token]
        return None
    
    username = token_data["username"]
    user = USERS_DB.get(username)
    if not user:
        return None
    
    return {
        "username": user["username"],
        "realname": user.get("realname", user.get("nickname", "")),
        "role": user["role"],
        "email": user.get("email", ""),
    }


def revoke_token(token: str) -> bool:
    """撤销 token（登出）"""
    if token in TOKENS:
        del TOKENS[token]
        return True
    return False


def get_user_info(username: str) -> Optional[Dict[str, Any]]:
    """获取用户信息"""
    user = USERS_DB.get(username)
    if not user:
        return None
    
    return {
        "username": user["username"],
        "nickname": user["nickname"],
        "role": user["role"],
        "email": user.get("email", ""),
    }


def cleanup_expired_tokens():
    """清理过期的 token"""
    now = datetime.now()
    expired = [token for token, data in TOKENS.items() if now > data["expire_at"]]
    for token in expired:
        del TOKENS[token]
    return len(expired)


# ============ 用户管理 CRUD ============

def get_user_list(page: int = 1, page_size: int = 10, keyword: str = "") -> Dict[str, Any]:
    """获取用户列表（分页）"""
    users = list(USERS_DB.values())
    
    # 关键词过滤
    if keyword:
        users = [u for u in users if keyword.lower() in u["username"].lower() 
                 or keyword.lower() in u.get("realname", "").lower()]
    
    total = len(users)
    start = (page - 1) * page_size
    end = start + page_size
    
    # 返回不含密码的用户列表
    result = []
    for u in users[start:end]:
        result.append({
            "id": u.get("id"),
            "username": u["username"],
            "realname": u.get("realname", u.get("nickname", "")),
            "email": u.get("email", ""),
            "status": u.get("status", 1),
            "created_at": u.get("created_at", ""),
        })
    
    return {"list": result, "total": total, "page": page, "pageSize": page_size}


def get_user_by_id(user_id: int) -> Optional[Dict[str, Any]]:
    """根据 ID 获取用户"""
    for u in USERS_DB.values():
        if u.get("id") == user_id:
            return {
                "id": u.get("id"),
                "username": u["username"],
                "realname": u.get("realname", u.get("nickname", "")),
                "email": u.get("email", ""),
                "status": u.get("status", 1),
                "created_at": u.get("created_at", ""),
            }
    return None


def create_user_full(username: str, password: str, realname: str, email: str) -> Dict[str, Any]:
    """创建新用户（完整版）"""
    global _user_id_counter
    
    if username in USERS_DB:
        return {"ok": False, "error": "用户名已存在"}
    
    user_id = _user_id_counter
    _user_id_counter += 1
    
    USERS_DB[username] = {
        "id": user_id,
        "password": hash_password(password),
        "username": username,
        "realname": realname,
        "role": "user",
        "email": email,
        "status": 1,
        "created_at": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    }
    
    return {"ok": True, "id": user_id}


def update_user(user_id: int, data: Dict[str, Any]) -> Dict[str, Any]:
    """更新用户信息"""
    for username, u in USERS_DB.items():
        if u.get("id") == user_id:
            if "realname" in data:
                u["realname"] = data["realname"]
            if "email" in data:
                u["email"] = data["email"]
            if "status" in data:
                u["status"] = data["status"]
            if "password" in data and data["password"]:
                u["password"] = hash_password(data["password"])
            return {"ok": True}
    return {"ok": False, "error": "用户不存在"}


def delete_user(user_id: int) -> Dict[str, Any]:
    """删除用户"""
    if user_id == 1:
        return {"ok": False, "error": "不能删除系统管理员"}
    
    for username, u in list(USERS_DB.items()):
        if u.get("id") == user_id:
            del USERS_DB[username]
            return {"ok": True}
    return {"ok": False, "error": "用户不存在"}


# ============ 角色管理 CRUD ============

def get_role_list(page: int = 1, page_size: int = 10, keyword: str = "") -> Dict[str, Any]:
    """获取角色列表（分页）"""
    roles = list(ROLES_DB.values())
    
    # 关键词过滤
    if keyword:
        roles = [r for r in roles if keyword.lower() in r["name"].lower() 
                 or keyword.lower() in r.get("code", "").lower()]
    
    total = len(roles)
    start = (page - 1) * page_size
    end = start + page_size
    
    return {"list": roles[start:end], "total": total, "page": page, "pageSize": page_size}


def get_role_by_id(role_id: int) -> Optional[Dict[str, Any]]:
    """根据 ID 获取角色"""
    return ROLES_DB.get(role_id)


def create_role(name: str, code: str, description: str = "") -> Dict[str, Any]:
    """创建新角色"""
    global _role_id_counter
    
    # 检查 code 是否重复
    for r in ROLES_DB.values():
        if r["code"] == code:
            return {"ok": False, "error": "角色标识已存在"}
    
    role_id = _role_id_counter
    _role_id_counter += 1
    
    ROLES_DB[role_id] = {
        "id": role_id,
        "name": name,
        "code": code,
        "description": description,
        "created_at": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    }
    
    return {"ok": True, "id": role_id}


def update_role(role_id: int, data: Dict[str, Any]) -> Dict[str, Any]:
    """更新角色信息"""
    if role_id not in ROLES_DB:
        return {"ok": False, "error": "角色不存在"}
    
    role = ROLES_DB[role_id]
    if "name" in data:
        role["name"] = data["name"]
    if "code" in data:
        # 检查 code 是否与其他角色重复
        for r in ROLES_DB.values():
            if r["id"] != role_id and r["code"] == data["code"]:
                return {"ok": False, "error": "角色标识已存在"}
        role["code"] = data["code"]
    if "description" in data:
        role["description"] = data["description"]
    
    return {"ok": True}


def delete_role(role_id: int) -> Dict[str, Any]:
    """删除角色"""
    if role_id in (1, 2):
        return {"ok": False, "error": "不能删除系统内置角色"}
    
    if role_id not in ROLES_DB:
        return {"ok": False, "error": "角色不存在"}
    
    del ROLES_DB[role_id]
    return {"ok": True}
