import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { authApi } from '../services/api';

// 创建认证上下文
const AuthContext = createContext(null);

// 认证提供者组件
export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);

  // 登出函数 - 使用 useCallback 避免循环依赖
  const logout = useCallback(() => {
    localStorage.removeItem('admin_token');
    localStorage.removeItem('admin_user');
    setUser(null);
    setIsAuthenticated(false);
  }, []);

  // 初始化时检查本地存储的 token
  useEffect(() => {
    // 使用 requestAnimationFrame 避免同步 setState
    const initAuth = () => {
      const token = localStorage.getItem('admin_token');
      const userData = localStorage.getItem('admin_user');
      
      if (token && userData) {
        try {
          setUser(JSON.parse(userData));
          setIsAuthenticated(true);
        } catch {
          // 解析失败，清除存储
          logout();
        }
      }
      setLoading(false);
    };
    
    requestAnimationFrame(initAuth);
  }, [logout]);

  // 登录函数
  const login = async (username, password) => {
    try {
      const response = await authApi.login(username, password);
      
      if (response.token) {
        localStorage.setItem('admin_token', response.token);
        localStorage.setItem('admin_user', JSON.stringify(response.user));
        setUser(response.user);
        setIsAuthenticated(true);
        return { success: true };
      }
      
      return { success: false, error: response.error || '登录失败' };
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '网络错误，请稍后重试' 
      };
    }
  };

  // 检查是否有管理员权限
  const isAdmin = () => {
    return user?.role === 'admin';
  };

  const value = {
    user,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    loading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

// 自定义 Hook 使用认证上下文
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
