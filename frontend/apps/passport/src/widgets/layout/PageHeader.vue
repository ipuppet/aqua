<template>
  <el-row align="middle" class="header" type="flex">
    <el-col :span="12" class="logo">
      {{ siteConfig.title }}
    </el-col>
    <el-col :span="12" class="text-right" v-if="route.path !== registerUrl">
      <span>还没有账户？</span>
      <router-link :to="registerUrl">点此注册</router-link>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router"
import { siteConfig } from "#/shared/config/site"

const route = useRoute()
const registerUrl = "/register"
</script>

<style scoped>
.header {
  padding: 2rem 1.5rem 1.5rem 1.5rem;
}

.text-right {
  text-align: right;
}

.logo {
  font-size: 26px;
}
</style>
