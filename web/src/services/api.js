/**
 * NATMap Admin API Service
 * 与 Go 后端 API 交互的服务层
 */

import axios from 'axios';

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器 - 添加认证信息
apiClient.interceptors.request.use(
  (config) => {
    // 从 localStorage 获取 token
    const token = localStorage.getItem('admin_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 统一错误处理
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const message = error.response?.data?.error || error.message || '请求失败';
    console.error('API Error:', message);
    return Promise.reject(new Error(message));
  }
);

// ==================== 认证 API ====================

export const authApi = {
  // 登录
  login: (username, password) => 
    apiClient.post('/api/auth/login', { username, password }),
  
  // 获取当前用户信息
  me: () => apiClient.get('/api/auth/me'),
};

// ==================== 租户管理 API ====================

export const tenantApi = {
  // 获取租户列表
  getList: () => apiClient.get('/api/admin?type=tenant'),
  
  // 创建租户
  create: (data) => apiClient.post('/api/admin?type=tenant', data),
  
  // 更新租户
  update: (id, data) => apiClient.put(`/api/admin?type=tenant&id=${id}`, data),
  
  // 删除租户
  delete: (id) => apiClient.delete(`/api/admin?type=tenant&id=${id}`),
};

// ==================== 应用管理 API ====================

export const appApi = {
  // 获取应用列表
  getList: () => apiClient.get('/api/admin?type=app'),
  
  // 创建应用
  create: (data) => apiClient.post('/api/admin?type=app', data),
  
  // 更新应用
  update: (id, data) => apiClient.put(`/api/admin?type=app&id=${id}`, data),
  
  // 删除应用
  delete: (id) => apiClient.delete(`/api/admin?type=app&id=${id}`),
};

// ==================== 映射管理 API ====================

export const mappingApi = {
  // 获取映射列表
  getList: () => apiClient.get('/api/admin?type=mapping'),
  
  // 创建映射
  create: (data) => apiClient.post('/api/admin?type=mapping', data),
  
  // 更新映射
  update: (id, data) => apiClient.put(`/api/admin?type=mapping&id=${id}`, data),
  
  // 删除映射
  delete: (id) => apiClient.delete(`/api/admin?type=mapping&id=${id}`),
};

// ==================== 查询 API ====================

export const queryApi = {
  // 查询映射
  getMapping: (tenantId, appId) => 
    apiClient.get(`/api/get?tenant_id=${tenantId}&app_id=${appId}`),
  
  // 健康检查
  healthCheck: () => apiClient.get('/api/test'),
};

export default apiClient;
