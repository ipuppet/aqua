import { createApp } from "vue"
import { createPinia } from "pinia"

import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import "#/app/styles/index.css"

import * as echarts from "echarts"
import { config } from "md-editor-v3"

import App from "#/app/App.vue"
import router from "#/app/router"
import { ElMessage } from "element-plus"

// 禁用短链接插件，解决图片链接无法编辑的问题
// 参考: https://github.com/imzbf/md-editor-v3/issues/972
config({
  codeMirrorExtensions(extensions) {
    // 过滤掉 linkShortener 扩展，防止 URL 被自动缩短
    return extensions.filter(item => item.type !== "linkShortener")
  }
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia).use(router).use(ElementPlus)

app.provide("$echarts", echarts)

app.config.errorHandler = (err: unknown, _instance: unknown, _info: string) => {
  const message = err instanceof Error ? err.message : String(err)
  ElMessage.error(message)
}

app.mount("#app")
