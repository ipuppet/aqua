<template>
  <div class="tags-view-container" :class="{ 'tags-view-active': menuStore.active }">
    <el-scrollbar ref="scrollContainer" :vertical="false" class="scroll-container" @wheel.prevent="handleScroll">
      <router-link
        v-for="tag in visitedViews"
        ref="tag"
        :key="tag.path"
        :class="isActive(tag) ? 'active' : ''"
        :to="{ path: tag.path, query: tag.query, fullPath: tag.fullPath }"
        class="tags-view-item"
        @click.middle="!isAffix(tag) ? closeSelectedTag(tag) : ''"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        {{ getTabTitle(tag) }}
        <span v-if="!isAffix(tag)" class="el-icon-close" @click.prevent.stop="closeSelectedTag(tag)">
          <component :is="iconsStore.Close" style="width: 10px; height: 10px" />
        </span>
      </router-link>
    </el-scrollbar>
    <ul v-show="visible" :style="{ left: left + 'px', top: top + 'px' }" class="contextmenu">
      <li v-if="!isAffix(selectedTag)" @click="closeSelectedTag(selectedTag)">关闭</li>
      <li @click="closeOthersTags">关闭其他</li>
      <li @click="closeAllTags(selectedTag)">关闭所有</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from "vue-router"
import type { RouteLocationNormalized } from "vue-router"
import { useTabsStore } from "#/shared/stores/tabs"
import { useIconsStore } from "#/shared/stores/icons"
import { useMenuStore } from "#/shared/stores/menu"
import { storeToRefs } from "pinia"

const route = useRoute()
const router = useRouter()
const tabsStore = useTabsStore()
const iconsStore = useIconsStore()
const menuStore = useMenuStore()
const { visitedViews } = storeToRefs(tabsStore)
const { menus } = storeToRefs(menuStore)

const visible = ref(false)
const top = ref(0)
const left = ref(0)
const selectedTag = ref<RouteLocationNormalized>({} as RouteLocationNormalized)

const isActive = (view: RouteLocationNormalized) => {
  return view.path === route.path
}

const isAffix = (tag: RouteLocationNormalized) => {
  return tag.meta && tag.meta.affix
}

const findMenuName = (menus: any[], path: string): string | undefined => {
  for (const menu of menus) {
    if (menu.url === path) {
      return menu.name
    }
    if (menu.children && menu.children.length > 0) {
      const name = findMenuName(menu.children, path)
      if (name) return name
    }
  }
  return undefined
}

const getTabTitle = (view: RouteLocationNormalized) => {
  const title = findMenuName(menus.value, view.path)
  if (title) {
    return title
  }
  return (view.meta.title as string) || "no-name"
}

const addTags = () => {
  const { name } = route
  if (name) {
    tabsStore.addView(route)
  }
  return false
}

const initTags = () => {
  const routes = router.getRoutes()
  const affixTags = routes.filter(r => r.meta && r.meta.affix)
  for (const tag of affixTags) {
    // Must have tag name
    if (tag.name) {
      tabsStore.addVisitedView(tag as unknown as RouteLocationNormalized)
    }
  }
}

const closeSelectedTag = (view: RouteLocationNormalized) => {
  if (isAffix(view)) return
  tabsStore.delView(view).then(({ visitedViews }: any) => {
    if (isActive(view)) {
      toLastView(visitedViews, view)
    }
  })
}

const closeOthersTags = () => {
  router.push(selectedTag.value)
  tabsStore.delOthersViews(selectedTag.value).then(() => {
    moveToCurrentTag()
  })
}

const closeAllTags = (view: RouteLocationNormalized) => {
  tabsStore.delAllViews().then(({ visitedViews }: any) => {
    if (visitedViews.find((v: any) => v.path === view.path)) {
      return
    }
    toLastView(visitedViews, view)
  })
}

const toLastView = (visitedViews: RouteLocationNormalized[], view: RouteLocationNormalized) => {
  const latestView = visitedViews.slice(-1)[0]
  if (latestView) {
    router.push(latestView.fullPath || latestView.path)
  } else {
    if (view.name === "首页") {
      router.replace({ path: "/redirect" + view.fullPath })
    } else {
      router.push("/")
    }
  }
}

