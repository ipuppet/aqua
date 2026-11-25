# Passport 服务专项规范

IAM 总图见 [ARCHITECTURE.md](../standards/ARCHITECTURE.md) §2.3。本服务是认证面。

## 1. 职责

负责：登录、注册、验证码、密码重置、token 刷新与登出；会话查询与撤销（前端直接调 passport）。

账号行用 `domain.Account`（凭证与 JWT claims 所需字段），不是 user 的档案实体。

不负责：用户资料、角色权限、管理后台聚合。

## 2. 文件

logic：入口 `auth.go`（`AuthService`）；子文件 `login.go`、`logout.go`、`registration.go`、`password_reset.go`、`token.go`、`turnstile.go`、`session.go`。禁止回到 `passport.go` 单文件。

- `auth.go`：构造器、通用 token 生命周期配置。
- `login.go`：账号密码登录。
- `token.go`：access/refresh 验证、轮换、撤销点。
- `session.go`：会话查询与撤销（admin gRPC 与用户自助共用）。

dao：`user.go`、`session.go`、`verification_code.go`、`password_reset_token.go`、`login_audit.go`。对外端口在 `ports.go` 的 `AuthRepository`。refresh 会话走 `session.go` 的 Redis 键，不要再拆 `refresh_token.go`。不要在 logic 拼 SQL。

grpcserver：`admin_server.go`（`ListUserSessions`、`RevokeUserSession`）、`passport_server.go`（登录/注册/密码重置/会话自助）。

proto：`admin.proto` → `AdminPassportService`；`passport.proto` → `PassportService`。

## 3. 安全

- access 短寿命；refresh 长寿命，必须能轮换与撤销。
- 退出默认只撤当前 refresh；`logout_all=true` 撤该用户全部历史 refresh。
- 登录与刷新都要检查账号锁定。
- 刷新时先校验撤销标记与 jti 存活，再签新对。撤销标记用 Unix 秒。

## 4. HTTP

路径不加服务名前缀：`/api/login/password`、`/api/token`、`/api/register` 等。

- 当前用户会话：`/api/sessions`，`middleware.Authenticate(publicKey)`。
- 管理会话：`/api/users/{uid}/sessions`，需 `user.manager`。

## 5. 本服务回归补充

调整构造函数签名时，同步 `logic/provider.go`、`wire_gen.go` 和对应测试。
