<template>
  <div
    :style="`
        width: calc(100% - ${props.space * 2}px);
        margin: ${props.space}px;
        columns: ${props.maxColumn !== undefined ? autoColumn : props.column};
        column-gap: ${props.space}px;
    `"
  >
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
interface Props {
  space?: number
  column?: number
  maxColumn?: number
  minColumn?: number
}

const props = withDefaults(defineProps<Props>(), {
  space: 30,
  column: 3,
  maxColumn: undefined,
  minColumn: 1
})

const screenWidth = ref(document.body.clientWidth)

const resize = () => {
  screenWidth.value = document.body.clientWidth
}

onMounted(() => {
  screenWidth.value = document.body.clientWidth
  window.onresize = resize
})

const autoColumn = computed(() => {
  const width = screenWidth.value
  let colNum = props.maxColumn
  if (width >= 1920) {
    // 保持最大
  } else if (width >= 1200) {
    colNum = props.maxColumn! - 1
  } else if (width >= 992) {
    colNum = props.maxColumn! - (colNum! <= 3 ? 2 : 1)
  } else if (width >= 768) {
    colNum = props.maxColumn! - 2
  } else {
    colNum = props.minColumn
  }
  if (colNum! < props.minColumn) colNum = props.minColumn
  return colNum
})

// 暴露给模板
defineExpose({
  autoColumn
})
</script>

<style scoped></style>
