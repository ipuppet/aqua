export const securityConfig = {
  turnstileSiteKey: import.meta.env.VITE_TURNSTILE_SITE_KEY ?? "0x4AAAAAAA05D9IlG8DuPYS2",
  turnstileScriptUrl: "https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit"
} as const
