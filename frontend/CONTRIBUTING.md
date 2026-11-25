# 开发约定总览

`frontend/` 下三个应用（`apps/admin` / `apps/blog` /
`apps/passport`）共用同一套 Vue 3 + Vite + TypeScript +
FSD 约定。目标是代码易读、易维护、易扩展，并通过层级隔离业务域。

详细规范见 [`docs/`](./docs/README.md)。本文件只列提交前必守的要点。

## FSD 分层

`src/` 六层：

| 层          | 职责                                                                |
| ----------- | ------------------------------------------------------------------- |
| `app/`      | 装配：`main.ts`、`App.vue`、路由、全局样式                          |
| `pages/`    | 路由页面（薄壳编排）；私有组件放 `pages/<name>/components/`         |
| `widgets/`  | 页面级布局壳，统一放 `widgets/layout/`                              |
| `features/` | 可复用业务能力；推荐 `{ui,model,api,lib}`；**feature 之间禁止互引** |
| `entities/` | 业务实体与实体 UI / 模型                                            |
| `shared/`   | 跨业务：`lib`（含 API client 装配）、`config`、`ui`、`stores`       |

依赖方向（自上而下）：

```text
app -> pages -> widgets -> features -> entities -> shared
```

- 路径别名统一 `#/`；HTTP 基础设施用 `@aqua/http`，不要在应用内再复制
- 层级约束已写入各 app 的 `eslint.config.ts`（`no-restricted-imports`）

## 路由

- 默认 **hash** 模式
- 路由集中在 `src/app/router/index.ts`

## API 与状态

- 页面经 `#/shared/lib/api` 的 `client` 调用生成路径，不要手写 Resource
  / 裸 fetch
- 跨业务状态 → `shared/stores/`；域内状态 → `features/*/model/` 或
  `entities/*/model/`

## 检查

在 `frontend/` 下：

```bash
npm run check -w @aqua/<app>   # type + lint + format
npm run check                  # 全部 app + @aqua/http 测试
```
