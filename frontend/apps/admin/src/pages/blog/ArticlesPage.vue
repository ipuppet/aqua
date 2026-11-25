<template>
  <div>
    <div class="articles" v-loading="loading" v-show="articles.length > 0">
      <div :key="article.aid + '.' + article.created_at" class="articles-item" v-for="article in articles">
        <div class="articles-item">
          <a @click="openDialog(article.aid)" class="title pointer" v-html="getArticleTitle(article)"></a>
          <div class="info">
            <span class="date">{{ new Date(article.created_at).toLocaleString() }}</span>
            <span :key="tag.id" v-for="(tag, i) in article.tags">
              {{ tag.name }}
              <span v-if="i !== article.tags.length - 1"> | </span>
            </span>
          </div>
          <render-markdown :model-value="article.summary" />
        </div>
        <hr />
      </div>
    </div>
    <div v-show="articles.length === 0 && !loading" class="empty-tip">暂无文章</div>
    <el-pagination
      v-show="total > 0"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus"
import { useNavStore } from "#/shared/stores/nav"
import { storeToRefs } from "pinia"
import { client, buildApiUrl } from "#/shared/lib/api"

const navStore = useNavStore()
const { searchKeyword } = storeToRefs(navStore)

const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const articles = ref([])
const isSearchMode = ref(false)

function escapeHtml(value: string) {
  return value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;")
}

function sanitizeHighlightHtml(value: string) {
  const emOpenPlaceholder = "__EM_OPEN__"
  const emClosePlaceholder = "__EM_CLOSE__"

  // Keep only <em> highlight tags, escape any other HTML.
  return escapeHtml(value.replace(/<\s*em\s*>/gi, emOpenPlaceholder).replace(/<\s*\/\s*em\s*>/gi, emClosePlaceholder))
    .replace(new RegExp(emOpenPlaceholder, "g"), "<em>")
    .replace(new RegExp(emClosePlaceholder, "g"), "</em>")
}

function getArticleTitle(article: any) {
  const rawTitle = String(article?._formatted?.title ?? article?.title ?? "")
  return isSearchMode.value ? sanitizeHighlightHtml(rawTitle) : escapeHtml(rawTitle)
}

async function loadArticles() {
  loading.value = true
  const { data: res, error } = await client.GET("/api/articles", {
    baseUrl: buildApiUrl("blog"),
    params: { query: { page: currentPage.value, pageSize: pageSize.value } }
  })
  if (!error && res) {
    total.value = Number(res.total ?? 0)
    articles.value = (res.articles ?? []) as any
    loading.value = false
  } else {
    loading.value = false
    ElMessage.error(error?.message ?? "获取文章列表失败")
  }
}

function handleSearch(keyword: string) {
  const trimmedKeyword = keyword.trim()
  if (trimmedKeyword === "") {
    clearSearch()
  } else {
    isSearchMode.value = true
    currentPage.value = 1
    loadSearchResults()
  }
}

async function loadSearchResults() {
  loading.value = true
  const { data: res, error } = await client.GET("/api/search/{keyword}", {
    baseUrl: buildApiUrl("blog"),
    params: {
      path: { keyword: searchKeyword.value },
      query: { page: currentPage.value, pageSize: pageSize.value }
    }
  })

  loading.value = false
  if (!error && res) {
    const sourceArticles = res.articles ?? []
    if (sourceArticles.length === 0) {
      ElMessage.warning("未找到相关内容")
      articles.value = []
      total.value = 0
    } else {
      total.value = Number(res.total ?? sourceArticles.length)
      articles.value = (sourceArticles as any[]).map((article: any) => ({
        ...article,
        title: article?._formatted?.title ?? article?.title ?? "",
        summary: article?._formatted?.summary ?? article?.summary ?? ""
      }))
    }
  } else {
    ElMessage.error(error?.message ?? "搜索失败")
  }
}

function clearSearch() {
  isSearchMode.value = false
  currentPage.value = 1
  pageSize.value = 10
  loadArticles()
}

function handleSizeChange(val: number) {
  pageSize.value = val
  currentPage.value = 1
  if (isSearchMode.value) {
    loadSearchResults()
  } else {
    loadArticles()
  }
}

function handleCurrentChange(val: number) {
  currentPage.value = val
  if (isSearchMode.value) {
    loadSearchResults()
  } else {
    loadArticles()
  }
}

function openDialog(aid: any = null) {
  if (aid !== null) window.open("/#/blog/editor/" + aid + "/", "Editor", "width=1280,height=960")
  else window.open("/#/blog/editor/", "Editor", "width=1280,height=960")
}

async function syncIndex() {
  ElMessage.info("正在发送同步索引请求...")
  const { error } = await client.POST("/api/meili/sync", {
    baseUrl: buildApiUrl("blog"),
    body: {}
  })
  if (!error) {
    ElMessage.success("同步索引成功")
  } else {
    ElMessage.error(error?.message ?? "同步索引失败")
  }
}

onMounted(() => {
  loadArticles()
  const navButtons = [
    {
      text: "同步索引",
      type: "warning",
      click: () => {
        syncIndex()
      }
    },
    {
      text: "添加文章",
      click: () => {
        openDialog()
      }
    }
  ]
  navStore.setNavButtons(navButtons)
  navStore.setSearch(true, handleSearch, clearSearch)
})
onUnmounted(() => {
  navStore.clearNavButtons()
  navStore.setSearch(false)
})
</script>

<style scoped>
.articles {
  background: #ffffff;
  padding: 20px;
}

.articles hr {
  width: 100px;
  height: 0;
  border-bottom: 1px solid #b5b5b5;
  margin-top: 30px;
  margin-bottom: 30px;
}

.articles-item :deep(img) {
  max-width: 100%;
  max-height: 300px;
}

.title {
  text-decoration: none;
  color: #000000;
  transition: 0.3s;
  font-size: 26px;
}

.title:hover {
  margin-left: 10px;
}

.articles-item :deep(em) {
  font-style: normal;
  color: #333333;
  background-color: #ffe58f;
  border-radius: 2px;
  padding: 0 2px;
}

.info,
.tags {
  color: #adadad;
  font-size: 12px;
}

.tags {
  margin: 0 5px 0 5px;
}

.tags:hover {
  color: #343434;
  margin: 0 10px 0 10px;
}

.date {
  margin-right: 10px;
}

.empty-tip {
  text-align: center;
  color: #adadad;
  padding: 40px;
}
.pagination {
  margin-top: 20px;
  justify-content: center;
}
</style>
