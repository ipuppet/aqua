<template>
  <section class="about-editor-page" v-loading="loading">
    <div class="about-editor-area">
      <MdEditor
        v-model="markdown"
        class="editor"
        :class="{ 'full-screen': fullScreen }"
        @on-fullscreen="(f: boolean) => (fullScreen = f)"
        @onUploadImg="onUploadImg"
        @onSave="saveToLocal"
      />
    </div>
    <div class="about-editor-toolbar">
      <el-row class="editor-control-strip">
        <el-col :span="4">
          <el-button @click="helpInfo" circle type="info">
            <Info />
          </el-button>
        </el-col>
        <el-col :span="20" class="text-align-right">
          <el-button-group>
            <el-button title="恢复" @click="loadFromLocal" type="info">
              <Refresh />
            </el-button>
            <el-button :loading="saveButtonLoading" @click="saveArticle" type="primary">
              <Check />
            </el-button>
          </el-button-group>
        </el-col>
      </el-row>
    </div>
  </section>
</template>

<script setup lang="ts">
import { MdEditor } from "md-editor-v3"
import "md-editor-v3/lib/style.css"
import Info from "virtual:icons/mdi/information"
import Refresh from "virtual:icons/mdi/refresh"
import Check from "virtual:icons/mdi/check"
import { ElMessage, ElMessageBox } from "element-plus"
import { useNavStore } from "#/shared/stores/nav"
import { client, buildApiUrl } from "#/shared/lib/api"
import { localData } from "#/shared/lib/storage"

const navStore = useNavStore()

const loading = ref(true)
const fullScreen = ref(false)
const saveButtonLoading = ref(false)
const markdown = ref("")

function message(text: any, type?: string) {
  ElMessage({
    showClose: true,
    message: String(text),
    type: type as any
  })
}

function helpInfo() {
  const info =
    '快捷键 "ctrl+s" 可保存在本地存储中。<br>您可以点击右下角恢复按钮恢复上一次保存在本地的信息。<br>您可以保存空内容来清除本地存储。'
  ElMessageBox.confirm(info, "帮助", {
    dangerouslyUseHTMLString: true,
    type: "info",
    center: true
  }).catch(() => {})
}

function saveToLocal() {
  if (loading.value) return
  localData.set("blog.about", markdown.value)
  message("本地保存成功", "success")
}

function loadFromLocal() {
  const msg = "恢复您保存在本地的文章？"
  ElMessageBox.confirm(msg, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "info",
    center: true
  })
    .then(() => {
      markdown.value = localData.get("blog.about")
      message("恢复成功", "success")
    })
    .catch(() => {})
}

async function saveArticle() {
  if (loading.value) return
  saveButtonLoading.value = true
  const { error } = await client.PUT("/api/about", {
    baseUrl: buildApiUrl("blog"),
    body: { content: markdown.value }
  })
  if (!error) {
    message("更新成功", "success")
  } else {
    message(error?.message ?? "更新失败", "error")
  }
  saveButtonLoading.value = false
}

async function onUploadImg(files: File[], callback: (urls: string[]) => void) {
  const urls = await Promise.all(
    files.map(async file => {
      const formData = new FormData()
      formData.append("image", file)
      const { data: res, error } = await (client as any).POST("/api/file/image/blog", {
        body: formData as any
      })
      if (!error && res) {
        return res.url
      } else {
        message(error?.message ?? "上传失败", "error")
        return ""
      }
    })
  )
  callback(urls.filter(url => url !== ""))
}

onMounted(async () => {
  const { data: res, error } = await client.GET("/api/about", {
    baseUrl: buildApiUrl("blog")
  })
  if (!error && res) {
    markdown.value = res.content ?? ""
    loading.value = false
  } else {
    message(error?.message ?? "获取失败", "error")
  }
})

onUnmounted(() => {
  navStore.clearNavButtons()
})
</script>

<style scoped>
.about-editor-area {
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
}

.editor {
  /* 
       Dynamic Height Calculation:
       Viewport Height (100vh)
       - Header Height (Nav + Tabs + Border)
       - App Content Padding (Top + Bottom)
       - Footer Height
       - Editor Toolbar Height
       - Safety Buffer (8px)
    */
  height: calc(
    100vh - var(--header-height) - (var(--app-content-padding) * 2) - var(--footer-height) -
      var(--editor-toolbar-height) - 8px
  );
}

.editor.full-screen {
  max-height: 100vh;
}

.about-editor-toolbar {
  background: #ffffff;
  border-radius: 12px;
  padding: 12px 16px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.04);
}

.editor-control-strip {
  margin: 0;
  width: 100%;
}

.text-align-right {
  text-align: right;
}
</style>
