# 通用架构标准

全仓微服务的结构基线。编码风格、错误响应、回归命令见 [ENGINEERING.md](./ENGINEERING.md)；新服务骨架见 [SERVICE_TEMPLATE.md](./SERVICE_TEMPLATE.md)。
服务特有规则放在 `docs/services/`。

## 1. 目录结构

每个标准微服务两个入口：

- `services/<service>/cmd/`：进程入口。
- `services/<service>/`：业务实现。

`services/<service>/` 标准子模块：

- `domain/`：领域模型。
- `logic/`：用例编排与业务规则。
- `dao/`：DB / Cache。
- `grpcserver/`：gRPC 协议接入（含 admin 与业务 service）。
- `gateway.go`：HTTP 挂载（薄路由，只挂 grpc-gateway mux）。
- `main.go`：`ServiceName`、`Assemble`、module 装配。
- `wire.go` / `wire_gen.go`：依赖注入。

按需引入：

- `types/`：跨层 DTO / 共享接口。
- `pkg/`：与本域业务解耦的可复用组件。
- `internal/`：仅服务内可见。

默认不新增 `handler/`。例外见 [services/README.md](../services/README.md)。

## 2. 分层与依赖

### 2.1 方向

```
gateway (route mount) -> grpcserver -> logic -> dao
```

- `grpcserver` 不做业务决策。
- 禁止跨服务引用他服务 `logic` / `internal`。
- `domain` 作为底层模型，被 dao / logic / grpcserver 依赖。

### 2.2 各层职责

- `gateway`：路径分组、中间件。鉴权写在中间件，不把 `public` / `private` 写进 URL。统一用 `gateway.GinHandler` 包 grpc-gateway mux。对外 HTTP 由 Nginx 暴露与转发；CORS 也在 Nginx，服务侧不挂 CORS middleware。
- `grpcserver`：解包 → 调 logic → 组包。
- `logic`：业务规则、默认值、分页归一化、缓存失效。
- `dao`：只做持久化与缓存访问。接口集中在 `dao/ports.go`，实现按能力拆文件；读 / 写 / 统计不要塞进同一个 Repository。

### 2.3 IAM（认证、授权）

passport + user 构成 IAM 域，但是两个服务：

- **passport（认证 / IdP）**：登录、注册、session、JWT 签发与撤销。账号行用 `services/passport/domain.Account`。
- **user（目录）**：用户资料、RBAC 管理。档案用 `services/user/domain.UserModel`。
- **各服务门口**：
  - AuthN：`middleware.Authenticate` 验 JWT，写入 `authctx.TokenUser`。
  - AuthZ：需要权限码时挂 `middleware.RequireAnyPermission` / `RequireAllPermissions`，经 `pkg/rbac.Reader` 在本地判定。禁止每个请求打 user gRPC 做权限检查。

共享包：

- 跨服务请求身份只有 `authctx.TokenUser`（uid + roles）。禁止把 `UserModel` / `Account` 放进 `pkg/`。
- `pkg/rbac`：角色/权限读模型、缓存键、失效事件；是 `ums.rbac_*` 读侧内核，不是 transport。
- 角色/权限记录用 `pkg/rbac.Role` / `pkg/rbac.Permission`，user API 视图用 `RoleResponse` / `PermissionResponse`。不要在各服务再复制表模型。
- `proto/aqua/user`：用户与 RBAC 的 API 契约。管理 HTTP 挂在 user 服务 gateway 上，不要在 admin 进程再包一层。
- HTTP 权限判定放在 `pkg/transport/middleware`，不要再拆 `pkg/transport/authz`。
- 需要「用户资料 + 角色/权限」时，组合业务仓储与 `pkg/rbac.Reader`，不要新增聚合 loader 包。
- 需要 RBAC 读能力的服务在注入图声明 `rbac.Set`，把 `rbac.Reader` 交给授权中间件。

### 2.4 服务发现与内部 gRPC

- `pkg/infra/discovery`：etcd 查询、实例选择、缓存。不承载鉴权密钥。
- `pkg/transport/grpcx`：入站内部鉴权拦截器与 gRPC 错误映射。当前没有出站 gRPC 客户端工厂。
- 入站内部 token 由 `grpcx.UnaryServerInternalAuth` 校验，值来自 `models.InternalServiceConfig`（`app.yaml` 的 `internal_service.token`）。由 bootstrap 装配，业务层不要再拆 token 依赖。

