# 开发关键信息与扩展指南

## 概览
- 框架：`gin`（Web）、`gorm` + `glebarez/sqlite`（DB）、`golang-jwt/jwt/v5`（JWT）、`bcrypt`（密码）
- 功能：注册、登录、刷新令牌、登出/登出全部、角色与权限、Swagger 测试页面
- 目标：为后续增量开发提供明确的入口与最佳实践

## 目录结构与职责
- `cmd/server/main.go`：应用入口、加载配置、DB 迁移、路由注册、示例受权资源
- `internal/config/config.go`：配置加载（支持 `.env` 与环境变量）
- `internal/db/db.go`：SQLite 初始化（自动创建数据目录）
- `internal/models/*`：数据模型（User、RefreshToken、Permission）
- `internal/auth/*`：密码哈希/校验、JWT 生成/解析、刷新令牌生成与哈希
- `internal/middleware/*`：认证中间件、角色中间件、权限中间件
- `internal/handlers/*`：业务控制器（Auth、User、Admin）
- `internal/docs/docs.go`：Swagger 文档（内置 JSON，免生成）
- `docs/summary.md`：功能概览说明
- `docs/dev-notes.md`：开发者信息与扩展指南（本文件）

## 配置关键项
- `HTTP_ADDR` 默认 `:8080`
- `SQLITE_PATH` 默认 `./data/app.db`
- `JWT_SECRET` 必填强随机值
- `TOKEN_TTL` 默认 `24h`
- `REFRESH_TTL` 默认 `168h`（7 天）
- `BCRYPT_COST` 默认 `12`
- `SEED_USER=1`、`SEED_USERNAME`、`SEED_PASSWORD`：启动时创建管理员种子用户

## 数据模型
- `internal/models/user.go:1`
  - 字段：`id`、`username`（唯一）、`password_hash`、`email?`、`role`、`is_active`、`created_at`、`updated_at`、`last_login_at?`
  - 自动时间：`created_at` 使用 `gorm:"autoCreateTime"`，`updated_at` 使用 `gorm:"autoUpdateTime"`
- `internal/models/refresh_token.go:1`
  - 字段：`id`、`user_id`、`token_hash`（sha256）、`expires_at`、`revoked`、时间戳
  - 仅存储哈希，原始值下发给客户端
- `internal/models/permission.go:1`
  - 字段：`id`、`role`、`resource`、`action`、`allowed`

## 认证与授权
- 认证中间件：`internal/middleware/auth.go:12`
  - 从 `Authorization: Bearer <token>` 解析 JWT，注入 `user_id`、`username`、`role`
- 角色中间件：`internal/middleware/roles.go:1`
  - 允许集合校验，如 `RequireRole("admin")`
- 权限中间件：`internal/middleware/perm.go:1`
  - DB 查询 `role+resource+action+allowed=true`，未命中返回 403

## JWT 与刷新令牌
- Claims：`internal/auth/jwt.go:1`
  - 含 `user_id`、`username`、`role`、`exp`、`iat`
- 生成与解析：`internal/auth/jwt.go:1`
- 刷新令牌：`internal/auth/refresh.go:1`
  - 生成 32 字节随机原始值与其哈希（sha256）
  - 存库时仅存哈希
- 生命周期（`internal/handlers/auth.go:31,76,105,141,155`）：
  1. 登录：签发访问令牌与刷新令牌（存库哈希）
  2. 刷新：校验有效→撤销旧刷新令牌→签发新访问与刷新令牌→存新刷新令牌哈希
  3. 登出：按哈希撤销刷新令牌
  4. 登出全部：撤销该用户所有未撤销的刷新令牌

## 路由与接口
- 注册：`POST /api/v1/auth/register`（`internal/handlers/auth.go:76`）
- 登录：`POST /api/v1/auth/login`（`internal/handlers/auth.go:31`）
- 刷新：`POST /api/v1/auth/refresh`（`internal/handlers/auth.go:105`）
- 登出：`POST /api/v1/auth/logout`（`internal/handlers/auth.go:141`）
- 登出全部：`POST /api/v1/auth/logout_all`（需认证）（`internal/handlers/auth.go:155`）
- 当前用户：`GET /api/v1/me`（需认证）（`internal/handlers/user.go:18`）
- 管理员接口（需 `RequireRole("admin")`）：`cmd/server/main.go:41`
  - 用户：`GET /api/v1/admin/users`、`PUT /api/v1/admin/users/:id/role`
  - 权限：`GET/POST/DELETE /api/v1/admin/permissions`
