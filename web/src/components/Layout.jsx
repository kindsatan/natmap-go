import React, { useState } from 'react';
import {
  Box,
  AppBar,
  Toolbar,
  Typography,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  IconButton,
  Divider,
  Container,
  Avatar,
  Menu,
  MenuItem,
  Tooltip,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Business as BusinessIcon,
  Apps as AppsIcon,
  Route as RouteIcon,
  Dashboard as DashboardIcon,
  Logout as LogoutIcon,
  AccountCircle as AccountIcon,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const drawerWidth = 240;

const menuItems = [
  { path: '/', label: '概览', icon: <DashboardIcon /> },
  { path: '/tenants', label: '租户管理', icon: <BusinessIcon /> },
  { path: '/apps', label: '应用管理', icon: <AppsIcon /> },
  { path: '/mappings', label: '映射管理', icon: <RouteIcon /> },
];

function Layout({ children }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleNavigation = (path) => {
    navigate(path);
    setMobileOpen(false);
  };

  const drawer = (
    <Box>
      <Toolbar sx={{ justifyContent: 'center' }}>
        <Typography variant="h6" noWrap component="div" sx={{ fontWeight: 'bold' }}>
          NATMap 管理
        </Typography>
      </Toolbar>
      <Divider />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.path} disablePadding>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => handleNavigation(item.path)}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: 'primary.main',
                  color: 'white',
                  '&:hover': {
                    backgroundColor: 'primary.dark',
                  },
                  '& .MuiListItemIcon-root': {
                    color: 'white',
                  },
                },
              }}
            >
              <ListItemIcon
                sx={{
                  color: location.pathname === item.path ? 'white' : 'inherit',
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      {/* 顶部导航栏 */}
      <AppBar
        position="fixed"
        sx={{
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            NATMap 管理后台
          </Typography>
          
          {/* 用户信息菜单 */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="body2" color="inherit" sx={{ opacity: 0.9 }}>
              {user?.username || 'Admin'}
            </Typography>
            <Tooltip title="账户菜单">
              <IconButton onClick={handleMenuOpen} color="inherit">
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'secondary.main' }}>
                  <AccountIcon />
                </Avatar>
              </IconButton>
            </Tooltip>
          </Box>
          
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            transformOrigin={{ vertical: 'top', horizontal: 'right' }}
          >
            <MenuItem disabled>
              <Typography variant="body2" color="text.secondary">
                角色: {user?.role === 'admin' ? '管理员' : '用户'}
              </Typography>
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>
              <LogoutIcon fontSize="small" sx={{ mr: 1 }} />
              退出登录
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* 侧边导航栏 - 移动端 */}
      <Box
        component="nav"
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', md: 'none' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>

        {/* 侧边导航栏 - 桌面端 */}
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', md: 'block' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: drawerWidth,
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      {/* 主内容区 */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          bgcolor: '#f5f5f5',
          p: 3,
          width: { md: `calc(100% - ${drawerWidth}px)` },
          minHeight: '100vh',
        }}
      >
        <Toolbar />
        <Container maxWidth="xl" sx={{ py: 2 }}>
          {children}
        </Container>
      </Box>
    </Box>
  );
}

export default Layout;
