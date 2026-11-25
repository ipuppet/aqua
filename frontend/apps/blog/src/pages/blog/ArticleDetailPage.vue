<template>
  <div :class="{ dark: darkModeStore.isDarkMode }" class="router-view-container">
    <div class="box">
      <loading-ball-atom v-show="!isLoaded" class="loading" />
      <h1 :key="'title' + article.aid">{{ article.title }}</h1>
      <article-info :key="'info' + article.aid" :date="article.createdAt" :tags="article.tags" />
      <div :key="'content' + article.aid" class="markdown"><render-markdown v-model="article.markdown" /></div>
      <div class="context">
        <div class="previous">
          Previous Post <br />
          <span v-if="!previous" :class="{ dark: darkModeStore.isDarkMode }" class="link">No previous post</span>
          <router-link
            v-else
            :to="'/detail/' + previous.aid + '/'"
            :class="{ dark: darkModeStore.isDarkMode }"
            class="link"
          >
            {{ previous.title }}
          </router-link>
        </div>
        <div class="next">
          Next Post <br />
          <span v-if="!next" :class="{ dark: darkModeStore.isDarkMode }" class="link">No next post</span>
          <router-link
            v-else
            :to="'/detail/' + next.aid + '/'"
            :class="{ dark: darkModeStore.isDarkMode }"
            class="link"
          >
            {{ next.title }}
          </router-link>
        </div>
      </div>
      <el-divider />
      <article-comments
        :key="'comments' + article.aid"
        :data="comments.comments"
        :total="comments.total"
        @commentPosted="getComments"
      />
    </div>
  </div>
</template>
<script>
import ArticleInfo from "#/entities/article/ui/ArticleInfo.vue"
import LoadingBallAtom from "#/shared/ui/LoadingBallAtom.vue"
import RenderMarkdown from "#/shared/ui/RenderMarkdown.vue"
import ArticleComments from "#/features/comment/ui/ArticleComments.vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { sessionData } from "#/shared/lib/storage"
import { client } from "#/shared/lib/api"
export default {
  components: { ArticleInfo, LoadingBallAtom, RenderMarkdown, ArticleComments },
  setup() {
    const darkModeStore = useDarkModeStore()
    return { darkModeStore }
  },
  data: () => ({
    article: {},
    comments: { total: 0, comments: [], page: 1, pageSize: 10 },
    previous: {},
    next: {},
    isLoaded: false
  }),
  methods: {
    async getArticle() {
      const aid = this.$route.params.aid
      const articles = sessionData.get("articles") ?? {}
      if (!articles[aid]) {
        const { data, error } = await client.GET("/api/articles/{aid}", {
          params: { path: { aid: Number(aid) } }
        })
        if (!error && data) {
          this.article = data.article ?? {}
          this.previous = data.previous ?? {}
          this.next = data.next ?? {}
          this.isLoaded = true
          document.title = this.article.title ?? ""
          const newArticle = {}
          newArticle[aid] = data
          sessionData.set("articles", Object.assign(articles, newArticle))
        } else {
          this.isLoaded = true
        }
      } else {
        this.article = articles[aid].article
        this.previous = articles[aid].previous
        this.next = articles[aid].next
        this.isLoaded = true
        document.title = this.article.title
      }
    },
    async getComments() {
      const aid = this.$route.params.aid
      const { data, error } = await client.GET("/api/comments/articles/{aid}", {
        params: {
          path: { aid: Number(aid) },
          query: { page: this.comments.page, pageSize: this.comments.pageSize }
        }
      })
      if (!error && data) {
        this.comments.total = data.total ?? 0
        this.comments.comments = data.comments ?? []
      } else {
        this.comments.total = 0
        this.comments.comments = []
      }
    }
  },
  created() {
    this.getArticle()
    this.getComments()
  },
  watch: {
    $route(to, from) {
      if (to !== from) {
        this.isLoaded = false
        this.getArticle()
        this.getComments()
      }
    }
  }
}
</script>
<style scoped>
.box {
  margin: 0 30px;
}

.box h1 {
  margin: 0;
}

.markdown {
  margin-top: 20px;
  margin-bottom: 20px;
}

.context {
  color: var(--text-color);
  opacity: 0.7;
  font-size: 12px;
  height: 35px;
  position: relative;
  transition:
    color 0.8s ease,
    opacity 0.8s ease;
}

.next {
  position: absolute;
  right: 0;
  height: 35px;
  width: 50%;
  text-align: right;
}

.previous {
  position: absolute;
  left: 0;
  height: 35px;
  width: 50%;
  text-align: left;
}
</style>
