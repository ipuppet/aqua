<template>
  <div>
    <div class="login-title">Login to Your Account</div>

    <login-by-password ref="byPasswordRef" v-show="loginMode === 2" @login="handlePasswordLogin" />
    <div v-show="loginMode === -2">
      {{ nickname }}，您已登陆。 <br />如果需要登陆的页面未正常显示您可以尝试刷新那个页面。 <br /><a
        @click="logout"
        class="pointer"
        >Logout</a
      >
    </div>

    <div class="msg-box" :class="messageType === 'error' ? 'text-error' : 'text-info'">
      {{ message }}
    </div>

    <el-row>
      <el-col :span="12">
        <div v-show="loginMode === 2">
          <router-link to="/forgot-password">忘记密码？</router-link>
        </div>
      </el-col>
      <el-col :span="12">
        <button :disabled="loginStatus > 0" @click="login" class="login-btn" v-show="loginMode > 0">
          <span v-show="loginStatus === -1">Login</span>
          <span v-show="loginStatus === 1">
            <el-icon><Check /></el-icon>
          </span>
          <span v-show="loginStatus === -2">Retry</span>
          <span v-show="loginStatus === 2">
            <el-icon class="is-loading"><Loading /></el-icon>
          </span>
        </button>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { Check, Loading } from "@element-plus/icons-vue"
import LoginByPassword from "#/pages/login/components/LoginByPassword.vue"
import { client } from "#/shared/lib/api"
import { ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY, getCookie, setCookie, delCookie } from "@aqua/http"

const message = ref("")
const messageType = ref<"info" | "error">("info")
const nickname = ref("")
// 登陆状态 -1初始状态 -2登陆失败 1登陆成功 2登陆中
const loginStatus = ref(-1)
// 登陆模式 -1初始状态 -2已经登陆 2账号密码模式
const loginMode = ref(-1)

const byPasswordRef = ref<{ submit: () => void } | null>(null)

function decodeJwtPayload(token: string): Record<string, any> | null {
  const parts = token.split(".")
  if (parts.length < 2 || !parts[1]) return null
  try {
    const base64 = parts[1].replace(/-/g, "+").replace(/_/g, "/")
    const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), "=")
    const decoded = decodeURIComponent(
      atob(padded)
        .split("")
        .map(char => `%${char.charCodeAt(0).toString(16).padStart(2, "0")}`)
        .join("")
    )
    const payload = JSON.parse(decoded)
    return payload && typeof payload === "object" ? payload : null
  } catch {
    return null
  }
}

function pickDisplayName(accessToken?: string): string {
  const payload = accessToken ? decodeJwtPayload(accessToken) : null
  return payload?.nk?.trim() ?? ""
}

function getQueryVariable(variable: string): string | false {
  const query = window.location.search.substring(1)
  const vars = query.split("&")
  for (let i = 0; i < vars.length; i++) {
    const item = vars[i]
    if (!item) continue
    const pair = item.split("=")
    if (pair[0] !== undefined && pair[0] === variable && pair[1] !== undefined) {
      return pair[1]
    }
  }
  return false
}

function isAllowedRedirect(url: string): boolean {
  if (url === "false") return false
  if (url.startsWith("/") && !url.startsWith("//")) return true
  try {
    const parsed = new URL(url, window.location.origin)
    if (parsed.protocol !== "https:" && parsed.protocol !== "http:") return false
    const host = parsed.hostname
    return host === "localhost" || host === "example.com" || host.endsWith(".example.com")
  } catch {
    return false
  }
}

function enterKey(event: KeyboardEvent) {
  const code = event.keyCode || event.which || event.charCode
  if (code === 13) {
    login()
  }
}

async function handlePasswordLogin(data: { account: string; password: string }) {
  await performLogin(client.POST("/api/login/password", { body: data }), data)
}

