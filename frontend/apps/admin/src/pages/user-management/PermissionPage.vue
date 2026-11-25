<template>
  <div>
    <el-dialog :title="editedPermission.isNew ? '添加权限' : '编辑权限'" v-model="dialogEditorVisible" width="60%">
      <el-form
        status-icon
        label-position="left"
        label-width="120px"
        :model="editedPermission"
        ref="editedPermissionForm"
        :rules="rules"
      >
        <el-form-item label="代码" prop="code">
          <el-input v-model="editedPermission.code" placeholder="权限代码"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="editedPermission.description" placeholder="权限的描述信息"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <span class="dialog-footer-left">
            <el-popconfirm title="确定删除吗？" v-if="!editedPermission.isNew" @confirm="deletePermission">
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
    <Waterfall-layout :max-column="4">
      <card-template v-loading="loading" v-for="permission of permissions" v-bind:key="permission.id">
        <div class="button-editor" @click="edit(permission)">编辑</div>
        <div class="module">{{ permission.id }}: {{ permission.code }}</div>
        <div class="description">描述：&nbsp;&nbsp;{{ permission.description }}</div>
      </card-template>
    </Waterfall-layout>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus"
import { useNavStore } from "#/shared/stores/nav"
import { client, buildApiUrl } from "#/shared/lib/api"

const navStore = useNavStore()

const editedPermissionForm = ref(null)
const loading = ref(true)
const buttons = ref({
  submit: false,
  delete: false
})
const dialogEditorVisible = ref(false)
const permissions = ref([
  {
    id: 0,
    code: "",
    description: ""
  }
])
const editedPermission = ref({
  isNew: true,
  code: "",
  description: ""
})
const rules = ref({
  code: [
    { required: true, message: "请输入权限代码", trigger: "blur" },
    { min: 1, max: 50, message: "长度在 1 到 50 个字符", trigger: "blur" }
  ],
  description: [{ min: 1, max: 300, message: "长度在 1 到 300 个字符", trigger: "blur" }]
})

function edit(permission: Record<string, any> = {}) {
  dialogEditorVisible.value = true
  editedPermissionForm.value?.resetFields()
  permission["isNew"] = permission.id === undefined
  editedPermission.value = JSON.parse(JSON.stringify(permission))
}
function editComplete() {
  // 锁定按钮
  buttons.value.submit = true
  // 验证数据
  editedPermissionForm.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        // 判断新增还是修改
        if (editedPermission.value.isNew) {
          const { id: _id, ...createData } = editedPermission.value
          const { error } = await client.POST("/api/access-control/permissions", {
            baseUrl: buildApiUrl("user"),
            body: createData as any
          })
          if (error) throw new Error(error.message ?? "请求失败")
          ElMessage.success("创建成功")
        } else {
          const { id, ...updateData } = editedPermission.value
          const { error } = await client.PUT("/api/access-control/permissions/{id}", {
            baseUrl: buildApiUrl("user"),
            params: { path: { id: Number(id) } },
            body: updateData as any
          })
          if (error) throw new Error(error.message ?? "请求失败")
          ElMessage.success("更新成功")
        }
        dialogEditorVisible.value = false
        await getPermissionsFromServer()
      } catch (error: any) {
        ElMessage.error(error.message ?? "请求失败")
      } finally {
        buttons.value.submit = false
      }
    } else {
      buttons.value.submit = false
    }
    return false
  })
}
async function deletePermission() {
  buttons.value.delete = true
  try {
    const { error } = await client.DELETE("/api/access-control/permissions/{id}", {
      baseUrl: buildApiUrl("user"),
      params: { path: { id: Number(editedPermission.value.id) } }
    })
    if (error) throw new Error(error.message ?? "请求失败")
    ElMessage.success("删除成功")
    dialogEditorVisible.value = false
    buttons.value.delete = false
    await getPermissionsFromServer()
  } catch (error: any) {
    buttons.value.delete = false
    ElMessage.error(error.message ?? "请求失败")
  }
}
async function getPermissionsFromServer(callback?: () => void) {
  try {
    const { data: res, error } = await client.GET("/api/access-control/permissions", {
      baseUrl: buildApiUrl("user")
    })
    if (error) throw new Error(error.message ?? "获取权限失败")
    permissions.value = (res?.permissions as any) ?? []
    if (typeof callback === "function") callback()
  } catch (error: any) {
    ElMessage.error(error.message ?? "获取权限失败")
  }
}

onMounted(() => {
  getPermissionsFromServer(() => {
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
.module {
  font-size: 16px;
}

.description {
  font-size: 14px;
  color: #adadad;
  margin-top: 10px;
}

.button-editor {
  float: right;
  padding: 5px;
  background-color: #e8e8e8;
  font-size: 12px;
  cursor: pointer;
  border-radius: 5px;
  height: 13px;
  line-height: 13px;
  user-select: none;
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
