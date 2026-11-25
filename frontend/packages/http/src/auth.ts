import { getCookie, delCookie } from "./cookie"

export const ACCESS_TOKEN_KEY = "access_token"
export const REFRESH_TOKEN_KEY = "refresh_token"

export function hasAuthToken(): boolean {
  const accessToken = getCookie(ACCESS_TOKEN_KEY)
  const refreshToken = getCookie(REFRESH_TOKEN_KEY)
  return Boolean(accessToken?.trim()) || Boolean(refreshToken?.trim())
}

export function clearAuthTokens(domain?: string | null): void {
  delCookie(ACCESS_TOKEN_KEY, domain)
  delCookie(REFRESH_TOKEN_KEY, domain)
}
