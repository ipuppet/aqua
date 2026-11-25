import createClient, { type Client, type Middleware } from "openapi-fetch"
import type { paths } from "./gen/api"
import { ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY, clearAuthTokens } from "./auth"
import { configureCookies, delCookie, getCookie, setCookie } from "./cookie"

export type AuthMode = "refresh" | "none"

export interface CreateAquaClientOptions {
  baseUrl: string
  /** Absolute base URL for the passport service, required when authMode is "refresh". */
  passportApiUrl?: string
  /**
   * - refresh: on 401, single-flight refresh then retry once (admin)
   * - none: attach token if present; never refresh (passport, blog)
   */
  authMode?: AuthMode
  /**
   * Cookie Domain attribute. Empty/undefined = host-only.
   * Production cross-subdomain: ".example.com"
   */
  cookieDomain?: string
  /** Called once when refresh fails or auth cannot be recovered (authMode=refresh only). */
  onAuthFailure?: (message: string) => void
}

export type AquaClient = Client<paths>

function isTokenRefreshRequest(request: Request): boolean {
  try {
    return new URL(request.url).pathname.endsWith("/api/token")
  } catch {
    return request.url.includes("/api/token")
  }
}

function toErrorMessage(error: unknown): string {
  if (error instanceof Error && error.message) return error.message
  return String(error)
}

async function refreshAccessToken(passportApiUrl: string, cookieDomain?: string): Promise<string> {
  const accessToken = getCookie(ACCESS_TOKEN_KEY)
  const refreshToken = getCookie(REFRESH_TOKEN_KEY)
  if (!refreshToken) {
    throw new Error("Missing refresh token")
  }

  const payload: Record<string, string> = { refreshToken }
  if (accessToken) payload.accessToken = accessToken

  const response = await fetch(`${passportApiUrl}/api/token`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  })

  const body = (await response.json().catch(() => ({}))) as Record<string, unknown>
  if (!response.ok || body.status === false) {
    const message = typeof body.message === "string" ? body.message : "Token refresh failed"
    throw new Error(message)
  }

  const data = ((body.data ?? body) as Record<string, unknown>) ?? {}
  const newAccessToken = data.accessToken
  if (typeof newAccessToken !== "string" || !newAccessToken) {
    throw new Error("No accessToken in refresh response")
  }

  const exp = data.exp
  const expiresInDays = typeof exp === "number" ? Math.max((exp - Math.floor(Date.now() / 1000)) / 86400, 1 / 1440) : 2
  setCookie(ACCESS_TOKEN_KEY, newAccessToken, expiresInDays, cookieDomain)

  const newRefreshToken = data.refreshToken
  if (typeof newRefreshToken === "string" && newRefreshToken) {
    setCookie(REFRESH_TOKEN_KEY, newRefreshToken, 30, cookieDomain)
  }
  return newAccessToken
}

function createAuthMiddleware(options: {
  authMode: AuthMode
  passportApiUrl: string
  cookieDomain?: string
  onAuthFailure?: (message: string) => void
}): Middleware {
  const retryCopies = new WeakMap<Request, Request>()
  let refreshPromise: Promise<string> | null = null
  let hasAuthFailureNotified = false

  return {
    async onRequest({ request }) {
      const accessToken = getCookie(ACCESS_TOKEN_KEY)
      if (accessToken) {
        request.headers.set("Authorization", `Bearer ${accessToken}`)
      }
      retryCopies.set(request, request.clone())
      return request
    },
    async onResponse({ request, response }) {
      if (options.authMode !== "refresh" || response.status !== 401 || isTokenRefreshRequest(request)) {
        return undefined
      }

      try {
        refreshPromise ??= refreshAccessToken(options.passportApiUrl, options.cookieDomain).finally(() => {
          refreshPromise = null
        })
        const newAccessToken = await refreshPromise

        const original = retryCopies.get(request) ?? request
        const retry = new Request(original, { headers: new Headers(original.headers) })
        retry.headers.set("Authorization", `Bearer ${newAccessToken}`)
        return await fetch(retry)
      } catch (refreshError: unknown) {
        const refreshMessage = toErrorMessage(refreshError)
        const isMissingToken = refreshMessage === "Missing refresh token"
        const message = isMissingToken ? "Missing token" : refreshMessage || "登录已过期，请重新登录"

        if (!hasAuthFailureNotified) {
          hasAuthFailureNotified = true
          clearAuthTokens(options.cookieDomain)
          options.onAuthFailure?.(message)
          setTimeout(() => {
            hasAuthFailureNotified = false
          }, 2000)
        } else {
          delCookie(ACCESS_TOKEN_KEY, options.cookieDomain)
        }
        refreshPromise = null
        return response
      }
    }
  }
}

export function createAquaClient(options: CreateAquaClientOptions): AquaClient {
  const authMode: AuthMode = options.authMode ?? "none"
  const cookieDomain = options.cookieDomain || undefined

  if (authMode === "refresh" && !options.passportApiUrl) {
    throw new Error('createAquaClient: passportApiUrl is required when authMode is "refresh"')
  }

  configureCookies({ domain: cookieDomain })

  const client = createClient<paths>({
    baseUrl: options.baseUrl
  })

  client.use(
    createAuthMiddleware({
      authMode,
      passportApiUrl: options.passportApiUrl ?? "",
      cookieDomain,
      onAuthFailure: options.onAuthFailure
    })
  )

  return client
}
