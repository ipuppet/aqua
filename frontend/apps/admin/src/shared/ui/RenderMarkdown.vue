<template>
  <MdPreview :modelValue="props.modelValue" theme="light" codeTheme="github" />
</template>

<script setup lang="ts">
import { MdPreview } from "md-editor-v3"
import "md-editor-v3/lib/preview.css"

const props = defineProps({
  modelValue: {
    type: String,
    default: ""
  }
})
</script>

<style scoped>
.md-preview {
  /* 确保过渡效果平滑 */
  transition: all 0.8s ease;
}
</style>
