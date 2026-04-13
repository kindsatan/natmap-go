# NATMap Go 版本

NATMap 应用的 Go 语言重构版本，从 Cloudflare Workers + D1 迁移到 Go + MySQL 架构。

## 特性

- **高性能**: 使用 Go 语言的高并发特性
- **内存缓存**: Cache-Aside 模式，支持 60 秒 TTL
- **MySQL 存储**: 本地数据库，数据自主可控
- **RESTful API**: 兼容原有 API 接口
- **管理后台**: 完整的 CRUD 功能
- **JWT 认证**: 安全的 API 认证机制

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 5.7+
- **缓存**: 内存缓存 (go-cache)
- **日志**: Zap
- **配置**: Viper

## 快速开始

### 1. 安装依赖

```bash
# 安装 Go 1.21+
# https://golang.org/dl/

# 克隆项目
git clone https://github.com/kindsatan/natmap-go.git
cd natmap-go

# 下载依赖
go mod download
```

### 2. 配置 MySQL

```bash
# 登录 MySQL
mysql -u root -p

# 执行初始化脚本
source scripts/init.sql
```

### 3. 配置应用

编辑 `config.yaml`:

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: "123456"
  dbname: natmap
```

### 4. 运行数据库迁移

```bash
go run cmd/natmap/main.go -migrate
```

### 5. 启动服务

```bash
go run cmd/natmap/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 公开接口

#### 健康检查
```bash
GET /api/test
```

#### 获取映射
```bash
GET /api/get?tenant_id=3&app_id=3
```

响应:
```json
{
  "public_ip": "36.96.128.250",
  "public_port": 10472,
  "updated_at": "2026-04-08 20:34:36",
  "_cache": "HIT"
}
```

#### 更新映射
```bash
POST /api/update
Content-Type: application/json

{
  "app": "http2tunnel1",
  "ip": "36.96.128.250",
  "port": 10472,
  "proto": "tcp",
  "local_ip": "192.168.1.100",
  "local_port": 8080
}
```

### 管理后台接口 (需要 Basic Auth)

默认账号: `admin` / `admin123`

#### 租户管理
```bash
# 列表
GET /api/admin?type=tenant

# 创建
POST /api/admin?type=tenant
{
  "tenant_name": "新租户"
}

# 更新
PUT /api/admin?type=tenant&id=1
{
  "tenant_name": "更新后的租户名"
}

# 删除
DELETE /api/admin?type=tenant&id=1
```

#### 应用管理
```bash
# 列表
GET /api/admin?type=app

# 创建
POST /api/admin?type=app
{
  "tenant_id": 1,
  "app_name": "新应用",
  "description": "应用描述"
}

# 更新
PUT /api/admin?type=app&id=1
{
  "app_name": "更新后的应用名",
  "description": "新描述"
}

# 删除
DELETE /api/admin?type=app&id=1
```

#### 映射管理
```bash
# 列表
GET /api/admin?type=mapping

# 创建
POST /api/admin?type=mapping
{
  "tenant_id": 1,
  "app_id": 1,
  "public_ip": "36.96.128.250",
  "public_port": 10472,
  "local_ip": "192.168.1.100",
  "local_port": 8080,
  "protocol": "tcp"
}

# 更新
PUT /api/admin?type=mapping&id=1
{
  "public_ip": "36.96.128.251",
  "public_port": 10473
}

# 删除
DELETE /api/admin?type=mapping&id=1
```

## 性能测试

```bash
# Linux/Mac
cd scripts
chmod +x benchmark-fast.sh
./benchmark-fast.sh --requests 100

# Windows PowerShell
.\scripts\benchmark-fast.ps1 -Requests 100
```

## 数据迁移

### 从 Cloudflare D1 迁移

1. 从 D1 导出数据为 JSON 格式
2. 使用迁移工具导入:

```bash
go run cmd/migrate/main.go \
  -tenants tenants.json \
  -apps apps.json \
  -mappings mappings.json
```

## 部署

### 使用 systemd (Linux)

创建 `/etc/systemd/system/natmap.service`:

```ini
[Unit]
Description=NATMap Service
After=network.target

[Service]
Type=simple
User=natmap
WorkingDirectory=/opt/natmap
ExecStart=/opt/natmap/natmap
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl enable natmap
sudo systemctl start natmap
```

### 使用 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o natmap cmd/natmap/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/natmap .
COPY config.yaml .
EXPOSE 8080
CMD ["./natmap"]
```

构建并运行:
```bash
docker build -t natmap .
docker run -d -p 8080:8080 natmap
```

## 配置说明

### config.yaml

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| server.port | 服务端口 | 8080 |
| server.mode | 运行模式 (debug/release) | release |
| database.host | 数据库主机 | localhost |
| database.port | 数据库端口 | 3306 |
| database.username | 数据库用户名 | root |
| database.password | 数据库密码 | - |
| database.dbname | 数据库名 | natmap |
| cache.enabled | 是否启用缓存 | true |
| cache.ttl | 缓存过期时间(秒) | 60 |
| jwt.secret | JWT 密钥 | - |
| jwt.expire_hours | Token 过期时间(小时) | 24 |

## 目录结构

```
natmap-go/
├── cmd/
│   └── natmap/          # 主程序入口
│       └── main.go
├── internal/
│   ├── cache/           # 缓存层
│   ├── config/          # 配置管理
│   ├── handlers/        # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── migrator/        # 数据迁移
│   └── models/          # 数据模型
├── scripts/             # 脚本文件
│   ├── init.sql         # 数据库初始化
│   └── benchmark*.sh    # 性能测试
├── web/                 # 静态文件
├── config.yaml          # 配置文件
├── go.mod               # Go 模块
└── README.md            # 本文档
```

## 与原 Cloudflare 版本对比

| 特性 | Cloudflare 版本 | Go 版本 |
|------|----------------|---------|
| 部署 | Cloudflare Pages | 本地/服务器 |
| 数据库 | Cloudflare D1 | MySQL |
| 缓存 | Cloudflare KV | 内存缓存 |
| 性能 | 边缘计算延迟 | 本地低延迟 |
| 数据控制 | 受限 | 完全自主 |
| 成本 | 按量付费 | 固定成本 |

## 许可证

MIT License
