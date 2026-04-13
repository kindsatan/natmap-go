import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Alert,
  CircularProgress,
  Container,
} from '@mui/material';
import { LockOutlined as LockIcon } from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    if (!username || !password) {
      setError('请输入用户名和密码');
      setLoading(false);
      return;
    }

    const result = await login(username, password);
    
    if (!result.success) {
      setError(result.error);
    }
    
    setLoading(false);
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <Container maxWidth="sm">
        <Card elevation={8} sx={{ borderRadius: 3 }}>
          <CardContent sx={{ p: 4 }}>
            <Box sx={{ textAlign: 'center', mb: 3 }}>
              <Box
                sx={{
                  backgroundColor: 'primary.main',
                  color: 'white',
                  borderRadius: '50%',
                  width: 64,
                  height: 64,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  margin: '0 auto 16px',
                }}
              >
                <LockIcon fontSize="large" />
              </Box>
              <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
                NATMap 管理后台
              </Typography>
              <Typography variant="body2" color="text.secondary">
                请输入管理员账号登录系统
              </Typography>
            </Box>

            {error && (
              <Alert severity="error" sx={{ mb: 3 }}>
                {error}
              </Alert>
            )}

            <Box component="form" onSubmit={handleSubmit}>
              <TextField
                fullWidth
                label="用户名"
                variant="outlined"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                margin="normal"
                autoComplete="username"
                autoFocus
              />
              <TextField
                fullWidth
                label="密码"
                type="password"
                variant="outlined"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                margin="normal"
                autoComplete="current-password"
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={loading}
                sx={{ mt: 3, mb: 2, py: 1.5 }}
              >
                {loading ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  '登录'
                )}
              </Button>
            </Box>

            <Typography variant="body2" color="text.secondary" align="center" sx={{ mt: 2 }}>
              默认账号: admin / admin123
            </Typography>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
}

export default Login;
