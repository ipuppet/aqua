# FSD 架构规范

定义 admin / blog / passport 共用的 Feature-Sliced Design 落地规则。

## 分层职责

- `app`：应用装配（`main.ts`、router、全局样式）
- `pages`：路由页面（编排 + 页面私有组件）
- `widgets`：页面级组合 UI（导航、壳层、布局）
- `features`：可复用业务能力（用户可触发的动作）
- `entities`：业务实体（模型 / UI）
- `shared`：跨业务通用能力（工具、基础 UI、配置、类型、公共 store）

## 依赖方向

```text
app -> pages -> widgets -> features -> entities -> shared
```

禁止反向依赖与跨层穿透。

## 目录约定

### app

- `app/router/*`：路由与守卫
- `app/styles/*`：全局样式与 CSS 变量

### pages

- 一个路由页面一个目录
- 页面私有组件放 `components/`，不要用含糊的 `ui/`

### widgets

- 统一 `widgets/layout/**`
- 不要把页面级布局塞进 `shared/ui`

### features

- 推荐：`features/<name>/{ui,model,api,lib}`
- 只放业务能力，不放 HTTP 客户端组装

### entities

- 推荐：`entities/<domain>/{model,ui}`
- 实体级 UI 可放 `entities/<domain>/ui/`
- 不要再写 `*Resource`；请求用 `#/shared/lib/api` 的 `client`

### shared

- `shared/lib/api.ts`：本端 `createAquaClient` 装配与 `buildApiUrl` 包装
- `shared/lib`：通用函数
- `shared/config`：跨层稳定配置；Token 键名从 `@aqua/http` 再导出即可
- `shared/ui`：纯展示 / 通用组件
- `shared/stores`：跨业务全局 store

## 命名

- 页面组件：`*Page.vue`（如 `ArticleDetailPage.vue`）
- 跨目录导入统一 `#/`，禁止相对路径跨层引用

## 配置放置

- `shared/config`：可被多层复用的稳定配置
- `app`：仅装配阶段使用、不应被下层复用的配置

## 工程化

- ESLint `no-restricted-imports` 固化层级（别名 `#/`）
- HTTP 用 `@aqua/http` 的
  `createAquaClient` + 生成类型；应用内不要再复制 URL 构建、token 刷新或 cookie 逻辑
- 禁止仅为兼容保留的转发层 / 旧目录入口
- 合入前：`npm run check -w <app>` 通过
