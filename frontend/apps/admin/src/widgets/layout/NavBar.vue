<template>
  <div :class="['nav', { active: menuStore.active }]">
    <functional-button class="functional-button" />
    <nav-search v-if="navStore.showSearch" class="nav-search" />
    <nav-user />
  </div>
</template>

<script setup lang="ts">
import { useMenuStore } from "#/shared/stores/menu"
import { useNavStore } from "#/shared/stores/nav"
import NavUser from "#/widgets/layout/NavUser.vue"
import FunctionalButton from "#/widgets/layout/FunctionalButton.vue"
import NavSearch from "#/widgets/layout/NavSearch.vue"

const menuStore = useMenuStore()
const navStore = useNavStore()
</script>

<style scoped>
.nav {
  height: var(--nav-height);
  font-size: 14px;
  background-color: #ffffff;
  box-shadow: 0 0 35px 0 rgba(154, 161, 171, 0.15);
}

.functional-button {
  float: left;
  height: var(--nav-height);
  overflow: hidden;
}

.nav-search {
  float: left;
  height: var(--nav-height);
  display: flex;
  align-items: center;
  margin-left: 20px;
}
</style>
