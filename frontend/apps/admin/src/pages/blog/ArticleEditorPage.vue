<template>
  <div class="editor-container" v-loading="loading">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
        <MdEditor
          class="editor"
          v-model="article.markdown"
          :tab-width="4"
          @onUploadImg="onUploadImg"
          @onSave="saveToLocal"
        />
      </el-col>
      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
        <el-row class="editor-info">
          标题：
          <el-input placeholder="请输入标题" v-model="article.title"></el-input>
          <div class="info-item">
            摘要：
            <el-input
              :autosize="{ minRows: 5, maxRows: 8 }"
              placeholder="摘要"
              type="textarea"
              v-model="article.summary"
            >
            </el-input>
          </div>
          <div class="info-item">
            <el-tag
              :closable="true"
              :key="tag.id"
              @close="handleTagClose(tag)"
              effect="dark"
              v-for="tag in article.tags"
            >
              {{ tag.name }}
            </el-tag>
            <el-dropdown
              @click="createNewTag"
              @command="addTag"
              plain
              size="small"
              split-button
              trigger="click"
              type="info"
            >
              New
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="tag" :key="tag.id" v-for="tag in tags">
                    {{ tag.name }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-row>
        <el-row class="editor-control-strip">
          <el-col :span="4">
            <el-button @click="helpInfo" circle type="info">
              <Info />
            </el-button>
          </el-col>
          <el-col :span="20" class="text-align-right">
            <el-button-group>
              <el-button v-if="!isNewArticle" :loading="buttons.delete.loading" @click="deleteArticle" type="danger">
                <Trash />
              </el-button>
              <el-button title="恢复" @click="loadFromLocal" type="info">
                <Refresh />
              </el-button>
              <el-button :loading="buttons.save.loading" @click="saveArticle" type="primary">
                <Check />
              </el-button>
            </el-button-group>
          </el-col>
        </el-row>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { MdEditor } from "md-editor-v3"
import "md-editor-v3/lib/style.css"
import { ElMessage, ElMessageBox } from "element-plus"
import Info from "virtual:icons/mdi/information"
import Trash from "virtual:icons/mdi/delete"
import Refresh from "virtual:icons/mdi/refresh"
import Check from "virtual:icons/mdi/check"
import { useRoute, useRouter } from "vue-router"
// useSiteStore removed; local message wrapper is used instead
import { useMenuStore } from "#/shared/stores/menu"
import { localData } from "#/shared/lib/storage"
import { client, buildApiUrl } from "#/shared/lib/api"

interface Tag {
  id: number
  name: string
}

interface Article {
  aid: number
  title: string
  summary: string
  markdown: string
  tags?: Tag[]
}

const route = useRoute()
const router = useRouter()
function message(text: any, type?: string) {
  ElMessage({
    showClose: true,
    message: String(text),
    type: type as any
  })
}

const loading = ref(true)
const article = ref<Article>({
  aid: 0,
  title: "",
  summary: "",
  markdown: "",
  tags: []
})
const tags = ref<Tag[]>([])
const isNewArticle = ref(true)
const buttons = reactive({
  save: { loading: false },
  delete: { loading: false }
})
const menuStore = useMenuStore()
let _prevMenuActive: boolean | null = null

function helpInfo() {
  ElMessageBox.alert(
    "Ctrl+S 保存到本地\n\n图片删除功能需要自行实现，目前版本(3.x)的 md-editor-v3 还没有实现。",
    "帮助",
    {
      confirmButtonText: "OK"
    }
  )
}

function saveToLocal() {
  localData.set("article.content", article.value)
  message("已保存到本地", "success")
}

function loadFromLocal() {
  const localArticle = localData.get("article.content")
  if (localArticle) {
    article.value = localArticle
    message("已从本地恢复", "success")
  } else {
    message("本地没有备份", "info")
  }
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

async function saveArticle() {
  buttons.save.loading = true

  if (isNewArticle.value) {
    const { data: res, error } = await client.POST("/api/articles", {
      baseUrl: buildApiUrl("blog"),
      body: article.value as any
    })
    buttons.save.loading = false
    if (!error && res) {
      message("保存成功", "success")
      router.push(`/blog/editor/${res.aid}/`)
    } else {
      message(error?.message ?? "保存失败", "error")
    }
  } else {
    const { data: res, error } = await client.PUT("/api/articles/{aid}", {
      baseUrl: buildApiUrl("blog"),
      params: { path: { aid: Number(article.value.aid) } },
      body: article.value as any
    })
    buttons.save.loading = false
    if (!error && res) {
      message("保存成功", "success")
    } else {
      message(error?.message ?? "保存失败", "error")
    }
  }
}

function deleteArticle() {
  ElMessageBox.confirm("此操作将永久删除该文章, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  })
    .then(async () => {
      buttons.delete.loading = true
      const { error } = await client.DELETE("/api/articles/{aid}", {
        baseUrl: buildApiUrl("blog"),
        params: { path: { aid: Number(article.value.aid) } }
      })
      buttons.delete.loading = false
      if (!error) {
        message("删除成功", "success")
        router.push("/blog/articles/")
      } else {
        message(error?.message ?? "删除失败", "error")
      }
    })
    .catch(() => {
      message("已取消删除", "info")
    })
}

function handleTagClose(tag: Tag) {
  if (article.value.tags) {
    article.value.tags = article.value.tags.filter((t: Tag) => t.id !== tag.id)
  }
}

function addTag(tag: Tag) {
  article.value.tags ??= []
  if (!article.value.tags.find((t: Tag) => t.id === tag.id)) {
    article.value.tags.push(tag)
  }
}

function createNewTag() {
  ElMessageBox.prompt("请输入新标签", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消"
  }).then(async ({ value }: { value: string }) => {
    const { data: res, error } = await client.POST("/api/tags", {
      baseUrl: buildApiUrl("blog"),
      body: { name: value }
    })
    if (!error && res) {
      message("创建成功", "success")
      // Refetch tags to get the newly created tag with its ID
      const { data: tagsRes } = await client.GET("/api/tags", {
        baseUrl: buildApiUrl("blog")
      })
      if (tagsRes?.tags) {
        tags.value = tagsRes.tags as any
        const newTag = tagsRes.tags.find((t: any) => t.name === value)
        if (newTag) {
          addTag(newTag as any)
        }
      }
    } else {
      message(error?.message ?? "创建失败", "error")
    }
  })
}

let _prevMenuWidth: string | null = null
let _prevMenuWidthActive: string | null = null
let _mainMenuEl: HTMLElement | null = null
let _mainMenuDisplay: string | null = null
let _hideStyleEl: HTMLStyleElement | null = null

onMounted(() => {
  // hide left menu by overriding CSS variables and hiding <main-menu>
  try {
    const rootStyles = getComputedStyle(document.documentElement)
    _prevMenuWidth = rootStyles.getPropertyValue("--menu-width") || null
    _prevMenuWidthActive = rootStyles.getPropertyValue("--menu-width-active") || null
    document.documentElement.style.setProperty("--menu-width", "0px")
    document.documentElement.style.setProperty("--menu-width-active", "0px")

    // add a root class and style rules to force-hide the menu (covers .menu.active override)
    document.documentElement.classList.add("hide-left-menu")
    const style = document.createElement("style")
    style.id = "hide-left-menu-style"
    style.innerHTML = `
            .hide-left-menu .menu,
            .hide-left-menu .menu.active {
                width: 0 !important;
                min-width: 0 !important;
                overflow: hidden !important;
            }
            .hide-left-menu .router-view {
                margin-left: 0 !important;
                width: 100% !important;
            }
        `
    document.head.appendChild(style)
    _hideStyleEl = style

    _mainMenuEl = document.querySelector(".menu") as HTMLElement | null
    if (_mainMenuEl) {
      _mainMenuDisplay = _mainMenuEl.style.display || null
      _mainMenuEl.style.display = "none"
    }
  } catch {
    // ignore DOM errors in SSR or tests
  }

  // store-driven hide: collapse the menu (remembers previous state)
  try {
    _prevMenuActive = menuStore.active
    menuStore.setMenuState(false)
  } catch {
    // ignore if store unavailable
  }

  // existing data loading logic
  const aid = Number(route.params.aid)
  if (aid) {
    isNewArticle.value = false
    client
      .GET("/api/articles/{aid}/edit", {
        baseUrl: buildApiUrl("blog"),
        params: { path: { aid } }
      })
      .then(({ data: res, error }) => {
        if (!error && res) {
          article.value = (res.article ?? res) as any
        } else {
          message(error?.message ?? "获取文章失败", "error")
        }
        loading.value = false
      })
  } else {
    loading.value = false
  }

  client
    .GET("/api/tags", {
      baseUrl: buildApiUrl("blog")
    })
    .then(({ data: res, error }) => {
      if (!error && res) {
        tags.value = res.tags as any
      }
    })
})

onUnmounted(() => {
  // restore menu styles
  try {
    if (_prevMenuWidth !== null) document.documentElement.style.setProperty("--menu-width", _prevMenuWidth)
    if (_prevMenuWidthActive !== null)
      document.documentElement.style.setProperty("--menu-width-active", _prevMenuWidthActive)
    if (_mainMenuEl) _mainMenuEl.style.display = _mainMenuDisplay ?? ""
    if (_hideStyleEl && _hideStyleEl.parentNode) _hideStyleEl.parentNode.removeChild(_hideStyleEl)
    document.documentElement.classList.remove("hide-left-menu")
    // restore previous store-driven menu state
    try {
      if (_prevMenuActive !== null) menuStore.setMenuState(_prevMenuActive)
    } catch {
      // ignore
    }
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.editor {
  /* 
        Dynamic Height Calculation:
        100vh - Header(121) - Padding(40) - Footer(61) - Buffer(varies)
    */
  height: calc(100vh - var(--header-height) - (var(--app-content-padding) * 2) - var(--footer-height) - 2px);
}

.editor-info {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.editor-control-strip {
  margin-top: 20px;
}
</style>
