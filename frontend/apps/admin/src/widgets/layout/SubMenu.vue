<template>
  <el-sub-menu :index="props.menu.name">
    <template #title>
      <el-icon>
        <component :is="iconsStore[props.menu.icon]" />
      </el-icon>
      <span>{{ props.menu.name }}</span>
    </template>
    <el-menu-item-group v-for="(item, index) of props.menu.children" :key="index">
      <sub-menu v-if="item.children" :menu="item" :key="index" />
      <el-menu-item v-else :key="item.url" :index="item.url">
        <el-icon>
          <component :is="iconsStore[item.icon]" />
        </el-icon>
        <span>{{ item.name }}</span>
      </el-menu-item>
    </el-menu-item-group>
  </el-sub-menu>
</template>

<script setup lang="ts">
// 使用 Pinia store
import { useIconsStore } from "#/shared/stores/icons"

interface MenuItem {
  name: string
  icon: string
  children: Array<{
    name: string
    icon: string
    router: string
  }>
}

interface Props {
  menu: MenuItem
}

const props = defineProps<Props>()
const iconsStore = useIconsStore()
</script>

<style scoped></style>
