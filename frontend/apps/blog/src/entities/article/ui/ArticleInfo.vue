<template>
  <div class="info">
    <span class="date">{{ new Date(date).toLocaleString() }}</span>
    <span v-for="(tag, i) in tags" :key="i">
      <router-link :to="'/tag/' + tag.name + '/'" :class="{ dark: darkModeStore.isDarkMode }" class="tags link">{{
        tag.name
      }}</router-link>
      <span v-if="i !== tags.length - 1">|</span>
    </span>
  </div>
</template>
<script setup>
import { useDarkModeStore } from "#/shared/stores/darkMode"
const darkModeStore = useDarkModeStore()
defineProps({ date: String, tags: Array })
</script>
<style scoped>
.info,
.tags {
  color: var(--text-color);
  opacity: 0.7;
  font-size: 12px;
}
.tags {
  margin: 0 5px;
}
.date {
  margin-right: 5px;
}
</style>
