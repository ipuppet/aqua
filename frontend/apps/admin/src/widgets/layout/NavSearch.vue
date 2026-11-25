<template>
  <div class="nav-search-container">
    <el-input
      clearable
      @clear="handleClear"
      @keyup.enter="handleSearch"
      placeholder="搜索"
      v-model="searchKeyword"
      style="width: 300px"
    >
      <template #append>
        <el-button @click="handleSearch">
          <Search />
        </el-button>
      </template>
    </el-input>
  </div>
</template>

<script setup lang="ts">
import Search from "virtual:icons/mdi/magnify"
import { useNavStore } from "#/shared/stores/nav"
import { storeToRefs } from "pinia"

const navStore = useNavStore()
const { searchKeyword } = storeToRefs(navStore)

function handleSearch() {
  navStore.triggerSearch()
}

function handleClear() {
  navStore.clearSearch()
}
</script>

<style scoped>
.nav-search-container {
  display: flex;
  align-items: center;
}
</style>