- 权限保护示例资源：`GET /api/v1/reports`（需 `user` 对 `reports:read` 权限）（`cmd/server/main.go:44`）

## Swagger 测试环境
- 入口集成：`cmd/server/main.go:37` 注册路由 `GET /swagger/*any`
- 文档源：`internal/docs/docs.go:1`（内置 JSON，定义 `BearerAuth`）
- 使用：打开 `http://localhost:<port>/swagger/index.html`，在右上角 “Authorize” 填写 `Bearer <token>`

## 构建与运行
- 开发运行：`go run cmd/server/main.go`
- 编译可执行：`go build -o bin/natmap.exe cmd/server/main.go`
- 运行可执行：`bin/natmap.exe`（读取 `.env`）
- 端口覆盖示例：`set HTTP_ADDR=:8080 && bin\natmap.exe`

## 测试与示例
- 单元测试：`go test ./...`
- Curl 登录（Windows PowerShell 使用 `curl.exe`）：
  - `curl.exe -s -X POST "http://localhost:8080/api/v1/auth/login" -H "Content-Type: application/json" -d "{\"username\":\"admin\",\"password\":\"Admin123!\"}"`
- Authorize 值示例：`Bearer <token>`（注意包含前缀）

## 常见问题与排查
- 401 Unauthorized：未设置或格式错误的 `Authorization` 头；令牌过期；端口不一致
- 403 Forbidden：角色不足或权限未配置；在管理员接口新增权限后重试
- Swagger 未生效：确认访问的是集成路由 `/swagger/index.html`；清除浏览器缓存或使用无痕模式

## 安全规范
- 密码只存储 `bcrypt` 哈希
- 刷新令牌仅存储哈希；原始值仅下发客户端
- 错误信息模糊化：避免暴露用户存在与否
- 严禁硬编码密钥：使用 `JWT_SECRET` 环境变量

## 扩展指南（建议路径）
1. 新增受权限保护的资源
   - 定义路由：在 `cmd/server/main.go` 中添加 `GET/POST` 路由
   - 附加中间件：`AuthMiddleware` + `RequirePermission(db, resource, action)`
   - 管理员创建权限规则：`POST /api/v1/admin/permissions`
2. 引入角色分级
   - 扩展 `Role` 取值与语义（例如 `editor`、`viewer`）
   - 在权限表中为新角色配置资源动作
3. 权限缓存
   - 通过 `sync.Map` 或本地 LRU 缓存查询结果，设置 TTL 与失效策略
4. 会话管控
   - 缩短访问令牌 TTL，增强刷新令牌旋转策略
   - 可选：访问令牌黑名单（需持久化与高性能匹配）
5. 审计与日志
   - 记录登录、刷新、登出与权限变更事件（建议结构化日志 + 可查询后端）

## 开发约定
- 统一使用 `UTC` 时间写入（注册时已采用 `time.Now().UTC()`）
- 路由前缀：API 统一以 `/api/v1` 开头
- 控制器只返回必要字段，避免输出敏感信息（如 `password_hash`）
- 代码风格：保持现有 import 与分层结构，避免跨层耦合

## 关键代码引用
- 入口与路由：`cmd/server/main.go:18`
- 认证中间件：`internal/middleware/auth.go:12`
- 角色中间件：`internal/middleware/roles.go:1`
- 权限中间件：`internal/middleware/perm.go:1`
- 登录：`internal/handlers/auth.go:31`
- 注册：`internal/handlers/auth.go:76`
- 刷新：`internal/handlers/auth.go:105`
- 登出/登出全部：`internal/handlers/auth.go:141,155`
- 用户模型：`internal/models/user.go:1`（自动时间标签在 `CreatedAt/UpdatedAt`）
- 刷新令牌模型：`internal/models/refresh_token.go:1`
- 权限模型：`internal/models/permission.go:1`
- Swagger 文档：`internal/docs/docs.go:1`

---
如需新增模块或重构，请在本文件追加“变更记录”与“迁移步骤”，确保后续协作可追踪与复现。
