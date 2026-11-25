<template>
  <div class="card-container" :style="spaceStyle">
    <div :id="id" class="box">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: "CardTemplate" })

const generateId = (length = 10) => {
  const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  let result = ""
  for (let i = 0; i < length; i++) {
    result += characters[Math.floor(Math.random() * characters.length)]
  }
  return result
}

const id = `card_${generateId()}`
const spaceStyle = ref("")

onBeforeMount(() => {
  const parent = getCurrentInstance()?.proxy?.$parent as { waterfall?: boolean; space?: number } | undefined
  if (parent?.waterfall && spaceStyle.value === "") {
    spaceStyle.value = `margin-bottom:${parent.space}px`
  }
})

const getId = () => id

defineExpose({ getId })
</script>

<style scoped>
.card-container {
  width: calc(100% - 20px);
  padding: 10px;
  break-inside: avoid;
}

.box {
  width: calc(100% - 40px);
  padding: 20px;
  box-shadow:
    0 0 10px rgba(0, 2, 4, 0.06),
    0 0 1px rgba(0, 2, 4, 0.11);
  border-radius: 10px;
  transition: background 0.8s;
  background: #ffffff;
}
</style>
