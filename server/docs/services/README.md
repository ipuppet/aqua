# 服务专项规范

本目录只放服务特有规则。分层、IAM、Wire、错误响应、回归命令见 `docs/standards/`，这里不要复述。

- [passport.md](./passport.md)
- [user.md](./user.md)
- [blog.md](./blog.md)

## 允许例外

须在对应服务文档写明原因：

- 无 JWT 校验需求时，可不引入 JWT key provider（纯 webhook 等）。
- 无内部 gRPC 暴露时，可省略 `MountGRPC`。
- `admin` 仅保留文件上传与清缓存的手写 Gin handler（multipart 不适合 grpc-gateway）；管理 CRUD HTTP 挂在所属服务 gateway。
- logic 实现方默认不定义 Service 接口；仅多 facet 绑定或消费方本地 mock 时保留。
