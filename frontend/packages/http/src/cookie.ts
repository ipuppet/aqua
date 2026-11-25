export type CookieOptions = {
  /** Empty/undefined = host-only cookie (no Domain attribute). */
  domain?: string
}

let defaultDomain: string | undefined

/**
 * Set the default Domain attribute for subsequent set/del cookie calls.
 * Pass "" or omit to use host-only cookies (local / preview).
 */
export function configureCookies(options: CookieOptions = {}): void {
  defaultDomain = options.domain || undefined
}

export function getCookieDomain(): string | undefined {
  return defaultDomain
}

function resolveDomain(domain?: string | null): string | undefined {
  if (domain === null) return undefined
  if (domain !== undefined) return domain || undefined
  return defaultDomain
}

function isLocalhost(): boolean {
  if (typeof window === "undefined") return true
  const hostname = window.location.hostname
  return hostname === "localhost" || hostname === "127.0.0.1"
}

function securityAttrs(): string {
  return isLocalhost() ? "; SameSite=Lax" : "; Secure; SameSite=Lax"
}

export function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()?.split(";").shift() ?? null
  return null
}

export function setCookie(name: string, value: string, days?: number, domain?: string | null): void {
  if (typeof document === "undefined") return
  let expires = ""
  if (days != null) {
    const date = new Date()
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000)
    expires = `; expires=${date.toUTCString()}`
  }
  const resolved = resolveDomain(domain)
  const domainStr = resolved ? `; domain=${resolved}` : ""
  document.cookie = `${name}=${value}${expires}; path=/${domainStr}${securityAttrs()}`
}

export function delCookie(name: string, domain?: string | null): void {
  if (typeof document === "undefined") return
  const resolved = resolveDomain(domain)
  const domainStr = resolved ? `; domain=${resolved}` : ""
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/${domainStr}${securityAttrs()}`
}
