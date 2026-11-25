<template>
  <div :class="{ dark: darkModeStore.isDarkMode }" class="container">
    <el-row align="middle" type="flex">
      <el-col :span="2"></el-col>
      <el-col :span="10"
        ><span>2018-{{ year }}</span
        ><i-custom-heart class="icon-heart" /><span>{{ siteConfig.author }}</span></el-col
      >
      <el-col :span="10" style="text-align: right"
        >Github
        <a :href="siteConfig.github?.link" :class="{ dark: darkModeStore.isDarkMode }" class="link">{{
          siteConfig.github?.name
        }}</a></el-col
      >
    </el-row>
  </div>
</template>
<script setup>
import { ref } from "vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { siteConfig } from "#/shared/config"
const darkModeStore = useDarkModeStore()
const year = ref(new Date().getFullYear())
</script>

<style scoped>
.container {
  position: relative;
  margin-top: 50px;
  padding-top: 20px;
  border-top: 1px solid rgba(152, 166, 173, 0.2);
  color: var(--text-color);
  background: var(--main-block-color);
  height: 100px;
  transition: background 0.8s;
  font-size: 16px;
}

.icon-heart {
  margin: 0 5px -2px 5px;
}
</style>
