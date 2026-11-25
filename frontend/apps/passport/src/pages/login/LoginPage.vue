<template>
  <el-row class="login-body">
    <el-col :span="12">
      <login-left-panel class="section hidden-xs-only" />
    </el-col>
    <el-col :span="12" :xs="24">
      <login-form-panel class="section" />
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import LoginLeftPanel from "#/pages/login/components/LoginLeftPanel.vue"
import LoginFormPanel from "#/pages/login/components/LoginFormPanel.vue"
</script>

<style>
.login-body {
  height: 80vh;
  margin-top: 15vh;
}

.section {
  text-align: center;
  padding-left: 20px;
  padding-right: 20px;
}

.pointer {
  cursor: pointer;
}
</style>