async function performLogin(loginPromise: Promise<any>, _data: any) {
  // 锁定按钮
  message.value = ""
  messageType.value = "info"
  loginStatus.value = 2

  const { data, error } = await loginPromise
  if (!error && data) {
    const accessToken = typeof data?.accessToken === "string" ? data.accessToken : ""
    const refreshToken = typeof data?.refreshToken === "string" ? data.refreshToken : ""
    if (!accessToken || !refreshToken) {
      loginStatus.value = -2
      message.value = "登录成功但未返回完整 token"
      messageType.value = "error"
      return
    }
    loginStatus.value = 1
    setCookie(ACCESS_TOKEN_KEY, accessToken, 1)
    setCookie(REFRESH_TOKEN_KEY, refreshToken, 15)
    nickname.value = pickDisplayName(accessToken)
    // 跳转
    const toUrl = getQueryVariable("from")
    const decoded = toUrl ? decodeURIComponent(toUrl) : ""
    if (decoded && isAllowedRedirect(decoded)) {
      message.value = "正在跳转: " + decoded
      messageType.value = "info"
      setTimeout(() => {
        window.location.href = decoded
      }, 1000)
    } else {
      window.location.reload()
    }
  } else {
    loginStatus.value = -2
    message.value = error?.message ?? "登录失败，请稍后重试"
    messageType.value = "error"
  }
}

function login() {
  // 检查登录模式是否有效
  if (loginMode.value <= 0) {
    message.value = "登录模式未初始化，请刷新页面重试"
    messageType.value = "error"
    return
  }

  // 检查按钮是否被禁用
  if (loginStatus.value > 0) {
    message.value = "登录中，请稍候..."
    messageType.value = "info"
    return
  }

  // 触发对应子组件的表单提交
  if (loginMode.value === 2) {
    if (!byPasswordRef.value) {
      message.value = "密码登录组件未加载，请刷新页面重试"
      messageType.value = "error"
      return
    }
    byPasswordRef.value.submit()
  } else {
    message.value = `未知的登录模式: ${loginMode.value}`
    messageType.value = "error"
  }
}

function logout() {
  delCookie(ACCESS_TOKEN_KEY)
  delCookie(REFRESH_TOKEN_KEY)
  location.reload()
}

onMounted(async () => {
  const accessToken = getCookie(ACCESS_TOKEN_KEY)
  const refreshToken = getCookie(REFRESH_TOKEN_KEY)
  if (accessToken !== null) {
    const { data, error } = await client.POST("/api/token", {
      body: { accessToken, refreshToken: refreshToken ?? undefined }
    })

    if (!error && data) {
      const newAccess = typeof data?.accessToken === "string" ? data.accessToken : ""
      const newRefresh = typeof data?.refreshToken === "string" ? data.refreshToken : ""
      if (newAccess && accessToken !== newAccess) {
        delCookie(ACCESS_TOKEN_KEY)
        setCookie(ACCESS_TOKEN_KEY, newAccess, 1)
      }
      if (newRefresh && refreshToken !== newRefresh) {
        delCookie(REFRESH_TOKEN_KEY)
        setCookie(REFRESH_TOKEN_KEY, newRefresh, 15)
      }
      const tokenForName = [newAccess, accessToken].find(token => typeof token === "string" && token.length > 0) ?? ""
      nickname.value = pickDisplayName(tokenForName)
      loginMode.value = -2
    } else {
      delCookie(ACCESS_TOKEN_KEY)
      delCookie(REFRESH_TOKEN_KEY)
      loginMode.value = 2
      if (error) {
        message.value = "登陆失效：" + (error?.message ?? "请重新登录")
      } else {
        message.value = "登陆失效：返回数据为空"
      }
      messageType.value = "error"
    }
  } else {
    loginMode.value = 2
  }

  document.addEventListener("keyup", enterKey)
})

onUnmounted(() => {
  document.removeEventListener("keyup", enterKey)
})
</script>

<style scoped>
/*Login to Your Account文本*/
.login-title {
  color: #fff;
  margin-bottom: 20px;
  font-size: 20px;
}

/*错误信息提示框*/
.msg-box {
  padding-left: 15px;
  text-align: left;
  margin-bottom: 10px;
  min-height: 20px;
}

.text-error {
  color: #ff6f77;
}

.text-info {
  color: #fff;
}

/*忘记密码*/
.login-forget {
  float: left;
  margin-left: 12%;
}

/*登录按钮*/
.login-btn {
  border-radius: 45px;
  padding: 10px 40px;
  background-color: #ff407c;
  font-size: 14px;
  float: right;
  margin-right: 10px;
}

.pointer {
  cursor: pointer;
}
</style>
