# 项目总结

## 概览
- 提供基于 Go 的 REST API：注册、登录、刷新令牌、登出与登出全部、角色与权限控制、受保护资源示例
- 技术栈：`gin`（Web）、`GORM` + 纯 Go 驱动 `github.com/glebarez/sqlite`（数据库）、`golang-jwt/jwt/v5`（访问令牌）、`bcrypt`（密码哈希）
- 已实现：用户模型与迁移、注册/登录/刷新/登出接口、JWT 携带角色、角色中间件与权限中间件、管理员接口（用户与权限管理）、受权限保护示例路由

## 数据模型
- User（`internal/models/user.go`）
  - `id`、`username`（唯一）、`password_hash`、`email?`、`role`（默认 `user`）、`is_active`、`created_at`、`updated_at`、`last_login_at?`
- RefreshToken（`internal/models/refresh_token.go`）
  - `id`、`user_id`、`token_hash`（`sha256`）、`expires_at`、`revoked`、时间戳
- Permission（`internal/models/permission.go`）
  - `id`、`role`、`resource`、`action`、`allowed`

## 配置项（`.env` 或环境变量）
- `HTTP_ADDR`（默认 `:8080`）
- `SQLITE_PATH`（默认 `./data/app.db`）
- `JWT_SECRET`（必须设置为强随机值）
- `TOKEN_TTL`（默认 `24h`）
- `REFRESH_TTL`（默认 `168h`，即 7 天）
- `BCRYPT_COST`（默认 `12`）
- `SEED_USER`（`1` 开启种子用户）
- `SEED_USERNAME`、`SEED_PASSWORD`（创建管理员种子用户）

## 接口与行为
- 注册 `POST /api/v1/auth/register`
  - Body：`{"username":"<string>","password":"<string>","email":"<string?>"}`
  - 返回：`{"id":<uint>,"username":"<string>"}`
- 登录 `POST /api/v1/auth/login`
  - Body：`{"username":"<string>","password":"<string>"}`
  - 返回：`{"token":"<jwt>","token_type":"Bearer","expires_in":<seconds>,"refresh_token":"<string>"}`
- 刷新 `POST /api/v1/auth/refresh`
  - Body：`{"refresh_token":"<string>"}`
  - 行为：验证未撤销且未过期；撤销旧刷新令牌；返回新访问令牌与新刷新令牌
- 登出 `POST /api/v1/auth/logout`（需携带刷新令牌原始值）
  - Body：`{"refresh_token":"<string>"}`
  - 行为：将对应刷新令牌标记为撤销
- 登出全部 `POST /api/v1/auth/logout_all`（需认证）
  - 行为：撤销当前用户所有未撤销的刷新令牌
- 当前用户 `GET /api/v1/me`（需认证）
  - 返回：`{"id","username","email","last_login_at"}`
- 管理员（需 `RequireRole("admin")`）
  - `GET /api/v1/admin/ping` 返回 `{ ok: true }`
  - 用户管理：`GET /api/v1/admin/users`，`PUT /api/v1/admin/users/:id/role`
  - 权限管理：`GET /api/v1/admin/permissions`，`POST /api/v1/admin/permissions`，`DELETE /api/v1/admin/permissions/:id`
- 受权限保护示例资源
  - `GET /api/v1/reports` 需权限 `role=user`、`resource=reports`、`action=read` 且 `allowed=true`

## 中间件
- 认证中间件 `AuthMiddleware`
  - 从 `Authorization: Bearer <token>` 解析 JWT，注入 `user_id`、`username`、`role`
- 角色中间件 `RequireRole("admin")`
  - 校验上下文中的 `role` 必须在允许集合内
- 权限中间件 `RequirePermission(db, resource, action)`
  - 基于 `role+resource+action+allowed=true` 查询权限表，未命中则 403

## 安全与规范
- 密码使用 `bcrypt` 哈希；从不返回明文或哈希
- 刷新令牌仅存储哈希（`sha256`）；原始值仅下发给客户端
- 登录与刷新失败返回模糊错误（避免暴露用户状态）
- JWT 携带 `role` 与过期；秘钥来源于环境变量
- 管理接口受角色保护；受权限资源路由按中间件检查

## 启动与验证
- 安装依赖：`go mod tidy`
- 启动：`go run cmd/server/main.go`
- 健康检查：`GET /health`
- 验证流程示例：
  1. 管理员登录获取 `token`
  2. 创建权限：`POST /api/v1/admin/permissions`（`{"role":"user","resource":"reports","action":"read","allowed":true}`）
  3. 普通用户登录获取 `token`，访问 `GET /api/v1/reports` 成功返回 `items`
  4. 刷新令牌：调用 `POST /api/v1/auth/refresh`，返回新访问与刷新令牌，旧刷新令牌被撤销
  5. 登出：`POST /api/v1/auth/logout` 撤销指定刷新令牌；`POST /api/v1/auth/logout_all` 撤销当前用户所有刷新令牌

## 关键代码位置
- 入口与路由：`cmd/server/main.go`
- 配置加载：`internal/config/config.go`
- 模型：`internal/models/user.go`、`internal/models/refresh_token.go`、`internal/models/permission.go`
- 认证工具：`internal/auth/password.go`、`internal/auth/jwt.go`、`internal/auth/refresh.go`
- 中间件：`internal/middleware/auth.go`、`internal/middleware/roles.go`、`internal/middleware/perm.go`
- 控制器：`internal/handlers/auth.go`、`internal/handlers/user.go`、`internal/handlers/admin.go`

## 单元测试
- 密码与 JWT 核心测试：`internal/auth/password_test.go`、`internal/auth/jwt_test.go`、`internal/auth/refresh_test.go`
- 运行：`go test ./...`

## 后续建议
- 权限通配支持（如资源前缀或 `*` 动作）
- 权限缓存（减少 DB 命中）与缓存失效策略
- 访问令牌黑名单或短 TTL + 刷新旋转加强会话控制
- 审计日志：登录、刷新、登出、权限变更等事件记录
