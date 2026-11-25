<template>
  <div class="form-card">
    <h2>创建账户</h2>
    <Form @submit="handleRegister">
      <div class="form-group">
        <label for="username">用户名</label>
        <Field name="username" rules="required" v-model="username" v-slot="{ field, errors: fieldErrors }">
          <input type="text" id="username" v-bind="field" required />
          <p class="error-message" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </Field>
      </div>
      <div class="form-group">
        <label for="nickname">昵称</label>
        <Field name="nickname" rules="required" v-model="nickname" v-slot="{ field, errors: fieldErrors }">
          <input type="text" id="nickname" v-bind="field" required />
          <p class="error-message" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </Field>
      </div>
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
        <label for="verificationCode">验证码</label>
        <div style="display: flex; align-items: center">
          <Field
            name="emailVerificationCode"
            rules="required"
            v-model="emailVerificationCode"
            v-slot="{ field, errors: fieldErrors }"
          >
            <input type="text" id="verificationCode" v-bind="field" required style="flex-grow: 1; margin-right: 10px" />
            <p class="error-message" v-if="fieldErrors.length > 0">
              {{ fieldErrors[0] }}
            </p>
          </Field>
          <button
            type="button"
            @click="sendVerificationCode"
            class="send-code-button"
            :disabled="isSendingVerificationCode"
          >
            {{ isSendingVerificationCode ? "发送中..." : "发送验证码" }}
          </button>
        </div>
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
      <div class="form-group">
        <label for="sex">性别</label>
        <Field name="sex" rules="required" v-model="sex" v-slot="{ field }">
          <select id="sex" v-bind="field" required>
            <option value="-1">其他</option>
            <option value="0">女</option>
            <option value="1">男</option>
          </select>
        </Field>
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <Field name="password" rules="required|password" v-model="password" v-slot="{ field, errors: fieldErrors }">
          <input type="password" id="password" v-bind="field" required />
          <p class="error-message" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </Field>
      </div>
      <div class="form-group">
        <label for="confirmPassword">确认密码</label>
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
      <button type="submit" class="submit-button">注册</button>
    </Form>
    <div class="login-link">
      <p>已有账户？ <router-link to="/">点击登录</router-link></p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { ElMessage } from "element-plus"
import { Form, Field } from "vee-validate"
import { client } from "#/shared/lib/api"
import { useTurnstile } from "#/features/auth/model/useTurnstile"
import { EMAIL_REGEX, validatePassword, validateUsername } from "#/shared/lib/validation/validation"

const { turnstileContainerId, turnstileToken, turnstileLoading, resetTurnstile } = useTurnstile()

const username = ref("")
const nickname = ref("")
const email = ref("")
const emailVerificationCode = ref("")
const isSendingVerificationCode = ref(false)
const sex = ref("-1")
const password = ref("")
const confirmPassword = ref("")

async function sendVerificationCode() {
  if (isSendingVerificationCode.value) return
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

  isSendingVerificationCode.value = true

  const { error } = await client.POST("/api/verify/email", {
    body: {
      email: email.value,
      turnstileToken: turnstileToken.value
    }
  })

  if (!error) {
    ElMessage({ type: "success", message: "验证码已发送，请查收邮箱。" })
  } else {
    ElMessage({
      type: "error",
      message: error?.message ?? "发送验证码失败，请稍后再试。"
    })
  }
  isSendingVerificationCode.value = false
  resetTurnstile()
}

async function handleRegister() {
  const usernameError = validateUsername(username.value)
  if (usernameError) {
    ElMessage({ type: "error", message: usernameError })
    return
  }
  const passwordError = validatePassword(password.value, confirmPassword.value)
  if (passwordError) {
    ElMessage({ type: "error", message: passwordError })
    return
  }

  const { error } = await client.POST("/api/register", {
    body: {
      username: username.value,
      nickname: nickname.value,
      gender: parseInt(sex.value),
      email: email.value,
      password: password.value,
      verificationCode: emailVerificationCode.value
    }
  })

  if (!error) {
    ElMessage({ type: "success", message: "注册成功！" })
  } else {
    ElMessage({
      type: "error",
      message: error?.message ?? "注册失败，请稍后再试。"
    })
  }
}
</script>

<style lang="scss" scoped>
@use "#/shared/styles/form-card.scss" as *;

form::-webkit-scrollbar {
  width: 8px;
}
form::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}
form::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 10px;
}
form::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}

.send-code-button {
  padding: 12px 15px;
  background-color: #727cf5;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 0.9rem;
  white-space: nowrap;
  transition: background-color 0.3s ease;
}
.send-code-button:hover {
  background-color: #5a63d3;
}
.send-code-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-card .form-group select {
  width: 100%;
  padding: 12px 15px;
  background-color: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 5px;
  box-sizing: border-box;
  color: #fff;
  font-size: 0.9rem;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22rgba(255,255,255,0.7)%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.4-5.4-12.8z%22%2F%3E%3C%2Fsvg%3E");
  background-repeat: no-repeat;
  background-position: right 15px center;
  background-size: 12px;
  padding-right: 40px;
}
.form-card .form-group select option {
  background-color: #5c539d;
  color: #fff;
}

.turnstile-wrapper {
  position: relative;
  min-height: 65px;
}

.turnstile-container {
  min-height: 65px;
}
</style>
