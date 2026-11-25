<template>
  <waterfall-layout :max-column="3">
    <hello-card />
    <cache-clear-card />
    <blog-access-chart-card />
  </waterfall-layout>
</template>

<script setup lang="ts">
import HelloCard from "#/pages/home/components/HelloCard.vue"
import BlogAccessChartCard from "#/pages/home/components/BlogAccessChartCard.vue"
import CacheClearCard from "#/pages/home/components/CacheClearCard.vue"

defineOptions({ name: "HomePage" })
</script>

<style scoped>
.card {
  margin-bottom: 30px;
}
</style>
