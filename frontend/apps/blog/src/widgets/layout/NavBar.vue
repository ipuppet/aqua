<template>
  <div class="container">
    <div class="nav" :class="{ dark: darkModeStore.isDarkMode }">
      <el-row justify="space-between">
        <el-col :span="16">
          <router-link :class="{ dark: darkModeStore.isDarkMode }" class="link" to="/"
            ><div class="nav-button"><i-tabler-home /><br />Home</div></router-link
          >
          <router-link :class="{ dark: darkModeStore.isDarkMode }" class="link" to="/archive/"
            ><div class="nav-button"><i-tabler-archive /><br />Archive</div></router-link
          >
          <router-link :class="{ dark: darkModeStore.isDarkMode }" class="link" to="/about/"
            ><div class="nav-button"><i-tabler-info-circle /><br />About</div></router-link
          >
        </el-col>
        <el-col :span="1"><i-tabler-sun :class="displayMode" class="dark-button" @click="toggleDarkMode" /></el-col>
      </el-row>
    </div>
  </div>
</template>
<script setup>
import { computed, onMounted, onUnmounted } from "vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
const darkModeStore = useDarkModeStore()
const displayMode = computed(() => (darkModeStore.isDarkMode ? "rotate" : "re-rotate"))
function toggleDarkMode() {
  darkModeStore.setDarkMode(!darkModeStore.isDarkMode)
}
function darkModeListener(e) {
  darkModeStore.setDarkMode(e.matches)
}
onMounted(() => {
  const m = window.matchMedia("(prefers-color-scheme: dark)")
  darkModeListener(m)
  m.addEventListener("change", darkModeListener)
})
onUnmounted(() => {
  window.matchMedia("(prefers-color-scheme: dark)").removeEventListener("change", darkModeListener)
})
</script>
<style scoped>
.container {
  width: 100%;
  height: 60px;
  overflow: hidden;
}
.nav {
  position: absolute;
  width: calc(100% - 60px);
  padding: 0 30px;
  background: var(--main-block-color);
  box-shadow:
    0 4px 10px rgba(0, 2, 4, 0.06),
    0 0 1px rgba(0, 2, 4, 0.11);
  transition: background 0.8s;
}
.nav-button {
  width: 50px;
  height: 34px;
  display: inline;
  text-align: center;
  margin: 10px 10px 10px 0;
  padding: 3px 0;
  float: left;
  font-size: 12px;
  color: var(--text-color);
}
.dark-button {
  float: right;
  margin: 15px 0;
  font-size: 25px;
  color: var(--text-color);
  cursor: pointer;
}
</style>
