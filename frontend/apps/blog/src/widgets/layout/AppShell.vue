<template>
  <div id="app" :class="{ dark: darkModeStore.isDarkMode }">
    <nav-bar />
    <div class="space-top"></div>
    <router-view
      :class="['router-view-' + contentColumn, { dark: darkModeStore.isDarkMode }]"
      :style="{ width: contentWidth }"
      class="router-view"
    />
    <div :class="{ 'root-widgets-bottom': contentColumn === 1 }" class="root-widgets root-widgets-left">
      <my-info :class="{ dark: darkModeStore.isDarkMode }" :style="{ width: widgetWidth }" />
      <template v-if="contentColumn <= 2">
        <search-box :class="{ dark: darkModeStore.isDarkMode }" :style="{ width: widgetWidth }" />
        <tag-cloud :class="{ dark: darkModeStore.isDarkMode }" :style="{ width: widgetWidth }" :tags="tags" />
      </template>
    </div>
    <div class="root-widgets root-widgets-right">
      <template v-if="contentColumn > 2">
        <search-box :class="{ dark: darkModeStore.isDarkMode }" :style="{ width: widgetWidth }" />
        <tag-cloud :class="{ dark: darkModeStore.isDarkMode }" :style="{ width: widgetWidth }" :tags="tags" />
      </template>
    </div>
    <page-footer />
    <back-top-button />
  </div>
</template>
<script setup lang="ts">
import NavBar from "#/widgets/layout/NavBar.vue"
import PageFooter from "#/widgets/layout/PageFooter.vue"
import MyInfo from "#/widgets/layout/MyInfo.vue"
import SearchBox from "#/widgets/layout/SearchBox.vue"
import TagCloud from "#/widgets/layout/TagCloud.vue"
import BackTopButton from "#/shared/ui/BackTopButton.vue"
import { ref, computed, watch, onMounted, onUnmounted } from "vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { localData } from "#/shared/lib/storage"
import { siteConfig } from "#/shared/config"
import { useTags } from "#/features/tag/model/useTags"

const darkModeStore = useDarkModeStore()
watch(
  () => darkModeStore.isDarkMode,
  isDark => {
    document.documentElement.classList.toggle("dark", isDark)
    document.body.classList.toggle("dark", isDark)
  },
  { immediate: true }
)
const screenWidth = ref<number>(document.body.clientWidth)
const { tags, loadTags } = useTags()
document.title = siteConfig.title
const contentColumnSpace = 30
const contentColumn = computed(() => (screenWidth.value >= 1200 ? 3 : screenWidth.value >= 768 ? 2 : 1))
const widgetWidth = computed(() => {
  const c = contentColumn.value
  let w = (screenWidth.value - contentColumnSpace * (c + 1)) / c
  if (c > 1) w -= w / 3
  return w - 40 + "px"
})
const contentWidth = computed(() => {
  const c = contentColumn.value
  let w = (screenWidth.value - contentColumnSpace * (c + 1)) / c
  if (c > 1) w += (w / 3) * (c - 1)
  return w + "px"
})
function handleScroll(e: Event): void {
  if (contentColumn.value <= 1) {
    document.querySelectorAll(".root-widgets").forEach(ele => ele.removeAttribute("style"))
    return
  }
  const top = 80
  const topOffset = 0
  const bottomOffset = 170
  const height = document.body.scrollHeight
  const target = e.target as Document
  const scrollTop = Number(target.documentElement.scrollTop) || Number(target.body.scrollTop)
  document.querySelectorAll(".root-widgets").forEach(ele => {
    if (ele.childNodes.length === 0) return
    const widgetHeight = ele.clientHeight + topOffset
    const offsetBottom = height - scrollTop - widgetHeight
    ele.removeAttribute("style")
    const htmlElement = ele as HTMLElement
    if (scrollTop <= top - topOffset) {
      htmlElement.style.position = "absolute"
      htmlElement.style.top = top + "px"
    } else if (offsetBottom <= bottomOffset) {
      htmlElement.style.position = "absolute"
      htmlElement.style.bottom = bottomOffset + "px"
      htmlElement.style.top = "auto"
    } else {
      htmlElement.style.position = "fixed"
      htmlElement.style.top = topOffset + "px"
    }
  })
}
onMounted(() => {
  window.addEventListener("scroll", handleScroll)
  window.onresize = () => {
    screenWidth.value = document.body.clientWidth
  }
  loadTags()
  if (localData.get("darkMode") === null) localData.set("darkMode", false)
  darkModeStore.setDarkMode(localData.get("darkMode") === "true" || localData.get("darkMode") === true)
})
onUnmounted(() => window.removeEventListener("scroll", handleScroll))
</script>
<style lang="scss">
.space-top {
  margin-top: 50px;
}
.router-view {
  position: relative;
  min-height: calc(100vh - 171px);
}
.router-view-1 {
  margin: 0 30px;
}
.router-view-2 {
  margin: 0 30px 0 auto;
}
.router-view-3 {
  margin: 0 auto;
}
.root-widgets {
  position: absolute;
  top: 80px;
}
.root-widgets-bottom {
  position: relative !important;
  top: 0 !important;
  margin: 0 30px !important;
}
.root-widgets-left {
  margin-left: 30px;
}
.root-widgets-right {
  margin-right: 30px;
  right: 0;
}
</style>
