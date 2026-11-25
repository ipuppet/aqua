# Aqua server

Go 业务后端（module `aqua`）。Compose 服务名仍为 `web_*`（运行时命名，与目录无关）。

## Usage

```shell
make build
make run target=mono
make check
```

Copy `config/app.example.yaml` to `config/app.yaml` and replace the placeholders
before running the full stack. See `config/README.md`.

