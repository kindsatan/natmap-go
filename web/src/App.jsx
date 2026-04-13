import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import TenantManagement from './pages/TenantManagement';
import AppManagement from './pages/AppManagement';
import MappingManagement from './pages/MappingManagement';

// 创建主题
const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#f5f5f5',
    },
  },
  typography: {
    fontFamily: [
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
    ].join(','),
  },
});

// 受保护的路由包装器
function ProtectedLayout({ children }) {
  return (
    <ProtectedRoute>
      <Layout>{children}</Layout>
    </ProtectedRoute>
  );
}

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <Routes>
            {/* 登录页面 - 公开访问 */}
            <Route path="/login" element={<Login />} />
            
            {/* 受保护的管理后台路由 */}
            <Route path="/" element={<ProtectedLayout><Dashboard /></ProtectedLayout>} />
            <Route path="/tenants" element={<ProtectedLayout><TenantManagement /></ProtectedLayout>} />
            <Route path="/apps" element={<ProtectedLayout><AppManagement /></ProtectedLayout>} />
            <Route path="/mappings" element={<ProtectedLayout><MappingManagement /></ProtectedLayout>} />
            
            {/* 默认重定向到登录页 */}
            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
