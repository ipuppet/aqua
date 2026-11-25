# aqua-frontend

前端 npm workspace：`apps/*` + `packages/*`。Go API 在仓库根目录的 `server/`。

| 路径 | 说明 |
| --- | --- |
| `apps/admin` | 管理后台（`@aqua/admin`） |
| `apps/blog` | 博客前台（`@aqua/blog`） |
| `apps/passport` | 登录 / 注册（`@aqua/passport`） |
| `packages/http` | `@aqua/http` |
| `packages/eslint-config` | `@aqua/eslint-config` |
| `packages/tsconfig` | `@aqua/tsconfig` |
| `packages/vite-config` | `@aqua/vite-config` |

## 开发

```bash
cd frontend
npm install

npm run dev:admin
npm run dev:blog
npm run dev:passport

npm run check
npm run build:admin
npm run build:blog
npm run build:passport
```

环境变量：复制 `.env.development.example` / `.env.production.example` 为
`.env.development` / `.env.production`（Vite `envDir` 指向 workspace 根）。

## 文档

- [CONTRIBUTING.md](./CONTRIBUTING.md) — 开发约定总览
- [docs/](./docs/README.md) — FSD / API / 状态 / 表单等规范
