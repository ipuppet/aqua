<template>
  <transition
    enter-active-class="animate__animated animate__fadeInRight animate__faster"
    leave-active-class="animate__animated animate__zoomOutRight"
  >
    <div v-show="isShowButton" class="back-top-button" @click="backToTop">
      <i-tabler-chevron-up :class="{ dark: darkModeStore.isDarkMode }" />
    </div>
  </transition>
</template>

<script setup>
import { ref } from "vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
const darkModeStore = useDarkModeStore()
const isShowButton = ref(false)
window.addEventListener(
  "scroll",
  () => {
    const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
    isShowButton.value = scrollTop > 300
  },
  true
)
function backToTop() {
  scrollToTop(0, 200)
}
function scrollToTop(number = 0, time) {
  if (!time) {
    document.body.scrollTop = document.documentElement.scrollTop = number
    return
  }
  const spacingTime = 20
  let spacingInex = time / spacingTime
  let nowTop = document.body.scrollTop + document.documentElement.scrollTop
  const everTop = (number - nowTop) / spacingInex
  const scrollTimer = setInterval(() => {
    if (spacingInex > 0) {
      spacingInex--
      scrollToTop((nowTop += everTop))
    } else {
      clearInterval(scrollTimer)
    }
  }, spacingTime)
}
</script>

<style scoped>
.back-top-button {
  position: fixed;
  z-index: 1000;
  width: 30px;
  height: 30px;
  right: 20px;
  bottom: 45px;
  font-size: 26px;
  text-align: center;
  line-height: 35px;
  color: var(--text-color);
  cursor: pointer;
}
</style>
