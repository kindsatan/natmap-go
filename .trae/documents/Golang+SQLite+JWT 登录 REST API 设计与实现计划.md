## 需求概述
- 使用 Golang 编写 RESTful API
- 基于 SQLite 存储用户数据
- 实现登录接口，返回 JWT token；后续业务均需携带 token 访问

## 技术选型
- Web 框架：gin（轻量、生态成熟）
- ORM：GORM + `gorm.io/driver/sqlite`（简化迁移与模型维护）
- Token：`github.com/golang-jwt/jwt/v5`
- 密码加密：`golang.org/x/crypto/bcrypt`
- 配置：从环境变量读取（如 `JWT_SECRET`, `TOKEN_TTL`, `BCRYPT_COST`, `SQLITE_PATH`）

## 库表设计（SQLite DDL）
- `users`
  - `id` INTEGER PRIMARY KEY AUTOINCREMENT
  - `username` TEXT UNIQUE NOT NULL
  - `password_hash` TEXT NOT NULL
  - `email` TEXT NULL
  - `is_active` INTEGER NOT NULL DEFAULT 1
  - `created_at` DATETIME NOT NULL
  - `updated_at` DATETIME NOT NULL
  - `last_login_at` DATETIME NULL
- 说明：登录仅需 `username + password`；邮箱可选；采用 `is_active` 控制停用；记录登录时间

## 认证与授权
- 登录成功后签发 JWT：
  - Claims：`sub`=用户ID，`username`，`exp`，`iat`
  - 签名算法：`HS256`，密钥从 `JWT_SECRET` 环境变量读取
  - Token 过期时间：默认 24h（可由 `TOKEN_TTL` 配置）
- 授权：
  - 中间件校验 `Authorization: Bearer <token>`
  - 拦截并解析 JWT，失败返回 401；成功将用户信息注入上下文

## API 设计
- `POST /api/v1/auth/login`
  - 入参：`{ "username": string, "password": string }`
  - 返回：`{ "token": string, "token_type": "Bearer", "expires_in": number }`
  - 错误：400（参数错误）、401（用户名或密码错误）、423（用户被禁用）、500（服务器错误）
- `GET /api/v1/me`（受保护示例）
  - 头部：`Authorization: Bearer <token>`
  - 返回：当前用户基础信息（`id, username, email, last_login_at`）
- 可选：`POST /api/v1/auth/register`（仅开发/演示用途）或提供种子脚本创建初始用户

## 项目结构
- `cmd/server/main.go`：入口，加载配置、初始化 DB、注册路由、中间件
- `internal/config`：配置加载（环境变量）
- `internal/db`：数据库连接与迁移
- `internal/models`：`User` 模型（GORM 标签）
- `internal/handlers`：`AuthHandler`（登录）、`UserHandler`（me）
- `internal/middleware`：`AuthMiddleware`（JWT 校验）
- `internal/auth`：JWT 生成/解析、密码比较封装

## 初始化与迁移
- 启动时自动执行 GORM 自动迁移：创建 `users` 表
- 若无用户且启用 `SEED_USER=1`，则根据 `SEED_USERNAME/SEED_PASSWORD` 创建一个初始用户（密码使用 bcrypt 加密）

## 错误处理与安全
- 登录统一返回模糊错误（不透露用户名是否存在）
- bcrypt 成本默认 12（`BCRYPT_COST` 可调），避免过低；禁止明文存储密码
- JWT 密钥必须来自安全环境变量；禁止硬编码
- 响应不包含敏感字段（如 `password_hash`）
- 考虑基础速率限制（后续可加），防爆破

## 测试与验证
- 单元测试：
  - 密码哈希与校验
  - 登录成功/失败路径
  - JWT 生成与过期校验
  - 受保护路由访问（有/无/过期 token）
- 集成验证：
  - 启动服务后，用 `curl` 调用 `login` 获取 token；携带 token 访问 `/me`

## 交付内容
- 可运行的服务（`go run cmd/server/main.go`）
- SQLite 数据文件路径由 `SQLITE_PATH` 指定（默认 `./data/app.db`）
- 示例 `.env.example`，说明启动方式与环境变量

## 后续扩展
- 刷新令牌/黑名单机制、角色与权限（RBAC）、审计日志
- 注册、重置密码、多因素认证

请确认是否按上述计划实施；确认后我将开始编写代码与测试，并在本仓库中新增相应文件与实现。