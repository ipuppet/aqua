# User 服务专项规范

IAM 总图见 [ARCHITECTURE.md](../standards/ARCHITECTURE.md) §2.3。本服务是目录面。

## 1. 职责

负责：用户资料查询/更新、用户管理、RBAC（角色/权限/关系）。

不负责：登录签发、token 生命周期、任何 session 接口。前端会话必须直连 passport。

## 2. 文件

logic / dao 都按三块：`user.go`（公开资料与当前用户）、`user_management.go`（管理侧列表/检索/更新）、`access_control.go`（角色、权限、用户角色）。

`user_management.go` 不要改成 `management.go`：本包没有管理聚合入口，management 只是 user 的子能力。

端口：`UserRepository`、`UserManagementRepository`、`AccessControlRepository`。

- `UserManagementService` 只依赖 `UserManagementRepository`。
- `AccessControlService` 依赖 `AccessControlRepository` + `rbac.EventBus`。

grpcserver：`admin_server.go`（管理用户与 RBAC）、`user_server.go`（自助与公开资料）。`Admin*` 可留在 gRPC 入口名（如 `AdminUserServiceServer`），包内实现仍用能力名。

## 3. 本服务的 RBAC 约束

读模型与 HTTP 判定见架构 §2.3。这里只补充本域：

- 表列定义在 `services/user/domain/schema`（如 `UserColumns`），不要放进 loader / dao 包。`UserColumns` 只能从该包引用。
- `services/*/dao` 不要自己写角色/权限读取 SQL，委托 `pkg/rbac.Reader`。
- 不要在 logic / grpcserver 里再写一套权限码判断。
- RBAC 缓存（`user:GetRoles:*`、`role:GetPermissions:*`）只通过 RBAC 事件消费器失效。logic 禁止对这组键 `Del` / `DelMatch`，但可以发事件（UserRoleAssigned、RoleUpdated、PermissionDeleted）。

## 4. HTTP

- 当前用户：`/api/users/me`、`/api/users/me/menu`、`/api/users/me/structure`。
- 公开资料：`GET /api/users/{uid}`（仅 GET；PUT 走管理权限）。
- 静态段 `me`、`search` 必须先于 `/{uid}` 注册。
- 管理会话不在本服务：`/api/users/{uid}/sessions` 由 passport 挂载。

挂载：`RegisterUserServiceHandlerServer`、`RegisterAdminUserServiceHandlerServer`。
