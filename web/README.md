# NATMap React 管理后台

基于 React + Material-UI 的现代化管理后台，用于管理 NATMap 的租户、应用和映射记录。

## 功能特性

- **概览面板**: 显示服务器状态、统计数据和快速操作指南
- **租户管理**: 创建、编辑、删除租户（公司/组织）
- **应用管理**: 创建、编辑、删除应用，支持关联租户
- **映射管理**: 创建、编辑、删除 NAT 映射记录
- **响应式设计**: 支持桌面端和移动端访问
- **现代化 UI**: 使用 Material-UI 组件库，美观大方

## 技术栈

- **React 18**: 前端框架
- **Material-UI 5**: UI 组件库
- **React Router 6**: 路由管理
- **Axios**: HTTP 客户端
- **Vite**: 构建工具

## 开发环境

### 前置要求

- Node.js 16+
- npm 或 yarn

### 安装依赖

```bash
cd web
npm install
```

### 启动开发服务器

```bash
npm run dev
```

开发服务器运行在 http://localhost:5173

### 构建生产版本

```bash
npm run build
```

构建后的文件位于 `dist/` 目录。

## 部署

### 与 Go 后端集成

1. 构建前端：
```bash
cd web
npm run build
```

2. 复制构建文件到 Go 项目的 web 目录：
```bash
# Windows
xcopy /E /I /Y dist\* ..\web\

# Linux/Mac
cp -r dist/* ../web/
```

3. 启动 Go 服务：
```bash
cd ..
./natmap.exe
```

4. 访问管理后台：
```
http://localhost:8080
```

## API 接口

管理后台通过以下 API 与后端通信：

### 租户管理
- `GET /api/admin?type=tenant` - 获取租户列表
- `POST /api/admin?type=tenant` - 创建租户
- `PUT /api/admin?type=tenant&id={id}` - 更新租户
- `DELETE /api/admin?type=tenant&id={id}` - 删除租户

### 应用管理
- `GET /api/admin?type=app` - 获取应用列表
- `POST /api/admin?type=app` - 创建应用
- `PUT /api/admin?type=app&id={id}` - 更新应用
- `DELETE /api/admin?type=app&id={id}` - 删除应用

### 映射管理
- `GET /api/admin?type=mapping` - 获取映射列表
- `POST /api/admin?type=mapping` - 创建映射
- `PUT /api/admin?type=mapping&id={id}` - 更新映射
- `DELETE /api/admin?type=mapping&id={id}` - 删除映射

## 项目结构

```
web/
├── src/
│   ├── components/       # 可复用组件
│   │   └── Layout.jsx   # 主布局组件
│   ├── pages/           # 页面组件
│   │   ├── Dashboard.jsx       # 概览页面
│   │   ├── TenantManagement.jsx # 租户管理
│   │   ├── AppManagement.jsx    # 应用管理
│   │   └── MappingManagement.jsx # 映射管理
│   ├── services/        # API 服务
│   │   └── api.js      # API 接口封装
│   ├── App.jsx         # 主应用组件
│   └── main.jsx        # 应用入口
├── dist/               # 构建输出
└── package.json
```

## 配置

### API 地址配置

在 `src/services/api.js` 中修改 API 基础地址：

```javascript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
```

或在运行构建时设置环境变量：

```bash
VITE_API_URL=http://your-api-server npm run build
```

## 浏览器支持

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## 许可证

MIT
