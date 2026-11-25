export const DEV_API_PORTS = {
  root: 8080,
  passport: 8081,
  admin: 8082,
  user: 8083,
  blog: 8084
} as const

export type AquaService = keyof typeof DEV_API_PORTS

function isDevApiHost(apiHost: string): boolean {
  return apiHost.includes("localhost") || apiHost.includes("127.0.0.1") || apiHost.includes("192.168.")
}

function resolveHostname(hostname?: string): string {
  if (hostname) return hostname
  if (typeof window !== "undefined") return window.location.hostname
  return ""
}

/**
 * Build the absolute base URL for a backend service.
 * `apiHost` is typically `import.meta.env.VITE_API` (hostname only, no scheme).
 */
export function buildApiUrl(sub: AquaService, apiHost: string, hostname?: string): string {
  if (isDevApiHost(apiHost)) {
    return `http://${apiHost}:${DEV_API_PORTS[sub]}`
  }

  let finalApi = apiHost
  const host = resolveHostname(hostname)
  if (host) {
    const hostnameParts = host.split(".")
    const eoIndex = hostnameParts.indexOf("eo")
    if (eoIndex !== -1) {
      const apiParts = finalApi.split(".")
      if (!apiParts.includes("eo")) {
        const fromEnd = hostnameParts.length - eoIndex
        const targetIndex = Math.max(0, apiParts.length - (fromEnd - 1))
        apiParts.splice(targetIndex, 0, "eo")
        finalApi = apiParts.join(".")
      }
    }
  }

  return `https://${finalApi}/${sub}`
}
