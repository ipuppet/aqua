<template>
  <div class="container">
    <el-button
      v-for="(button, index) in buttons"
      :key="index"
      :size="button.size"
      :icon="button.icon"
      @click="button.click"
      >{{ button.text }}
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { useNavStore } from "#/shared/stores/nav"
import { storeToRefs } from "pinia"

const navStore = useNavStore()
const { buttons } = storeToRefs(navStore)
</script>

<style scoped>
.container {
  padding-left: 20px;
  line-height: var(--nav-height);
}
</style>
