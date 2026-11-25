export const securityConfig = {
  turnstileSiteKey: import.meta.env.VITE_TURNSTILE_SITE_KEY ?? "0x4AAAAAAA05D9IlG8DuPYS2",
  turnstileScriptUrl: "https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit",
  passwordMinLength: 6,
  passwordMaxLength: 50,
  usernameMinLength: 3,
  usernameMaxLength: 20
} as const
