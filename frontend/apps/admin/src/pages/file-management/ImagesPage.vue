<template>
  <div>
    <div class="toolbar">
      <el-select v-model="selectedModule" placeholder="选择模块" style="width: 200px" @change="handleModuleChange">
        <el-option v-for="module in modules" :key="module" :label="module" :value="module" />
      </el-select>
      <el-button @click="loadImages(true)" :icon="Refresh" :loading="loading">刷新列表</el-button>
      <el-upload action="" :show-file-list="false" :http-request="customUpload" :before-upload="beforeUpload">
        <el-button type="primary" :icon="Upload" :loading="uploading">上传图片</el-button>
      </el-upload>
    </div>

    <div v-if="loadError" class="error-tip">
      <span>{{ loadError }}</span>
      <el-button type="primary" :icon="Refresh" @click="loadImages(true)" size="small">重试</el-button>
    </div>

    <div v-loading="loading" v-show="!loadError" class="content-area">
      <div v-if="images.length === 0 && !loading" class="empty-tip">暂无图片</div>

      <div v-for="(group, date) in groupedImages" :key="date" class="date-group">
        <h3 class="date-header">{{ date }}</h3>
        <div class="image-grid">
          <div v-for="image in group" :key="image.path" class="image-card">
            <div class="image-wrapper">
              <img
                :src="image.url"
                :data-zoom-src="image.url"
                :alt="image.name"
                class="image-thumbnail zoomable-image"
                @error="handleImageError"
              />
            </div>
            <div class="image-info">
              <div class="image-name" :title="image.name">{{ image.name }}</div>
              <div class="image-meta">
                <span>{{ formatSize(image.size) }}</span>
                <span>{{ formatDate(image.modTime) }}</span>
              </div>
              <div class="image-actions">
                <el-button @click.stop="copyUrl(image)" :icon="Link" link size="small" title="复制URL">
                  复制URL
                </el-button>
                <el-button @click.stop="deleteImage(image)" :icon="Trash" link type="danger" size="small" title="删除">
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from "element-plus"
import Refresh from "virtual:icons/mdi/refresh"
import Trash from "virtual:icons/mdi/delete"
import Link from "virtual:icons/mdi/link"
import Upload from "virtual:icons/mdi/upload"
import { client } from "#/shared/lib/api"
import mediumZoom from "medium-zoom"

type ImageInfo = {
  path: string
  name: string
  url: string
  size: number
  modTime: string
}

const loading = ref(true)
const uploading = ref(false)
const loadError = ref("")
const modules = ref(["blog", "user"])
const selectedModule = ref(localStorage.getItem("file_image_last_module") ?? modules.value[0])
const images = ref<ImageInfo[]>([])
let zoomInstance: ReturnType<typeof mediumZoom> | null = null

function handleModuleChange() {
  localStorage.setItem("file_image_last_module", selectedModule.value)
  loadImages()
}

function beforeUpload(file: File) {
  const isImage = file.type.startsWith("image/")
  if (!isImage) {
    ElMessage.error("只能上传图片文件!")
  }
  return isImage
}

async function customUpload(options: any) {
  const { file } = options
  uploading.value = true
  const formData = new FormData()
  formData.append("image", file)

  const { error } = await (client as any).POST(`/api/file/image/${selectedModule.value}`, {
    body: formData as any
  })

  if (!error) {
    ElMessage.success("上传成功")
    loadImages()
  } else {
    ElMessage.error(error?.message ?? "上传失败")
  }
  uploading.value = false
}

// 按日期分组图片
const groupedImages = computed(() => {
  const groups: Record<string, ImageInfo[]> = {}
  images.value.forEach((image: ImageInfo) => {
    const date = new Date(image.modTime).toLocaleDateString("zh-CN")
    groups[date] ??= []
    groups[date].push(image)
  })
  // 按日期倒序排列
  return Object.fromEntries(
    Object.entries(groups).sort((a, b) => {
      return new Date(b[0]).getTime() - new Date(a[0]).getTime()
    })
  )
})

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + " B"
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + " KB"
  return (bytes / (1024 * 1024)).toFixed(2) + " MB"
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString("zh-CN")
}

function handleImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.style.display = "none"
  const errorDiv = document.createElement("div")
  errorDiv.className = "image-error"
  errorDiv.textContent = "加载失败"
  target.parentElement?.appendChild(errorDiv)
}

function initZoom() {
  // 销毁之前的实例
  if (zoomInstance) {
    zoomInstance.detach()
    zoomInstance = null
  }
  // 等待DOM更新后初始化，medium-zoom会自动处理图片加载
  nextTick(() => {
    zoomInstance = mediumZoom(".zoomable-image", {
      background: "rgba(0, 0, 0, 0.8)",
      margin: 24
    })
  })
}

async function loadImages(refresh = false) {
  loading.value = true
  loadError.value = ""

  const queryStr = refresh ? "?refresh=true" : ""
  const { data: res, error } = await (client as any).GET(`/api/file/image/${selectedModule.value}${queryStr}`)

  if (!error && res) {
    images.value = res
    loading.value = false
    loadError.value = ""
    initZoom()
  } else {
    loading.value = false
    loadError.value = error?.message ?? "获取图片列表失败"
    ElMessage.error(loadError.value)
  }
}

function copyUrl(image: ImageInfo) {
  navigator.clipboard
    .writeText(image.url)
    .then(() => {
      ElMessage.success("URL已复制到剪贴板")
    })
    .catch(() => {
      // 降级方案
      const textarea = document.createElement("textarea")
      textarea.value = image.url
      document.body.appendChild(textarea)
      textarea.select()
      try {
        document.execCommand("copy")
        ElMessage.success("URL已复制到剪贴板")
      } catch {
        ElMessage.error("复制失败")
      }
      document.body.removeChild(textarea)
    })
}

function deleteImage(image: ImageInfo) {
  ElMessageBox.confirm(`确定要删除图片 "${image.name}" 吗？`, "确认删除", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  })
    .then(async () => {
      const { error } = await (client as any).DELETE(`/api/file/image/${image.path}`)
      if (!error) {
        ElMessage.success("删除成功")
        loadImages()
      } else {
        ElMessage.error(error?.message ?? "删除失败")
      }
    })
    .catch(() => {})
}

onMounted(() => {
  loadImages()
})

onUnmounted(() => {
  if (zoomInstance) {
    zoomInstance.detach()
  }
})
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  align-items: center;
}

.empty-tip {
  text-align: center;
  color: #adadad;
  padding: 40px;
}

.error-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 24px;
  margin-bottom: 20px;
  background-color: #fef0f0;
  color: #f56c6c;
  border-radius: 8px;
}

.date-group {
  margin-bottom: 30px;
}

.date-header {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 15px;
  color: #333;
  border-bottom: 2px solid #e4e7ed;
  padding-bottom: 8px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.image-card {
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.image-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.image-wrapper {
  width: 100%;
  height: 200px;
  overflow: hidden;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-thumbnail {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: zoom-in;
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.image-info {
  padding: 12px;
}

.image-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
}

.image-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}
</style>

<style>
/* 确保 medium-zoom 遮罩层和图片在 nav (98) 和 menu (1501) 之上 */
.medium-zoom-overlay {
  z-index: 2000 !important;
}

.medium-zoom-image--opened {
  z-index: 2001 !important;
  object-fit: contain !important;
  max-width: 90vw !important;
  max-height: 90vh !important;
  width: auto !important;
  height: auto !important;
}
</style>
