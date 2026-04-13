import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { CircularProgress, Box } from '@mui/material';
import { useAuth } from '../contexts/AuthContext';

// 受保护路由组件
function ProtectedRoute({ children, requireAdmin = false }) {
  const { isAuthenticated, isAdmin, loading } = useAuth();
  const location = useLocation();

  // 加载中显示加载动画
  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  // 未认证，重定向到登录页
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // 需要管理员权限但当前用户不是管理员
  if (requireAdmin && !isAdmin()) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          flexDirection: 'column',
        }}
      >
        <h1>403 - 权限不足</h1>
        <p>您需要管理员权限才能访问此页面</p>
      </Box>
    );
  }

  // 认证通过，渲染子组件
  return children;
}

export default ProtectedRoute;
