# Aqua frontend

基于 Vue 3、TypeScript 和 Vite 的 Aqua 前端 npm
workspace，按 FSD 分层，包含三个应用和一组共享包：

| 路径                     | 包名                  | 说明                                   |
| ------------------------ | --------------------- | -------------------------------------- |
| `apps/admin`             | `@aqua/admin`         | 管理后台：文章、用户与权限、系统配置   |
| `apps/blog`              | `@aqua/blog`          | 博客前台：文章、搜索、评论、暗色模式   |
| `apps/passport`          | `@aqua/passport`      | 登录 / 注册 / 邮箱验证 / Turnstile     |
| `packages/http`          | `@aqua/http`          | HTTP client、认证刷新、各服务 URL 组装 |
| `packages/eslint-config` | `@aqua/eslint-config` | 共享 ESLint 配置                       |
| `packages/tsconfig`      | `@aqua/tsconfig`      | 共享 TypeScript 配置                   |
| `packages/vite-config`   | `@aqua/vite-config`   | 共享 Vite 配置                         |

对应 Go API 在 [`../server/`](../server/README.md)。

## 环境要求

- Node.js 22（与部署镜像一致）
- npm 10+

## 安装

```bash
cd frontend
npm ci
```

## 环境变量

Vite 的 `envDir` 指向 `frontend/` 根目录：

```bash
cp .env.example .env.development
cp .env.example .env.production
```

| 变量                      | 说明                                                            |
| ------------------------- | --------------------------------------------------------------- |
| `VITE_API`                | API 主机名。开发环境只填主机名，客户端会按服务拼接端口          |
| `VITE_AUTH_COOKIE_DOMAIN` | 需要跨子域共享登录 Cookie 时填写，例如 `.example.com`；本地留空 |

开发模式下，`VITE_API=localhost` 会按 `packages/http/src/url.ts` 的
`DEV_API_PORTS` 访问各服务：

| 服务       | 端口   |
| ---------- | ------ |
| `passport` | `8081` |
| `admin`    | `8082` |
| `user`     | `8083` |
| `blog`     | `8084` |

也就是说，本地前端开发需要把对应服务分别启动在以上端口。连接单一反向代理网关时，把
`VITE_API` 设为网关域名（不要填
`localhost`，否则客户端会进入多端口开发模式），生产环境示例见
[`../deploy/README.md`](../deploy/README.md)。

## 开发

```bash
npm run dev:admin
npm run dev:blog
npm run dev:passport
```

默认使用 Vite 开发服务器。三个应用的路由默认是 hash 模式，路由集中在各应用的
`src/app/router/index.ts`。

## 检查与构建

```bash
npm run check                 # 三个 app 的 type/lint/format + @aqua/http 测试
npm run test                  # 只运行 @aqua/http 测试

npm run check -w @aqua/admin  # 只检查一个应用
npm run build:admin
npm run build:blog
npm run build:passport
npm run gen                   # 重新生成 @aqua/http 的接口定义
```

构建产物在 `apps/<app>/dist/`。部署时需要为三个 SPA 配置前端路由回退；
`../deploy/` 中的示例已经处理了这一点。

## 应用差异

| 应用       | `authMode` | 说明                                         |
| ---------- | ---------- | -------------------------------------------- |
| `admin`    | `refresh`  | 使用 refresh token；刷新失败后提示并跳转登录 |
| `blog`     | `none`     | 公开站点；需要登录的操作由页面自行引导       |
| `passport` | `none`     | 本站是登录页；鉴权失败只抛错，不自我跳转     |

业务请求统一通过 `#/shared/lib/api` 的 `client`
发起，不要在页面里手写 Resource 或裸 `fetch`。跨应用状态放
`shared/stores/`，业务域状态放 `features/*/model/` 或 `entities/*/model/`。

## 代码组织

`src/` 按 FSD 分为六层，依赖方向自上而下：

```text
app -> pages -> widgets -> features -> entities -> shared
```

- `app/`：入口、路由、全局样式
- `pages/`：路由页面，薄壳编排
- `widgets/`：页面级布局壳
- `features/`：可复用业务能力，feature 之间禁止互引
- `entities/`：业务实体、实体 UI 与模型
- `shared/`：跨业务工具、配置、UI 和 stores

路径别名统一使用 `#/`；层级约束写在各自的 `eslint.config.ts` 中。完整约定见
[`CONTRIBUTING.md`](./CONTRIBUTING.md) 和 [`docs/`](./docs/README.md)。

## 部署

- 三个应用分别构建为静态站点，由 Nginx 或对象存储托管。
- API 与前端分属不同子域；登录 Cookie 的域由 `VITE_AUTH_COOKIE_DOMAIN` 控制。
- 仓库根目录的 [`deploy/`](../deploy/README.md)
  提供 Compose 部署示例，包含三端构建、Nginx 路由和 API 反向代理。

## 文档

- [开发约定](./CONTRIBUTING.md)
- [前端规范](./docs/README.md)
- [HTTP 与错误处理](./docs/standards/api-and-error-handling.md)
- [FSD 架构](./docs/standards/fsd-architecture.md)
- [状态管理](./docs/standards/state-management.md)
- [表单与校验](./docs/standards/validation-and-forms.md)
