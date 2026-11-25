<template>
  <el-container :class="{ dark: darkModeStore.isDarkMode }" class="article-comments-section">
    <el-header v-if="total > 0">Comments ({{ total }})</el-header>
    <el-header v-else>No comments yet. Be the first to comment!</el-header>
    <el-main class="new-comment-area">
      <el-input
        v-model="userComment.content"
        type="textarea"
        :autosize="{ minRows: 3, maxRows: 5 }"
        placeholder="Share your thoughts..."
        clearable
      />
      <el-button type="primary" :icon="Check" @click="handleUserComment" style="margin-top: 10px"
        >Post Comment</el-button
      >
    </el-main>
    <el-container class="comment-list-container" direction="vertical">
      <el-container class="comment-container" v-for="(comment, i) in comments" :key="comment.id || i">
        <el-aside>
          <el-avatar :size="50" :src="users[comment.uid]?.avatar"
            ><img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"
          /></el-avatar>
        </el-aside>
        <el-main>
          <div class="comment-meta">
            <el-link :href="users[comment.uid]?.website || '#'" target="_blank" :underline="false">{{
              users[comment.uid]?.nickname || "Anonymous"
            }}</el-link>
            <el-text class="comment-time" size="small">{{ new Date(comment.createdAt).toLocaleString() }}</el-text>
          </div>
          <div class="comment-content">
            <el-text tag="p">{{ comment.content }}</el-text>
          </div>
          <div v-if="children[comment.id] && children[comment.id].total > 0" class="child-comment-actions">
            <el-button
              link
              type="primary"
              size="small"
              @click="toggleChildren(comment.id)"
              class="toggle-replies-button"
            >
              {{
                children[comment.id].isCollapsed
                  ? `View ${children[comment.id].total} ${children[comment.id].total === 1 ? "reply" : "replies"}`
                  : "Hide replies"
              }}
            </el-button>
          </div>
          <div
            v-if="children[comment.id] && children[comment.id].total > 0 && !children[comment.id].isCollapsed"
            class="child-comments-wrapper"
          >
            <div
              v-if="
                children[comment.id].isLoading &&
                children[comment.id].page === 0 &&
                children[comment.id].comments.length === 0
              "
              class="loading-indicator centered-loading"
            >
              <el-text size="small">Loading replies...</el-text>
            </div>
            <template v-if="children[comment.id].comments.length > 0">
              <el-container
                class="child-container"
                v-for="(child, childIdx) in children[comment.id].comments"
                :key="child.id || childIdx"
              >
                <el-aside
                  ><el-avatar :size="35" :src="users[child.uid]?.avatar"
                    ><img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png" /></el-avatar
                ></el-aside>
                <el-main class="child-comment-main">
                  <div class="comment-meta">
                    <el-link :href="users[child.uid]?.website || '#'" target="_blank" :underline="false">{{
                      users[child.uid]?.nickname || "Anonymous"
                    }}</el-link>
                    <el-text class="comment-time" size="small" v-if="child.createdAt">{{
                      new Date(child.createdAt).toLocaleString()
                    }}</el-text>
                  </div>
                  <div class="comment-content child-comment-content">
                    <el-text tag="p" size="small">{{ child.content }}</el-text>
                  </div>
                </el-main>
              </el-container>
            </template>
            <div
              v-else-if="
                !children[comment.id].isLoading &&
                children[comment.id].page > 0 &&
                children[comment.id].comments.length === 0
              "
            >
              <el-text size="small">No replies to display.</el-text>
            </div>
            <div
              class="load-more-action"
              v-if="
                !children[comment.id].isLoading &&
                children[comment.id].comments.length > 0 &&
                children[comment.id].comments.length < children[comment.id].total
              "
            >
              <el-button
                link
                type="primary"
                size="small"
                @click="handleLoadMoreChildren(comment.id)"
                class="load-more-button"
                >Load more replies ({{
                  children[comment.id].total - children[comment.id].comments.length
                }}
                remaining)</el-button
              >
            </div>
            <div
              class="load-more-action loading-indicator"
              v-if="children[comment.id].isLoading && children[comment.id].page > 0"
            >
              <el-text size="small">Loading more replies...</el-text>
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </el-container>
</template>

<script setup>
import { Check } from "@element-plus/icons-vue"
import { ref, watch } from "vue"
import { useRoute } from "vue-router"
import { ElMessage } from "element-plus"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { client } from "#/shared/lib/api"
import { getCookie, ACCESS_TOKEN_KEY } from "@aqua/http"

const route = useRoute()
const darkModeStore = useDarkModeStore()
const emit = defineEmits(["commentPosted"])
const INITIAL_CHILD_PAGE_SIZE = 5

const { total, data } = defineProps({
  total: { type: Number, required: true },
  data: { type: Array, required: true }
})

const users = ref({})
const comments = ref([])
const children = ref({})
const userComment = ref({ content: "" })

