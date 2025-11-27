import api from './api';

export interface User {
  username: string;
  nickname: string;
  role: string;
  email: string;
}

export interface LoginResponse {
  code: number;
  message: string;
  data: {
    token: string;
    user: User;
    expiresIn: number;
  };
}

export interface RegisterResponse {
  code: number;
  message: string;
  data: {
    message: string;
  };
}

class AuthService {
  private readonly TOKEN_KEY = 'auth_token';
  private readonly USER_KEY = 'user_info';

  /**
   * 用户登录
   */
  async login(username: string, password: string): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/api/auth/login', {
      username,
      password,
    });

    if (response.data.code === 0) {
      const { token, user } = response.data.data;
      this.setToken(token);
      this.setUser(user);
    }

    return response.data;
  }

  /**
   * 用户注册
   */
  async register(
    username: string,
    password: string,
    nickname: string,
    email?: string
  ): Promise<RegisterResponse> {
    const response = await api.post<RegisterResponse>('/api/auth/register', {
      username,
      password,
      nickname,
      email,
    });

    return response.data;
  }

  /**
   * 用户登出
   */
  async logout(): Promise<void> {
    try {
      await api.post('/api/auth/logout');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      this.clearAuth();
    }
  }

  /**
   * 获取当前用户信息
   */
  async getCurrentUser(): Promise<User | null> {
    try {
      const response = await api.get<{ code: number; data: User }>('/api/auth/me');
      if (response.data.code === 0) {
        this.setUser(response.data.data);
        return response.data.data;
      }
    } catch (error) {
      console.error('Get current user error:', error);
      this.clearAuth();
    }
    return null;
  }

  /**
   * 保存 token
   */
  setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
  }

  /**
   * 获取 token
   */
  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  /**
   * 保存用户信息
   */
  setUser(user: User): void {
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
  }

  /**
   * 获取用户信息
   */
  getUser(): User | null {
    const userStr = localStorage.getItem(this.USER_KEY);
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch {
        return null;
      }
    }
    return null;
  }

  /**
   * 清除认证信息
   */
  clearAuth(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
  }

  /**
   * 检查是否已登录
   */
  isAuthenticated(): boolean {
    return !!this.getToken();
  }
}

export const authService = new AuthService();
