import { createApp } from "vue"
import { createPinia } from "pinia"
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"

import App from "#/app/App.vue"
import router from "#/app/router"

import "#/shared/lib/validation/vee-validate-rules"
import "#/app/styles/index.scss"

const app = createApp(App)
const pinia = createPinia()

app.use(pinia).use(router).use(ElementPlus)

app.mount("#app")
