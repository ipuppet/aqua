import { createApp } from "vue"
import { createPinia } from "pinia"
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import "md-editor-v3/lib/style.css"
import "#/app/styles/index.scss"
import App from "#/app/App.vue"
import router from "#/app/router"

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus)

app.mount("#app")
