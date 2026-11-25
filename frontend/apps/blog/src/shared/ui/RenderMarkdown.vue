<template>
  <MdPreview
    :modelValue="markdownContent"
    :theme="darkModeStore.isDarkMode ? 'dark' : 'light'"
    :codeTheme="darkModeStore.isDarkMode ? 'github-dark' : 'github'"
    class="markdown-preview"
  />
</template>

<script setup>
import { ref, watch } from "vue"
import { MdPreview } from "md-editor-v3"
import { useDarkModeStore } from "#/shared/stores/darkMode"

const darkModeStore = useDarkModeStore()
const props = defineProps({ markdown: String })
const markdownContent = ref(props.markdown ?? "")

watch(
  () => props.markdown,
  newValue => {
    markdownContent.value = newValue ?? ""
  }
)
</script>

<style scoped>
.markdown-preview {
  transition: all 0.8s ease;
}
</style>
