<template>
  <div class="router-view-container">
    <render-markdown v-model="about" class="content" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue"
import RenderMarkdown from "#/shared/ui/RenderMarkdown.vue"
import { client } from "#/shared/lib/api"
import { sessionData } from "#/shared/lib/storage"
const about = ref("")
onMounted(async () => {
  document.title = "About"
  const cachedAbout = sessionData.get("about")
  if (!cachedAbout) {
    const { data, error } = await client.GET("/api/about")
    if (!error && data) {
      about.value = data.content ?? ""
      sessionData.set("about", data.content ?? "")
    }
  } else {
    about.value = cachedAbout
  }
})
</script>
<style scoped>
.content {
  padding: 0 30px;
  overflow-x: hidden;
  word-break: break-word;
}

.content :deep(*) {
  max-width: 100%;
  box-sizing: border-box;
}

.content :deep(pre) {
  overflow-x: auto;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
}

.content :deep(table) {
  display: block;
  overflow-x: auto;
}
</style>
