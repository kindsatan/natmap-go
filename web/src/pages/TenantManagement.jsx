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
  Alert,
  Snackbar,
  Tooltip,
  Chip,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Business as BusinessIcon,
} from '@mui/icons-material';
import { tenantApi } from '../services/api';

function TenantManagement() {
  const [tenants, setTenants] = useState([]);
  const [loading, setLoading] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [editingTenant, setEditingTenant] = useState(null);
  const [formData, setFormData] = useState({ tenant_name: '' });
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // 加载租户列表
  const loadTenants = async () => {
    setLoading(true);
    try {
      const data = await tenantApi.getList();
      setTenants(Array.isArray(data) ? data : []);
    } catch (error) {
      showSnackbar(error.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTenants();
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
    setEditingTenant(null);
    setFormData({ tenant_name: '' });
    setOpenDialog(true);
  };

  // 打开编辑对话框
  const handleOpenEdit = (tenant) => {
    setEditingTenant(tenant);
    setFormData({ tenant_name: tenant.tenant_name });
    setOpenDialog(true);
  };

  // 关闭对话框
  const handleCloseDialog = () => {
    setOpenDialog(false);
    setEditingTenant(null);
  };

  // 保存租户
  const handleSave = async () => {
    if (!formData.tenant_name.trim()) {
      showSnackbar('请输入租户名称', 'warning');
      return;
    }

    try {
      if (editingTenant) {
        await tenantApi.update(editingTenant.id, formData);
        showSnackbar('租户更新成功');
      } else {
        await tenantApi.create(formData);
        showSnackbar('租户创建成功');
      }
      handleCloseDialog();
      loadTenants();
    } catch (error) {
      showSnackbar(error.message, 'error');
    }
  };

  // 删除租户
  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除这个租户吗？相关的应用和映射也将被删除。')) {
      return;
    }

    try {
      await tenantApi.delete(id);
      showSnackbar('租户删除成功');
      loadTenants();
    } catch (error) {
      showSnackbar(error.message, 'error');
    }
  };

  return (
    <Box>
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h5" component="h1" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <BusinessIcon color="primary" />
              租户管理
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleOpenCreate}
              disabled={loading}
            >
              添加租户
            </Button>
          </Box>

          <TableContainer component={Paper} variant="outlined">
            <Table>
              <TableHead>
                <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
                  <TableCell>ID</TableCell>
                  <TableCell>租户名称</TableCell>
                  <TableCell>创建时间</TableCell>
                  <TableCell align="right">操作</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {tenants.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} align="center" sx={{ py: 4, color: 'text.secondary' }}>
                      暂无租户数据
                    </TableCell>
                  </TableRow>
                ) : (
                  tenants.map((tenant) => (
                    <TableRow key={tenant.id} hover>
                      <TableCell>
                        <Chip label={tenant.id} size="small" variant="outlined" />
                      </TableCell>
                      <TableCell>
                        <Typography fontWeight="medium">{tenant.tenant_name}</Typography>
                      </TableCell>
                      <TableCell sx={{ color: 'text.secondary' }}>
                        {tenant.created_at}
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="编辑">
                          <IconButton
                            size="small"
                            onClick={() => handleOpenEdit(tenant)}
                            color="primary"
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="删除">
                          <IconButton
                            size="small"
                            onClick={() => handleDelete(tenant.id)}
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
        <DialogTitle>
          {editingTenant ? '编辑租户' : '添加租户'}
        </DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="租户名称"
            fullWidth
            variant="outlined"
            value={formData.tenant_name}
            onChange={(e) => setFormData({ ...formData, tenant_name: e.target.value })}
            placeholder="请输入租户名称"
            sx={{ mt: 1 }}
          />
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={handleCloseDialog}>取消</Button>
          <Button onClick={handleSave} variant="contained">
            {editingTenant ? '更新' : '创建'}
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

export default TenantManagement;