const scrollContainer = ref()

const handleScroll = (e: WheelEvent) => {
  const eventDelta = -e.deltaY
  const $scrollWrapper = scrollContainer.value?.$refs.wrap
  if ($scrollWrapper) {
    $scrollWrapper.scrollLeft = $scrollWrapper.scrollLeft + eventDelta / 4
  }
}

const openMenu = (tag: RouteLocationNormalized, e: MouseEvent) => {
  const menuMinWidth = 105
  // use the container's rect for relative positioning
  const containerRect = (scrollContainer.value?.$el.parentNode as HTMLElement).getBoundingClientRect()

  const offsetLeft = containerRect.left
  const offsetWidth = containerRect.width
  const maxLeft = offsetWidth - menuMinWidth // left boundary
  const left1 = e.clientX - offsetLeft + 15 // 15: margin right

  if (left1 > maxLeft) {
    left.value = maxLeft
  } else {
    left.value = left1
  }

  top.value = e.clientY - containerRect.top
  visible.value = true
  selectedTag.value = tag
}

const closeMenu = () => {
  visible.value = false
}

const moveToCurrentTag = () => {
  // scroll logic if needed, simplify for now
  /*
  const tags = document.querySelectorAll(".tags-view-item")
  nextTick(() => {
    for (const tag of tags as any) {
      if (tag.getAttribute('href') === route.path) {
        // scrollContainer.value.moveToTarget(tag)
        // break
      }
    }
  })
  */
}

watch(
  () => route.path,
  () => {
    addTags()
    moveToCurrentTag()
  }
)

watch(visible, (value: boolean) => {
  if (value) {
    document.body.addEventListener("click", closeMenu)
  } else {
    document.body.removeEventListener("click", closeMenu)
  }
})

onMounted(() => {
  initTags()
  addTags()
})
</script>

<style scoped>
.tags-view-container {
  height: var(--tags-height);
  width: 100%;
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  display: flex;
  align-items: center;
  padding: 0 10px;
  box-sizing: border-box;
  position: relative;
  z-index: 10;
}

.scroll-container {
  white-space: nowrap;
  position: relative;
  overflow: hidden;
  width: 100%;
  display: flex;
  align-items: center;
}

/* Ensure the scrollbar content is aligned */
:deep(.el-scrollbar__view) {
  display: inline-flex;
  align-items: center;
  height: 100%;
}

.tags-view-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  cursor: pointer;
  height: 36px;
  line-height: 36px;
  border: 1px solid transparent;
  border-radius: 12px;
  color: #555;
  background: #f7f8fa;
  padding: 0 16px;
  font-size: 13px;
  margin-right: 8px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
  user-select: none;
  font-weight: 500;
}

.tags-view-item:first-of-type {
  margin-left: 5px;
}

.tags-view-item:last-of-type {
  margin-right: 15px;
}

.tags-view-item:hover {
  background: #eef0f5;
  color: #333;
  transform: translateY(-1px);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}

.tags-view-item.active {
  background-color: #ffffff;
  color: #42b983;
  border-color: rgba(66, 185, 131, 0.3);
  box-shadow: 0 4px 12px rgba(66, 185, 131, 0.15);
  font-weight: 600;
}

.tags-view-item.active::before {
  content: "";
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 0; /* Removing the bottom line idea, sticking to card style */
}

/* Close icon style */
.tags-view-item .el-icon-close {
  width: 16px;
  height: 16px;
  margin-left: 8px;
  border-radius: 50%;
  text-align: center;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
  transform-origin: 100% 50%;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  font-size: 10px;
  color: #aeb4bf;
}

.tags-view-item .el-icon-close:hover {
  background-color: #ff4d4f;
  color: #fff;
  transform: scale(1.1);
}

.contextmenu {
  margin: 0;
  background: #fff;
  z-index: 3000;
  position: absolute;
  list-style-type: none;
  padding: 5px 0;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  color: #333;
  box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);
}

.contextmenu li {
  margin: 0;
  padding: 7px 16px;
  cursor: pointer;
}

.contextmenu li:hover {
  background: #eee;
}
</style>
