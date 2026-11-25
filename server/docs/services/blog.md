# Blog 服务专项规范

## 1. 职责

- 读：文章详情、列表、标签、归档、about。
- 写：文章、标签、about、Meili 索引同步。
- 访问：访问日志与按天统计。统计与文章读模型解耦，查询接口不要和统计接口绑死。

## 2. 文件

logic：

- 读入口 `blog.go`（`ArticleService`）：查询、分页、上下文拼装。
- 写入口 `admin.go`（`ManagementService`）：写操作、缓存失效、索引联动。
- 写子文件：`article.go`、`tag.go`、`about.go`、`index.go`。

dao：

- 读：`blog.go`、`comment.go`。
- 写：`article_write.go`、`tag_write.go`、`transaction.go`。
- 访问：`access.go`（独立 `AccessRepository`）。
- 端口在 `ports.go`。

`Admin*` 只出现在对外协议（gRPC `AdminServer` / proto service）。包内实现用能力名。

## 3. 缓存与索引

写操作后必须清对应缓存键。文章变更要更新索引，删除要删索引。

## 4. HTTP

读写共用 `/api/articles`、`/api/tags`、`/api/about`；GET 公开，写入需 `blog.manager`。

- 编辑读取：`GET /api/articles/{aid}/edit`（`AdminBlogService.GetArticle`），不要走公开详情（公开详情会记访问量）。
- `AdminBlogService` 只保留写入、编辑读取、访问统计、Meili 同步。列表 / 搜索 / 标签 / about 读取走 `ArticleService`。
- 评论：`GET /api/comments/articles/{aid}`、`GET /api/comments/{id}/children`；`POST /api/comments/articles/{aid}` 需登录。

挂载：`RegisterArticleServiceHandlerServer`、`RegisterCommentServiceHandlerServer`、`RegisterAdminBlogServiceHandlerServer`。

生成物：`gen/go/aqua/blog/v1/blog.pb.gw.go`、`gen/openapi/aqua/blog/v1/*.openapi.json`。
