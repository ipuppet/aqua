<template>
  <div>
    <el-dialog :title="editedRole.isNew ? '添加角色' : '编辑角色'" v-model="dialogEditorVisible" width="60%">
      <el-form
        status-icon
        label-position="left"
        label-width="80px"
        :model="editedRole"
        ref="editedRoleForm"
        :rules="rules"
      >
        <el-form-item label="代码" prop="code">
          <el-input v-model="editedRole.code" placeholder="角色代码"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="editedRole.description" placeholder="角色的描述信息"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <span class="dialog-footer-left">
            <el-popconfirm title="确定删除吗？" v-if="!editedRole.isNew" @confirm="deleteRole">
              <template #reference>
                <el-button type="danger" :disabled="buttons.delete" :loading="buttons.delete">删 除</el-button>
              </template>
            </el-popconfirm>
          </span>
          <span class="dialog-footer-right">
            <el-button @click="dialogEditorVisible = false">取 消</el-button>
            <el-button type="primary" @click="editComplete" :disabled="buttons.submit" :loading="buttons.submit">
              确 定
            </el-button>
          </span>
        </span>
      </template>
    </el-dialog>
    <el-dialog :title="'编辑权限: ' + selectedRoleForPerm.code" v-model="permissionDialogVisible">
      <el-table v-loading="permissionDialogLoading" :data="editedRolePermissions" stripe style="width: 100%">
        <el-table-column prop="id" label="id" />
        <el-table-column prop="code" label="代码" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-popconfirm title="确定移除此权限？" @confirm="deleteRolePermission(scope)">
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  :disabled="buttons.deletePermission"
                  :loading="buttons.deletePermission"
                >
                  删 除
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <el-button class="space-top" type="primary" circle @click="openSelectPermission">
        <Plus />
      </el-button>
    </el-dialog>
    <el-dialog title="选择权限" v-model="allPermissionsVisible">
      <div v-loading="allPermissionsLoading">
        <el-select v-model="selectedPermissionId" filterable placeholder="请选择">
          <el-option
            v-for="permission in allPermissions"
            :key="permission.id"
            :label="`[${permission.id}][${permission.code}]: ${permission.description}`"
            :value="permission.id"
          />
        </el-select>
        <el-button
          :loading="buttons.addPermission"
          :disabled="selectedPermissionId === undefined"
          style="float: right"
          type="success"
          circle
          @click="addPermissionToRole"
        >
          <Check />
        </el-button>
      </div>
    </el-dialog>
    <waterfall-layout :max-column="4">
      <card-template v-loading="loading" v-for="role of roles" v-bind:key="role.id">
        <div class="card-header-row">
          <div class="name">{{ role.id }}: {{ role.code }}</div>
          <div class="card-actions">
            <span class="button-editor" @click="edit(role)">编辑</span>
            <span class="button-editor" @click="openPermissionDialog(role)">分配权限</span>
          </div>
        </div>
        <div class="description">描述：&nbsp;&nbsp;{{ role.description }}</div>
      </card-template>
    </waterfall-layout>
  </div>
</template>

<script setup lang="ts">
import Plus from "virtual:icons/mdi/plus"
import Check from "virtual:icons/mdi/check"
import { ElMessage } from "element-plus"
import { useNavStore } from "#/shared/stores/nav"
import { client, buildApiUrl } from "#/shared/lib/api"

const navStore = useNavStore()

const editedRoleForm = ref(null)
const loading = ref(true)
const buttons = ref({
  submit: false,
  delete: false,
  deletePermission: false,
  addPermission: false
})
const dialogEditorVisible = ref(false)
const roles = ref([
  {
    id: 0,
    code: "",
    description: ""
  }
])
const editedRole = ref({
  isNew: true,
  code: "",
  description: ""
})
const rules = ref({
  code: [
    { required: true, message: "请输入角色代码", trigger: "blur" },
    { min: 1, max: 50, message: "长度在 1 到 50 个字符", trigger: "blur" }
  ],
  description: [{ min: 1, max: 300, message: "长度在 1 到 300 个字符", trigger: "blur" }]
})

// 分配权限相关
const selectedRoleForPerm = ref<{ id: number | string; code: string }>({ id: 0, code: "" })
const permissionDialogVisible = ref(false)
const permissionDialogLoading = ref(true)
const editedRolePermissions = ref<Record<string, unknown>[]>([])
const allPermissionsVisible = ref(false)
const allPermissionsLoading = ref(true)
const allPermissions = ref<Record<string, unknown>[]>([])
const selectedPermissionId = ref<number | undefined>(undefined)

function edit(role: Record<string, any> = {}) {
  dialogEditorVisible.value = true
  editedRoleForm.value?.resetFields()
  role["isNew"] = role.id === undefined
  editedRole.value = JSON.parse(JSON.stringify(role))
}

