<template>
  <div class="menu" :class="[{ active }]" v-loading="!loaded">
    <div @click="menuControl" class="pointer menu-control-button">
      <Fold v-if="active" />
      <Unfold v-else />
    </div>
    <div class="menu-logo">
      <span class="menu-logo-title">{{ siteConfig.title }}</span>
    </div>
    <div class="menu-header">
      <div class="menu-navigation">NAVIGATION</div>
      <div class="pointer lock" v-show="active" @click="lockMenu">
        <el-tooltip effect="dark" :content="uniqueOpenedInformation" placement="top">
          <Lock v-if="uniqueOpened" />
          <Unlock v-else />
        </el-tooltip>
      </div>
    </div>

    <el-menu
      class="el-menu-vertical"
      :default-active="$route.path"
      background-color="rgba(0,0,0,0)"
      text-color="#c2cfd6"
      active-text-color="#ffffff"
      :router="true"
      :collapse="!active"
      :unique-opened="uniqueOpened"
    >
      <template v-for="(item, index) of menus">
        <sub-menu v-if="item.children" :menu="item" :key="index" />
        <el-menu-item v-else :key="item.url" :index="item.url">
          <el-icon>
            <component :is="iconsStore[item.icon]" />
          </el-icon>
          <span>{{ item.name }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import Fold from "virtual:icons/mdi/menu-open"
import Unfold from "virtual:icons/mdi/menu"
import Lock from "virtual:icons/mdi/lock"
import Unlock from "virtual:icons/mdi/lock-open"
import { useMenuStore } from "#/shared/stores/menu"
import { useIconsStore } from "#/shared/stores/icons"
import { siteConfig } from "#/shared/config/site"
import { storeToRefs } from "pinia"
import { localData } from "#/shared/lib/storage"

const menuStore = useMenuStore()
const iconsStore = useIconsStore()
const { active, loaded, menus } = storeToRefs(menuStore)

const uniqueOpened = ref(true)
const uniqueOpenedInformation = `当锁关闭时，只能一个菜单是开启状态。\n当锁开启时，可同时开启多个菜单。`

function menuControl(states?: boolean) {
  if (typeof states !== "boolean") states = active.value !== true
  menuStore.setMenuState(states)
  localData.set("menu.menuState", states)
}
function lockMenu() {
  uniqueOpened.value = !uniqueOpened.value
  localData.set("menu.uniqueOpened", uniqueOpened.value)
}

onMounted(() => {
  uniqueOpened.value = localData.get("menu.uniqueOpened") === "true"
  const menuState = localData.get("menu.menuState")
  if (menuState !== null) menuStore.setMenuState(menuState === "true")
})
</script>

<style scoped>
.el-menu-vertical:not(.el-menu--collapse) {
  width: var(--menu-width-active);
  min-height: 400px;
}

.menu {
  /* position: fixed; removed for flex layout */
  z-index: 1501;
  height: 100%;
  width: var(--menu-width);
  background: linear-gradient(135deg, #8f75da 0, #727cf5 60%);
  transition: 0.3s ease-in-out;
  font-size: 16px;
  overflow-x: hidden;
  user-select: none;
  color: #c2cfd6;
}

.menu-control-button {
  position: relative;
  margin: 10px 20px 0 auto;
  width: 20px;
  height: 20px;
  line-height: 20px;
  text-align: center;
}

/*logo部分*/
.menu-logo {
  position: relative;
  width: 100%;
  height: 0;
  opacity: 0;
  line-height: 90px;
  text-align: center;
  transition: 0.3s;
}

.menu-logo-title {
  font-weight: lighter;
  font-size: 24px;
  color: #f5f5f5;
}

.menu-header {
  position: relative;
  height: 0;
  line-height: 0;
  overflow: hidden;
  margin-bottom: 20px;
}

.menu-navigation {
  float: left;
  margin: 0 auto 0 30px;
  height: 20px;
  font-size: 14px;
}

.lock {
  float: right;
  margin: 0 15px 0 auto;
  height: 20px;
}

/*开启菜单*/
.menu.active {
  width: var(--menu-width-active);
}

.menu.active .menu-logo {
  height: 90px;
  opacity: 1;
}

.menu.active .menu-header {
  height: 20px;
  line-height: 20px;
}
</style>
