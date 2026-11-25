import { securityConfig } from "#/shared/config/security"

export const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function validatePassword(password: string, confirmPassword: string): string | null {
  if (password.length < securityConfig.passwordMinLength) {
    return "密码长度至少为6位。"
  }
  if (password.length > securityConfig.passwordMaxLength) {
    return "密码长度不能超过50位。"
  }
  if (password !== confirmPassword) {
    return "两次输入的密码不一致。"
  }
  return null
}

export function validateUsername(username: string): string | null {
  if (username.length < securityConfig.usernameMinLength || username.length > securityConfig.usernameMaxLength) {
    return "用户名长度应在3到20个字符之间。"
  }
  return null
}
