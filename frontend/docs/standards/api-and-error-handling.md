# API 与错误处理规范

HTTP 基础设施在 workspace 包
`@aqua/http`（`frontend/packages/http`）。类型与路径来自 buf 的
`grpc-ecosystem/openapiv3` 产物（`npm run gen`），请求走 `openapi-fetch`，不要再手写
Resource / axios。

## 职责划分

**共享包 `@aqua/http`：**

- `createAquaClient`：openapi-fetch 工厂（`authMode` / `cookieDomain` /
  `onAuthFailure`）
- `buildApiUrl`：按服务名拼开发端口或网关 URL（`apiHost` 由各 app 注入）
- Cookie：`getCookie` / `setCookie` / `configureCookies`
- Token：`ACCESS_TOKEN_KEY` / `REFRESH_TOKEN_KEY` / `clearAuthTokens`
- 生成类型：`paths` / `components` / `operations`
- `getErrorMessage`

**各应用 `src/shared/lib/api.ts`：**

- 用本端 `VITE_API` 包一层 `buildApiUrl`
- 调用 `createAquaClient` 组装本端策略（UI 提示、跳转登录）
- 导出本端 `client` 单例

## 基本原则

1. 请求一律经 `#/shared/lib/api` 的 `client`，使用生成路径：

   ```ts
   import { client, buildApiUrl } from "#/shared/lib/api"

   const { data, error } = await client.GET("/api/articles", {
     params: { query: { page: 1, pageSize: 10 } }
   })
   if (error) {
     ElMessage.error(error.message ?? "获取文章列表失败")
     return
   }
   ```

2. 跨服务调用显式传
   `baseUrl: buildApiUrl("user" | "passport" | …)`，不要再 new 第二个 client。
3. 错误对象用 `error.message`；需要兜底文案时再用 `getErrorMessage`。
4. 响应字段以 OpenAPI 生成为准（camelCase），不要再做 Resource 里的 snake_case 映射。

## 各端 createAquaClient 策略

| 应用     | `authMode`  | 说明                                                              |
| -------- | ----------- | ----------------------------------------------------------------- |
| admin    | `"refresh"` | 401 单飞刷新并重试一次；失败走 `onAuthFailure`（提示 + 跳转登录） |
| blog     | `"none"`    | 有 token 则带上；不静默刷新、不跳转登录                           |
| passport | `"none"`    | 已在登录站；鉴权失败只返回 error，**不** `redirectToLogin`        |

共同点：

- 有 access token 时都会带 `Authorization`
- `cookieDomain`：生产多为 `.example.com`；本地为空（host-only）
- 包内不依赖 Element Plus，不直接做 UI / 跳转
- `/api/token` 自身失败不会递归刷新

Passport 登录成功后的 `from` 回跳须走白名单（仅 `localhost` / `example.com`
及其子域，且 http(s)）。
