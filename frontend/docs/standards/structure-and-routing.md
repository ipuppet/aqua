# 目录结构与路由规范

以 FSD 为基础的稳定分层：

- `app`：应用装配（入口、路由、全局样式）
- `pages`：路由页面与页面私有组件
- `widgets`：页面级布局 / 壳层
- `features`：可复用业务能力
- `entities`：业务实体与实体 API
- `shared`：跨层通用能力（lib / ui / config / types / stores）

## 结构（简化）

```text
src/
  app/
  pages/
  widgets/
  features/
  entities/
  shared/
```

## 依赖方向

推荐：`app -> pages -> widgets -> features -> entities -> shared`

强制：

- `shared` 不能依赖 `entities` / `features` / `widgets` / `pages` / `app`
- `entities` 不能依赖 `features` / `widgets` / `pages` / `app`
- `features` 不能依赖 `widgets` / `pages` / `app`
- `widgets` 不能依赖 `pages` / `app`
- `pages` 禁止依赖 `app`

约束已写入各 app `eslint.config.ts` 的 `no-restricted-imports`。

## 路由

- 统一维护在 `src/app/router/index.ts`
- **hash** 模式
- 页面组件在 `src/pages/**`

## 命名与目录

- `widgets` → `widgets/layout/**`
- 页面私有组件 → `pages/<name>/components/**`
- 别名 `#/`
- 业务能力 → `features/*/model|ui|api`

## config

- `shared/config`：跨层通用配置（站点、安全等）
- `app`：仅启动 / 装配阶段使用的配置
