<template>
  <Form @submit="handleSubmit" @invalid="handleInvalid" v-slot="{ submitForm, errors, meta, values }">
    <div v-show="false">{{ submitFormRef = submitForm }}</div>
    <div v-show="false">{{ formMeta = meta }}</div>
    <div v-show="false">{{ formErrors = errors }}</div>
    <div v-show="false">{{ formValues = values }}</div>
    <button type="submit" ref="submitButtonRef" style="display: none">Submit</button>
    <div class="login-mode">
      <!--账号密码登陆-->
      <Field name="account" rules="required" v-slot="{ field, errors: fieldErrors }">
        <div class="input-group">
          <label for="account" v-show="false">Account</label>
          <input class="login-input" id="account" placeholder="Account" type="text" v-bind="field" />
          <p class="feedback" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </div>
      </Field>
      <Field name="password" rules="required" v-slot="{ field, errors: fieldErrors }">
        <div class="input-group">
          <label for="password" v-show="false">Password</label>
          <input
            :type="isShowPassword ? 'text' : 'password'"
            class="login-input"
            id="password"
            placeholder="Password"
            v-bind="field"
          />
          <span @click="lockPassword" @mouseout="hidePassword" @mouseover="showPassword" class="eye-box">
            <el-icon :class="{ active: isShowPassword }"><View /></el-icon>
          </span>
          <p class="feedback" v-if="fieldErrors.length > 0">
            {{ fieldErrors[0] }}
          </p>
        </div>
      </Field>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { ref, nextTick } from "vue"
import { View } from "@element-plus/icons-vue"
import { Form, Field } from "vee-validate"

interface Emits {
  (e: "login", data: { account: string; password: string }): void
}

const emit = defineEmits<Emits>()

const eye_open_lock = ref(false)
const isShowPassword = ref(false)
const submitFormRef = ref<(() => void) | null>(null)
const formMeta = ref<any>(null)
const formErrors = ref<any>(null)
const formValues = ref<any>(null)
const submitButtonRef = ref<HTMLButtonElement | null>(null)

function handleSubmit(values: any) {
  emit("login", {
    account: values.account ?? "",
    password: values.password ?? ""
  })
}

function handleInvalid(_values: any, _errors: any) {
  // 表单验证失败时的处理
}

// 暴露方法供父组件调用
defineExpose({
  submit: async () => {
    await nextTick()

    if (submitButtonRef.value) {
      submitButtonRef.value.click()
      return
    }

    if (submitFormRef.value) {
      submitFormRef.value()
      await nextTick()

      setTimeout(async () => {
        await nextTick()
        if (formMeta.value && formMeta.value.valid && formValues.value) {
          handleSubmit(formValues.value)
        }
      }, 100)
    }
  }
})

function showPassword() {
  if (!eye_open_lock.value) isShowPassword.value = true
}

function hidePassword() {
  if (!eye_open_lock.value) isShowPassword.value = false
}

function lockPassword() {
  if (!eye_open_lock.value) {
    eye_open_lock.value = true
    isShowPassword.value = true
  } else {
    eye_open_lock.value = false
    isShowPassword.value = false
  }
}
</script>

<style scoped>
/*输入框*/
.login-input {
  padding: 20px;
  padding-left: 5%;
  padding-right: 10%;
  width: 80%;
  border-radius: 45px;
  margin: 10px 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background-color: rgba(255, 255, 255, 0);
  color: #fff;
  -webkit-text-fill-color: rgba(255, 255, 255, 0.6);
  outline: none;
  transition: 0.3s;
}

.feedback {
  text-align: left;
  color: #ff6f77;
  margin: 0 0 0 15%;
}

.input-group {
  position: relative;
}

/*显示密码按钮*/
.eye-box {
  height: 20px;
  width: 20px;
  line-height: 20px;
  text-align: center;
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 35px;
  cursor: pointer;
  font-size: 18px;
}

.eye-box :deep(.el-icon.active) {
  color: #535858;
}
</style>
