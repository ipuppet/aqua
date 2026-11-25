import { createRouter, createWebHashHistory } from "vue-router"
import type { RouteRecordRaw } from "vue-router"

const HomePage = () => import("#/pages/home/HomePage.vue")
const ArticleDetailPage = () => import("#/pages/blog/ArticleDetailPage.vue")
const SearchPage = () => import("#/pages/blog/SearchPage.vue")
const TagPage = () => import("#/pages/blog/TagPage.vue")
const ArchivePage = () => import("#/pages/blog/ArchivePage.vue")
const AboutPage = () => import("#/pages/blog/AboutPage.vue")

const routes: RouteRecordRaw[] = [
  { path: "/", name: "HomePage", component: HomePage },
  { path: "/detail/:aid/", name: "ArticleDetail", component: ArticleDetailPage },
  { path: "/search/:count/:kw/", name: "SearchPage", component: SearchPage },
  { path: "/tag/:tag/", name: "TagPage", component: TagPage },
  { path: "/archive/", name: "ArchivePage", component: ArchivePage },
  { path: "/about/", name: "AboutPage", component: AboutPage }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
