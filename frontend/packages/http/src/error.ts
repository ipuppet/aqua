export function getErrorMessage(error: unknown, fallbackMessage = "请求失败"): string {
  if (typeof error === "string" && error.trim()) return error
  if (error instanceof Error && error.message) return error.message
  if (error && typeof error === "object") {
    const withMessage = error as { message?: unknown }
    if (typeof withMessage.message === "string" && withMessage.message.trim()) {
      return withMessage.message
    }
  }
  return fallbackMessage
}
