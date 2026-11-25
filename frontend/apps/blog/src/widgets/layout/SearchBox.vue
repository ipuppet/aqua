<template>
  <div class="card-container">
    <div class="search-box">
      <input
        id="kw"
        v-model="kw"
        :class="{ dark: darkModeStore.isDarkMode }"
        class="input-box"
        placeholder="Search..."
        type="text"
        @keyup.enter="search"
      />
      <button :class="{ dark: darkModeStore.isDarkMode }" class="search-button" @click="search">
        <i-custom-search />
      </button>
    </div>
  </div>
</template>
<script setup>
import { ref, watch } from "vue"
import { useRouter, useRoute } from "vue-router"
import { useDarkModeStore } from "#/shared/stores/darkMode"
const router = useRouter()
const route = useRoute()
const darkModeStore = useDarkModeStore()
const kw = ref("")
watch(
  () => route.path,
  () => {
    kw.value = ""
  }
)
function search() {
  if (kw.value === "") {
    alert("Please enter a keyword")
  } else {
    router.push(`/search/30/${kw.value}/`)
  }
}
</script>
<style scoped>
.search-box {
  height: 36px;
}
.input-box {
  background: rgba(0, 0, 0, 0);
  height: 100%;
  width: calc(70% - 20px);
  margin-left: 15%;
  padding: 0 10px;
  border-radius: 18px;
  border: 1px solid var(--input-border-color);
  font-size: 16px;
  color: var(--text-color);
}
.search-button {
  cursor: pointer;
  background: none;
  border: none;
  position: relative;
  left: -25px;
  opacity: 0.5;
  color: var(--text-color);
}
</style>
