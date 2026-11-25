<template>
  <card-template>
    <div class="cache-card">
      <h3>系统缓存</h3>
      <p>清理系统后台缓存，用于更新配置或修复数据同步问题。</p>
      <div class="actions">
        <el-button type="danger" :loading="loading" @click="handleClearCache">
          <template #icon>
            <Delete />
          </template>
          清理缓存
        </el-button>
      </div>
    </div>
  </card-template>
</template>

<script setup lang="ts">
import { client } from "#/shared/lib/api"
import { ElMessage } from "element-plus"
import Delete from "virtual:icons/mdi/delete"

defineOptions({ name: "CacheClearCard" })

const loading = ref(false)

async function handleClearCache() {
  loading.value = true
  const { error } = await (client as any).DELETE("/api/cache")
  if (!error) {
    ElMessage.success("缓存清理成功")
  } else {
    ElMessage.error(error?.message ?? "缓存清理失败")
  }
  loading.value = false
}
</script>

<style scoped>
.cache-card {
  padding: 10px;
}

h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

p {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
  line-height: 1.5;
}

.actions {
  display: flex;
  justify-content: flex-end;
}
</style>
