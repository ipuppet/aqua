<template>
  <article-list
    :is-loaded="isLoaded"
    :page="page"
    @current-change="currentChange"
    @prev-click="prevClick"
    @next-click="nextClick"
  />
</template>
<script>
import ArticleList from "#/entities/article/ui/ArticleList.vue"
import { sessionData } from "#/shared/lib/storage"
import { client } from "#/shared/lib/api"
import { siteConfig } from "#/shared/config"
export default {
  name: "HomePage",
  components: { ArticleList },
  data: () => ({
    articles: {},
    page: { articles: [], pageSize: 5, total: 1, current: 1 },
    isLoaded: false,
    currentRequestPage: null
  }),
  methods: {
    async getArticleList() {
      if (this.currentRequestPage === this.page.current) return
      window.scrollTo(0, 0)
      if (this.page.current in this.articles) {
        this.page.articles = this.articles[this.page.current]
        this.isLoaded = true
        return
      }
      const homePage = sessionData.get("homePage")
      if (homePage?.articles?.[this.page.current]) {
        this.articles[this.page.current] = homePage.articles[this.page.current]
        this.page = Object.assign(this.page, { total: homePage.page.total, articles: this.articles[this.page.current] })
        this.isLoaded = true
        return
      }
      this.isLoaded = false
      this.currentRequestPage = this.page.current
      const requestPage = this.page.current

      const { data, error } = await client.GET("/api/articles", {
        params: { query: { page: this.page.current, pageSize: this.page.pageSize } }
      })

      if (!error && data) {
        if (requestPage === this.page.current) {
          this.articles[this.page.current] = data.articles
          this.page = Object.assign(this.page, { total: data.total, articles: this.articles[this.page.current] })
          const updatedHomePage = homePage ?? { articles: {}, page: {} }
          updatedHomePage.articles = this.articles
          updatedHomePage.page.total = this.page.total
          sessionData.set("homePage", updatedHomePage)
        }
        this.isLoaded = true
        this.currentRequestPage = null
      } else {
        if (requestPage === this.page.current) this.isLoaded = true
        this.currentRequestPage = null
      }
    },
    currentChange(page) {
      this.page.current = page
      this.getArticleList()
    },
    prevClick() {
      if (this.currentRequestPage !== null) return
      if (this.page.current > 1) {
        this.page.current -= 1
        this.getArticleList()
      }
    },
    nextClick() {
      if (this.currentRequestPage !== null) return
      if (this.page.current < Math.ceil(this.page.total / this.page.pageSize)) {
        this.page.current += 1
        this.getArticleList()
      }
    }
  },
  mounted() {
    this.getArticleList()
  },
  created() {
    document.title = siteConfig.title
    if (!sessionData.get("homePage")) sessionData.set("homePage", { articles: {}, page: { total: 0 } })
  }
}
</script>
