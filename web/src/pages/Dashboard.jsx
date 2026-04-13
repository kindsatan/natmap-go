import React, { useState, useEffect } from 'react';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  Chip,
  LinearProgress,
} from '@mui/material';
import {
  Business as BusinessIcon,
  Apps as AppsIcon,
  Route as RouteIcon,
  CheckCircle as CheckCircleIcon,
} from '@mui/icons-material';
import { tenantApi, appApi, mappingApi, queryApi } from '../services/api';

function StatCard({ title, count, icon, color, subtitle }) {
  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <Box
            sx={{
              backgroundColor: `${color}.main`,
              color: 'white',
              borderRadius: 2,
              p: 1,
              mr: 2,
              display: 'flex',
            }}
          >
            {icon}
          </Box>
          <Box>
            <Typography color="textSecondary" variant="body2">
              {title}
            </Typography>
            <Typography variant="h4" component="div" sx={{ fontWeight: 'bold' }}>
              {count}
            </Typography>
          </Box>
        </Box>
        {subtitle && (
          <Typography variant="body2" color="textSecondary">
            {subtitle}
          </Typography>
        )}
      </CardContent>
    </Card>
  );
}

function Dashboard() {
  const [stats, setStats] = useState({
    tenants: 0,
    apps: 0,
    mappings: 0,
  });
  const [loading, setLoading] = useState(true);
  const [serverStatus, setServerStatus] = useState({ online: false, message: '' });

  useEffect(() => {
    loadStats();
    checkServerStatus();
  }, []);

  const loadStats = async () => {
    setLoading(true);
    try {
      const [tenants, apps, mappings] = await Promise.all([
        tenantApi.getList(),
        appApi.getList(),
        mappingApi.getList(),
      ]);
      setStats({
        tenants: Array.isArray(tenants) ? tenants.length : 0,
        apps: Array.isArray(apps) ? apps.length : 0,
        mappings: Array.isArray(mappings) ? mappings.length : 0,
      });
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const checkServerStatus = async () => {
    try {
      const result = await queryApi.healthCheck();
      setServerStatus({
        online: true,
        message: result.message || '服务正常运行',
      });
    } catch {
      setServerStatus({
        online: false,
        message: '无法连接到服务器',
      });
    }
  };

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 'bold' }}>
        概览
      </Typography>

      {/* 服务器状态 */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <CheckCircleIcon
              color={serverStatus.online ? 'success' : 'error'}
              sx={{ fontSize: 40 }}
            />
            <Box>
              <Typography variant="h6">
                服务器状态
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
                <Chip
                  label={serverStatus.online ? '在线' : '离线'}
                  color={serverStatus.online ? 'success' : 'error'}
                  size="small"
                />
                <Typography variant="body2" color="textSecondary">
                  {serverStatus.message}
                </Typography>
              </Box>
            </Box>
          </Box>
        </CardContent>
      </Card>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* 统计卡片 */}
      <Grid container spacing={3}>
        <Grid item xs={12} sm={6} md={4}>
          <StatCard
            title="租户数量"
            count={stats.tenants}
            icon={<BusinessIcon />}
            color="primary"
            subtitle="已注册的公司/组织"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <StatCard
            title="应用数量"
            count={stats.apps}
            icon={<AppsIcon />}
            color="secondary"
            subtitle="已配置的应用服务"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <StatCard
            title="映射数量"
            count={stats.mappings}
            icon={<RouteIcon />}
            color="success"
            subtitle="活跃的 NAT 映射"
          />
        </Grid>
      </Grid>

      {/* 快速操作提示 */}
      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            快速开始
          </Typography>
          <Typography variant="body1" color="textSecondary" paragraph>
            欢迎使用 NATMap 管理后台！您可以：
          </Typography>
          <Box component="ul" sx={{ pl: 2 }}>
            <Typography component="li" variant="body1" sx={{ mb: 1 }}>
              在<strong>租户管理</strong>中添加新的公司或组织
            </Typography>
            <Typography component="li" variant="body1" sx={{ mb: 1 }}>
              在<strong>应用管理</strong>中配置需要映射的应用服务
            </Typography>
            <Typography component="li" variant="body1" sx={{ mb: 1 }}>
              在<strong>映射管理</strong>中查看和管理 NAT 映射记录
            </Typography>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}

export default Dashboard;
