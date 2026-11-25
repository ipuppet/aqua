import { buildApiUrl as buildAquaApiUrl, createAquaClient, type AquaService } from "@aqua/http"

const apiHost = import.meta.env.VITE_API ?? ""

export function buildApiUrl(sub: AquaService): string {
  return buildAquaApiUrl(sub, apiHost)
}

export const client = createAquaClient({
  baseUrl: buildApiUrl("passport"),
  authMode: "none",
  cookieDomain: import.meta.env.VITE_AUTH_COOKIE_DOMAIN
})
