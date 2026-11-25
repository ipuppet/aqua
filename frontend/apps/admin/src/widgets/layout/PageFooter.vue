<template>
  <div class="footer">
    <el-row align="middle" type="flex">
      <el-col :span="12">
        <span>2018-{{ year }}</span>
        <el-icon class="icon-heart">
          <Heart />
        </el-icon>
        <span class="bottom-link">{{ siteConfig.title }}</span>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import Heart from "virtual:icons/mdi/heart"
import { siteConfig } from "#/shared/config/site"

const year = ref(new Date().getFullYear())
</script>

<style scoped>
.footer {
  height: var(--footer-height);
  display: flex;
  align-items: center;
  padding: 0 20px; /* Horizontal padding only, let height handle vertical */
  border-top: 1px solid rgba(152, 166, 173, 0.2);
  color: #98a6ad;
  background-color: #fafbfe;
  transition: 0.3s;
  box-sizing: border-box;
  width: 100%;
}

/* Ensure the inner el-row takes full width */
.footer :deep(.el-row) {
  width: 100%;
}

.footer span {
  margin-left: 10px;
}

.icon-heart {
  margin: 0 5px 0 5px;
  vertical-align: middle;
}

.bottom-link {
  vertical-align: middle;
}
</style>