## 3. 文件组织

- logic 按业务能力拆分，不要按技术动作碎片化。一个能力一个入口文件（类型、构造器、依赖）；大了再拆子文件，同一能力命名风格必须统一。禁止 `logic.go` / `<service>.go` 当长期 God File。
- logic 默认导出具体类型（`type UserService struct`），`NewXxxService` 返回 `*UserService`。禁止 `*Impl`，禁止默认再包一层同名接口。接口定义在消费方。
- dao 实现默认非导出（`type repository struct`），通过 `XxxRepository` 对外暴露。
- grpcserver 结构体可导出（`type UserServer struct`）。
- 可空字符串用 `pkg/types.NullString`，domain 禁止依赖 `infra/database`。

复杂服务防 God Package / 循环依赖：

- `types/`：跨层 DTO、共享 payload。不要把所有结构体都搬进去。
- `pkg/`：业务无关可复用能力。
- `internal/`：服务私有实现。
- 协议 DTO 留在 proto / grpcserver；logic / dao 结构体留在本层。

## 4. 启动与编排

### 4.1 `cmd/main.go`

默认：

```go
func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    <service>.ServiceName,
		Factory: <service>.Assemble,
	})
}
```


### 4.2 `services/<service>/main.go`

- `Assemble` 只做 module assembly：返回 `MountHTTP` / `MountGRPC`。监听器由 `bootstrap.Run`（单服务）或 `bootstrap.RunAll`（单体）创建。
- 禁止在 `Assemble` 手工加载 JWT key、filestore 等共享基础依赖。
- 可复用清理（如 registry close）走 infra helper，禁止复制 cleanup 闭包。
- 单体入口：`services/mono/cmd` 调用 `bootstrap.RunAll`（admin / passport / user / blog）。

## 5. Wire

### 5.1 命名

每层在 `provider.go` 暴露 set，包内只有一个时可直接叫 `Set`：

```go
var Set = wire.NewSet(
	NewRepository,
	wire.Bind(new(AuthRepository), new(*repository)),
)
```

- 构造函数统一 `NewXxx`，不要和 `ProvideXxx` 混用。
- set 用 `XxxSet`，或包内单 set 时用 `Set`。

### 5.2 `buildModule`

默认签名，且不要另套容器结构：

```go
func buildModule(rt *bootstrap.Runtime) (*bootstrap.Module, error)
```

`bootstrap.Runtime` 是进程级上下文（env / config / logger / cleanup），不是 IoC 容器。只有无法用 pkg provider 表达的临时参数，才允许未导出 injector，并在服务文档写明原因。

推荐组合顺序：

1. `bootstrap.RuntimeBaseSet`
2. pkg 通用（`bootstrap.FileStoreSet`、`config.*ConfigSet`、`token.VerifierSet`、`token.PrivateKeySet`）
3. pkg 基础设施（`persist.Set`、discovery、rbac 等）
4. 本域（dao → logic → grpcserver）
5. `provideModule`

### 5.3 归属与禁止项

- pkg provider：跨服务、无领域语义（config 字段提取、JWT key、filestore、registry cleanup）。
- service provider：仅本域（dao / logic / grpcserver / root middleware）。
- 构造器禁止直接依赖 `*models.AppConfig`，除非确实要跨多个无关配置域。优先注入具体结构体（`models.EtcdConfig`、`models.InternalServiceConfig` 等），用 `pkg/infra/config` 的 `*Set` 提取。
- 禁止在 Wire 图里直接提供 `string` / `int` / `bool`。
- 禁止多个服务重复声明同构 provider（`newPublicKey`、`newFileStore`）。两处及以上复用必须下沉到 pkg，以 `*Set` 暴露。
- JWT 校验用 `token.VerifierSet`（含 `config.JWTConfigSet` + `token.PublicKeySet`）；签发服务再加 `token.PrivateKeySet`。禁止在 `services/*/main.go` 手工 `Load*Key`。
- 持久化用 `persist.Set`（database + redis）。
- filestore 用 `bootstrap.FileStoreSet`。
- 启用 discovery 的服务，cleanup 必须走 shared helper。

## 6. 结构变更

- 重命名 / 拆分必须一次改完：gateway / grpcserver / logic / dao / wire / tests。
- 构造函数签名变化必须同步 `provider` 与 `wire_gen.go`。
- 回归命令见 [ENGINEERING.md](./ENGINEERING.md)。
