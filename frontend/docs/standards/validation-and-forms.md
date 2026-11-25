# 表单与校验规范

使用 Vee-Validate 时，校验逻辑集中复用，页面只绑规则名与文案。

## 目录组织

有表单重校验的应用（当前主要是 passport）放在：

- `src/shared/lib/validation/validation.ts`：基础正则与通用函数（邮箱、用户名、密码等）
- `src/shared/lib/validation/vee-validate-rules.ts`：规则注册与封装

其他应用若引入同类能力，沿用同一路径，不要在页面里散落实现。

## 原则

1. **复用现有函数**  
   如 `validatePassword`、`validateUsername`。新增规则先写进 `validation.ts`，再注册到
   `vee-validate-rules.ts`。

2. **UI 与规则解耦**  
   规则返回可读文案或错误 key；页面决定 `ElMessage` 或 inline 展示。

3. **Turnstile 等人机验证**  
   逻辑集中在 `features/auth/`（或等价 feature）；表单只关心 token 是否存在。

## 示例

```ts
import { Form, Field } from "vee-validate"
import {
  validatePassword,
  validateUsername
} from "#/shared/lib/validation/validation"

async function handleRegister() {
  const usernameError = validateUsername(username.value)
  if (usernameError) {
    ElMessage({ type: "error", message: usernameError })
    return
  }

  const passwordError = validatePassword(password.value, confirmPassword.value)
  if (passwordError) {
    ElMessage({ type: "error", message: passwordError })
    return
  }
}
```
