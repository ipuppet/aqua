<template>
  <div class="page">
    <div class="header">
      <div class="title">会话管理</div>
      <div class="actions">
        <el-button
          :disabled="selectedSessions.length === 0"
          :loading="loading"
          @click="confirmBulkRevokeSelected"
          plain
        >
          批量吊销所选
        </el-button>
        <el-button :loading="loading" @click="confirmRevokeOtherSessions" plain>退出其他会话</el-button>
        <el-button :loading="loading" @click="refresh" plain>刷新</el-button>
      </div>
    </div>

    <el-table
      v-loading="loading"
      :data="sortedSessions"
      stripe
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="48" :selectable="isSelectableSession" />
      <el-table-column prop="deviceInfo" label="设备信息" min-width="220">
        <template #default="{ row }">
          <div class="device-cell">
            <span>{{ row.deviceInfo || "-" }}</span>
            <el-tag v-if="isCurrentSession(row)" size="small" effect="plain" type="primary">当前设备</el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{ row }">
          <span>{{ formatTime(row.createdAt) }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="sessionId" label="Session ID" min-width="240">
        <template #default="{ row }">
          <el-tooltip placement="top" :content="row.sessionId">
            <span class="mono ellipsis">{{ row.sessionId }}</span>
          </el-tooltip>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button
            link
            type="danger"
            size="small"
            :disabled="isCurrentSession(row)"
            @click="confirmRevoke(row.sessionId)"
          >
            吊销
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from "element-plus"
import { client, buildApiUrl } from "#/shared/lib/api"
import { ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY, getCookie, type components } from "@aqua/http"

type UserSession = components["schemas"]["aqua.passport.v1.UserSession"]

const loading = ref(false)
const sessions = ref<UserSession[]>([])
const selectedSessions = ref<UserSession[]>([])
const currentSessionId = computed(() => getCurrentSessionId())

const sortedSessions = computed(() => {
  const currentId = currentSessionId.value
  const list = [...sessions.value]
  const rank = (session: UserSession) => (currentId && (session.sessionId ?? "").trim() === currentId ? 0 : 1)
  return list.sort((a, b) => {
    const rankDiff = rank(a) - rank(b)
    if (rankDiff !== 0) return rankDiff
    const timeA = Date.parse(a.createdAt ?? "")
    const timeB = Date.parse(b.createdAt ?? "")
    if (Number.isNaN(timeA) || Number.isNaN(timeB)) return 0
    return timeB - timeA
  })
})

function parseSessionIdFromJwt(token: string | null | undefined): string | null {
  if (!token) return null
  const parts = token.split(".")
  if (parts.length < 2) return null
  const payloadPart = parts[1]
  if (!payloadPart) return null
  try {
    const payloadB64 = payloadPart.replace(/-/g, "+").replace(/_/g, "/")
    const padded = payloadB64.padEnd(Math.ceil(payloadB64.length / 4) * 4, "=")
    const json = atob(padded)
    const payload = JSON.parse(json) as { session_id?: unknown }
    const sessionId = typeof payload?.session_id === "string" ? payload.session_id : null
    return sessionId ? sessionId.trim() : null
  } catch {
    return null
  }
}

/** 会话 session_id 仅在 refresh_token 中可靠，优先读取 rt。 */
function getCurrentSessionId(): string | null {
  return parseSessionIdFromJwt(getCookie(REFRESH_TOKEN_KEY)) ?? parseSessionIdFromJwt(getCookie(ACCESS_TOKEN_KEY))
}

function formatTime(value?: string | null) {
  if (!value) return "-"
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toLocaleString()
}

function handleSelectionChange(rows: UserSession[]) {
  selectedSessions.value = Array.isArray(rows) ? rows.filter(isSelectableSession) : []
}

function isCurrentSession(session: UserSession) {
  const currentId = currentSessionId.value
  if (!currentId) return false
  return (session.sessionId ?? "").trim() === currentId
}

function isSelectableSession(session: UserSession) {
  return !isCurrentSession(session)
}

async function refresh() {
  loading.value = true
  const { data: res, error } = await client.GET("/api/sessions", {
    baseUrl: buildApiUrl("passport")
  })
  if (!error && res) {
    sessions.value = res.sessions ?? []
  } else {
    ElMessage.error(error?.message ?? "获取会话失败")
  }
  loading.value = false
}

async function revokeMany(sessionIds: Array<string | undefined | null>) {
  const unique = Array.from(new Set(sessionIds.map(s => (s ?? "").trim()).filter(Boolean)))
  if (unique.length === 0) return { ok: 0, fail: 0 }

  let ok = 0
  let fail = 0
  for (const sessionId of unique) {
    const { error } = await client.DELETE("/api/sessions/{session_id}", {
      baseUrl: buildApiUrl("passport"),
      params: { path: { session_id: sessionId } }
    })
    if (!error) {
      ok++
    } else {
      fail++
    }
  }
  return { ok, fail }
}

async function confirmBulkRevokeSelected() {
  const sessionIds = selectedSessions.value.filter(isSelectableSession).map(s => s.sessionId)
  if (sessionIds.length === 0) return

  try {
    await ElMessageBox.confirm(`将吊销所选 ${sessionIds.length} 个会话。是否继续？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
  } catch {
    return
  }

  loading.value = true
  try {
    const { ok, fail } = await revokeMany(sessionIds)
    if (fail === 0) ElMessage.success(`已吊销 ${ok} 个会话`)
    else ElMessage.warning(`已吊销 ${ok} 个，会话失败 ${fail} 个`)
    await refresh()
  } finally {
    loading.value = false
  }
}

async function confirmRevokeOtherSessions() {
  const currentId = getCurrentSessionId()
  if (!currentId) {
    ElMessage.error("无法识别当前会话（令牌中缺少 session_id），请先使用“批量吊销所选”")
    return
  }

  const other = sessions.value
    .filter(s => (s.sessionId ?? "").trim() && (s.sessionId ?? "").trim() !== currentId)
    .map(s => s.sessionId)
  if (other.length === 0) {
    ElMessage.info("没有可退出的其他会话")
    return
  }

  try {
    await ElMessageBox.confirm(`将退出其他 ${other.length} 个会话（保留当前会话）。是否继续？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
  } catch {
    return
  }

  loading.value = true
  try {
    const { ok, fail } = await revokeMany(other)
    if (fail === 0) ElMessage.success(`已退出其他 ${ok} 个会话`)
    else ElMessage.warning(`已退出 ${ok} 个，会话失败 ${fail} 个`)
    await refresh()
  } finally {
    loading.value = false
  }
}

async function confirmRevoke(sessionId: string) {
  const trimmed = (sessionId ?? "").trim()
  if (!trimmed) return
  if (trimmed === currentSessionId.value) {
    ElMessage.info("当前会话不可在此处吊销，请使用退出登录")
    return
  }

  try {
    await ElMessageBox.confirm("吊销该会话后，该设备将需要重新登录。是否继续？", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
  } catch {
    return
  }

  loading.value = true
  const { error } = await client.DELETE("/api/sessions/{session_id}", {
    baseUrl: buildApiUrl("passport"),
    params: { path: { session_id: trimmed } }
  })

  if (!error) {
    ElMessage.success("已吊销")
    await refresh()
  } else {
    ElMessage.error(error?.message ?? "吊销失败")
  }
  loading.value = false
}

onMounted(() => refresh())
</script>

<style scoped>
.page {
  padding: 4px 2px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.actions {
  display: flex;
  gap: 8px;
}

.title {
  font-size: 14px;
  font-weight: 600;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

.ellipsis {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.device-cell {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
</style>
