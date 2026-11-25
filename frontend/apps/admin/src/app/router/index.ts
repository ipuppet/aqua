import { createRouter, createWebHashHistory } from "vue-router"
import type { RouteRecordRaw } from "vue-router"

const HomePage = () => import("#/pages/home/HomePage.vue")
const Me = {
  InfoPage: () => import("#/pages/me/InfoPage.vue"),
  SecurityPage: () => import("#/pages/me/SecurityPage.vue")
}
const User = {
  UserInfoPage: () => import("#/pages/user-management/UserInfoPage.vue"),
  RolePage: () => import("#/pages/user-management/RolePage.vue"),
  PermissionPage: () => import("#/pages/user-management/PermissionPage.vue")
}
const Blog = {
  TagsPage: () => import("#/pages/blog/TagsPage.vue"),
  ArticlesPage: () => import("#/pages/blog/ArticlesPage.vue"),
  AboutEditorPage: () => import("#/pages/blog/AboutEditorPage.vue"),
  ArticleEditorPage: () => import("#/pages/blog/ArticleEditorPage.vue")
}
const File = {
  ImagesPage: () => import("#/pages/file-management/ImagesPage.vue")
}

const routes: Array<RouteRecordRaw> = [
  { path: "/", name: "home", component: HomePage, meta: { title: "首页", affix: true } },
  { path: "/me/info", name: "me-info", component: Me.InfoPage, meta: { title: "个人信息" } },
  { path: "/me/security", name: "me-security", component: Me.SecurityPage, meta: { title: "安全设置" } },
  { path: "/user-manager/info/", name: "user-info", component: User.UserInfoPage, meta: { title: "用户管理" } },
  { path: "/user-manager/role/", name: "user-role", component: User.RolePage, meta: { title: "角色管理" } },
  {
    path: "/user-manager/permission/",
    name: "user-permission",
    component: User.PermissionPage,
    meta: { title: "权限管理" }
  },
  { path: "/blog/articles/", name: "blog-articles", component: Blog.ArticlesPage, meta: { title: "文章管理" } },
  {
    path: "/blog/editor/:aid/",
    name: "blog-article-editor",
    component: Blog.ArticleEditorPage,
    meta: { hideMenu: true, title: "编辑文章" }
  },
  {
    path: "/blog/editor/",
    name: "blog-article-editor-create",
    component: Blog.ArticleEditorPage,
    meta: { hideMenu: true, title: "创建文章" }
  },
  { path: "/blog/tags/", name: "blog-tags", component: Blog.TagsPage, meta: { title: "标签管理" } },
  { path: "/blog/about/", name: "blog-about", component: Blog.AboutEditorPage, meta: { title: "关于页面" } },
  { path: "/file/images/", name: "file-images", component: File.ImagesPage, meta: { title: "图片管理" } }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
