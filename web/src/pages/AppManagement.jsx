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
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Apps as AppsIcon,
} from '@mui/icons-material';
import { appApi, tenantApi } from '../services/api';

function AppManagement() {
  const [apps, setApps] = useState([]);
  const [tenants, setTenants] = useState([]);
  const [loading, setLoading] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [editingApp, setEditingApp] = useState(null);
  const [formData, setFormData] = useState({
    tenant_id: '',
    app_name: '',
    description: '',
  });
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    try {
      const [appsData, tenantsData] = await Promise.all([
        appApi.getList(),
        tenantApi.getList(),
      ]);
      setApps(Array.isArray(appsData) ? appsData : []);
      setTenants(Array.isArray(tenantsData) ? tenantsData : []);
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
    setEditingApp(null);
    setFormData({
      tenant_id: tenants[0]?.id || '',
      app_name: '',
      description: '',
    });
    setOpenDialog(true);
  };

  // 打开编辑对话框
  const handleOpenEdit = (app) => {
    setEditingApp(app);
    setFormData({
      tenant_id: app.tenant_id,
      app_name: app.app_name,
      description: app.description || '',
    });
    setOpenDialog(true);
  };

  // 关闭对话框
  const handleCloseDialog = () => {
    setOpenDialog(false);
    setEditingApp(null);
  };

  // 保存应用
  const handleSave = async () => {
    if (!formData.tenant_id) {
      showSnackbar('请选择所属租户', 'warning');
      return;
    }
    if (!formData.app_name.trim()) {
      showSnackbar('请输入应用名称', 'warning');
      return;
    }

    try {
      if (editingApp) {
        await appApi.update(editingApp.id, formData);
        showSnackbar('应用更新成功');
      } else {
        await appApi.create(formData);
        showSnackbar('应用创建成功');
      }
      handleCloseDialog();
      loadData();
    } catch (error) {
      showSnackbar(error.message, 'error');
    }
  };

  // 删除应用
  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除这个应用吗？相关的映射也将被删除。')) {
      return;
    }

    try {
      await appApi.delete(id);
      showSnackbar('应用删除成功');
      loadData();
    } catch (error) {
      showSnackbar(error.message, 'error');
    }
  };

  // 获取租户名称
  const getTenantName = (tenantId) => {
    const tenant = tenants.find((t) => t.id === tenantId);
    return tenant ? tenant.tenant_name : '-';
  };

  return (
    <Box>
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h5" component="h1" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <AppsIcon color="primary" />
              应用管理
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleOpenCreate}
              disabled={loading || tenants.length === 0}
            >
              添加应用
            </Button>
          </Box>

          {tenants.length === 0 && (
            <Alert severity="warning" sx={{ mb: 3 }}>
              请先创建租户，然后再添加应用。
            </Alert>
          )}

          <TableContainer component={Paper} variant="outlined">
            <Table>
              <TableHead>
                <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
                  <TableCell>ID</TableCell>
                  <TableCell>应用名称</TableCell>
                  <TableCell>所属租户</TableCell>
                  <TableCell>描述</TableCell>
                  <TableCell>创建时间</TableCell>
                  <TableCell align="right">操作</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {apps.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} align="center" sx={{ py: 4, color: 'text.secondary' }}>
                      暂无应用数据
                    </TableCell>
                  </TableRow>
                ) : (
                  apps.map((app) => (
                    <TableRow key={app.id} hover>
                      <TableCell>
                        <Chip label={app.id} size="small" variant="outlined" />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={app.app_name}
                          size="small"
                          color="primary"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>{getTenantName(app.tenant_id)}</TableCell>
                      <TableCell sx={{ color: 'text.secondary', maxWidth: 200 }}>
                        {app.description || '-'}
                      </TableCell>
                      <TableCell sx={{ color: 'text.secondary' }}>
                        {app.created_at}
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="编辑">
                          <IconButton
                            size="small"
                            onClick={() => handleOpenEdit(app)}
                            color="primary"
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="删除">
                          <IconButton
                            size="small"
                            onClick={() => handleDelete(app.id)}
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
      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
        <DialogTitle>{editingApp ? '编辑应用' : '添加应用'}</DialogTitle>
        <DialogContent>
          <FormControl fullWidth margin="dense" sx={{ mt: 1, mb: 2 }}>
            <InputLabel>所属租户</InputLabel>
            <Select
              value={formData.tenant_id}
              onChange={(e) => setFormData({ ...formData, tenant_id: e.target.value })}
              label="所属租户"
            >
              {tenants.map((tenant) => (
                <MenuItem key={tenant.id} value={tenant.id}>
                  {tenant.tenant_name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <TextField
            margin="dense"
            label="应用名称"
            fullWidth
            variant="outlined"
            value={formData.app_name}
            onChange={(e) => setFormData({ ...formData, app_name: e.target.value })}
            placeholder="例如: vpn, ssh, rdp"
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="描述"
            fullWidth
            variant="outlined"
            multiline
            rows={3}
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="可选，用于说明该应用的用途"
          />
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={handleCloseDialog}>取消</Button>
          <Button onClick={handleSave} variant="contained">
            {editingApp ? '更新' : '创建'}
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

export default AppManagement;
