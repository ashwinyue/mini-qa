"""用户认证模块

提供用户登录、注册、token 验证等功能
"""
import hashlib
import secrets
import time
from typing import Optional, Dict, Any
from datetime import datetime, timedelta

# 简单的内存存储（生产环境应使用数据库）
USERS_DB = {
    "admin": {
        "password": hashlib.sha256("admin123".encode()).hexdigest(),
        "username": "admin",
        "nickname": "管理员",
        "role": "admin",
        "email": "admin@example.com",
    },
    "demo": {
        "password": hashlib.sha256("demo123".encode()).hexdigest(),
        "username": "demo",
        "nickname": "演示用户",
        "role": "user",
        "email": "demo@example.com",
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
        "nickname": user["nickname"],
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
        "nickname": user["nickname"],
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