async function loadUserInfo(uid) {
  if (!users.value[uid]) {
    users.value[uid] = { nickname: "Loading...", avatar: "" }
    const { data: res, error } = await client.GET("/api/users/{uid}", {
      params: { path: { uid: Number(uid) } }
    })
    if (!error && res) {
      users.value[uid] = res
    } else {
      users.value[uid] = { nickname: "User", avatar: "" }
    }
  }
}

async function loadChildren(parentId, pageToLoad = 1, pageSize = INITIAL_CHILD_PAGE_SIZE) {
  if (!children.value[parentId]) {
    children.value[parentId] = { comments: [], total: 0, page: 0, pageSize, isCollapsed: true, isLoading: false }
  } else {
    if (typeof children.value[parentId].isCollapsed === "undefined") children.value[parentId].isCollapsed = true
    if (typeof children.value[parentId].isLoading === "undefined") children.value[parentId].isLoading = false
    if (typeof children.value[parentId].pageSize === "undefined") children.value[parentId].pageSize = pageSize
  }
  const childData = children.value[parentId]
  if (childData.isLoading) return
  if (pageToLoad > 1 && pageToLoad <= childData.page && childData.comments.length >= childData.total) return
  childData.isLoading = true

  const {
    data: res,
    error,
    response
  } = await client.GET("/api/comments/{id}/children", {
    params: { path: { id: Number(parentId) }, query: { page: pageToLoad, pageSize } }
  })

  if (!error && res) {
    if (pageToLoad === 1) childData.comments = res.comments ?? []
    else {
      const existingIds = new Set(childData.comments.map(c => c.id))
      childData.comments.push(...(res.comments ?? []).filter(c => !existingIds.has(c.id)))
    }
    childData.total = res.total ?? 0
    childData.page = pageToLoad
    childData.pageSize = pageSize
    ;(res.comments ?? []).forEach(child => loadUserInfo(child.uid))
  } else if (response.status === 404 && pageToLoad === 1) {
    childData.total = 0
    childData.comments = []
  }
  childData.isLoading = false
}

function toggleChildren(id) {
  if (!children.value[id]) return
  children.value[id].isCollapsed = !children.value[id].isCollapsed
  const childData = children.value[id]
  if (!childData.isCollapsed && childData.comments.length === 0 && childData.total > 0 && !childData.isLoading) {
    loadChildren(id, 1, childData.pageSize ?? INITIAL_CHILD_PAGE_SIZE)
  }
}

function handleLoadMoreChildren(commentId) {
  const childData = children.value[commentId]
  if (childData && !childData.isLoading && childData.comments.length < childData.total) {
    loadChildren(commentId, childData.page + 1, childData.pageSize)
  }
}

function loadComments(commentData) {
  comments.value = commentData
  commentData.forEach(comment => {
    setTimeout(() => {
      loadChildren(comment.id, 1, INITIAL_CHILD_PAGE_SIZE)
      loadUserInfo(comment.uid)
    }, 1)
  })
}

async function handleUserComment() {
  if (getCookie(ACCESS_TOKEN_KEY) === null) {
    ElMessage.error("Please login to comment.")
    return
  }
  if (userComment.value.content.trim() === "") {
    ElMessage.warning("Comment content cannot be empty.")
    return
  }
  const articleId = Number(route.params.aid)
  if (!articleId) {
    ElMessage.error("Article ID is missing. Cannot post comment.")
    return
  }

  const { error } = await client.POST("/api/comments/articles/{aid}", {
    params: { path: { aid: articleId } },
    body: { content: userComment.value.content, parentId: -1 }
  })

  if (!error) {
    ElMessage.success("Comment posted successfully!")
    userComment.value.content = ""
    emit("commentPosted")
  } else {
    ElMessage.error(error?.message ?? "Failed to post comment.")
  }
}

watch(
  () => data,
  current => loadComments(current),
  { immediate: true }
)
</script>

<style scoped>
.article-comments-section {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: var(--main-block-color);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}
.new-comment-area {
  margin-bottom: 30px;
  padding: 0 !important;
}
.comment-list-container {
  padding: 0;
}
.comment-container {
  display: flex;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--input-border-color);
}
.comment-container > .el-aside {
  width: 50px !important;
  height: 50px !important;
  margin-right: 15px;
  flex-shrink: 0;
  overflow: hidden;
}
.comment-container > .el-main {
  padding: 0 !important;
  flex-grow: 1;
  margin-left: 0 !important;
}
.comment-meta {
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.comment-content {
  background-color: var(--main-block-color);
  padding: 10px 15px;
  border-radius: 6px;
  line-height: 1.6;
  border: 1px solid var(--input-border-color);
}
.child-comments-wrapper {
  margin-top: 15px;
  padding-left: 20px;
  border-left: 2px solid var(--input-border-color);
}
.child-container {
  display: flex;
  margin-top: 15px;
}
.child-container > .el-aside {
  width: 35px !important;
  height: 35px !important;
  margin-right: 10px;
  flex-shrink: 0;
  overflow: hidden;
}
.child-container > .el-main {
  padding: 0 !important;
  flex-grow: 1;
  margin-left: 0 !important;
}
.load-more-action {
  margin-top: 10px;
  text-align: left;
}
</style>