async function editComplete() {
  // 锁定按钮
  buttons.value.submit = true
  // 验证数据
  editedRoleForm.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        // 判断新增还是修改
        if (editedRole.value.isNew) {
          const { id: _id, ...createData } = editedRole.value
          const { error } = await client.POST("/api/access-control/roles", {
            baseUrl: buildApiUrl("user"),
            body: createData as any
          })
          if (error) throw new Error(error.message ?? "请求失败")
          ElMessage.success("创建成功")
        } else {
          const { id, ...updateData } = editedRole.value
          const { error } = await client.PUT("/api/access-control/roles/{id}", {
            baseUrl: buildApiUrl("user"),
            params: { path: { id: Number(id) } },
            body: updateData as any
          })
          if (error) throw new Error(error.message ?? "请求失败")
          ElMessage.success("更新成功")
        }
        dialogEditorVisible.value = false
        await getRolesFromServer()
      } catch (error: any) {
        ElMessage.error(error.message ?? "请求失败")
      } finally {
        buttons.value.submit = false
      }
    }
    return false
  })
}

async function deleteRole() {
  buttons.value.delete = true
  try {
    const { error } = await client.DELETE("/api/access-control/roles/{id}", {
      baseUrl: buildApiUrl("user"),
      params: { path: { id: Number(editedRole.value.id) } }
    })
    if (error) throw new Error(error.message ?? "请求失败")
    ElMessage.success("删除成功")
    dialogEditorVisible.value = false
    await getRolesFromServer()
  } catch (error: any) {
    ElMessage.error(error.message ?? "请求失败")
  } finally {
    buttons.value.delete = false
  }
}

async function getRolesFromServer(callback?: () => void) {
  try {
    const { data: res, error } = await client.GET("/api/access-control/roles", {
      baseUrl: buildApiUrl("user")
    })
    if (error) throw new Error(error.message ?? "获取角色列表失败")
    roles.value = (res?.roles as any) ?? []
    if (typeof callback === "function") callback()
  } catch (error: any) {
    ElMessage.error(error.message ?? "获取角色列表失败")
  }
}

function openPermissionDialog(role: { id: number | string; code: string }) {
  selectedRoleForPerm.value = role
  editedRolePermissions.value = []
  permissionDialogVisible.value = true
  loadRolePermissions()
}

async function loadRolePermissions() {
  permissionDialogLoading.value = true
  const { data: res, error } = await client.GET("/api/access-control/roles/{id}/permissions", {
    baseUrl: buildApiUrl("user"),
    params: { path: { id: Number(selectedRoleForPerm.value.id) } }
  })
  if (!error && res) {
    editedRolePermissions.value = (res.permissions as any) ?? []
    permissionDialogLoading.value = false
  } else {
    permissionDialogLoading.value = false
    ElMessage.error(error?.message ?? "获取权限列表失败")
  }
}

async function openSelectPermission() {
  if (allPermissions.value.length === 0) {
    allPermissionsLoading.value = true
    const { data: res, error } = await client.GET("/api/access-control/permissions", {
      baseUrl: buildApiUrl("user")
    })
    if (!error && res) {
      allPermissions.value = (res.permissions as any) ?? []
      allPermissionsLoading.value = false
      allPermissionsVisible.value = true
    } else {
      allPermissionsLoading.value = false
      ElMessage.error(error?.message ?? "获取权限列表失败")
    }
  } else {
    allPermissionsVisible.value = true
  }
  selectedPermissionId.value = undefined
}

async function addPermissionToRole() {
  if (selectedPermissionId.value === undefined) return
  buttons.value.addPermission = true

  const { error } = await client.POST("/api/access-control/roles/{role_id}/permissions/{permission_id}", {
    baseUrl: buildApiUrl("user"),
    params: {
      path: {
        role_id: Number(selectedRoleForPerm.value.id),
        permission_id: Number(selectedPermissionId.value)
      }
    }
  })

  if (!error) {
    ElMessage.success("添加成功")
    buttons.value.addPermission = false
    allPermissionsVisible.value = false
    loadRolePermissions()
    selectedPermissionId.value = undefined
  } else {
    buttons.value.addPermission = false
    ElMessage.error(error?.message ?? "添加失败")
  }
}

async function deleteRolePermission(scope: { row: { id: number | string } }) {
  buttons.value.deletePermission = true

  const { error } = await client.DELETE("/api/access-control/roles/{role_id}/permissions/{permission_id}", {
    baseUrl: buildApiUrl("user"),
    params: {
      path: {
        role_id: Number(selectedRoleForPerm.value.id),
        permission_id: Number(scope.row.id)
      }
    }
  })

  if (!error) {
    ElMessage.success("移除成功")
    buttons.value.deletePermission = false
    loadRolePermissions()
  } else {
    buttons.value.deletePermission = false
    ElMessage.error(error?.message ?? "移除失败")
  }
}

onMounted(() => {
  getRolesFromServer(() => {
    // 添加nav按钮
    navStore.setNavButton({
      text: "添加",
      click: () => {
        edit()
      }
    })
    // 结束加载动画
    loading.value = false
  })
})

onUnmounted(() => {
  navStore.clearNavButtons()
})
</script>

<style scoped>
.name {
  font-size: 16px;
}

.description {
  font-size: 14px;
  color: #adadad;
  margin-top: 10px;
}

.card-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.name {
  min-width: 0;
  flex: 1;
}

.button-editor {
  padding: 5px;
  background-color: #e8e8e8;
  font-size: 12px;
  cursor: pointer;
  border-radius: 5px;
  height: 13px;
  line-height: 13px;
  user-select: none;
}

.space-top {
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.dialog-footer-right {
  margin-left: auto;
}
</style>
