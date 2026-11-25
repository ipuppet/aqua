/// <reference types="vite/client" />

declare module "*.vue" {
  import type { DefineComponent } from "vue"
  const component: DefineComponent<object, object, unknown>
  export default component
}

interface ImportMetaEnv {
  readonly VITE_API?: string
  readonly VITE_AUTH_COOKIE_DOMAIN?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
