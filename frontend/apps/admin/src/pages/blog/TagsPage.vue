<template>
  <div v-loading="loading">
    <el-table :data="tags" stripe style="width: 100%">
      <el-table-column fixed prop="id" label="ID" width="100"> </el-table-column>
      <el-table-column fixed prop="name" label="标签" width="180"> </el-table-column>
      <el-table-column label="操作" width="180">
        <template v-slot="scope">
          <el-button @click="editTag(scope.row)" link size="small">编辑</el-button>
          <el-button
            :loading="showArticlesButtonLoading[String(scope.row.id)]"
            @click="getArticlesByTag(scope.row.id)"
            link
            size="small"
            >查看文章
          </el-button>
          <el-button @click="deleteTag(scope.row)" link type="danger" size="small">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
// 使用自动导入的 API
import { ElMessage, ElMessageBox } from "element-plus"
import { useNavStore } from "#/shared/stores/nav"
import { client, buildApiUrl } from "#/shared/lib/api"

const navStore = useNavStore()

const loading = ref(true)
const showArticlesButtonLoading = ref({})
const tags = ref([])

async function getTags() {
  loading.value = true
  const { data: res, error } = await client.GET("/api/tags", {
    baseUrl: buildApiUrl("blog")
  })
  if (!error && res) {
    tags.value = res.tags as any
    loading.value = false
  } else {
    loading.value = false
    ElMessage.error(error?.message ?? "获取失败")
  }
}

function addTag() {
  ElMessageBox.prompt("标签名称:", "添加标签", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    inputPattern: /[a-zA-Z0-9\u4e00-\u9fa5]+/,
    inputErrorMessage: "只能输入数字、英文字母和汉字"
  })
    .then(async ({ value }: { value: string }) => {
      value = value.trim()
      if (value === "") {
        ElMessage({
          type: "info",
          message: "未作修改"
        })
      } else {
        const { error } = await client.POST("/api/tags", {
          baseUrl: buildApiUrl("blog"),
          body: { name: value }
        })
        if (!error) {
          ElMessage({
            showClose: true,
            message: `标签: "${value}" 已添加`,
            type: "success"
          })
          getTags()
        } else {
          ElMessage({
            showClose: true,
            message: error?.message ?? "添加失败",
            type: "error"
          })
        }
      }
    })
    .catch(() => {})
}

function editTag(tag: any) {
  ElMessageBox.prompt("标签名称:", "编辑", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    inputValue: tag.name,
    inputPattern: /[a-zA-Z0-9\u4e00-\u9fa5]+/,
    inputErrorMessage: "只能输入数字、英文字母和汉字"
  })
    .then(async ({ value }: { value: string }) => {
      value = value.trim()
      if (value === "" || value === tag.name) {
        ElMessage({
          type: "info",
          message: "未作修改"
        })
      } else {
        const { error } = await client.PUT("/api/tags/{id}", {
          baseUrl: buildApiUrl("blog"),
          params: { path: { id: tag.id } },
          body: { name: value }
        })
        if (!error) {
          ElMessage({
            showClose: true,
            message: `标签: "${tag.name}" 已修改为 "${value}"`,
            type: "success"
          })
          getTags()
        } else {
          ElMessage({
            showClose: true,
            message: error?.message ?? "修改失败",
            type: "error"
          })
        }
      }
    })
    .catch(() => {})
}

function deleteTag(tag: any) {
  ElMessageBox({
    type: "warning",
    title: "删除",
    message: `将会删除标签: "${tag.name}" 是否确认？`
  })
    .then(async () => {
      const { error } = await client.DELETE("/api/tags/{id}", {
        baseUrl: buildApiUrl("blog"),
        params: { path: { id: tag.id } }
      })
      if (!error) {
        ElMessage({
          showClose: true,
          message: `标签: "${tag.name}" 已被删除`,
          type: "success"
        })
        getTags()
      } else {
        ElMessage({
          showClose: true,
          message: error?.message ?? "删除失败",
          type: "error"
        })
      }
    })
    .catch(() => {})
}

async function getArticlesByTag(id: any) {
  showArticlesButtonLoading.value[String(id)] = true
  const { data: res, error } = await client.GET("/api/tags/{id}/articles", {
    baseUrl: buildApiUrl("blog"),
    params: { path: { id: Number(id) } }
  })

  if (!error && res) {
    const articlesList = res.articles ?? []
    const articles = [h("span", { style: "font-weight: bold" }, `共 ${articlesList.length} 篇文章`), h("hr")]
    articlesList.forEach((article: any) => {
      articles.push(h("span", { style: "font-size: 14px" }, article.title))
      articles.push(h("br"))
    })
    ElMessageBox({
      title: "文章列表",
      message: h("div", null, articles)
    }).catch(() => {})
  } else {
    ElMessage({
      showClose: true,
      message: error?.message ?? "获取失败",
      type: "error"
    })
  }
  showArticlesButtonLoading.value[String(id)] = false
}

onMounted(() => {
  // 添加刷新按钮
  const addButton = {
    text: "新建标签",
    click: () => addTag()
  }
  const refreshButton = {
    text: "刷新",
    click: () => getTags()
  }
  navStore.setNavButtons([addButton, refreshButton])
  getTags()
})
onUnmounted(() => {
  navStore.clearNavButtons()
})
</script>

<style scoped></style>
