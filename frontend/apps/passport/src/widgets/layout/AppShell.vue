<template>
  <div class="app-shell">
    <page-header />
    <router-view />
    <page-footer />
  </div>
</template>

<script setup lang="ts">
import PageHeader from "#/widgets/layout/PageHeader.vue"
import PageFooter from "#/widgets/layout/PageFooter.vue"
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
}
</style>
