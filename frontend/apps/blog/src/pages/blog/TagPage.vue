<template>
  <div>
    <div :class="{ dark: darkModeStore.isDarkMode }" v-show="isLoaded" class="header">{{ title }}</div>
    <article-list :is-loaded="isLoaded" :page="page" />
  </div>
</template>
<script>
import ArticleList from "#/entities/article/ui/ArticleList.vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { sessionData } from "#/shared/lib/storage"
import { client } from "#/shared/lib/api"
export default {
  components: { ArticleList },
  setup() {
    const darkModeStore = useDarkModeStore()
    return { darkModeStore }
  },
  data: () => ({ title: "", page: { articles: [], pageSize: 5, total: 1, current: 1 }, tag: "", isLoaded: false }),
  methods: {
    async getArticleListByTag() {
      this.tag = this.$route.params.tag
      this.title = `Articles of '${this.tag}'`
      document.title = this.title
      const tags = sessionData.get("tags") ?? []
      const tagObj = tags.find(t => t.name === this.tag)
      if (!tagObj) {
        this.isLoaded = true
        return
      }
      const tagPage = sessionData.get("tagPage") ?? {}
      if (tagPage[this.tag]) {
        this.page.articles = tagPage[this.tag].articles
        this.page.total = tagPage[this.tag].total
        this.page.pageSize = tagPage[this.tag].total
        this.isLoaded = true
      } else {
        const { data, error } = await client.GET("/api/tags/{id}/articles", {
          params: { path: { id: tagObj.id } }
        })
        if (!error && data) {
          this.page.articles = data.articles ?? []
          this.page.total = this.page.articles.length
          this.page.pageSize = this.page.articles.length
          tagPage[this.tag] = this.page
          sessionData.set("tagPage", tagPage)
          this.isLoaded = true
        } else {
          this.isLoaded = true
        }
      }
    }
  },
  created() {
    this.getArticleListByTag()
  },
  watch: {
    $route(to, from) {
      if (to !== from) {
        this.isLoaded = false
        this.getArticleListByTag()
      }
    }
  }
}
</script>
<style scoped>
.header {
  font-size: 36px;
  margin-bottom: 30px;
  color: var(--text-color);
}
</style>
