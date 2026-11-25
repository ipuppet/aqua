<template>
  <div class="app-shell" v-if="userStore.isLogin">
    <main-menu v-if="!hideMenu" class="sidebar-container" />
    <div class="main-container">
      <div class="header-container" v-if="!hideMenu">
        <nav-bar />
        <tabs-bar />
      </div>
      <div class="header-container" v-else>
        <nav-bar />
      </div>
      <div class="app-main">
        <div class="app-content-box">
          <router-view v-slot="{ Component }">
            <keep-alive :include="tabsStore.cachedViews">
              <component :is="Component" />
            </keep-alive>
          </router-view>
        </div>
        <page-footer class="page-footer" />
        <el-backtop target=".app-main" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from "#/shared/stores/user"
import { useRoute } from "vue-router"
import TabsBar from "#/widgets/layout/TabsBar.vue"
import MainMenu from "#/widgets/layout/MainMenu.vue"
import NavBar from "#/widgets/layout/NavBar.vue"
import PageFooter from "#/widgets/layout/PageFooter.vue"
import { useTabsStore } from "#/shared/stores/tabs"

const userStore = useUserStore()
const tabsStore = useTabsStore()
const route = useRoute()
const hideMenu = computed(() => Boolean((route.meta as any)?.hideMenu))
</script>

<style scoped>
.app-shell {
  display: flex;
  width: 100%;
  height: 100%;
}

.sidebar-container {
  flex-shrink: 0;
  transition: width 0.3s;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  height: 100%;
  overflow: hidden;
  background-color: #ffffff;
}

.header-container {
  flex-shrink: 0;
  background-color: #fff;
  z-index: 9;
}

.app-main {
  flex: 1;
  overflow-y: auto;
  position: relative;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.app-content-box {
  flex: 1;
  width: 100%;
  padding: var(--app-content-padding);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.page-footer {
  padding: 0;
  margin-top: auto;
  flex-shrink: 0;
}
</style>
