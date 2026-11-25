<template>
  <div>
    <loading-ball-atom v-show="!isLoaded" class="loading" />
    <div
      v-for="article in page.articles"
      :key="article.aid"
      :class="{ dark: darkModeStore.isDarkMode }"
      class="router-view-container"
    >
      <div class="margin-box">
        <router-link :to="'/detail/' + article.aid + '/'" :class="{ dark: darkModeStore.isDarkMode }" class="title link"
          ><span v-html="article.title"></span
        ></router-link>
        <div class="summary"><render-markdown v-model="article.summary" /></div>
        <article-info :date="article.createdAt" :tags="article.tags" class="info" />
      </div>
    </div>
    <div class="block">
      <el-pagination
        :hide-on-single-page="true"
        :page-size="page.pageSize"
        :pager-count="5"
        :total="page.total"
        :disabled="!isLoaded"
        background
        layout="prev, pager, next"
        @current-change="p => $emit('current-change', p)"
        @prev-click="$emit('prev-click')"
        @next-click="$emit('next-click')"
      />
    </div>
  </div>
</template>
<script setup>
import ArticleInfo from "#/entities/article/ui/ArticleInfo.vue"
import LoadingBallAtom from "#/shared/ui/LoadingBallAtom.vue"
import RenderMarkdown from "#/shared/ui/RenderMarkdown.vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
const darkModeStore = useDarkModeStore()
defineProps({ page: Object, isLoaded: Boolean })
</script>
<style scoped>
.router-view-container {
  overflow-x: hidden;
  margin-bottom: 30px;
  padding: 0;
}
.margin-box {
  margin: 30px;
  width: calc(100% - 60px);
}
.title {
  font-size: 26px;
}
.summary {
  margin-top: 10px;
}
.block {
  text-align: center;
}
</style>
