<template>
  <div>
    <div :class="{ dark: darkModeStore.isDarkMode }" v-show="isLoaded" class="header">{{ title }}</div>
    <article-list v-if="page.articles.length > 0" :is-loaded="isLoaded" :page="page" @current-change="changePage" />
    <span v-else v-show="isLoaded">Nothing found</span>
  </div>
</template>
<script>
import ArticleList from "#/entities/article/ui/ArticleList.vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { client } from "#/shared/lib/api"
export default {
  components: { ArticleList },
  setup() {
    const darkModeStore = useDarkModeStore()
    return { darkModeStore }
  },
  data: () => ({ title: "", page: { articles: [], pageSize: 10, total: 1, current: 1 }, kw: "", isLoaded: false }),
  methods: {
    changePage(page) {
      this.page.current = page
      this.getSearchResult()
    },
    async getSearchResult() {
      this.kw = this.$route.params.kw
      this.title = `Search result for '${this.kw}'`
      document.title = this.title
      this.isLoaded = false
      window.scrollTo(0, 0)
      const { data, error } = await client.GET("/api/search/{keyword}", {
        params: {
          path: { keyword: this.kw },
          query: { page: this.page.current, pageSize: this.page.pageSize }
        }
      })
      if (!error && data) {
        this.page.articles = data.articles ?? []
        this.page.total = data.total ?? 0
        this.isLoaded = true
      } else {
        this.isLoaded = true
      }
    }
  },
  created() {
    this.getSearchResult()
  },
  watch: {
    $route(to, from) {
      if (to !== from) {
        this.isLoaded = false
        this.getSearchResult()
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
