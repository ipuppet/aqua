<template>
  <div>
    <el-table v-loading="loading" :data="users" stripe style="width: 100%">
      <el-table-column v-for="item of column" :key="item.key" prop="uid" :label="item.name" :min-width="item.width">
        <template #default="scope">
          <img v-if="item.key === 'avatar'" class="avatar" :src="scope.row[item.key]" alt="User avatar" />
          <span v-else-if="item.type === 'text'">{{ scope.row[item.key] }}</span>
          <span v-else-if="item.type === 'radio'">
            {{ item.items.find((i: any) => i.value === scope.row[item.key])?.name }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="110">
        <template #default="scope">
          <el-row>
            <el-button @click="editUser(scope)" link size="small">信息管理</el-button>
          </el-row>
          <el-row>
            <el-button @click="editUserRole(scope)" link size="small">角色管理</el-button>
          </el-row>
          <el-row>
            <el-button @click="openUserSessions(scope)" link size="small">会话管理</el-button>
          </el-row>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="编辑信息" v-model="editedUserVisible">
      <el-form :model="users[editedUser]" label-position="left" label-width="70px">
        <el-form-item v-for="item of canEdit" :key="item.key" :label="item.name">
          <el-input v-if="item.type === 'text'" v-model="users[editedUser][item.key]" />
          <el-switch v-else-if="item.type === 'switch'" v-model="users[editedUser][item.key]" />
          <el-radio-group v-else-if="item.type === 'radio'" v-model="users[editedUser][item.key]">
            <el-radio-button v-for="i of item.items" :key="i.value" :value="i.value">
              {{ i.name }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editedUserVisible = false">取 消</el-button>
          <el-button :loading="editedUserSubmitButtonLoading" type="primary" @click="saveEditedUser">确 定</el-button>
        </span>
      </template>
    </el-dialog>
    <el-dialog :title="'编辑角色: ' + users[editedUser]?.username" v-model="editedUserRoleVisible">
      <el-table v-loading="editedUserRolesLoading" :data="editedUserRoles" stripe style="width: 100%">
        <el-table-column prop="id" label="id" />
        <el-table-column prop="code" label="代码" />
        <el-table-column label="操作" width="60">
          <template #default="scope">
            <el-button @click="deleteUserRole(scope)" link type="danger" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-button class="space-top" type="primary" circle @click="selectRole">
        <Plus />
      </el-button>
    </el-dialog>
    <el-dialog title="选择角色" v-model="allRolesVisible">
      <div v-loading="allRolesLoading">
        <el-select v-model="selectedRole" filterable placeholder="请选择">
          <el-option v-for="role in allRoles" :key="role.id" :label="`[${role.id}] ${role.code}`" :value="role.id">
          </el-option>
        </el-select>
        <el-button v-loading="addRoleButtonLoading" style="float: right" type="success" circle @click="addRoleToUser">
          <Check />
        </el-button>
      </div>
    </el-dialog>

    <el-dialog :title="'会话管理: ' + (users[editedUser]?.username ?? '')" v-model="userSessionsVisible" width="900px">
      <div class="sessions-header">
        <el-button :loading="userSessionsLoading" @click="refreshUserSessions" plain size="small">刷新</el-button>
      </div>
      <el-table v-loading="userSessionsLoading" :data="userSessions" stripe style="width: 100%">
        <el-table-column prop="deviceInfo" label="设备信息" min-width="220">
          <template #default="{ row }">
            <span>{{ row.deviceInfo || "-" }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            <span>{{ formatTime(row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sessionId" label="Session ID" min-width="260">
          <template #default="{ row }">
            <el-tooltip placement="top" :content="row.sessionId">
              <span class="mono ellipsis">{{ row.sessionId }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="90">
          <template #default="scope">
            <el-button
              @click="confirmRevokeUserSession(userSessions[scope.$index]?.sessionId)"
              link
              type="danger"
              size="small"
            >
              吊销
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import Plus from "virtual:icons/mdi/plus"
import Check from "virtual:icons/mdi/check"
import { ElMessage, ElMessageBox } from "element-plus"
import { useUserStore } from "#/shared/stores/user"
import { useNavStore } from "#/shared/stores/nav"
import { storeToRefs } from "pinia"
import { useRoute } from "vue-router"
import { client, buildApiUrl } from "#/shared/lib/api"
import type { components } from "@aqua/http"

type UserSession = components["schemas"]["aqua.passport.v1.UserSession"]

const userStore = useUserStore()
const navStore = useNavStore()
const { searchKeyword } = storeToRefs(navStore)
const $route = useRoute()

const loading = ref(true) // 用户列表加载动画
let page = 1 // 用户列表页码
const pageSize = 10
let end = false // 是否已加载服务器内所有用户
const users = ref([]) // 所有用户
// 编辑用户信息
const editedUser = ref(0) // 被编辑用户索引
const editedUserVisible = ref(false) // 被编辑用户信息框
const editedUserSubmitButtonLoading = ref(false) // 编辑用户信息保存按钮状态控制
// 编辑角色
const editedUserRoleVisible = ref(false) // 被编辑用户的角色编辑框
const editedUserRolesLoading = ref(true) // 被编辑用户的所有角色加载动画
const editedUserRoles = ref([]) // 被编辑用户的角色
// 添加新角色
const allRolesLoading = ref(true) // 所有角色加载动画
const allRolesVisible = ref(false) // 所有角色选取框
const allRoles = ref([]) // 用于添加角色时选取角色
const selectedRole = ref(undefined) // 添加角色时被选中的角色
const addRoleButtonLoading = ref(false)
// 会话管理
const userSessionsVisible = ref(false)
const userSessionsLoading = ref(false)
const userSessions = ref<UserSession[]>([])
let usersTemp: any[] = [] // 暂存用户信息，主界面用于显示搜索结果

const column = computed(() => {
  return userStore.structure.column
})
const canEdit = computed(() => {
  const result: Record<string, any> = {}
  const canEdit = userStore.structure.canEdit.concat(userStore.structure.adminEdit)
  for (const c of column.value) {
    if (canEdit.includes(c.key)) {
      result[c.key] = c
    }
  }
  return result
})

function clearSearch() {
  if (usersTemp.length !== 0) {
    users.value = usersTemp
    usersTemp = []
  }
}
async function handleSearch(keyword: string) {
  const trimmedKeyword = keyword.trim()
  if (trimmedKeyword === "") {
    clearSearch()
  } else {
    loading.value = true
    const { data: res, error } = await client.GET("/api/users/search/{keyword}", {
      baseUrl: buildApiUrl("user"),
      params: { path: { keyword: trimmedKeyword } }
    })

    loading.value = false
    if (!error && res) {
      const found = res.users ?? []
      if (found.length === 0) {
        ElMessage.error("未找到相关内容")
      } else {
        if (usersTemp.length === 0) usersTemp = users.value
        users.value = found as any
      }
    } else {
      ElMessage.error(error?.message ?? "搜索失败")
    }
  }
}
async function getUserInfoFromServer() {
  loading.value = true
  const { data: res, error } = await client.GET("/api/users", {
    baseUrl: buildApiUrl("user"),
    params: { query: { page, pageSize } }
  })

  loading.value = false
  if (!error && res) {
    if (res.users?.length === 0) {
      end = true
    } else {
      users.value = users.value.concat(res.users as any)
      usersTemp = users.value
    }
  } else {
    ElMessage.error(error?.message ?? "获取失败")
  }
}
async function getUserRole(uid: any) {
  editedUserRolesLoading.value = true
  const { data: res, error } = await client.GET("/api/users/{uid}/roles", {
    baseUrl: buildApiUrl("user"),
    params: { path: { uid: Number(uid) } }
  })

  editedUserRolesLoading.value = false
  if (!error && res) {
    editedUserRoles.value = (res.roles as any) ?? []
  } else {
    ElMessage.error(error?.message ?? "获取失败")
  }
}
function editUser(user: any) {
  editedUser.value = user.$index
  editedUserVisible.value = true
}
function editUserRole(user: any) {
  editedUser.value = user.$index
  editedUserRoleVisible.value = true
  editedUserRolesLoading.value = true
  getUserRole(users.value[editedUser.value].uid)
}

function formatTime(value?: string | null) {
  if (!value) return "-"
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toLocaleString()
}

async function refreshUserSessions() {
  const uid = users.value?.[editedUser.value]?.uid
  if (uid === undefined || uid === null) return

  userSessionsLoading.value = true
  const { data: res, error } = await client.GET("/api/users/{uid}/sessions", {
    baseUrl: buildApiUrl("passport"),
    params: { path: { uid: Number(uid) } }
  })
  if (!error && res) {
    userSessions.value = res.sessions ?? []
  } else {
    ElMessage.error(error?.message ?? "获取会话失败")
  }
  userSessionsLoading.value = false
}

function openUserSessions(user: any) {
  editedUser.value = user.$index
  userSessionsVisible.value = true
  refreshUserSessions()
}

async function confirmRevokeUserSession(sessionId?: string) {
  const trimmed = (sessionId ?? "").trim()
  if (!trimmed) return

  try {
    await ElMessageBox.confirm("吊销该会话后，该设备将需要重新登录。是否继续？", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
  } catch {
    return
  }

  const uid = users.value?.[editedUser.value]?.uid
  if (uid === undefined || uid === null) return

  userSessionsLoading.value = true

  const { error } = await client.DELETE("/api/users/{uid}/sessions/{session_id}", {
    baseUrl: buildApiUrl("passport"),
    params: { path: { uid: Number(uid), session_id: trimmed } }
  })

  if (!error) {
    ElMessage.success("已吊销")
    await refreshUserSessions()
  } else {
    ElMessage.error(error?.message ?? "吊销失败")
  }
  userSessionsLoading.value = false
}

async function selectRole() {
  allRolesVisible.value = true
  if (allRoles.value.length > 0) return
  allRolesLoading.value = true

  const { data: res, error } = await client.GET("/api/access-control/roles", {
    baseUrl: buildApiUrl("user")
  })

  allRolesLoading.value = false
  if (!error && res) {
    allRoles.value = (res.roles as any) ?? []
  } else {
    ElMessage.error(error?.message ?? "获取失败")
  }
}

async function addRoleToUser() {
  addRoleButtonLoading.value = true
  const uid = users.value[editedUser.value].uid
  const { error } = await client.POST("/api/users/{uid}/roles", {
    baseUrl: buildApiUrl("user"),
    params: { path: { uid: Number(uid) } },
    body: { id: selectedRole.value } as any
  })

  addRoleButtonLoading.value = false
  if (!error) {
    ElMessage.success("Success")
    getUserRole(uid)
    setTimeout(() => (allRolesVisible.value = false), 100)
  } else {
    ElMessage.error(error?.message ?? "添加失败")
  }
}

function deleteUserRole(role: any) {
  role = editedUserRoles.value[role.$index]
  ElMessageBox.confirm(
    `此操作将永久删除'${users.value[editedUser.value].username}'的角色: ${role.code}, 是否继续?`,
    "提示",
    {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }
  )
    .then(async () => {
      editedUserRolesLoading.value = true
      const uid = users.value[editedUser.value].uid
      const { error } = await client.DELETE("/api/users/{uid}/roles/{id}", {
        baseUrl: buildApiUrl("user"),
        params: { path: { uid: Number(uid), id: Number(role.id) } }
      })

      if (!error) {
        getUserRole(uid)
      } else {
        editedUserRolesLoading.value = false
        ElMessage.error(error?.message ?? "删除失败")
      }
    })
    .catch(() => {})
}

async function saveEditedUser() {
  editedUserSubmitButtonLoading.value = true
  const uid = users.value[editedUser.value].uid

  const { error } = await client.PUT("/api/users/{uid}", {
    baseUrl: buildApiUrl("user"),
    params: { path: { uid: Number(uid) } },
    body: users.value[editedUser.value] as any
  })

  editedUserSubmitButtonLoading.value = false
  if (!error) {
    ElMessage.success("Success")
    editedUserVisible.value = false
  } else {
    ElMessage.error(error?.message ?? "保存失败")
  }
}
function handleScroll() {
  if ($route.path === "/user-manager/info/") {
    const ele = document.querySelector(".router-view")
    if (ele && ele.scrollTop + ele.clientHeight >= ele.scrollHeight) {
      const trimmedKeyword = searchKeyword.value.trim()
      if (!loading.value && !end && trimmedKeyword === "") {
        page++
        getUserInfoFromServer()
      }
    }
  }
}

onMounted(() => {
  getUserInfoFromServer()
  const routerView = document.querySelector(".router-view")
  if (routerView) {
    routerView.addEventListener("scroll", handleScroll)
  }
  navStore.setSearch(true, handleSearch, clearSearch)
})
onUnmounted(() => {
  navStore.setSearch(false)
})
</script>

<style scoped>
.avatar {
  width: 40px;
  height: 40px;
}

.space-top {
  margin-top: 20px;
}

.sessions-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
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
</style>
