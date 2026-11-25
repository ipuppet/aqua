<template>
  <div class="form-card">
    <h2>设置新密码</h2>
    <div v-if="!token">
      <p class="error-message">重置链接无效或已过期，请重新发送重置邮件。</p>
      <div class="login-link">
        <router-link to="/forgot-password">重新发送重置邮件</router-link>
      </div>
    </div>
    <Form v-else @submit="handleSubmit">
      <div class="form-group">
        <label for="password">新密码</label>
        <Field name="password" rules="required|password" v-model="password" v-slot="{ field, errors: fieldErrors }">
          <input type="password" id="password" v-bind="field" required />
          <p class="error-message" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </Field>
      </div>
      <div class="form-group">
        <label for="confirmPassword">确认新密码</label>
        <Field
          name="confirmPassword"
          rules="required"
          v-model="confirmPassword"
          v-slot="{ field, errors: fieldErrors }"
        >
          <input type="password" id="confirmPassword" v-bind="field" required />
          <p class="error-message" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </Field>
      </div>

      <button type="submit" class="submit-button" :disabled="isSubmitting">
        {{ isSubmitting ? "提交中..." : "重置密码" }}
      </button>
    </Form>
    <div class="login-link">
      <p><router-link to="/">返回登录</router-link></p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue"
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from "element-plus"
import { Form, Field } from "vee-validate"
import { client } from "#/shared/lib/api"
import { validatePassword } from "#/shared/lib/validation/validation"

const route = useRoute()
const router = useRouter()

const token = computed(() => {
  const value = route.query.token
  return typeof value === "string" ? value : undefined
})

const password = ref("")
const confirmPassword = ref("")
const isSubmitting = ref(false)

async function handleSubmit() {
  if (!token.value) {
    ElMessage({ type: "error", message: "重置链接无效或已过期，请重新发起“忘记密码”操作。" })
    return
  }

  const passwordError = validatePassword(password.value, confirmPassword.value)
  if (passwordError) {
    ElMessage({ type: "error", message: passwordError })
    return
  }
  if (isSubmitting.value) return

  isSubmitting.value = true
  const { error } = await client.POST("/api/password/reset", {
    body: {
      token: token.value,
      password: password.value
    }
  })

  if (!error) {
    ElMessage({ type: "success", message: "密码重置成功，请使用新密码登录。" })
    router.push("/")
  } else {
    ElMessage({
      type: "error",
      message: error?.message ?? "密码重置失败，请稍后再试。"
    })
  }
  isSubmitting.value = false
}
</script>

<style lang="scss" scoped>
@use "#/shared/styles/form-card.scss" as *;
</style>
