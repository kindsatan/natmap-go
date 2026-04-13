import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Snackbar,
  Tooltip,
  Chip,
  Grid,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Route as RouteIcon,
} from '@mui/icons-material';
import { mappingApi, tenantApi, appApi } from '../services/api';

function MappingManagement() {
  const [mappings, setMappings] = useState([]);
  const [tenants, setTenants] = useState([]);
  const [apps, setApps] = useState([]);
  const [loading, setLoading] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [editingMapping, setEditingMapping] = useState(null);
  const [formData, setFormData] = useState({
    tenant_id: '',
    app_id: '',
    public_ip: '',
    public_port: '',
    local_ip: '',
    local_port: '',
    protocol: 'tcp',
  });
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    try {
      const [mappingsData, tenantsData, appsData] = await Promise.all([
        mappingApi.getList(),
        tenantApi.getList(),
        appApi.getList(),
      ]);
      setMappings(Array.isArray(mappingsData) ? mappingsData : []);
      setTenants(Array.isArray(tenantsData) ? tenantsData : []);
      setApps(Array.isArray(appsData) ? appsData : []);
    } catch (error) {
      showSnackbar(error.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  // 显示提示消息
  const showSnackbar = (message, severity = 'success') => {
    setSnackbar({ open: true, message, severity });
  };

  // 关闭提示
  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  // 打开创建对话框
  const handleOpenCreate = () => {
    setEditingMapping(null);
    setFormData({
      tenant_id: tenants[0]?.id || '',
      app_id: apps[0]?.id || '',
      public_ip: '',
      public_port: '',
      local_ip: '',
      local_port: '',
      protocol: 'tcp',
    });
    setOpenDialog(true);
  };

  // 打开编辑对话框
  const handleOpenEdit = (mapping) => {
    setEditingMapping(mapping);
    setFormData({
      tenant_id: mapping.tenant_id,
      app_id: mapping.app_id,
      public_ip: mapping.public_ip,
      public_port: mapping.public_port,
      local_ip: mapping.local_ip || '',
      local_port: mapping.local_port || '',
      protocol: mapping.protocol || 'tcp',
    });
    setOpenDialog(true);
  };

  // 关闭对话框
  const handleCloseDialog = () => {
    setOpenDialog(false);
    setEditingMapping(null);
  };

  // 保存映射
  const handleSave = async () => {
    if (!formData.tenant_id || !formData.app_id) {
      showSnackbar('请选择租户和应用', 'warning');
      return;
    }
    if (!formData.public_ip.trim()) {
      showSnackbar('请输入公网 IP', 'warning');
      return;
    }
    if (!formData.public_port || formData.public_port < 1 || formData.public_port > 65535) {
      showSnackbar('请输入有效的端口号 (1-65535)', 'warning');
      return;
    }

    try {
      const payload = {
        ...formData,
        tenant_id: parseInt(formData.tenant_id),
        app_id: parseInt(formData.app_id),
        public_port: parseInt(formData.public_port),
        local_port: formData.local_port ? parseInt(formData.local_port) : 0,
      };

      if (editingMapping) {
        await mappingApi.update(editingMapping.id, payload);
        showSnackbar('映射更新成功');
      } else {
        await mappingApi.create(payload);
        showSnackbar('映射创建成功');
      }
      handleCloseDialog();
      loadData();
    } catch (error) {
      showSnackbar(error.message, 'error');
    }
  };

  // 删除映射
  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除这条映射记录吗？')) {
      return;
    }

    try {
      await mappingApi.delete(id);
      showSnackbar('映射删除成功');
      loadData();
    } catch (error) {
      showSnackbar(error.message, 'error');
    }
  };

  // 获取应用名称
  const getAppName = (appId) => {
    const app = apps.find((a) => a.id === appId);
    return app ? app.app_name : '-';
  };

  // 获取租户名称
  const getTenantName = (tenantId) => {
    const tenant = tenants.find((t) => t.id === tenantId);
    return tenant ? tenant.tenant_name : '-';
  };

  // 过滤出选中租户的应用
  const getFilteredApps = () => {
    if (!formData.tenant_id) return apps;
    return apps.filter((app) => app.tenant_id === parseInt(formData.tenant_id));
  };

  return (
    <Box>
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h5" component="h1" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <RouteIcon color="primary" />
              映射管理
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleOpenCreate}
              disabled={loading || tenants.length === 0 || apps.length === 0}
            >
              添加映射
            </Button>
          </Box>

          {(tenants.length === 0 || apps.length === 0) && (
            <Alert severity="warning" sx={{ mb: 3 }}>
              请先创建租户和应用，然后再添加映射。
            </Alert>
          )}

          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
                  <TableCell>ID</TableCell>
                  <TableCell>租户</TableCell>
                  <TableCell>应用</TableCell>
                  <TableCell>公网地址</TableCell>
                  <TableCell>本地地址</TableCell>
                  <TableCell>协议</TableCell>
                  <TableCell>更新时间</TableCell>
                  <TableCell align="right">操作</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {mappings.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} align="center" sx={{ py: 4, color: 'text.secondary' }}>
                      暂无映射数据
                    </TableCell>
                  </TableRow>
                ) : (
                  mappings.map((mapping) => (
                    <TableRow key={mapping.id} hover>
                      <TableCell>
                        <Chip label={mapping.id} size="small" variant="outlined" />
                      </TableCell>
                      <TableCell>{getTenantName(mapping.tenant_id)}</TableCell>
                      <TableCell>
                        <Chip
                          label={getAppName(mapping.app_id)}
                          size="small"
                          color="primary"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={`${mapping.public_ip}:${mapping.public_port}`}
                          size="small"
                          color="success"
                          variant="outlined"
                          sx={{ fontFamily: 'monospace' }}
                        />
                      </TableCell>
                      <TableCell sx={{ fontFamily: 'monospace', color: 'text.secondary' }}>
                        {mapping.local_ip ? `${mapping.local_ip}:${mapping.local_port}` : '-'}
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={mapping.protocol?.toUpperCase() || 'TCP'}
                          size="small"
                          color="info"
                        />
                      </TableCell>
                      <TableCell sx={{ color: 'text.secondary' }}>
                        {mapping.updated_at}
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="编辑">
                          <IconButton
                            size="small"
                            onClick={() => handleOpenEdit(mapping)}
                            color="primary"
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="删除">
                          <IconButton
                            size="small"
                            onClick={() => handleDelete(mapping.id)}
                            color="error"
                          >
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      {/* 创建/编辑对话框 */}
      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="md" fullWidth>
        <DialogTitle>{editingMapping ? '编辑映射' : '添加映射'}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 0.5 }}>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>所属租户</InputLabel>
                <Select
                  value={formData.tenant_id}
                  onChange={(e) =>
                    setFormData({ ...formData, tenant_id: e.target.value, app_id: '' })
                  }
                  label="所属租户"
                >
                  {tenants.map((tenant) => (
                    <MenuItem key={tenant.id} value={tenant.id}>
                      {tenant.tenant_name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>所属应用</InputLabel>
                <Select
                  value={formData.app_id}
                  onChange={(e) => setFormData({ ...formData, app_id: e.target.value })}
                  label="所属应用"
                  disabled={!formData.tenant_id}
                >
                  {getFilteredApps().map((app) => (
                    <MenuItem key={app.id} value={app.id}>
                      {app.app_name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="公网 IP"
                value={formData.public_ip}
                onChange={(e) => setFormData({ ...formData, public_ip: e.target.value })}
                placeholder="例如: 192.168.1.100"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="公网端口"
                type="number"
                value={formData.public_port}
                onChange={(e) => setFormData({ ...formData, public_port: e.target.value })}
                placeholder="1-65535"
                inputProps={{ min: 1, max: 65535 }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="本地 IP (可选)"
                value={formData.local_ip}
                onChange={(e) => setFormData({ ...formData, local_ip: e.target.value })}
                placeholder="例如: 10.0.0.1"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="本地端口 (可选)"
                type="number"
                value={formData.local_port}
                onChange={(e) => setFormData({ ...formData, local_port: e.target.value })}
                placeholder="1-65535"
                inputProps={{ min: 1, max: 65535 }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>协议</InputLabel>
                <Select
                  value={formData.protocol}
                  onChange={(e) => setFormData({ ...formData, protocol: e.target.value })}
                  label="协议"
                >
                  <MenuItem value="tcp">TCP</MenuItem>
                  <MenuItem value="udp">UDP</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={handleCloseDialog}>取消</Button>
          <Button onClick={handleSave} variant="contained">
            {editingMapping ? '更新' : '创建'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* 提示消息 */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={3000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <Alert severity={snackbar.severity} onClose={handleCloseSnackbar}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}

export default MappingManagement;
