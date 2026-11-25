# 服务脚手架

```bash
make new-service name=<service>
# 等价于：
./scaffold/new_service.sh <service> [--no-mono] [--no-gen] [--force]
```

## 布局

```
scaffold/
  new_service.sh          # 生成器：拷贝模板、登记 mono、可选 codegen
  templates/
    service/              # grpc + grpc-gateway 最小服务（SayHello）
```

新增一种骨架时，在 `templates/<kind>/` 放树，再让 `new_service.sh` 按 kind 选择模板根目录。模板占位符：`__SERVICE__`（小写）与 `__Service__`（Pascal）。

`proto/` 下的文件会写到 `proto/aqua/<name>/v1/<name>.proto`，其余路径相对仓库根原样展开。
