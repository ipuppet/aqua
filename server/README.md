# Aqua server

Aqua 的 Go 后端，module 名为 `aqua`。默认提供四个业务服务和一个单体入口：

| 服务       | 说明                                           |
| ---------- | ---------------------------------------------- |
| `admin`    | 管理后台 API                                   |
| `blog`     | 博客前台 API                                   |
| `passport` | 登录、注册、会话与 JWT 签发                    |
| `user`     | 用户资料与 RBAC                                |
| `mono`     | 在同一个进程内挂载以上全部服务，供生产部署使用 |

服务之间通过 gRPC 通信，对外 HTTP 由 Gin + grpc-gateway 提供；服务发现使用 etcd，数据层使用 PostgreSQL、Redis 和 Meilisearch。CORS 和 TLS 由反向代理负责，Go 服务本身不挂 CORS 中间件。

## 环境要求

- Go 1.27.3（以 `go.mod` 为准）
- `buf`（生成 protobuf）
- `goimports`、`golangci-lint`（格式化和静态检查）
- PostgreSQL、Redis、etcd、Meilisearch（运行完整链路时）

`make tools` 会同步 Go module 工具、安装 `gopls` 并检查 `buf` 是否在 PATH 中。

## 快速开始

```bash
cd server
make tools

cp config/app.example.yaml config/app.yaml
cp config/admin.example.yaml config/admin.yaml
```

按 [`config/README.md`](./config/README.md) 生成 JWT 密钥并填写数据库、Redis、Meilisearch、邮件等配置。公开仓库只提供示例配置，不包含任何可用的生产密钥。

启动单体服务：

```bash
make run target=mono
```

默认监听 `:8080`（HTTP）和 `:18080`（gRPC），配置文件目录为 `./config`。健康检查：

```bash
curl http://127.0.0.1:8080/ping
```

也可以单独运行某个服务：

```bash
make run target=passport
```

## 构建与检查

```bash
make build                  # 构建全部服务
make build target=passport  # 只构建一个服务
make gen                    # buf generate + go generate
make fmt                    # goimports + buf format
make lint                   # golangci-lint + buf lint + 自定义检查
make test                   # 单元测试
make check                  # gen + fmt + lint + test
make check target=blog      # 只检查 blog 及其 proto
```

构建产物在 `bin/<service>/server`。可用 `GOOS`、`GOARCH`、`CGO_ENABLED`、`LDFLAGS` 控制构建，例如：

```bash
make build target=mono GOOS=linux GOARCH=arm64
```

新增服务：

```bash
make new-service name=comments
```

脚手架会生成 `services/comments` 和对应的 proto 目录；具体骨架见 [`docs/standards/SERVICE_TEMPLATE.md`](./docs/standards/SERVICE_TEMPLATE.md)。

## 目录结构

```text
server/
├── services/                 # 业务服务
│   ├── admin/
│   ├── blog/
│   ├── passport/
│   ├── user/
│   └── mono/                 # 单体入口
├── proto/                    # protobuf 契约
├── gen/                      # 生成代码
├── pkg/                      # 跨服务基础设施与共享包
├── scaffold/                 # 新服务脚手架
├── scripts/                  # 构建辅助与自定义检查
├── config/                   # 示例配置
└── docs/                     # 架构、工程规范、服务专项
```

每个标准服务包含：

```text
services/<service>/
├── cmd/main.go               # 进程入口
├── domain/                   # 领域模型
├── logic/                    # 业务规则与用例编排
├── dao/                      # 数据库 / 缓存访问
├── grpcserver/               # gRPC 接入
├── gateway.go                # HTTP 路由挂载
├── main.go                   # ServiceName / Assemble
└── wire.go / wire_gen.go     # 依赖注入
```

依赖方向为 `gateway -> grpcserver -> logic -> dao`；`grpcserver` 不做业务决策，`dao` 只负责持久化。详细约束见 [`docs/standards/ARCHITECTURE.md`](./docs/standards/ARCHITECTURE.md) 和 [`docs/standards/ENGINEERING.md`](./docs/standards/ENGINEERING.md)。

## 配置

服务只从 YAML 读取业务配置，不使用通用环境变量覆盖 YAML 字段。常用入口：

| 配置                              | 说明                                                        |
| --------------------------------- | ----------------------------------------------------------- |
| `-config-dir` / `AQUA_CONFIG_DIR` | 配置目录，默认 `./config`                                   |
| `app.yaml`                        | 数据库、Redis、etcd、Meilisearch、JWT、邮件、存储等共享配置 |
| `admin.yaml`                      | 管理后台的邮箱、Turnstile 等配置                            |
| `AQUA_STORE_DIR`                  | 文件存储目录                                                |
| `AQUA_HTTP_ADDR`                  | HTTP 监听地址                                               |
| `AQUA_GRPC_ADDR`                  | gRPC 监听地址                                               |

配置字段说明和本地密钥生成方式见 [`config/README.md`](./config/README.md)；容器化部署示例见 [`../deploy/README.md`](../deploy/README.md)。

## 数据库

建表脚本在 [`../postgres/structure/web-structure.sql`](../postgres/structure/web-structure.sql)。本地需要先创建 `web` 数据库和角色，再导入该文件；Docker 部署示例会自动完成这一步。数据库说明见 [`../postgres/README.md`](../postgres/README.md)。

## 文档

- [服务架构与规范](./docs/README.md)
- [Passport](./docs/services/passport.md)
- [User](./docs/services/user.md)
- [Blog](./docs/services/blog.md)
- [部署示例](../deploy/README.md)
- [前端](../frontend/README.md)
