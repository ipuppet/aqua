<template>
  <div class="card-container">
    <div class="title">Tags</div>
    <router-link
      v-for="tag in tags"
      :key="tag.id"
      :to="'/tag/' + tag.name + '/'"
      :class="{ dark: darkModeStore.isDarkMode }"
      class="link tag"
    >
      {{ tag.name }}
    </router-link>
  </div>
</template>
<script setup>
import { useDarkModeStore } from "#/shared/stores/darkMode"
const darkModeStore = useDarkModeStore()
defineProps({ tags: Array })
</script>
<style scoped>
.card-container {
  text-align: center;
}
.title {
  margin-bottom: 10px;
}
.tag {
  margin-left: 5px;
}
</style>
