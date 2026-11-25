# 通用工程规范

编码、命名、日志、错误响应与回归检查。分层、IAM、Wire 见 [ARCHITECTURE.md](./ARCHITECTURE.md)。服务差异见 `docs/services/`。

## 1. 注释

解释 Why（意图、约束、取舍），不要把代码再翻译一遍。命名能说清时就少写注释。

```go
// 反例：如果 err 不为 nil，返回 err
// 正例：提前返回，避免超时后继续写入已关闭 channel 导致 goroutine 泄漏。
if err != nil {
	return err
}
```

函数体内注释只用于：非直觉算法、并发/锁、协议边界、魔法数字、workaround。

### 1.1 godoc

- 导出名称必须有 doc comment，与声明之间不得空行，以对象名开头的完整句子。
- Package 注释放 `doc.go` 或主文件顶部。
- Type 注释写语义与零值；并发安全要显式写。
- Func 注释写返回值或副作用；实现细节放函数体内。
- 布尔函数推荐 `reports whether ...`。
- 优先用 `[Name]`、`[pkg.Name]` 做符号跳转。

标记：`TODO(name):`、`BUG(name):`、`Deprecated:`（原因 + 替代）。注释过时比没有更糟，改代码必须改注释。

## 2. 日志

- 注入字段统一 `logger *slog.Logger`，禁止 `_logger` / `log`。
- service 保留该字段，推荐 `func (s *svc) log(ctx context.Context) *slog.Logger`，调用点 `s.log(ctx)`，需要时再 `.With(...)`。
- HTTP client / registry 等用内部 accessor 给默认 logger；外部用构造器或 Option 注入，不直接改字段。
- root logger 只在 bootstrap 初始化一次。业务层不要再搞零参全局 `logger()`。

## 3. 命名

- 同一业务域只允许一套词根。禁止 `administration` 与 `management` 混用。
- grpcserver 注入字段禁止 `service` / `repository` 这种泛化名，用 `userService`、`articleService`。构造器参数名与结构体字段一致。
- Go 文件 `snake_case.go`。测试文件与被测能力一一映射。禁止 `common.go` / `constants.go` 承载单一业务。
- YAML 配置键统一 `snake_case`（`app.yaml` / `admin.yaml` / `root.yaml` / `clash.yaml`），Go `mapstructure` tag 对齐。禁止同一配置树混用驼峰。
- 入口文件放类型、构造器、依赖字段；子文件按职责拆。可见性规则见 [ARCHITECTURE.md](./ARCHITECTURE.md) §3。

### 3.1 缩写

跨层标识符只允许行业缩写：`API`、`ASCII`、`CPU`、`CSS`、`DNS`、`EOF`、`HTML`、`HTTP`、`HTTPS`、`ID`（`IDs`）、`IP`、`JSON`、`QPS`、`RAM`、`RPC`、`SQL`、`SSH`、`TCP`、`TLS`、`TTL`、`UDP`、`UI`、`UID`、`UUID`、`URI`、`URL`、`UTF8`、`XML`、`YAML`、`DB`、`JWT`、`RBAC`、`JTI`、`DTO`。必须全大写语义，禁止 `Api` / `Id` / `Http` 这类写法。复合名用 `roleID`、`tokenID`。

protobuf `GetId`、jwt `StandardClaims.Id` 等第三方固定 API 不强改；业务侧用 `TokenID()` 等本地名挡住。

仅局部变量可用：`ctx`、`err`、`ok`、`i`、`j`、`n`；`req` / `resp` / `msg` 仅函数内临时变量，禁止做跨层字段。`cfg` / `repo` 可用于服务内部私有字段与构造器参数，对外接口优先 `config` / `repository`。`svc` 禁止用于结构体字段和构造器参数，测试局部变量可用。禁止 `pwd` / `usr` 进入跨层命名。

### 3.2 管理侧文件

有管理聚合入口时可用 `admin.go`；子实现按职责命名（`article.go`），不强制统一前缀。无聚合入口时保留领域名（如 `user_management.go`），不要机械改成裸 `management.go`。同一包内不要混用 `admin_*.go` / `management_*.go` / `*_management.go`。

## 4. gRPC 与转换

协议职责见 [ARCHITECTURE.md](./ARCHITECTURE.md) §2.2。文件建议：`*_server.go` + `converter.go` + `generated.go`。`generated.go` 不手改；规则变了重新生成。分页、默认值、缓存策略在 logic。

Facade 只在有明确收益时引入（迁移兼容、边界裁剪）。纯转发 facade 删掉。

## 5. 错误处理与 HTTP 响应

### 5.1 三条路径

- **RPC 映射的 HTTP**（blog / user / passport）：grpc-gateway + `pkg/transport/gateway` 写 JSON envelope。gRPC code → HTTP status 只在 `grpcx.HTTPStatusFromCode`。gateway 不得依赖 `pkg/transport/response`。
- **Gin 边沿**（`Authenticate` / authz、`ErrorHandler`、admin 文件/缓存 handler）：`response.Fail(c, err)`。`Fail` 直接写 JSON 并 `Abort`。`ErrorHandler` 只做 panic recovery，不消费 `c.Errors`。业务返回 `apperr`，handler 不挑 HTTP 状态码。`response` 不得 import gRPC / `grpcx`，也不解析 `status.FromError`。存储 miss 在 gRPC 边沿用 `apperr.FromDAO`。绑定失败构造 `apperr.BadRequest`（admin `httpbind`）。

### 5.2 禁止与允许

- 禁止在 blog / user / passport 新增与 grpc-gateway 同路径的手写 HTTP handler。
- 禁止在非白名单路径用 `c.String` / `c.JSON` / `c.DataFromReader` 当通用错误返回。
- 禁止在协议层用 `errors.As` / 类型断言手工决定 HTTP 状态码。
- 禁止把 Gin 绑定助手放回 `pkg/`；admin 只用 `services/admin/internal/httpbind`。
- admin 仅保留文件上传与清缓存的手写 Gin handler，使用 `response.OK*` / `Fail`。

## 6. 回归

提交前：

```bash
go test ./services/<service>/...
```

涉及全仓或共享包时用 `go test ./...`。协议变更执行 `buf generate` 并提交产物。Wire / goverter 接口变更必须同步生成物。
