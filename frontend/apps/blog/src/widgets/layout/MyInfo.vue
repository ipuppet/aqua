<template>
  <div class="card-container">
    <div class="image"><img :alt="siteConfig.author" :src="siteConfig.authorImage" /></div>
    <div class="my-info">
      <span>{{ siteConfig.author }}</span>
      <p class="sentence">{{ siteConfig.description }}</p>
    </div>
  </div>
</template>
<script setup>
import { siteConfig } from "#/shared/config"
</script>
<style scoped>
.card-container {
  text-align: center;
}
.image {
  width: 150px;
  margin: 10px auto 0;
  overflow: hidden;
}
.image img {
  width: 100%;
  transition: transform 0.3s ease;
}
.my-info {
  margin-bottom: 20px;
}
.sentence {
  font-size: 12px;
  color: var(--text-color);
  opacity: 0.7;
  transition: color 0.8s ease;
}
</style>
