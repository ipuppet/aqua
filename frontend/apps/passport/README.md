# Passport Authentication

登录 / 注册站：密码与验证码登录、邮箱注册、Turnstile、Cookie 发 token。

包名：`@aqua/passport`（目录 `apps/passport`）。

## 功能

- 登录（密码 / 手机验证码）
- 注册（邮箱验证）
- Cloudflare Turnstile
- Cookie 认证；登录成功后按白名单 `from` 回跳

## 本端差异

- HTTP：`authMode: "none"`（本站即登录页，鉴权失败只抛错、不自我跳转）
- 表单校验：`src/shared/lib/validation/`（见
  [validation-and-forms](../../docs/standards/validation-and-forms.md)）

## 开发

在 `frontend/` 根目录：

```bash
npm install
npm run dev:passport
npm run check -w @aqua/passport
npm run build:passport
```

环境变量：`frontend/.env.development` / `frontend/.env.production`。

规范：[CONTRIBUTING](../../CONTRIBUTING.md) · [docs](../../docs/README.md)
