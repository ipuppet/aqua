<template>
  <div class="form-card">
    <h2>重置密码</h2>
    <Form @submit="handleSubmit">
      <div class="form-group">
        <label for="email">邮箱</label>
        <Field name="email" rules="required|email" v-model="email" v-slot="{ field, errors: fieldErrors }">
          <input type="email" id="email" v-bind="field" required />
          <p class="error-message" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </Field>
      </div>

      <div class="form-group">
        <div
          class="turnstile-wrapper"
          v-loading="turnstileLoading"
          element-loading-background="rgba(255, 255, 255, 0.28)"
        >
          <div :id="turnstileContainerId" class="turnstile-container"></div>
        </div>
      </div>

      <button type="submit" class="submit-button" :disabled="isSubmitting">
        {{ isSubmitting ? "发送中..." : "发送重置邮件" }}
      </button>
    </Form>
    <div class="login-link">
      <p>想起密码了？ <router-link to="/">返回登录</router-link></p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { ElMessage } from "element-plus"
import { Form, Field } from "vee-validate"
import { client } from "#/shared/lib/api"
import { useTurnstile } from "#/features/auth/model/useTurnstile"
import { EMAIL_REGEX } from "#/shared/lib/validation/validation"

const { turnstileContainerId, turnstileToken, turnstileLoading, resetTurnstile } = useTurnstile()

const email = ref("")
const isSubmitting = ref(false)

async function handleSubmit() {
  if (isSubmitting.value) return

  if (!email.value) {
    ElMessage({ type: "error", message: "请输入邮箱地址。" })
    return
  }
  if (!EMAIL_REGEX.test(email.value)) {
    ElMessage({ type: "error", message: "请输入有效的邮箱地址。" })
    return
  }
  if (!turnstileToken.value) {
    ElMessage({ type: "error", message: "请先完成人机验证。" })
    return
  }

  isSubmitting.value = true
  const { error } = await client.POST("/api/password/forgot", {
    body: {
      email: email.value,
      turnstileToken: turnstileToken.value
    }
  })

  if (!error) {
    ElMessage({
      type: "success",
      message: "如果该邮箱已注册，我们已经发送了一封重置密码邮件，请查收。"
    })
  } else {
    ElMessage({
      type: "error",
      message: error?.message ?? "发送重置邮件失败，请稍后再试。"
    })
  }

  isSubmitting.value = false
  resetTurnstile()
}
</script>

<style lang="scss" scoped>
@use "#/shared/styles/form-card.scss" as *;

.turnstile-wrapper {
  position: relative;
  min-height: 65px;
}

.turnstile-container {
  min-height: 65px;
}
</style>
