# 状态管理规范

使用 Pinia，按 FSD 分成两类 store。

## 分类

### 跨业务（`shared/stores/`）

多个业务域共用的全局状态。各端现有文件不同，以该 app 的 `src/shared/stores/`
为准（例如 admin 有 tabs / menu；blog 有 darkMode；passport 可能很少或没有）。

### 业务域（`features/*/model/` 或 `entities/*/model/`）

仅某一域使用的状态，放在对应切片内。

## 原则

1. **最小化全局状态**  
   页面内优先 `ref` / `reactive`；真正跨页共享再进 Pinia。

2. **逻辑分离**  
   Store 以存状态为主；异步请求优先放在视图或 Composable，经 `#/shared/lib/api`
   更新store。

3. **类型安全**  
   state 必须有明确 TypeScript 类型。

4. **持久化**  
   偏好类用 `localStorage`；敏感或短期数据优先 `sessionStorage`。

## 导入规则

- 上层可以依赖 `shared/stores/`
- **禁止跨 feature 导入另一个 feature 的 store**；若多域需要，提升到
  `shared/stores/`
