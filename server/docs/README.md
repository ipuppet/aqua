# Aqua 文档总览

- 通用标准：`docs/standards/`
- 服务专项：`docs/services/`

一条规则只写一处。冲突时，服务专项覆盖该服务相关的通用条目。

## 1. 通用标准

| 文档 | 管什么 |
| --- | --- |
| [架构](./standards/ARCHITECTURE.md) | 分层、依赖方向、IAM、启动、Wire |
| [工程规范](./standards/ENGINEERING.md) | 注释、日志、命名、错误响应、回归检查 |
| [服务模板](./standards/SERVICE_TEMPLATE.md) | `Assemble` 骨架、logic/dao 书写顺序 |

## 2. 服务专项

只写本服务特有的职责、文件布局、路由和约束。跨服务强制条款不要再抄一遍。

- [索引与例外](./services/README.md)
- [Passport](./services/passport.md)
- [User](./services/user.md)
- [Blog](./services/blog.md)

## 3. 怎么用

- 新建服务：`make new-service name=<service>`，脚手架在 `scaffold/`，约定见 [服务模板](./standards/SERVICE_TEMPLATE.md)。
- 改代码：先读架构与工程规范，再读对应服务专项。
- 提交前检查：见 [工程规范 · 回归](./standards/ENGINEERING.md)。
