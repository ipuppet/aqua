/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API?: string
  readonly VITE_AUTH_COOKIE_DOMAIN?: string
  readonly VITE_TURNSTILE_SITE_KEY?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

// Third-party library declarations
declare module "mavon-editor"

// Virtual modules
declare module "virtual:icons/*"

// Augment Vue runtime-core for injected plugin properties
declare module "@vue/runtime-core" {
  interface ComponentCustomProperties {
    $echarts: typeof import("echarts")
    $route: import("vue-router").RouteLocationNormalizedLoaded
    $router: import("vue-router").Router
  }
}
