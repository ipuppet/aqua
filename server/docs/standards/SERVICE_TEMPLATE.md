# 微服务模板与 main 规范

`blog` / `user` / `passport` 已落地的实现模板。分层与目录见 [ARCHITECTURE.md](./ARCHITECTURE.md)，回归见 [ENGINEERING.md](./ENGINEERING.md)。

## 0. 脚手架

```bash
make new-service name=<service>
make build target=<service>
make run target=<service>
# curl -s 'http://127.0.0.1:8080/api/<service>/hello?name=aqua'
```

生成 `proto/aqua/<service>/v1` 与 `services/<service>`，并默认登记到 `services/mono/cmd`。Makefile 通过 `services/*/cmd` 发现可构建服务。不改 Docker / Compose。脚手架独立在 `scaffold/`（生成器 `new_service.sh`，模板 `templates/service/`）。

## 1. 与架构的对应

标准子模块见架构 §1。这三服务都是 grpc + grpc-gateway，`gateway.go` 不写业务。除非明确是协议透传，不要新增 `handler/`。

grpcserver 按能力命名 `*_server.go`：

- `admin_server.go`：Admin gRPC。
- `<domain>_server.go`：对外/业务 gRPC（`user_server.go`、`article_server.go`）。
- `converter.go` + `generated.go`：goverter 纯函数转换。

不要用 `server.go` / `operations.go` 当长期能力文件。

## 2. `Assemble`

固定三步：构造 gRPC server → 注册 `Register*HandlerServer` → `MountHTTP` 里挂 `registerGatewayRoutes`。

约束：

- `main.go` 只组装，不写业务，不手工构造 JWT key / filestore / db / redis。
- 不用服务私有 `xxxContainer`；Wire 直接产出 `*bootstrap.Module`。
- `buildModule(rt *bootstrap.Runtime)`。
- gRPC 变量按能力：`adminServer`、`userServer`。

```go
func Assemble(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	return buildModule(rt)
}

func provideModule(
	...,
	publicKey *rsa.PublicKey,
) (*bootstrap.Module, error) {
	adminServer := grpcserver.NewAdminServer(...)
	xxxServer := grpcserver.NewXxxServer(...)

	gatewayMux := gateway.NewServeMux()
	if err := v1.RegisterXxxServiceHandlerServer(context.Background(), gatewayMux, xxxServer); err != nil {
		return nil, err
	}

	gatewayHandler := gateway.GinHandler(gatewayMux)

	return &bootstrap.Module{
		MountHTTP: func(engine *gin.Engine) {
			registerGatewayRoutes(engine, publicKey, gatewayHandler)
		},
		MountGRPC: func(server *grpc.Server) {
			v1.RegisterAdminXxxServiceServer(server, adminServer)
			v1.RegisterXxxServiceServer(server, xxxServer)
		},
	}, nil
}
```

`gateway.go` 只挂路径、分组和认证/授权中间件。禁止写业务、DAO 或复杂校验。

## 3. logic / dao 书写顺序

### logic

核心文件（`user.go`、`admin.go`）按此顺序：

1. `package` / `import`
2. `type XxxService struct`
3. `func NewXxxService(...) *XxxService`
4. 方法（可按实体拆到同包其他文件）
5. 私有辅助（`log()`、清缓存等）

默认不在实现方定义 Service 接口。需要可替换依赖时，在消费方定义小接口。

同一实体方法按 `Create` → `Get` → `List/Search` → `Update` → `Delete`。

类型名：`UserService`、`ManagementService`。禁止 `*Impl`、禁止再套同包大接口。

保留接口的例外：同一具体类型要绑多个消费侧 facet；消费方测试要本地 mock（接口定义在 handler/grpcserver）；多方法只读端口且测试大量替换（如某服务的 `ConfigProvider`）。

### dao

- Repository 接口集中在 `dao/ports.go`。
- `ports.go` 底部：包级私有 `repository`、`NewRepository(...) *repository`，以及 `var _ XxxRepository = (*repository)(nil)`。
- 方法顺序与 logic 相同。
- 接口名 `XxxRepository`；实现固定不导出 `repository`。
