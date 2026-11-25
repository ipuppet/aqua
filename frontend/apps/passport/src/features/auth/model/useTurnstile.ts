import { ref, onMounted, onUnmounted, nextTick } from "vue"
import { ElMessage } from "element-plus"
import { securityConfig } from "#/shared/config/security"

const TURNSTILE_CONTAINER_ID = "turnstile-widget-container"

interface TurnstileRenderOptions {
  sitekey: string
  callback?: (token: string) => void
  "error-callback"?: () => void
  "expired-callback"?: () => void
  theme?: "light" | "dark" | "auto"
  size?: "normal" | "compact" | "flexible"
}

let scriptLoadingPromise: Promise<void> | null = null
const MAX_RENDER_RETRY_COUNT = 20
const RENDER_RETRY_DELAY_MS = 150
const SCRIPT_LOAD_TIMEOUT_MS = 10000

export function useTurnstile(containerId: string = TURNSTILE_CONTAINER_ID) {
  const turnstileToken = ref("")
  const turnstileWidgetId = ref<string | null>(null)
  const turnstileLoading = ref(true)
  let renderRetryTimer: number | null = null

  function handleSuccess(token: string) {
    turnstileToken.value = token
  }

  function handleError() {
    turnstileToken.value = ""
    ElMessage({ type: "error", message: "人机验证失败，请刷新页面或稍后再试。" })
  }

  function handleExpired() {
    turnstileToken.value = ""
    ElMessage({ type: "warning", message: "人机验证已过期，请重新验证。" })
    if (typeof window !== "undefined" && (window as any).turnstile && turnstileWidgetId.value) {
      ;(window as any).turnstile.reset(turnstileWidgetId.value)
    }
  }

  function clearRenderRetryTimer() {
    if (renderRetryTimer !== null) {
      window.clearTimeout(renderRetryTimer)
      renderRetryTimer = null
    }
  }

  function setTurnstileLoading(loading: boolean) {
    turnstileLoading.value = loading
  }

  function renderWidget(retryCount: number = 0) {
    nextTick(() => {
      if (typeof window === "undefined" || !(window as any).turnstile) return
      const el = document.getElementById(containerId)
      if (turnstileWidgetId.value !== null) {
        setTurnstileLoading(false)
        return
      }
      if (!el) {
        if (retryCount >= MAX_RENDER_RETRY_COUNT) {
          setTurnstileLoading(false)
          return
        }
        clearRenderRetryTimer()
        renderRetryTimer = window.setTimeout(() => {
          renderWidget(retryCount + 1)
        }, RENDER_RETRY_DELAY_MS)
        return
      }

      const options: TurnstileRenderOptions = {
        sitekey: securityConfig.turnstileSiteKey,
        callback: handleSuccess,
        "error-callback": handleError,
        "expired-callback": handleExpired,
        theme: "auto",
        size: "normal"
      }
      const widgetId = (window as any).turnstile.render(
        "#" + containerId,
        options as unknown as Record<string, unknown>
      )
      turnstileWidgetId.value = widgetId
      setTurnstileLoading(false)
      clearRenderRetryTimer()
    })
  }

  function loadScript() {
    if (typeof window !== "undefined" && (window as any).turnstile) {
      return Promise.resolve()
    }

    scriptLoadingPromise ??= new Promise((resolve, reject) => {
      const scriptUrl = securityConfig.turnstileScriptUrl
      const existing = document.querySelector(`script[src="${scriptUrl}"]`)
      const failByTimeout = () => {
        scriptLoadingPromise = null
        reject(new Error("Turnstile script load timeout"))
      }
      const timeoutId = window.setTimeout(failByTimeout, SCRIPT_LOAD_TIMEOUT_MS)

      const onLoaded = () => {
        window.clearTimeout(timeoutId)
        resolve()
      }

      const onFailed = () => {
        window.clearTimeout(timeoutId)
        scriptLoadingPromise = null
        reject(new Error("Turnstile script load failed"))
      }

      if (existing) {
        existing.addEventListener("load", onLoaded, { once: true })
        existing.addEventListener("error", onFailed, { once: true })
        if (typeof window !== "undefined" && (window as any).turnstile) {
          onLoaded()
        }
        return
      }

      const script = document.createElement("script")
      script.src = scriptUrl
      script.async = true
      script.defer = true
      script.onload = onLoaded
      script.onerror = onFailed
      document.head.appendChild(script)
    })

    return scriptLoadingPromise
  }

  function loadScriptAndRender() {
    setTurnstileLoading(true)
    loadScript()
      .then(() => renderWidget())
      .catch(() => {
        setTurnstileLoading(false)
        ElMessage({ type: "error", message: "人机验证加载失败，请稍后重试。" })
      })
  }

  function ensureTurnstileWidget() {
    if (turnstileWidgetId.value !== null) {
      setTurnstileLoading(false)
      return
    }
    setTurnstileLoading(true)
    if (typeof window !== "undefined" && (window as any).turnstile) {
      renderWidget()
      return
    }
    loadScriptAndRender()
  }

  function resetTurnstile() {
    turnstileToken.value = ""
    if (typeof window !== "undefined" && (window as any).turnstile) {
      if (turnstileWidgetId.value !== null) {
        ;(window as any).turnstile.reset(turnstileWidgetId.value)
      } else {
        const el = document.getElementById(containerId)
        if (el) {
          try {
            ;(window as any).turnstile.reset("#" + containerId)
          } catch {
            // ignore
          }
        }
      }
    }
  }

  onMounted(() => {
    loadScriptAndRender()
  })

  onUnmounted(() => {
    clearRenderRetryTimer()
    if (typeof window !== "undefined" && (window as any).turnstile && turnstileWidgetId.value !== null) {
      try {
        ;(window as any).turnstile.remove(turnstileWidgetId.value)
      } catch {
        // ignore
      }
      turnstileWidgetId.value = null
    }
  })

  return {
    turnstileContainerId: containerId,
    turnstileToken,
    turnstileLoading,
    resetTurnstile,
    ensureTurnstileWidget
  }
}
