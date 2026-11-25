import { ElMessage } from "element-plus"
import { buildApiUrl as buildAquaApiUrl, createAquaClient, type AquaService } from "@aqua/http"
import { redirectToLogin } from "#/shared/lib/auth"

const apiHost = import.meta.env.VITE_API ?? ""

export function buildApiUrl(sub: AquaService): string {
  return buildAquaApiUrl(sub, apiHost)
}

export const client = createAquaClient({
  baseUrl: buildApiUrl("admin"),
  passportApiUrl: buildApiUrl("passport"),
  authMode: "refresh",
  cookieDomain: import.meta.env.VITE_AUTH_COOKIE_DOMAIN,
  onAuthFailure: (message: string) => {
    ElMessage.error(message)
    redirectToLogin()
  }
})
