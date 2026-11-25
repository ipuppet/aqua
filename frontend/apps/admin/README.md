# Admin Dashboard

管理后台：博客内容、用户权限、系统配置。

包名：`@aqua/admin`（目录 `apps/admin`）。

## 功能

- 文章 / 标签 / 关于页（Markdown）
- 用户、角色、权限
- 数据统计（ECharts）、个人信息

## 本端差异

- HTTP：`authMode: "refresh"`，刷新失败提示并跳转登录（见
  [api-and-error-handling](../../docs/standards/api-and-error-handling.md)）
- 全局 store：`user` / `menu` / `nav` / `tabs` / `icons`

## 开发

在 `frontend/` 根目录：

```bash
npm install
npm run dev:admin
npm run check -w @aqua/admin
npm run build:admin
```

环境变量：`frontend/.env.development` / `frontend/.env.production`。

规范：[CONTRIBUTING](../../CONTRIBUTING.md) · [docs](../../docs/README.md)
