# Blog Frontend

博客前台：文章浏览、搜索、评论、暗色模式。

包名：`@aqua/blog`（目录 `apps/blog`）。

## 功能

- 列表 / 详情 / 搜索 / 标签 / 归档 / 关于
- 评论
- 暗色模式

## 本端差异

- HTTP：`authMode: "none"`（公开站，需登录时页面自行引导）
- 全局 store：`darkMode`

## 开发

在 `frontend/` 根目录：

```bash
npm install
npm run dev:blog
npm run check -w @aqua/blog
npm run build:blog
```

环境变量：`frontend/.env.development` / `frontend/.env.production`。

规范：[CONTRIBUTING](../../CONTRIBUTING.md) · [docs](../../docs/README.md)
