<template>
  <div class="info-page">
    <el-form label-position="right" label-width="80px" :model="userInfo">
      <div v-if="accountFields.length" class="section">
        <div class="section-title">账号信息</div>
        <el-form-item v-for="item of accountFields" :key="item.key" :label="item.name">
          <el-input v-if="item.type === 'text'" :disabled="item.key === 'username'" v-model="userInfo[item.key]" />
          <el-switch v-else-if="item.type === 'switch'" v-model="userInfo[item.key]" />
          <el-radio-group v-else-if="item.type === 'radio'" v-model="userInfo[item.key]">
            <el-radio-button v-for="i of item.items" :key="i.value" :value="i.value">
              {{ i.name }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
      </div>

      <div v-if="profileFields.length" class="section">
        <div class="section-title">基本资料</div>
        <el-form-item v-for="item of profileFields" :key="item.key" :label="item.name">
          <template v-if="item.key === 'avatar'">
            <div class="avatar-editor">
              <el-avatar v-if="avatarDisplayUrl" :src="avatarDisplayUrl" :size="72" />
              <el-avatar v-else :size="72">?</el-avatar>
              <el-upload :show-file-list="false" :before-upload="beforeAvatarSelect">
                <el-button>选择头像</el-button>
              </el-upload>
            </div>
          </template>
          <el-input v-else-if="item.type === 'text'" v-model="userInfo[item.key]" />
          <el-switch v-else-if="item.type === 'switch'" v-model="userInfo[item.key]" />
          <el-radio-group v-else-if="item.type === 'radio'" v-model="userInfo[item.key]">
            <el-radio-button v-for="i of item.items" :key="i.value" :value="i.value">
              {{ i.name }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item class="submit-row">
          <el-button :loading="buttonLoading" type="primary" @click="onSubmit">保存</el-button>
        </el-form-item>
      </div>

      <div v-if="otherFields.length" class="section">
        <div class="section-title">其他信息</div>
        <el-form-item v-for="item of otherFields" :key="item.key" :label="item.name">
          <el-input v-if="item.type === 'text'" v-model="userInfo[item.key]" />
          <el-switch v-else-if="item.type === 'switch'" v-model="userInfo[item.key]" />
          <el-radio-group v-else-if="item.type === 'radio'" v-model="userInfo[item.key]">
            <el-radio-button v-for="i of item.items" :key="i.value" :value="i.value">
              {{ i.name }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
      </div>

      <div class="section">
        <div class="section-title">安全操作</div>
        <el-form-item label="密码">
          <el-button :loading="forgotPasswordLoading" @click="openResetPasswordDialog">发送重置邮件</el-button>
        </el-form-item>
      </div>
    </el-form>

    <el-dialog v-model="resetPasswordDialogVisible" title="重置密码" width="420px">
      <el-form label-position="right" label-width="90px">
        <el-form-item label="邮箱">
          <el-input v-model="resetPasswordEmail" disabled />
        </el-form-item>
        <el-form-item label="人机验证">
          <div
            class="turnstile-wrapper"
            v-loading="turnstileLoading"
            element-loading-background="rgba(255, 255, 255, 0.7)"
          >
            <div :id="turnstileContainerId" class="turnstile-container"></div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span>
          <el-button @click="resetPasswordDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="forgotPasswordLoading" @click="submitForgotPassword"
            >发送重置邮件</el-button
          >
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus"
import type { UploadProps } from "element-plus"
import { useUserStore } from "#/shared/stores/user"
import { storeToRefs } from "pinia"
import { client, buildApiUrl } from "#/shared/lib/api"
import { sessionData } from "#/shared/lib/storage"
import { useTurnstile } from "#/features/auth/model/useTurnstile"

const userStore = useUserStore()
const { info, structure } = storeToRefs(userStore)

const userInfo = ref<any>({})
const buttonLoading = ref(false)
const forgotPasswordLoading = ref(false)
const resetPasswordDialogVisible = ref(false)
const resetPasswordEmail = ref("")
const pendingAvatarFile = ref<File | null>(null)
const previewAvatarUrl = ref("")
const { turnstileContainerId, turnstileToken, turnstileLoading, resetTurnstile, ensureTurnstileWidget } =
  useTurnstile("me-info-turnstile")

const column = computed(() => structure.value.column)
const canEdit = computed(() => structure.value.canEdit)
const avatarDisplayUrl = computed(() =>
  previewAvatarUrl.value !== "" ? previewAvatarUrl.value : String(userInfo.value?.avatar ?? "")
)
const knownGroupedKeys = ["username", "email", "phone", "avatar", "nickname", "gender"]

function isEditableField(item: Record<string, any>) {
  return item?.key === "username" || canEdit.value.indexOf(item?.key) !== -1
}

const editableFields = computed<Record<string, any>[]>(() =>
  Array.isArray(column.value) ? column.value.filter((item: Record<string, any>) => isEditableField(item)) : []
)

function pickFields(keys: string[]) {
  return editableFields.value.filter(item => keys.includes(String(item.key)))
}

const accountFields = computed(() => pickFields(["username", "email", "phone"]))
const profileFields = computed(() => pickFields(["avatar", "nickname", "gender"]))
const otherFields = computed(() => editableFields.value.filter(item => !knownGroupedKeys.includes(String(item.key))))

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const beforeAvatarSelect: UploadProps["beforeUpload"] = file => {
  if (!file.type.startsWith("image/")) {
    ElMessage.error("只能上传图片文件")
    return false
  }
  if (previewAvatarUrl.value) {
    URL.revokeObjectURL(previewAvatarUrl.value)
  }
  pendingAvatarFile.value = file
  previewAvatarUrl.value = URL.createObjectURL(file)
  return false
}

function openResetPasswordDialog() {
  resetPasswordEmail.value = String(userInfo.value?.email ?? "")
  resetPasswordDialogVisible.value = true
  nextTick(() => {
    ensureTurnstileWidget()
    resetTurnstile()
  })
}

async function submitForgotPassword() {
  const email = resetPasswordEmail.value.trim()
  if (!EMAIL_REGEX.test(email)) {
    ElMessage.error("请输入有效邮箱")
    return
  }
  if (!turnstileToken.value) {
    ElMessage.error("请先完成人机验证")
    return
  }

  forgotPasswordLoading.value = true
  const { error } = await client.POST("/api/password/forgot", {
    baseUrl: buildApiUrl("passport"),
    body: {
      email,
      turnstileToken: turnstileToken.value
    }
  })

  if (!error) {
    ElMessage.success("如果该邮箱已注册，我们已经发送了一封重置密码邮件，请查收。")
    resetPasswordDialogVisible.value = false
  } else {
    ElMessage.error(error?.message ?? "发送失败")
  }
  forgotPasswordLoading.value = false
  resetTurnstile()
}

function clearAvatarPreview() {
  if (previewAvatarUrl.value) {
    URL.revokeObjectURL(previewAvatarUrl.value)
    previewAvatarUrl.value = ""
  }
  pendingAvatarFile.value = null
}

function loadUserInfo() {
  // 从 store 获取用户信息
  userInfo.value = { ...info.value }

  // 如果 store 中没有用户信息，单次从 API 拉取
  if (Object.keys(userInfo.value).length === 0 || !userInfo.value.uid) {
    client
      .GET("/api/users/me", { baseUrl: buildApiUrl("user") })
      .then(({ data: res, error }) => {
        if (!error && res) {
          userInfo.value = res as any
          userStore.setUserInfo(res as any)
          sessionData.set("user.info", res as any)
        }
      })
      .catch(() => {})
  }
}
async function onSubmit() {
  if (
    !Object.prototype.hasOwnProperty.call(userInfo.value, "uid") ||
    !Object.prototype.hasOwnProperty.call(userInfo.value, "nickname") ||
    !Object.prototype.hasOwnProperty.call(userInfo.value, "gender")
  ) {
    ElMessage.error("userInfo incomplete")
    return false
  }
  buttonLoading.value = true
  try {
    const submitData = { ...userInfo.value }
    if (pendingAvatarFile.value) {
      const formData = new FormData()
      formData.append("image", pendingAvatarFile.value)
      const { data: uploadRes, error: uploadError } = await (client as any).POST("/api/file/image/user", {
        body: formData as any
      })
      if (uploadError) throw new Error(uploadError.message ?? "上传失败")
      submitData.avatar = uploadRes?.url ?? submitData.avatar
    }
    const { error: updateError } = await client.PUT("/api/users/me", {
      baseUrl: buildApiUrl("user"),
      body: submitData as any
    })
    if (updateError) throw new Error(updateError.message ?? "修改失败")

    userInfo.value = submitData
    userStore.setUserInfo(submitData)
    sessionData.set("user.info", submitData)
    clearAvatarPreview()
    ElMessage.success("修改成功")
  } catch (error: any) {
    ElMessage.error(error.message ?? "修改失败")
  } finally {
    buttonLoading.value = false
  }
}

onMounted(() => loadUserInfo())
onUnmounted(() => {
  clearAvatarPreview()
})
</script>

<style scoped>
.info-page {
  max-width: 720px;
}

.section {
  margin-bottom: 20px;
  padding: 14px 14px 6px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.section-title {
  margin-bottom: 10px;
  font-size: 14px;
  font-weight: 600;
}

.avatar-editor {
  display: flex;
  align-items: center;
  gap: 12px;
}

.submit-row {
  margin-top: 8px;
}

.turnstile-container {
  min-height: 65px;
}

.turnstile-wrapper {
  position: relative;
  min-height: 65px;
}
</style>
