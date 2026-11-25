import { ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY, clearAuthTokens, getCookie } from "@aqua/http"
import { redirectToLogin } from "#/shared/lib/auth"
import { client, buildApiUrl } from "#/shared/lib/api"

export function useSession() {
  function parseSessionIdFromJwt(token: string | null | undefined): string | null {
    if (!token) return null
    const parts = token.split(".")
    if (parts.length < 2) return null
    const payloadPart = parts[1]
    if (!payloadPart) return null
    try {
      const payloadB64 = payloadPart.replace(/-/g, "+").replace(/_/g, "/")
      const padded = payloadB64.padEnd(Math.ceil(payloadB64.length / 4) * 4, "=")
      const json = atob(padded)
      const payload = JSON.parse(json) as { session_id?: unknown }
      const sessionId = typeof payload?.session_id === "string" ? payload.session_id : null
      return sessionId ? sessionId.trim() : null
    } catch {
      return null
    }
  }

  function getCurrentSessionId(): string | null {
    return parseSessionIdFromJwt(getCookie(REFRESH_TOKEN_KEY)) ?? parseSessionIdFromJwt(getCookie(ACCESS_TOKEN_KEY))
  }

  async function logout() {
    const currentSessionId = getCurrentSessionId()
    try {
      if (currentSessionId) {
        await client.DELETE("/api/sessions/{session_id}", {
          baseUrl: buildApiUrl("passport"),
          params: { path: { session_id: currentSessionId } }
        })
      }
    } catch {
      // Ignore API errors and continue clearing local auth state.
    } finally {
      clearAuthTokens()
      redirectToLogin()
    }
  }

  return { logout }
}
